package aclidentitydynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamotestutils"
	dynamoutils "github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutilsv2"
	"sort"
	"testing"
)

func Test_repository_FindIdentity(t *testing.T) {
	ctx := context.Background()
	dyn := dynamotestutils.NewTestContext(ctx, t)
	r := Must(New(dyn.Cfg, dyn.Table))

	dyn.Must(dyn.WithDbContent(ctx, []map[string]types.AttributeValue{
		{
			"PK":      dynamoutils.AttributeValueMemberS("USER#tony@stark.com"),
			"SK":      dynamoutils.AttributeValueMemberS("IDENTITY#"),
			"Email":   dynamoutils.AttributeValueMemberS("tony+other@stark.com"),
			"Name":    dynamoutils.AttributeValueMemberS("Tony Stark"),
			"Picture": dynamoutils.AttributeValueMemberS("/you/know/me.jpg"),
		},
	}))

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

func Test_repository_FindIdentities(t *testing.T) {
	ctx := context.Background()
	dyn := dynamotestutils.NewTestContext(ctx, t)
	r := Must(New(dyn.Cfg, dyn.Table))

	dyn.Must(dyn.WithDbContent(ctx, []map[string]types.AttributeValue{
		{
			"PK":      dynamoutils.AttributeValueMemberS("USER#tony@stark.com"),
			"SK":      dynamoutils.AttributeValueMemberS("IDENTITY#"),
			"Email":   dynamoutils.AttributeValueMemberS("tony@stark.com"),
			"Name":    dynamoutils.AttributeValueMemberS("Tony Stark"),
			"Picture": dynamoutils.AttributeValueMemberS("/you/know/me.jpg"),
		},
		{
			"PK":      dynamoutils.AttributeValueMemberS("USER#natasha@banner.com"),
			"SK":      dynamoutils.AttributeValueMemberS("IDENTITY#"),
			"Email":   dynamoutils.AttributeValueMemberS("natasha@banner.com"),
			"Name":    dynamoutils.AttributeValueMemberS("Natasha"),
			"Picture": dynamoutils.AttributeValueMemberS("/black-widow.jpg"),
		},
	}))

	type args struct {
		email []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*aclcore.Identity
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should de-duplicate input list of emails",
			args: args{[]string{"tony@stark.com", "tony@stark.com"}},
			want: []*aclcore.Identity{
				{Email: "tony@stark.com", Name: "Tony Stark", Picture: "/you/know/me.jpg"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should find several identities",
			args: args{[]string{"tony@stark.com", "natasha@banner.com"}},
			want: []*aclcore.Identity{
				{Email: "natasha@banner.com", Name: "Natasha", Picture: "/black-widow.jpg"},
				{Email: "tony@stark.com", Name: "Tony Stark", Picture: "/you/know/me.jpg"},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "it should return an empty list when not found",
			args:    args{[]string{"pepper@stark.com"}},
			want:    nil,
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.FindIdentities(tt.args.email)
			sort.Slice(got, func(i, j int) bool {
				return got[i].Email < got[j].Email
			})
			if tt.wantErr(t, err) && err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_repository_StoreIdentity(t *testing.T) {
	ctx := context.Background()
	dyn := dynamotestutils.NewTestContext(ctx, t)
	r := Must(New(dyn.Cfg, dyn.Table))

	type args struct {
		identity aclcore.Identity
	}
	tests := []struct {
		name      string
		args      args
		wantErr   assert.ErrorAssertionFunc
		wantAfter []map[string]types.AttributeValue
	}{
		{
			name: "it should create a brand-new identity details",
			args: args{identity: aclcore.Identity{
				Email:   "pepper@stark.com",
				Name:    "Pepper",
				Picture: "/pepper.jpg",
			}},
			wantErr: assert.NoError,
			wantAfter: []map[string]types.AttributeValue{
				{
					"PK":      dynamoutils.AttributeValueMemberS("USER#pepper@stark.com"),
					"SK":      dynamoutils.AttributeValueMemberS("IDENTITY#"),
					"Email":   dynamoutils.AttributeValueMemberS("pepper@stark.com"),
					"Name":    dynamoutils.AttributeValueMemberS("Pepper"),
					"Picture": dynamoutils.AttributeValueMemberS("/pepper.jpg"),
				},
				{
					"PK":      dynamoutils.AttributeValueMemberS("USER#tony@stark.com"),
					"SK":      dynamoutils.AttributeValueMemberS("IDENTITY#"),
					"Email":   dynamoutils.AttributeValueMemberS("tony+other@stark.com"),
					"Name":    dynamoutils.AttributeValueMemberS("Tony Stark"),
					"Picture": dynamoutils.AttributeValueMemberS("/you/know/me.jpg"),
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
			wantAfter: []map[string]types.AttributeValue{
				{
					"PK":      dynamoutils.AttributeValueMemberS("USER#tony@stark.com"),
					"SK":      dynamoutils.AttributeValueMemberS("IDENTITY#"),
					"Email":   dynamoutils.AttributeValueMemberS("tony@stark.com"),
					"Name":    dynamoutils.AttributeValueMemberS("Ironman"),
					"Picture": dynamoutils.AttributeValueMemberS("/ironman-3.jpg"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dyn = dyn.Subtest(t)
			dyn.Must(dyn.WithDbContent(ctx, []map[string]types.AttributeValue{
				{
					"PK":      dynamoutils.AttributeValueMemberS("USER#tony@stark.com"),
					"SK":      dynamoutils.AttributeValueMemberS("IDENTITY#"),
					"Email":   dynamoutils.AttributeValueMemberS("tony+other@stark.com"),
					"Name":    dynamoutils.AttributeValueMemberS("Tony Stark"),
					"Picture": dynamoutils.AttributeValueMemberS("/you/know/me.jpg"),
				},
			}))

			err := r.StoreIdentity(tt.args.identity)
			if tt.wantErr(t, err) && err == nil {
				dyn.MustBool(dyn.EqualContent(ctx, tt.wantAfter))
			}
		})
	}
}
