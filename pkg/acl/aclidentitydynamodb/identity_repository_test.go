package aclidentitydynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	dynamotestutils "github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamotestutils"
	dynamoutils2 "github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
	"testing"
)

func awsSession() *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		CredentialsChainVerboseErrors: aws.Bool(true),
		Endpoint:                      aws.String("http://localhost:4566"),
		Credentials:                   credentials.NewStaticCredentials("localstack", "localstack", ""),
		Region:                        aws.String("eu-west-1"),
	}))
}

func Test_repository_FindIdentity(t *testing.T) {
	sess, db, table := dynamotestutils.NewDbContext(t)
	r := Must(New(sess, table)).(*repository)

	dynamotestutils.SetContent(t, db, table, []map[string]*dynamodb.AttributeValue{
		{
			"PK":      {S: aws.String("USER#tony@stark.com")},
			"SK":      {S: aws.String("IDENTITY#")},
			"Email":   {S: aws.String("tony+other@stark.com")},
			"Name":    {S: aws.String("Tony Stark")},
			"Picture": {S: aws.String("/you/know/me.jpg")},
		},
	})

	type args struct {
		email string
	}
	tests := []struct {
		name    string
		args    args
		want    *aclcore.Identity
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should find identity that exists",
			args: args{"tony@stark.com"},
			want: &aclcore.Identity{
				Email:   "tony+other@stark.com",
				Name:    "Tony Stark",
				Picture: "/you/know/me.jpg",
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should throw a Not found if identity doesn't exist",
			args: args{"pepper@stark.com"},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, aclcore.IdentityDetailsNotFoundError, i)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.FindIdentity(tt.args.email)
			if tt.wantErr(t, err) && err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_repository_StoreIdentity(t *testing.T) {
	sess, _, table := dynamotestutils.NewDbContext(t)
	r := Must(New(sess, table)).(*repository)

	type args struct {
		identity aclcore.Identity
	}
	tests := []struct {
		name      string
		args      args
		wantErr   assert.ErrorAssertionFunc
		wantAfter []map[string]*dynamodb.AttributeValue
	}{
		{
			name: "it should create a brand-new identity details",
			args: args{identity: aclcore.Identity{
				Email:   "pepper@stark.com",
				Name:    "Pepper",
				Picture: "/pepper.jpg",
			}},
			wantErr: assert.NoError,
			wantAfter: []map[string]*dynamodb.AttributeValue{
				{
					"PK":      {S: aws.String("USER#tony@stark.com")},
					"SK":      {S: aws.String("IDENTITY#")},
					"Email":   {S: aws.String("tony+other@stark.com")},
					"Name":    {S: aws.String("Tony Stark")},
					"Picture": {S: aws.String("/you/know/me.jpg")},
				},
				{
					"PK":      {S: aws.String("USER#pepper@stark.com")},
					"SK":      {S: aws.String("IDENTITY#")},
					"Email":   {S: aws.String("pepper@stark.com")},
					"Name":    {S: aws.String("Pepper")},
					"Picture": {S: aws.String("/pepper.jpg")},
				},
			},
		},
		{
			name: "it should override an existing identity",
			args: args{identity: aclcore.Identity{
				Email:   "tony@stark.com",
				Name:    "Ironman",
				Picture: "/ironman-3.jpg",
			}},
			wantErr: assert.NoError,
			wantAfter: []map[string]*dynamodb.AttributeValue{
				{
					"PK":      {S: aws.String("USER#tony@stark.com")},
					"SK":      {S: aws.String("IDENTITY#")},
					"Email":   {S: aws.String("tony@stark.com")},
					"Name":    {S: aws.String("Ironman")},
					"Picture": {S: aws.String("/ironman-3.jpg")},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dynamotestutils.SetContent(t, r.db, table, []map[string]*dynamodb.AttributeValue{
				{
					"PK":      {S: aws.String("USER#tony@stark.com")},
					"SK":      {S: aws.String("IDENTITY#")},
					"Email":   {S: aws.String("tony+other@stark.com")},
					"Name":    {S: aws.String("Tony Stark")},
					"Picture": {S: aws.String("/you/know/me.jpg")},
				},
			})

			err := r.StoreIdentity(tt.args.identity)
			if tt.wantErr(t, err) && err == nil {
				after, err := dynamoutils2.AsSlice(dynamoutils2.NewScanStream(r.db, r.table))
				if assert.NoError(t, err) {
					assert.Equal(t, tt.wantAfter, after)
				}
			}
		})
	}
}
