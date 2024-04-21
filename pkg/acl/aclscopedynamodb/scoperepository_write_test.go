package aclscopedynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamotestutils"
	dynamoutils "github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutilsv2"
	"testing"
	"time"
)

func Test_repository_DeleteScopes(t *testing.T) {
	ironmanVisitorToPepper := map[string]types.AttributeValue{
		"PK":            dynamoutils.AttributeValueMemberS("USER#ironman@stark.com"),
		"SK":            dynamoutils.AttributeValueMemberS("SCOPE#album:visitor#pepperpotts@stark.com#wedding"),
		"Type":          dynamoutils.AttributeValueMemberS("album:visitor"),
		"GrantedAt":     dynamoutils.AttributeValueMemberS("2006-01-02T15:04:05.000000000Z"),
		"GrantedTo":     dynamoutils.AttributeValueMemberS("ironman@stark.com"),
		"ResourceOwner": dynamoutils.AttributeValueMemberS("pepperpotts@stark.com"),
		"ResourceId":    dynamoutils.AttributeValueMemberS("wedding"),
		"ResourceName":  dynamoutils.AttributeValueMemberS("Wedding Before EndGame"),
	}
	ironmanAOwner := map[string]types.AttributeValue{
		"PK":            dynamoutils.AttributeValueMemberS("USER#ironman@stark.com"),
		"SK":            dynamoutils.AttributeValueMemberS("SCOPE#owner:main#ironman@stark.com#"),
		"Type":          dynamoutils.AttributeValueMemberS("owner:main"),
		"GrantedAt":     dynamoutils.AttributeValueMemberS("2006-01-01T15:04:05.000000000Z"),
		"GrantedTo":     dynamoutils.AttributeValueMemberS("ironman@stark.com"),
		"ResourceOwner": dynamoutils.AttributeValueMemberS("ironman@stark.com"),
	}

	type args struct {
		ids []aclcore.ScopeId
	}
	tests := []struct {
		name        string
		args        args
		givenBefore []map[string]types.AttributeValue
		wantAfter   []map[string]types.AttributeValue
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name:        "it should delete a record that exists, without resource ID",
			givenBefore: []map[string]types.AttributeValue{ironmanAOwner, ironmanVisitorToPepper},
			args: args{
				ids: []aclcore.ScopeId{
					{Type: aclcore.MainOwnerScope, GrantedTo: ironmanEmail, ResourceOwner: ironmanEmail},
				},
			},
			wantAfter: []map[string]types.AttributeValue{ironmanVisitorToPepper},
			wantErr:   assert.NoError,
		},
		{
			name:        "it should delete a record that exists with a resource id",
			givenBefore: []map[string]types.AttributeValue{ironmanAOwner, ironmanVisitorToPepper},
			args: args{
				ids: []aclcore.ScopeId{
					{Type: aclcore.AlbumVisitorScope, GrantedTo: ironmanEmail, ResourceOwner: pepperEmail, ResourceId: "wedding"},
				},
			},
			wantAfter: []map[string]types.AttributeValue{ironmanAOwner},
			wantErr:   assert.NoError,
		},
		{
			name:        "it should not delete a record that doesn't exist",
			givenBefore: []map[string]types.AttributeValue{ironmanAOwner, ironmanVisitorToPepper},
			args: args{
				ids: []aclcore.ScopeId{
					{Type: aclcore.MainOwnerScope, GrantedTo: pepperEmail, ResourceOwner: pepperEmail},
				},
			},
			wantAfter: []map[string]types.AttributeValue{ironmanVisitorToPepper, ironmanAOwner},
			wantErr:   assert.NoError,
		},
	}

	dyn := dynamotestutils.NewTestContext(context.Background(), t)
	repo := Must(New(dyn.Cfg, dyn.Table)).(*repository)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dyn := dyn.Subtest(t)
			dyn.Must(dyn.WithDbContent(dyn.Ctx, tt.givenBefore))

			err := repo.DeleteScopes(tt.args.ids...)
			if tt.wantErr(t, err) && err == nil {
				dyn.MustBool(dyn.EqualContent(dyn.Ctx, tt.wantAfter))
			}
		})
	}
}

