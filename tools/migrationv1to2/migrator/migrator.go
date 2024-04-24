package migrator

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	sessionv1 "github.com/aws/aws-sdk-go/aws/session"
	dynamodbv1 "github.com/aws/aws-sdk-go/service/dynamodb"
	dynamodbattributev1 "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/archive"
	"github.com/thomasduchatelle/dphoto/pkg/archiveadapters/archivedynamo"
	"github.com/thomasduchatelle/dphoto/pkg/archiveadapters/asyncjobadapter"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/catalogadapters/catalogdynamo"
	"github.com/thomasduchatelle/dphoto/tools/migrationv1to2/datamodelv1"
	"path"
	"strings"
)

func Migrate(tableName string, arn string, repopulateCache bool) (int, error) {
	log.Infoln("Updating indexes ...")

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return 0, err
	}

	_, err = catalogdynamo.NewRepository(cfg, tableName)
	if err != nil {
		return 0, errors.Wrapf(err, "updating indexes failed ...")
	}

	awsSession := sessionv1.Must(sessionv1.NewSession())
	client := dynamodbv1.New(awsSession)

	log.Infof("Scanning through the all records")

	types := make(map[string]int)
	var patches []*dynamodbv1.WriteRequest

	var imageToResizes []*archive.ImageToResize

	err = client.ScanPages(&dynamodbv1.ScanInput{
		TableName: &tableName,
	}, func(output *dynamodbv1.ScanOutput, _ bool) bool {
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
						if repopulateCache {
							imageToResizes = foundNewLocation(imageToResizes, item)
						}
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

	asyncAdapter := asyncjobadapter.New(cfg, arn, "", 300)
	if len(imageToResizes) > 0 {
		log.Infof("%d medias to miniaturise ...", len(imageToResizes))
		err = asyncAdapter.LoadImagesInCache(imageToResizes...)
		if err != nil {
			return 0, errors.Wrapf(err, "sending SNS messages")
		}
	}

	return total, nil
}

func foundNewLocation(images []*archive.ImageToResize, item map[string]*dynamodbv1.AttributeValue) []*archive.ImageToResize {
	locationKey := *item["LocationKey"].S
	locationId := *item["LocationId"].S

	if !archive.SupportResize(locationKey) {
		return images
	}

	widths := []int{archive.MiniatureCachedWidth}
	if path.Base(locationKey)[:4] == "2022" {
		widths = archive.CacheableWidths
	}

	return append(images, &archive.ImageToResize{
		Owner:    path.Dir(path.Dir(locationKey)),
		MediaId:  locationId,
		StoreKey: locationKey,
		Widths:   widths,
	})
}

func migrateOldAlbum(patches []*dynamodbv1.WriteRequest, item map[string]*dynamodbv1.AttributeValue) ([]*dynamodbv1.WriteRequest, error) {
	owner := *item["PK"].S
	previousKey := map[string]*dynamodbv1.AttributeValue{
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
		&dynamodbv1.WriteRequest{
			PutRequest: &dynamodbv1.PutRequest{
				Item: item,
			},
		},
		&dynamodbv1.WriteRequest{
			DeleteRequest: &dynamodbv1.DeleteRequest{
				Key: previousKey,
			},
		},
	)
	return patches, nil
}

func migrateOldMedia(patches []*dynamodbv1.WriteRequest, item map[string]*dynamodbv1.AttributeValue) ([]*dynamodbv1.WriteRequest, error) {
	data := &datamodelv1.MediaData{}
	err := dynamodbattributev1.UnmarshalMap(item, data)
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
	newItem, err := dynamodbattributev1.MarshalMap(media)
	if err != nil {
		return nil, errors.Wrapf(err, "migrateOldMedia/marshal")
	}

	patches = append(
		patches,
		&dynamodbv1.WriteRequest{
			PutRequest: &dynamodbv1.PutRequest{
				Item: newItem,
			},
		},
		&dynamodbv1.WriteRequest{
			DeleteRequest: &dynamodbv1.DeleteRequest{Key: map[string]*dynamodbv1.AttributeValue{
				"PK": item["PK"],
				"SK": item["SK"],
			}},
		},
	)

	return patches, nil
}

func migrateOldLocation(patches []*dynamodbv1.WriteRequest, item map[string]*dynamodbv1.AttributeValue) ([]*dynamodbv1.WriteRequest, error) {
	data := &datamodelv1.MediaLocationData{}
	err := dynamodbattributev1.UnmarshalMap(item, data)
	if err != nil {
		return nil, errors.Wrapf(err, "migrateOldLocation/unmarshal")
	}

	owner := data.PK[:strings.Index(data.PK, "#")]
	mediaId, _ := catalog.GenerateMediaId(catalog.MediaSignature{
		SignatureSha256: data.SignatureHash,
		SignatureSize:   data.SignatureSize,
	})

	location := archivedynamo.MediaLocationRecord{
		TablePk:           archivedynamo.MediaLocationPk(owner, mediaId),
		LocationKeyPrefix: path.Dir(path.Join(owner, data.FolderName, data.Filename)),
		LocationId:        mediaId,
		LocationKey:       path.Join(owner, data.FolderName, data.Filename),
	}
	newItem, err := dynamodbattributev1.MarshalMap(location)
	if err != nil {
		return nil, errors.Wrapf(err, "migrateOldLocation/marshal")
	}

	patches = append(
		patches,
		&dynamodbv1.WriteRequest{
			PutRequest: &dynamodbv1.PutRequest{
				Item: newItem,
			},
		},
		&dynamodbv1.WriteRequest{
			DeleteRequest: &dynamodbv1.DeleteRequest{Key: map[string]*dynamodbv1.AttributeValue{
				"PK": item["PK"],
				"SK": item["SK"],
			}},
		},
	)

	return patches, nil
}

func stringAttribute(value string) *dynamodbv1.AttributeValue {
	return &dynamodbv1.AttributeValue{
		S: &value,
	}
}
