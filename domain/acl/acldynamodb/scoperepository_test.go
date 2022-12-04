package acldynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/domain/acl/aclcore"
	"sort"
	"testing"
	"time"
)

func awsSession() *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		CredentialsChainVerboseErrors: aws.Bool(true),
		Endpoint:                      aws.String("http://localhost:4566"),
		Credentials:                   credentials.NewStaticCredentials("localstack", "localstack", ""),
		Region:                        aws.String("eu-west-1"),
	}))
}

func Test_repository_ListUserScopes(t *testing.T) {
	type args struct {
		email string
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
			args:    args{"ironman@stark.com", []aclcore.ScopeType{aclcore.MediaVisitorScope}},
			wantErr: assert.NoError,
		},
		{
			name:    "it should not find grants for all requested scopes",
			args:    args{"ironman@stark.com", []aclcore.ScopeType{aclcore.ApiScope, aclcore.MainOwnerScope, aclcore.AlbumVisitorScope, aclcore.MediaVisitorScope}},
			wantErr: assert.NoError,
			want: []*aclcore.Scope{
				{
					Type:          aclcore.MainOwnerScope,
					GrantedAt:     time.Date(2006, 1, 1, 15, 4, 5, 0, time.UTC),
					GrantedTo:     "ironman@stark.com",
					ResourceOwner: "ironman@stark.com",
				},
				{
					Type:          aclcore.AlbumVisitorScope,
					GrantedAt:     time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
					GrantedTo:     "ironman@stark.com",
					ResourceOwner: "pepperpotts@stark.com",
					ResourceId:    "wedding",
					ResourceName:  "Wedding Before EndGame",
				},
				{
					Type:       aclcore.ApiScope,
					GrantedTo:  "ironman@stark.com",
					GrantedAt:  time.Date(2006, 1, 3, 15, 4, 5, 0, time.UTC),
					ResourceId: "usermanagement",
				},
			},
		},
		{
			name:    "it should not find grants for 2 specific scopes",
			args:    args{"ironman@stark.com", []aclcore.ScopeType{aclcore.ApiScope, aclcore.MainOwnerScope}},
			wantErr: assert.NoError,
			want: []*aclcore.Scope{
				{
					Type:          aclcore.MainOwnerScope,
					GrantedAt:     time.Date(2006, 1, 1, 15, 4, 5, 0, time.UTC),
					GrantedTo:     "ironman@stark.com",
					ResourceOwner: "ironman@stark.com",
				},
				{
					Type:       aclcore.ApiScope,
					GrantedAt:  time.Date(2006, 1, 3, 15, 4, 5, 0, time.UTC),
					GrantedTo:  "ironman@stark.com",
					ResourceId: "usermanagement",
				},
			},
		},
		{
			name:    "it should not find grants for 2 specific scopes",
			args:    args{"ironman@stark.com", []aclcore.ScopeType{aclcore.MainOwnerScope}},
			wantErr: assert.NoError,
			want: []*aclcore.Scope{
				{
					Type:          aclcore.MainOwnerScope,
					GrantedAt:     time.Date(2006, 1, 1, 15, 4, 5, 0, time.UTC),
					GrantedTo:     "ironman@stark.com",
					ResourceOwner: "ironman@stark.com",
				},
			},
		},
	}

	r := Must(New(awsSession(), "accesscontroladapter-scoperepository-listuserscopes-"+time.Now().Format("20060102150405.000"), true)).(*repository)
	defer r.db.DeleteTable(&dynamodb.DeleteTableInput{TableName: &r.table})

	_, err := r.db.BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			r.table: {
				{
					PutRequest: &dynamodb.PutRequest{Item: map[string]*dynamodb.AttributeValue{
						"PK":            {S: aws.String("USER#ironman@stark.com")},
						"SK":            {S: aws.String("SCOPE#album:visitor#pepperpotts@stark.com#wedding")},
						"Type":          {S: aws.String("album:visitor")},
						"GrantedAt":     {S: aws.String("2006-01-02T15:04:05.000000000Z")},
						"GrantedTo":     {S: aws.String("ironman@stark.com")},
						"ResourceOwner": {S: aws.String("pepperpotts@stark.com")},
						"ResourceId":    {S: aws.String("wedding")},
						"ResourceName":  {S: aws.String("Wedding Before EndGame")},
					}},
				},
				{
					PutRequest: &dynamodb.PutRequest{Item: map[string]*dynamodb.AttributeValue{
						"PK":         {S: aws.String("USER#ironman@stark.com")},
						"SK":         {S: aws.String("SCOPE#api##usermanagement")},
						"Type":       {S: aws.String("api")},
						"GrantedAt":  {S: aws.String("2006-01-03T15:04:05.000000000Z")},
						"GrantedTo":  {S: aws.String("ironman@stark.com")},
						"ResourceId": {S: aws.String("usermanagement")},
					}},
				},
				{
					PutRequest: &dynamodb.PutRequest{Item: map[string]*dynamodb.AttributeValue{
						"PK":            {S: aws.String("USER#ironman@stark.com")},
						"SK":            {S: aws.String("SCOPE#owner:main#ironman@stark.com#")},
						"Type":          {S: aws.String("owner:main")},
						"GrantedAt":     {S: aws.String("2006-01-01T15:04:05.000000000Z")},
						"GrantedTo":     {S: aws.String("ironman@stark.com")},
						"ResourceOwner": {S: aws.String("ironman@stark.com")},
					}},
				},
				{
					PutRequest: &dynamodb.PutRequest{Item: map[string]*dynamodb.AttributeValue{
						"PK":            {S: aws.String("USER#pepperpotts@stark.com")},
						"SK":            {S: aws.String("SCOPE#owner:main#pepperpotts@stark.com#")},
						"Type":          {S: aws.String("owner:main")},
						"GrantedAt":     {S: aws.String("2006-01-05T15:04:05.000000000Z")},
						"GrantedTo":     {S: aws.String("pepperpotts@stark.com")},
						"ResourceOwner": {S: aws.String("pepperpotts@stark.com")},
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
		email string
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

	r := Must(New(awsSession(), "accesscontroladapter-scoperepository-listuserscopes-"+time.Now().Format("20060102150405.000"), true)).(*repository)
	defer r.db.DeleteTable(&dynamodb.DeleteTableInput{TableName: &r.table})

	_, err := r.db.BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			r.table: {
				{
					PutRequest: &dynamodb.PutRequest{Item: map[string]*dynamodb.AttributeValue{
						"PK":            {S: aws.String("USER#ironman@stark.com")},
						"SK":            {S: aws.String("SCOPE#album:visitor#pepperpotts@stark.com#wedding")},
						"Type":          {S: aws.String("album:visitor")},
						"GrantedAt":     {S: aws.String("2006-01-02T15:04:05.000000000Z")},
						"GrantedTo":     {S: aws.String("ironman@stark.com")},
						"ResourceOwner": {S: aws.String("pepperpotts@stark.com")},
						"ResourceId":    {S: aws.String("wedding")},
						"ResourceName":  {S: aws.String("Wedding Before EndGame")},
					}},
				},
				{
					PutRequest: &dynamodb.PutRequest{Item: map[string]*dynamodb.AttributeValue{
						"PK":            {S: aws.String("USER#pepperpotts@stark.com")},
						"SK":            {S: aws.String("SCOPE#owner:main#pepperpotts@stark.com#")},
						"Type":          {S: aws.String("owner:main")},
						"GrantedAt":     {S: aws.String("2006-01-05T15:04:05.000000000Z")},
						"GrantedTo":     {S: aws.String("pepperpotts@stark.com")},
						"ResourceOwner": {S: aws.String("pepperpotts@stark.com")},
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

	r := Must(New(awsSession(), "accesscontroladapter-scoperepository-findscopesbyid-"+time.Now().Format("20060102150405.000"), true)).(*repository)
	defer r.db.DeleteTable(&dynamodb.DeleteTableInput{TableName: &r.table})

	_, err := r.db.BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			r.table: {
				{
					PutRequest: &dynamodb.PutRequest{Item: map[string]*dynamodb.AttributeValue{
						"PK":            {S: aws.String("USER#pepperpotts@stark.com")},
						"SK":            {S: aws.String("SCOPE#owner:main#pepperpotts@stark.com#")},
						"Type":          {S: aws.String("owner:main")},
						"GrantedAt":     {S: aws.String("2022-12-24T00:00:00.000000000Z")},
						"GrantedTo":     {S: aws.String("pepperpotts@stark.com")},
						"ResourceOwner": {S: aws.String("pepperpotts@stark.com")},
					}},
				},
				{
					PutRequest: &dynamodb.PutRequest{Item: map[string]*dynamodb.AttributeValue{
						"PK":            {S: aws.String("USER#ironman@stark.com")},
						"SK":            {S: aws.String("SCOPE#album:visitor#pepperpotts@stark.com#wedding")},
						"Type":          {S: aws.String("album:visitor")},
						"GrantedAt":     {S: aws.String("2022-12-24T01:00:00.000000000Z")},
						"GrantedTo":     {S: aws.String("ironman@stark.com")},
						"ResourceOwner": {S: aws.String("pepperpotts@stark.com")},
						"ResourceId":    {S: aws.String("wedding")},
						"ResourceName":  {S: aws.String("Wedding Before EndGame")},
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
