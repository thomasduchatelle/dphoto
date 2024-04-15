package aclscopedynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamotestutils"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
	"testing"
	"time"
)

func Test_repository_DeleteScopes(t *testing.T) {
	ironmanVisitorToPepper := map[string]*dynamodb.AttributeValue{
		"PK":            {S: aws.String("USER#ironman@stark.com")},
		"SK":            {S: aws.String("SCOPE#album:visitor#pepperpotts@stark.com#wedding")},
		"Type":          {S: aws.String("album:visitor")},
		"GrantedAt":     {S: aws.String("2006-01-02T15:04:05.000000000Z")},
		"GrantedTo":     {S: aws.String("ironman@stark.com")},
		"ResourceOwner": {S: aws.String("pepperpotts@stark.com")},
		"ResourceId":    {S: aws.String("wedding")},
		"ResourceName":  {S: aws.String("Wedding Before EndGame")},
	}
	ironmanAOwner := map[string]*dynamodb.AttributeValue{
		"PK":            {S: aws.String("USER#ironman@stark.com")},
		"SK":            {S: aws.String("SCOPE#owner:main#ironman@stark.com#")},
		"Type":          {S: aws.String("owner:main")},
		"GrantedAt":     {S: aws.String("2006-01-01T15:04:05.000000000Z")},
		"GrantedTo":     {S: aws.String("ironman@stark.com")},
		"ResourceOwner": {S: aws.String("ironman@stark.com")},
	}

	type args struct {
		ids []aclcore.ScopeId
	}
	tests := []struct {
		name        string
		args        args
		givenBefore []map[string]*dynamodb.AttributeValue
		wantAfter   []map[string]*dynamodb.AttributeValue
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name:        "it should delete a record that exists, without resource ID",
			givenBefore: []map[string]*dynamodb.AttributeValue{ironmanAOwner, ironmanVisitorToPepper},
			args: args{
				ids: []aclcore.ScopeId{
					{Type: aclcore.MainOwnerScope, GrantedTo: ironmanEmail, ResourceOwner: ironmanEmail},
				},
			},
			wantAfter: []map[string]*dynamodb.AttributeValue{ironmanVisitorToPepper},
			wantErr:   assert.NoError,
		},
		{
			name:        "it should delete a record that exists with a resource id",
			givenBefore: []map[string]*dynamodb.AttributeValue{ironmanAOwner, ironmanVisitorToPepper},
			args: args{
				ids: []aclcore.ScopeId{
					{Type: aclcore.AlbumVisitorScope, GrantedTo: ironmanEmail, ResourceOwner: pepperEmail, ResourceId: "wedding"},
				},
			},
			wantAfter: []map[string]*dynamodb.AttributeValue{ironmanAOwner},
			wantErr:   assert.NoError,
		},
		{
			name:        "it should not delete a record that doesn't exist",
			givenBefore: []map[string]*dynamodb.AttributeValue{ironmanAOwner, ironmanVisitorToPepper},
			args: args{
				ids: []aclcore.ScopeId{
					{Type: aclcore.MainOwnerScope, GrantedTo: pepperEmail, ResourceOwner: pepperEmail},
				},
			},
			wantAfter: []map[string]*dynamodb.AttributeValue{ironmanVisitorToPepper, ironmanAOwner},
			wantErr:   assert.NoError,
		},
	}

	awsSession, _, tableName := dynamotestutils.NewClientV1(t)
	repo := Must(New(awsSession, tableName)).(*repository)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dynamotestutils.SetContent(t, repo.db, repo.table, tt.givenBefore)
			err := repo.DeleteScopes(tt.args.ids...)
			if tt.wantErr(t, err) && err == nil {
				got, err := dynamoutils.AsSlice(dynamoutils.NewScanStream(repo.db, repo.table))
				if assert.NoError(t, err) {
					assert.Equal(t, tt.wantAfter, got)
				}
			}
		})
	}
}

