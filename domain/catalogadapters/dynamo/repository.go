package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Rep struct {
	db                      *dynamodb.DynamoDB
	table                   string
	localDynamodb           bool // localDynamodb is set to true to disable some feature - not available on localstack - like tagging
	findMovedMediaBatchSize int64
}

// NewRepository creates the repository and connect to the database
func NewRepository(awsSession *session.Session, tableName string) (*Rep, error) {
	rep := &Rep{
		db:                      dynamodb.New(awsSession),
		table:                   tableName,
		localDynamodb:           false,
		findMovedMediaBatchSize: int64(dynamoWriteBatchSize),
	}

	err := rep.CreateTableIfNecessary()
	return rep, err
}

// Must panics if there is an error
func Must(rep *Rep, err error) *Rep {
	if err != nil {
		panic(err)
	}

	return rep
}

// CreateTableIfNecessary creates the table if it doesn't exists ; or update it.
func (r *Rep) CreateTableIfNecessary() error {
	mdc := log.WithFields(log.Fields{
		"TableBackup":  r.table,
		"TableVersion": tableVersion,
	})
	mdc.Debugf("CreateTableIfNecessary > describe table '%s'", r.table)

	table, err := r.db.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: &r.table,
	})

	if aerr, ok := err.(awserr.Error); ok && aerr.Code() == dynamodb.ErrCodeResourceNotFoundException {
		table = nil
	} else if err != nil {
		return errors.Wrap(err, "failed to find existing table")
	}

	s := aws.String(dynamodb.ScalarAttributeTypeS)
	createTableInput := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("PK"), AttributeType: s},
			{AttributeName: aws.String("SK"), AttributeType: s},
			{AttributeName: aws.String("AlbumIndexPK"), AttributeType: s},
			{AttributeName: aws.String("AlbumIndexSK"), AttributeType: s},
			{AttributeName: aws.String("MoveTransactionStatus"), AttributeType: s},
			{AttributeName: aws.String("MoveTransaction"), AttributeType: s},
		},
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("AlbumIndex"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{AttributeName: aws.String("AlbumIndexPK"), KeyType: aws.String(dynamodb.KeyTypeHash)},
					{AttributeName: aws.String("AlbumIndexSK"), KeyType: aws.String(dynamodb.KeyTypeRange)},
				},
				Projection: &dynamodb.Projection{ProjectionType: aws.String(dynamodb.ProjectionTypeAll)},
			},
			{
				IndexName: aws.String("MoveTransaction"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{AttributeName: aws.String("MoveTransactionStatus"), KeyType: aws.String(dynamodb.KeyTypeHash)},
					{AttributeName: aws.String("PK"), KeyType: aws.String(dynamodb.KeyTypeRange)},
				},
				Projection: &dynamodb.Projection{ProjectionType: aws.String(dynamodb.ProjectionTypeKeysOnly)},
			},
			{
				IndexName: aws.String("MoveOrder"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{AttributeName: aws.String("MoveTransaction"), KeyType: aws.String(dynamodb.KeyTypeHash)},
					{AttributeName: aws.String("PK"), KeyType: aws.String(dynamodb.KeyTypeRange)},
				},
				Projection: &dynamodb.Projection{ProjectionType: aws.String(dynamodb.ProjectionTypeAll)},
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("PK"), KeyType: aws.String(dynamodb.KeyTypeHash)},
			{AttributeName: aws.String("SK"), KeyType: aws.String(dynamodb.KeyTypeRange)},
		},
		TableName: &r.table,
		Tags: []*dynamodb.Tag{
			{Key: aws.String(versionTagName), Value: aws.String(tableVersion)},
		},
	}

	if table == nil {
		mdc.Infoln("Creating dynamodb table...")
		_, err = r.db.CreateTable(createTableInput)

		return errors.Wrapf(err, "failed to create table %s", r.table)

	} else if r.localDynamodb {
		mdc.Debugln("Local dynamodb update is not supported due to lack of AWS Tag support.")

	} else {
		resource, err := r.db.ListTagsOfResource(&dynamodb.ListTagsOfResourceInput{
			ResourceArn: table.Table.TableArn,
		})
		if err != nil {
			return err
		}

		version := ""
		for _, t := range resource.Tags {
			if *t.Key == versionTagName {
				version = *t.Value
			}
		}

		if version != tableVersion {
			mdc.WithFields(log.Fields{
				"Table":           r.table,
				"Version":         tableVersion,
				"PreviousVersion": version,
			}).Errorln("Dynamodb table exists but must be updated...")
		}
	}

	return nil
}
