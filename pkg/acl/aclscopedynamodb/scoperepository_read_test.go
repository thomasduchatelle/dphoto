package aclscopedynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamotestutils"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"sort"
	"testing"
	"time"
)

const (
	ironmanEmail = "ironman@stark.com"
	pepperEmail  = "pepperpotts@stark.com"
)

func Test_repository_ListUserScopes(t *testing.T) {
	type args struct {
		email usermodel.UserId
		types []aclcore.ScopeType
	}
	tests := []struct {
		name    string
		args    args
		want    []*aclcore.Scope
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "it should not find any scope when user doesn't exist",
			args:    args{"batman@wayne.com", []aclcore.ScopeType{aclcore.MainOwnerScope}},
			wantErr: assert.NoError,
		},
		{
			name:    "it should not find any scope when user has no grant of that type",
			args:    args{ironmanEmail, []aclcore.ScopeType{aclcore.MediaVisitorScope}},
			wantErr: assert.NoError,
		},
		{
			name:    "it should not find grants for all requested scopes",
			args:    args{ironmanEmail, []aclcore.ScopeType{aclcore.ApiScope, aclcore.MainOwnerScope, aclcore.AlbumVisitorScope, aclcore.MediaVisitorScope}},
			wantErr: assert.NoError,
			want: []*aclcore.Scope{
				{
					Type:          aclcore.MainOwnerScope,
					GrantedAt:     time.Date(2006, 1, 1, 15, 4, 5, 0, time.UTC),
					GrantedTo:     ironmanEmail,
					ResourceOwner: ironmanEmail,
				},
				{
					Type:          aclcore.AlbumVisitorScope,
					GrantedAt:     time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
					GrantedTo:     ironmanEmail,
					ResourceOwner: pepperEmail,
					ResourceId:    "wedding",
					ResourceName:  "Wedding Before EndGame",
				},
				{
					Type:       aclcore.ApiScope,
					GrantedTo:  ironmanEmail,
					GrantedAt:  time.Date(2006, 1, 3, 15, 4, 5, 0, time.UTC),
					ResourceId: "usermanagement",
				},
			},
		},
		{
			name:    "it should not find grants for 2 specific scopes",
			args:    args{ironmanEmail, []aclcore.ScopeType{aclcore.ApiScope, aclcore.MainOwnerScope}},
			wantErr: assert.NoError,
			want: []*aclcore.Scope{
				{
					Type:          aclcore.MainOwnerScope,
					GrantedAt:     time.Date(2006, 1, 1, 15, 4, 5, 0, time.UTC),
					GrantedTo:     ironmanEmail,
					ResourceOwner: ironmanEmail,
				},
				{
					Type:       aclcore.ApiScope,
					GrantedAt:  time.Date(2006, 1, 3, 15, 4, 5, 0, time.UTC),
					GrantedTo:  ironmanEmail,
					ResourceId: "usermanagement",
				},
			},
		},
		{
			name:    "it should not find grants for 2 specific scopes",
			args:    args{ironmanEmail, []aclcore.ScopeType{aclcore.MainOwnerScope}},
			wantErr: assert.NoError,
			want: []*aclcore.Scope{
				{
					Type:          aclcore.MainOwnerScope,
					GrantedAt:     time.Date(2006, 1, 1, 15, 4, 5, 0, time.UTC),
					GrantedTo:     ironmanEmail,
					ResourceOwner: ironmanEmail,
				},
			},
		},
		{
			name:    "it should not find the 'owner:main' scope during authentication [bug: ApiScope being empty, the query stream was returning directly]",
			args:    args{pepperEmail, []aclcore.ScopeType{aclcore.ApiScope, aclcore.MainOwnerScope}},
			wantErr: assert.NoError,
			want: []*aclcore.Scope{
				{
					Type:          aclcore.MainOwnerScope,
					GrantedAt:     time.Date(2006, 1, 5, 15, 4, 5, 0, time.UTC),
					GrantedTo:     pepperEmail,
					ResourceOwner: pepperEmail,
				},
			},
		},
	}

	dyn := dynamotestutils.NewTestContext(context.Background(), t)
	r := Must(New(dyn.Cfg, dyn.Table)).(*repository)

	_, err := r.client.BatchWriteItem(dyn.Ctx, &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			r.table: {
				{
					PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
						"PK":            dynamoutils.AttributeValueMemberS("USER#ironman@stark.com"),
						"SK":            dynamoutils.AttributeValueMemberS("SCOPE#album:visitor#pepperpotts@stark.com#wedding"),
						"Type":          dynamoutils.AttributeValueMemberS("album:visitor"),
						"GrantedAt":     dynamoutils.AttributeValueMemberS("2006-01-02T15:04:05.000000000Z"),
						"GrantedTo":     dynamoutils.AttributeValueMemberS(ironmanEmail),
						"ResourceOwner": dynamoutils.AttributeValueMemberS(pepperEmail),
						"ResourceId":    dynamoutils.AttributeValueMemberS("wedding"),
						"ResourceName":  dynamoutils.AttributeValueMemberS("Wedding Before EndGame"),
					}},
				},
				{
					PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
						"PK":         dynamoutils.AttributeValueMemberS("USER#ironman@stark.com"),
						"SK":         dynamoutils.AttributeValueMemberS("SCOPE#api##usermanagement"),
						"Type":       dynamoutils.AttributeValueMemberS("api"),
						"GrantedAt":  dynamoutils.AttributeValueMemberS("2006-01-03T15:04:05.000000000Z"),
						"GrantedTo":  dynamoutils.AttributeValueMemberS(ironmanEmail),
						"ResourceId": dynamoutils.AttributeValueMemberS("usermanagement"),
					}},
				},
				{
					PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
						"PK":            dynamoutils.AttributeValueMemberS("USER#ironman@stark.com"),
						"SK":            dynamoutils.AttributeValueMemberS("SCOPE#owner:main#ironman@stark.com#"),
						"Type":          dynamoutils.AttributeValueMemberS("owner:main"),
						"GrantedAt":     dynamoutils.AttributeValueMemberS("2006-01-01T15:04:05.000000000Z"),
						"GrantedTo":     dynamoutils.AttributeValueMemberS(ironmanEmail),
						"ResourceOwner": dynamoutils.AttributeValueMemberS(ironmanEmail),
					}},
				},
				{
					PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
						"PK":            dynamoutils.AttributeValueMemberS("USER#pepperpotts@stark.com"),
						"SK":            dynamoutils.AttributeValueMemberS("SCOPE#owner:main#pepperpotts@stark.com#"),
						"Type":          dynamoutils.AttributeValueMemberS("owner:main"),
						"GrantedAt":     dynamoutils.AttributeValueMemberS("2006-01-05T15:04:05.000000000Z"),
						"GrantedTo":     dynamoutils.AttributeValueMemberS(pepperEmail),
						"ResourceOwner": dynamoutils.AttributeValueMemberS(pepperEmail),
					}},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.ListUserScopes(tt.args.email, tt.args.types...)
			if tt.wantErr(t, err) && err == nil {
				sort.Slice(got, func(i, j int) bool {
					return got[i].GrantedAt.Before(got[j].GrantedAt)
				})

				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_repository_ListOwnerScopes(t *testing.T) {
	type args struct {
		email ownermodel.Owner
		types []aclcore.ScopeType
	}
	tests := []struct {
		name    string
		args    args
		want    []*aclcore.Scope
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "it should not find any scope when the owner doesn't exist",
			args:    args{"batman@wayne.com", []aclcore.ScopeType{aclcore.MainOwnerScope}},
			wantErr: assert.NoError,
		},
		{
			name:    "it should not find any scope when user has no grant of that type",
			args:    args{"pepperpotts@stark.com", []aclcore.ScopeType{aclcore.MediaVisitorScope}},
			wantErr: assert.NoError,
		},
		{
			name:    "it should find to whom a user shares its content",
			args:    args{"pepperpotts@stark.com", []aclcore.ScopeType{aclcore.AlbumVisitorScope}},
			wantErr: assert.NoError,
			want: []*aclcore.Scope{
				{
					Type:          aclcore.AlbumVisitorScope,
					GrantedAt:     time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
					GrantedTo:     "ironman@stark.com",
					ResourceOwner: "pepperpotts@stark.com",
					ResourceId:    "wedding",
					ResourceName:  "Wedding Before EndGame",
				},
			},
		},
	}

	dyn := dynamotestutils.NewTestContext(context.Background(), t)
	r := Must(New(dyn.Cfg, dyn.Table)).(*repository)

	_, err := r.client.BatchWriteItem(dyn.Ctx, &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			r.table: {
				{
					PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
						"PK":            dynamoutils.AttributeValueMemberS("USER#ironman@stark.com"),
						"SK":            dynamoutils.AttributeValueMemberS("SCOPE#album:visitor#pepperpotts@stark.com#wedding"),
						"Type":          dynamoutils.AttributeValueMemberS("album:visitor"),
						"GrantedAt":     dynamoutils.AttributeValueMemberS("2006-01-02T15:04:05.000000000Z"),
						"GrantedTo":     dynamoutils.AttributeValueMemberS("ironman@stark.com"),
						"ResourceOwner": dynamoutils.AttributeValueMemberS("pepperpotts@stark.com"),
						"ResourceId":    dynamoutils.AttributeValueMemberS("wedding"),
						"ResourceName":  dynamoutils.AttributeValueMemberS("Wedding Before EndGame"),
					}},
				},
				{
					PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
						"PK":            dynamoutils.AttributeValueMemberS("USER#pepperpotts@stark.com"),
						"SK":            dynamoutils.AttributeValueMemberS("SCOPE#owner:main#pepperpotts@stark.com#"),
						"Type":          dynamoutils.AttributeValueMemberS("owner:main"),
						"GrantedAt":     dynamoutils.AttributeValueMemberS("2006-01-05T15:04:05.000000000Z"),
						"GrantedTo":     dynamoutils.AttributeValueMemberS("pepperpotts@stark.com"),
						"ResourceOwner": dynamoutils.AttributeValueMemberS("pepperpotts@stark.com"),
					}},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.ListOwnerScopes(tt.args.email, tt.args.types...)
			if tt.wantErr(t, err) && err == nil {
				sort.Slice(got, func(i, j int) bool {
					return got[i].GrantedAt.Before(got[j].GrantedAt)
				})

				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_repository_FindScopesById(t *testing.T) {
	tests := []struct {
		name     string
		scopeIds []aclcore.ScopeId
		want     []*aclcore.Scope
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "it should find all grants by their IDs",
			scopeIds: []aclcore.ScopeId{
				{
					Type:          aclcore.MainOwnerScope,
					GrantedTo:     "pepperpotts@stark.com",
					ResourceOwner: "pepperpotts@stark.com",
				},
				{
					Type:          aclcore.AlbumVisitorScope,
					GrantedTo:     "ironman@stark.com",
					ResourceOwner: "pepperpotts@stark.com",
					ResourceId:    "wedding",
				},
			},
			want: []*aclcore.Scope{
				{
					Type:          aclcore.MainOwnerScope,
					GrantedAt:     time.Date(2022, 12, 24, 0, 0, 0, 0, time.UTC),
					GrantedTo:     "pepperpotts@stark.com",
					ResourceOwner: "pepperpotts@stark.com",
				},
				{
					Type:          aclcore.AlbumVisitorScope,
					GrantedAt:     time.Date(2022, 12, 24, 1, 0, 0, 0, time.UTC),
					GrantedTo:     "ironman@stark.com",
					ResourceOwner: "pepperpotts@stark.com",
					ResourceId:    "wedding",
					ResourceName:  "Wedding Before EndGame",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should find a single grant with an ID that doesn't exists",
			scopeIds: []aclcore.ScopeId{
				{
					Type:          aclcore.MainOwnerScope,
					GrantedTo:     "pepperpotts@stark.com",
					ResourceOwner: "pepperpotts@stark.com",
				},
				{
					Type:          aclcore.AlbumVisitorScope,
					GrantedTo:     "bruce@wayne.com",
					ResourceOwner: "joker@?",
				},
			},
			want: []*aclcore.Scope{
				{
					Type:          aclcore.MainOwnerScope,
					GrantedAt:     time.Date(2022, 12, 24, 0, 0, 0, 0, time.UTC),
					GrantedTo:     "pepperpotts@stark.com",
					ResourceOwner: "pepperpotts@stark.com",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should return an empty list when no IDs are found",
			scopeIds: []aclcore.ScopeId{
				{
					Type:          aclcore.AlbumVisitorScope,
					GrantedTo:     "bruce@wayne.com",
					ResourceOwner: "joker@?",
				},
			},
			want:    nil,
			wantErr: assert.NoError,
		},
	}

	dyn := dynamotestutils.NewTestContext(context.Background(), t)
	r := Must(New(dyn.Cfg, dyn.Table)).(*repository)

	_, err := r.client.BatchWriteItem(dyn.Ctx, &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			r.table: {
				{
					PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
						"PK":            dynamoutils.AttributeValueMemberS("USER#pepperpotts@stark.com"),
						"SK":            dynamoutils.AttributeValueMemberS("SCOPE#owner:main#pepperpotts@stark.com#"),
						"Type":          dynamoutils.AttributeValueMemberS("owner:main"),
						"GrantedAt":     dynamoutils.AttributeValueMemberS("2022-12-24T00:00:00.000000000Z"),
						"GrantedTo":     dynamoutils.AttributeValueMemberS("pepperpotts@stark.com"),
						"ResourceOwner": dynamoutils.AttributeValueMemberS("pepperpotts@stark.com"),
					}},
				},
				{
					PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
						"PK":            dynamoutils.AttributeValueMemberS("USER#ironman@stark.com"),
						"SK":            dynamoutils.AttributeValueMemberS("SCOPE#album:visitor#pepperpotts@stark.com#wedding"),
						"Type":          dynamoutils.AttributeValueMemberS("album:visitor"),
						"GrantedAt":     dynamoutils.AttributeValueMemberS("2022-12-24T01:00:00.000000000Z"),
						"GrantedTo":     dynamoutils.AttributeValueMemberS("ironman@stark.com"),
						"ResourceOwner": dynamoutils.AttributeValueMemberS("pepperpotts@stark.com"),
						"ResourceId":    dynamoutils.AttributeValueMemberS("wedding"),
						"ResourceName":  dynamoutils.AttributeValueMemberS("Wedding Before EndGame"),
					}},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.FindScopesById(tt.scopeIds...)
			if tt.wantErr(t, err) && err == nil {
				sort.Slice(got, func(i, j int) bool {
					return got[i].GrantedAt.Before(got[j].GrantedAt)
				})

				assert.Equal(t, tt.want, got)
			}
		})
	}
}