func Test_repository_SaveIfNewScope(t *testing.T) {
	ironmanAOwner := map[string]*dynamodb.AttributeValue{
		"PK":            {S: aws.String("USER#ironman@stark.com")},
		"SK":            {S: aws.String("SCOPE#owner:main#ironman@stark.com#")},
		"Type":          {S: aws.String("owner:main")},
		"GrantedAt":     {S: aws.String("2006-01-01T15:04:05Z")},
		"GrantedTo":     {S: aws.String("ironman@stark.com")},
		"ResourceOwner": {S: aws.String("ironman@stark.com")},
		"ResourceId":    {NULL: aws.Bool(true)},
		"ResourceName":  {NULL: aws.Bool(true)},
	}

	type args struct {
		scope aclcore.Scope
	}
	tests := []struct {
		name        string
		args        args
		givenBefore []map[string]*dynamodb.AttributeValue
		wantAfter   []map[string]*dynamodb.AttributeValue
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name: "it should insert a scope if it doesn't already exist",
			args: args{scope: aclcore.Scope{
				Type:          aclcore.MainOwnerScope,
				GrantedAt:     time.Date(2006, 1, 1, 15, 4, 5, 0, time.UTC),
				GrantedTo:     ironmanEmail,
				ResourceOwner: ironmanEmail,
			}},
			wantAfter: []map[string]*dynamodb.AttributeValue{ironmanAOwner},
			wantErr:   assert.NoError,
		},
		{
			name:        "it should not insert it if it already exists",
			givenBefore: []map[string]*dynamodb.AttributeValue{ironmanAOwner},
			args: args{scope: aclcore.Scope{
				Type:          aclcore.MainOwnerScope,
				GrantedAt:     time.Now(), // not overridden
				GrantedTo:     ironmanEmail,
				ResourceOwner: ironmanEmail,
				ResourceName:  "something that won't be saved", // not updated
			}},
			wantAfter: []map[string]*dynamodb.AttributeValue{ironmanAOwner},
			wantErr:   assert.NoError,
		},
		{
			name:        "it should insert it without disturbing another existing scope",
			givenBefore: []map[string]*dynamodb.AttributeValue{ironmanAOwner},
			args: args{scope: aclcore.Scope{
				Type:          aclcore.MainOwnerScope,
				GrantedAt:     time.Date(2006, 1, 1, 15, 4, 5, 0, time.UTC),
				GrantedTo:     ironmanEmail,
				ResourceOwner: ironmanEmail,
				ResourceId:    ironmanEmail,
			}},
			wantAfter: []map[string]*dynamodb.AttributeValue{
				ironmanAOwner,
				{
					"PK":            {S: aws.String("USER#ironman@stark.com")},
					"SK":            {S: aws.String("SCOPE#owner:main#ironman@stark.com#ironman@stark.com")},
					"Type":          {S: aws.String("owner:main")},
					"GrantedAt":     {S: aws.String("2006-01-01T15:04:05Z")},
					"GrantedTo":     {S: aws.String("ironman@stark.com")},
					"ResourceOwner": {S: aws.String("ironman@stark.com")},
					"ResourceId":    {S: aws.String("ironman@stark.com")},
					"ResourceName":  {NULL: aws.Bool(true)},
				},
			},
			wantErr: assert.NoError,
		},
	}

	awsSession, _, tableName := dynamotestutils.NewClientV1(t)
	repo := Must(New(awsSession, tableName)).(*repository)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dynamotestutils.SetContent(t, repo.db, repo.table, tt.givenBefore)

			err := repo.SaveIfNewScope(tt.args.scope)
			if tt.wantErr(t, err) && err == nil {
				got, err := dynamoutils.AsSlice(dynamoutils.NewScanStream(repo.db, repo.table))
				if assert.NoError(t, err) {
					assert.Equal(t, tt.wantAfter, got)
				}
			}
		})
	}
}