func Test_repository_SaveIfNewScope(t *testing.T) {
	ironmanAOwner := map[string]types.AttributeValue{
		"PK":            dynamoutils.AttributeValueMemberS("USER#ironman@stark.com"),
		"SK":            dynamoutils.AttributeValueMemberS("SCOPE#owner:main#ironman@stark.com#"),
		"Type":          dynamoutils.AttributeValueMemberS("owner:main"),
		"GrantedAt":     dynamoutils.AttributeValueMemberS("2006-01-01T15:04:05Z"),
		"GrantedTo":     dynamoutils.AttributeValueMemberS("ironman@stark.com"),
		"ResourceOwner": dynamoutils.AttributeValueMemberS("ironman@stark.com"),
		"ResourceId":    dynamoutils.EmptyString,
		"ResourceName":  dynamoutils.EmptyString,
	}

	type args struct {
		scope aclcore.Scope
	}
	tests := []struct {
		name        string
		args        args
		givenBefore []map[string]types.AttributeValue
		wantAfter   []map[string]types.AttributeValue
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
			wantAfter: []map[string]types.AttributeValue{ironmanAOwner},
			wantErr:   assert.NoError,
		},
		{
			name:        "it should not insert it if it already exists",
			givenBefore: []map[string]types.AttributeValue{ironmanAOwner},
			args: args{scope: aclcore.Scope{
				Type:          aclcore.MainOwnerScope,
				GrantedAt:     time.Now(), // not overridden
				GrantedTo:     ironmanEmail,
				ResourceOwner: ironmanEmail,
				ResourceName:  "something that won't be saved", // not updated
			}},
			wantAfter: []map[string]types.AttributeValue{ironmanAOwner},
			wantErr:   assert.NoError,
		},
		{
			name:        "it should insert it without disturbing another existing scope",
			givenBefore: []map[string]types.AttributeValue{ironmanAOwner},
			args: args{scope: aclcore.Scope{
				Type:          aclcore.MainOwnerScope,
				GrantedAt:     time.Date(2006, 1, 1, 15, 4, 5, 0, time.UTC),
				GrantedTo:     ironmanEmail,
				ResourceOwner: ironmanEmail,
				ResourceId:    ironmanEmail,
			}},
			wantAfter: []map[string]types.AttributeValue{
				ironmanAOwner,
				{
					"PK":            dynamoutils.AttributeValueMemberS("USER#ironman@stark.com"),
					"SK":            dynamoutils.AttributeValueMemberS("SCOPE#owner:main#ironman@stark.com#ironman@stark.com"),
					"Type":          dynamoutils.AttributeValueMemberS("owner:main"),
					"GrantedAt":     dynamoutils.AttributeValueMemberS("2006-01-01T15:04:05Z"),
					"GrantedTo":     dynamoutils.AttributeValueMemberS("ironman@stark.com"),
					"ResourceOwner": dynamoutils.AttributeValueMemberS("ironman@stark.com"),
					"ResourceId":    dynamoutils.AttributeValueMemberS("ironman@stark.com"),
					"ResourceName":  dynamoutils.EmptyString,
				},
			},
			wantErr: assert.NoError,
		},
	}

	dyn := dynamotestutils.NewTestContext(context.Background(), t)
	repo := Must(New(dyn.Cfg, dyn.Table)).(*repository)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dyn := dyn.Subtest(t)
			dyn.Must(dyn.WithDbContent(dyn.Ctx, tt.givenBefore))

			err := repo.SaveIfNewScope(tt.args.scope)
			if tt.wantErr(t, err) && err == nil {
				dyn.MustBool(dyn.EqualContent(dyn.Ctx, tt.wantAfter))
			}
		})
	}
}
