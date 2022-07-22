package migrator

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/domain/archiveadapters/archivedynamo"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"github.com/thomasduchatelle/dphoto/domain/catalogadapters/catalogdynamo"
	"github.com/thomasduchatelle/dphoto/domain/catalogadapters/dynamoutils"
	"github.com/thomasduchatelle/dphoto/tools/migrationv1to2/datamodelv1"
	"path"
	"strings"
)

func Migrate(tableName string) (int, error) {
	log.Infoln("Updating indexes ...")
	awsSession := session.Must(session.NewSession())
	_, err := catalogdynamo.NewRepository(awsSession, tableName)
	if err != nil {
		return 0, errors.Wrapf(err, "updating indexes failed ...")
	}

	client := dynamodb.New(awsSession)

	log.Infof("Scanning through the all records")

	types := make(map[string]int)
	var patches []*dynamodb.WriteRequest

	err = client.ScanPages(&dynamodb.ScanInput{
		TableName: &tableName,
	}, func(output *dynamodb.ScanOutput, _ bool) bool {
		for _, item := range output.Items {
			pk := *item["PK"].S
			if idx := strings.Index(pk, "#"); idx > 0 {
				if pk[idx+1:] == "ALBUM" {
					types["ALBUM_NEW"]++

				} else {
					sk := *item["SK"].S

					if sk == "#METADATA" {
						if strings.Index(strings.TrimPrefix(pk[idx+1:], "MEDIA#"), "#") > 0 {
							types["MEDIA_OLD"]++
							patches, err = migrateOldMedia(patches, item)
						} else {
							types["MEDIA_NEW"]++
						}

					} else if sk == "LOCATION" {
						types["LOCATION_OLD"]++
						patches, err = migrateOldLocation(patches, item)
					} else if sk == "LOCATION#" {
						types["LOCATION_NEW"]++
					} else {
						types["UNKNOWN"]++
					}
				}
			} else {
				types["ALBUM_OLD"]++
				patches, err = migrateOldAlbum(patches, item)
			}
		}

		return true
	})
	if err != nil {
		return 0, err
	}

	log.Infof("Types count: %+v\n", types)

	if len(patches) > 0 {
		log.Infof("Running %d updates ... ", len(patches))
		err = dynamoutils.BufferedWriteItems(client, patches, tableName, dynamoutils.DynamoWriteBatchSize)
		if err != nil {
			return 0, err
		}
	} else {
		log.Infof("Nothing to migrate, everything is already up to date.")
	}

	total := 0
	for _, count := range types {
		total += count
	}

	return total, nil
}

func migrateOldAlbum(patches []*dynamodb.WriteRequest, item map[string]*dynamodb.AttributeValue) ([]*dynamodb.WriteRequest, error) {
	owner := *item["PK"].S
	previousKey := map[string]*dynamodb.AttributeValue{
		"PK": item["PK"],
		"SK": item["SK"],
	}

	tablePK := catalogdynamo.AlbumPrimaryKey(owner, *item["AlbumFolderName"].S)
	tableAlbumIndex := catalogdynamo.AlbumIndexedKey(owner, *item["AlbumFolderName"].S)

	item["PK"] = stringAttribute(tablePK.PK)
	item["SK"] = stringAttribute(tablePK.SK)
	item["AlbumIndexPK"] = stringAttribute(tableAlbumIndex.AlbumIndexPK)
	item["AlbumIndexSK"] = stringAttribute(tableAlbumIndex.AlbumIndexSK)

	patches = append(
		patches,
		&dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: item,
			},
		},
		&dynamodb.WriteRequest{
			DeleteRequest: &dynamodb.DeleteRequest{
				Key: previousKey,
			},
		},
	)
	return patches, nil
}

func migrateOldMedia(patches []*dynamodb.WriteRequest, item map[string]*dynamodb.AttributeValue) ([]*dynamodb.WriteRequest, error) {
	data := &datamodelv1.MediaData{}
	err := dynamodbattribute.UnmarshalMap(item, data)
	if err != nil {
		return nil, errors.Wrapf(err, "migrateOldMedia/unmarshal")
	}

	owner := data.PK[:strings.Index(data.PK, "#")]
	folderName := data.AlbumIndexPK[strings.Index(data.AlbumIndexPK, "#")+1:]

	mediaId, _ := catalog.GenerateMediaId(catalog.MediaSignature{
		SignatureSha256: data.SignatureHash,
		SignatureSize:   data.SignatureSize,
	})

	media := &catalogdynamo.MediaRecord{
		TablePk:       catalogdynamo.MediaPrimaryKey(owner, mediaId),
		AlbumIndexKey: catalogdynamo.MediaAlbumIndexedKey(owner, folderName, data.DateTime, mediaId),
		Id:            mediaId,
		Type:          data.Type,
		DateTime:      data.DateTime,
		Details:       data.Details,
		Filename:      data.Filename,
		SignatureSize: data.SignatureSize,
		SignatureHash: data.SignatureHash,
	}
	newItem, err := dynamodbattribute.MarshalMap(media)
	if err != nil {
		return nil, errors.Wrapf(err, "migrateOldMedia/marshal")
	}

	patches = append(
		patches,
		&dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: newItem,
			},
		},
		&dynamodb.WriteRequest{
			DeleteRequest: &dynamodb.DeleteRequest{Key: map[string]*dynamodb.AttributeValue{
				"PK": item["PK"],
				"SK": item["SK"],
			}},
		},
	)

	return patches, nil
}

func migrateOldLocation(patches []*dynamodb.WriteRequest, item map[string]*dynamodb.AttributeValue) ([]*dynamodb.WriteRequest, error) {
	data := &datamodelv1.MediaLocationData{}
	err := dynamodbattribute.UnmarshalMap(item, data)
	if err != nil {
		return nil, errors.Wrapf(err, "migrateOldLocation/unmarshal")
	}

	owner := data.PK[:strings.Index(data.PK, "#")]
	mediaId, _ := catalog.GenerateMediaId(catalog.MediaSignature{
		SignatureSha256: data.SignatureHash,
		SignatureSize:   data.SignatureSize,
	})

	location := archivedynamo.MediaLocation{
		TablePk: archivedynamo.MediaLocationPk(owner, mediaId),
		Id:      mediaId,
		Key:     path.Join(owner, data.FolderName, data.Filename),
	}
	newItem, err := dynamodbattribute.MarshalMap(location)
	if err != nil {
		return nil, errors.Wrapf(err, "migrateOldLocation/marshal")
	}

	patches = append(
		patches,
		&dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: newItem,
			},
		},
		&dynamodb.WriteRequest{
			DeleteRequest: &dynamodb.DeleteRequest{Key: map[string]*dynamodb.AttributeValue{
				"PK": item["PK"],
				"SK": item["SK"],
			}},
		},
	)

	return patches, nil
}

func stringAttribute(value string) *dynamodb.AttributeValue {
	return &dynamodb.AttributeValue{
		S: &value,
	}
}
