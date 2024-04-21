package aclrefreshdynamodb

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

func Test_repository_DeleteRefreshToken(t *testing.T) {
	ctx := context.Background()
	dyn := dynamotestutils.NewTestContext(ctx, t)
	r := Must(New(dyn.Cfg, dyn.Table))

	const secretRefreshToken = "1234567890qwertyuiop"

	type args struct {
		token string
	}
	tests := []struct {
		name        string
		args        args
		givenBefore []map[string]types.AttributeValue
		wantAfter   []map[string]types.AttributeValue
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name: "it should delete an existing token",
			args: args{secretRefreshToken},
			givenBefore: []map[string]types.AttributeValue{
				{
					"PK": dynamoutils.AttributeValueMemberS("REFRESH#" + secretRefreshToken),
					"SK": dynamoutils.AttributeValueMemberS("#REFRESH_SPEC"),
				},
			},
			wantAfter: nil,
			wantErr:   assert.NoError,
		},
		{
			name:        "it should ignore when the token has already been deleted",
			args:        args{secretRefreshToken},
			givenBefore: nil,
			wantAfter:   nil,
			wantErr:     assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dyn = dyn.Subtest(t)
			dyn.Must(dyn.WithDbContent(ctx, tt.givenBefore))

			err := r.DeleteRefreshToken(tt.args.token)
			if tt.wantErr(t, err, "DeleteRefreshToken() error = %v, wantErr %v", err, tt.wantErr) && err == nil {
				dyn.MustBool(dyn.EqualContent(ctx, tt.wantAfter))
			}
		})
	}
}

func Test_repository_FindRefreshToken(t *testing.T) {
	ctx := context.Background()
	dyn := dynamotestutils.NewTestContext(ctx, t)
	r := Must(New(dyn.Cfg, dyn.Table))

	const secretRefreshToken = "1234567890qwertyuiop"

	type args struct {
		token string
	}
	someday := time.Date(2021, 12, 24, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		givenBefore []map[string]types.AttributeValue
		args        args
		want        *aclcore.RefreshTokenSpec
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name: "it should find a token that exists",
			args: args{secretRefreshToken},
			givenBefore: []map[string]types.AttributeValue{
				{
					"PK":                  dynamoutils.AttributeValueMemberS("REFRESH#" + secretRefreshToken),
					"SK":                  dynamoutils.AttributeValueMemberS("#REFRESH_SPEC"),
					"Email":               dynamoutils.AttributeValueMemberS("tony@stark.com"),
					"RefreshTokenPurpose": dynamoutils.AttributeValueMemberS("web"),
					"AbsoluteExpiryTime":  dynamoutils.AttributeValueMemberS(someday.Format(time.RFC3339Nano)),
					"Scopes":              dynamoutils.AttributeValueMemberS("foo bar"),
				},
			},
			want: &aclcore.RefreshTokenSpec{
				Email:               "tony@stark.com",
				RefreshTokenPurpose: aclcore.RefreshTokenPurposeWeb,
				AbsoluteExpiryTime:  someday,
				Scopes:              []string{"foo", "bar"},
			},
			wantErr: assert.NoError,
		},
		{
			name:        "it should not find a token that doesn't exist",
			args:        args{secretRefreshToken},
			givenBefore: nil,
			want:        nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, aclcore.InvalidRefreshTokenError, i)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dyn := dyn.Subtest(t)
			dyn.Must(dyn.WithDbContent(ctx, tt.givenBefore))

			got, err := r.FindRefreshToken(tt.args.token)
			if tt.wantErr(t, err) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_repository_HouseKeepRefreshToken(t *testing.T) {
	ctx := context.Background()
	dyn := dynamotestutils.NewTestContext(ctx, t)
	r := Must(New(dyn.Cfg, dyn.Table))

	const firstRefreshToken = "1234567890qwertyuiop"
	const secondRefreshToken = "poiuytrewq0987654321"
	const thirdRefreshToken = "asdfghjklzxcvbnm"
	someday := time.Date(2021, 12, 24, 0, 0, 0, 0, time.UTC)

	aclcore.TimeFunc = func() time.Time {
		return someday
	}

	tests := []struct {
		name        string
		givenBefore []map[string]types.AttributeValue
		wantAfter   []map[string]types.AttributeValue
		wantErr     assert.ErrorAssertionFunc
		want        int
	}{
		{
			name: "it should delete any token that already have expired",
			givenBefore: []map[string]types.AttributeValue{
				{
					"PK":                 dynamoutils.AttributeValueMemberS("REFRESH#" + firstRefreshToken),
					"SK":                 dynamoutils.AttributeValueMemberS("#REFRESH_SPEC"),
					"AbsoluteExpiryTime": dynamoutils.AttributeValueMemberS(someday.Add(-1 * time.Hour).Format(time.RFC3339Nano)),
				},
				{
					"PK":                 dynamoutils.AttributeValueMemberS("REFRESH#" + secondRefreshToken),
					"SK":                 dynamoutils.AttributeValueMemberS("#REFRESH_SPEC"),
					"AbsoluteExpiryTime": dynamoutils.AttributeValueMemberS(someday.Format(time.RFC3339Nano)),
				},
				{
					"PK":                 dynamoutils.AttributeValueMemberS("REFRESH#" + thirdRefreshToken),
					"SK":                 dynamoutils.AttributeValueMemberS("#REFRESH_SPEC"),
					"AbsoluteExpiryTime": dynamoutils.AttributeValueMemberS(someday.Add(1 * time.Hour).Format(time.RFC3339Nano)),
				},
			},
			wantAfter: []map[string]types.AttributeValue{
				{
					"PK":                 dynamoutils.AttributeValueMemberS("REFRESH#" + thirdRefreshToken),
					"SK":                 dynamoutils.AttributeValueMemberS("#REFRESH_SPEC"),
					"AbsoluteExpiryTime": dynamoutils.AttributeValueMemberS(someday.Add(1 * time.Hour).Format(time.RFC3339Nano)),
				},
			},
			want:    2,
			wantErr: assert.NoError,
		},
		{
			name: "it should not delete anything if nothing has expired",
			givenBefore: []map[string]types.AttributeValue{
				{
					"PK":                 dynamoutils.AttributeValueMemberS("REFRESH#" + thirdRefreshToken),
					"SK":                 dynamoutils.AttributeValueMemberS("#REFRESH_SPEC"),
					"AbsoluteExpiryTime": dynamoutils.AttributeValueMemberS(someday.Add(1 * time.Hour).Format(time.RFC3339Nano)),
				},
			},
			wantAfter: []map[string]types.AttributeValue{
				{
					"PK":                 dynamoutils.AttributeValueMemberS("REFRESH#" + thirdRefreshToken),
					"SK":                 dynamoutils.AttributeValueMemberS("#REFRESH_SPEC"),
					"AbsoluteExpiryTime": dynamoutils.AttributeValueMemberS(someday.Add(1 * time.Hour).Format(time.RFC3339Nano)),
				},
			},
			want:    0,
			wantErr: assert.NoError,
		},
		{
			name:        "it should not fail when the DB is empty",
			givenBefore: nil,
			wantAfter:   nil,
			want:        0,
			wantErr:     assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dyn := dyn.Subtest(t)
			dyn.Must(dyn.WithDbContent(ctx, tt.givenBefore))

			got, err := r.HouseKeepRefreshToken()
			if tt.wantErr(t, err) {
				assert.Equal(t, tt.want, got)
				dyn.MustBool(dyn.EqualContent(ctx, tt.wantAfter))
			}
		})
	}
}

func Test_repository_StoreRefreshToken(t *testing.T) {
	ctx := context.Background()
	dyn := dynamotestutils.NewTestContext(ctx, t)
	r := Must(New(dyn.Cfg, dyn.Table))

	const firstRefreshToken = "1234567890qwertyuiop"
	const secondRefreshToken = "poiuytrewq0987654321"
	someday := time.Date(2021, 12, 24, 0, 0, 0, 0, time.UTC)

	type args struct {
		token string
		spec  aclcore.RefreshTokenSpec
	}
	tests := []struct {
		name        string
		args        args
		givenBefore []map[string]types.AttributeValue
		wantAfter   []map[string]types.AttributeValue
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name: "it should create a new token",
			args: args{firstRefreshToken, aclcore.RefreshTokenSpec{
				Email:               "tony@stark.com",
				RefreshTokenPurpose: aclcore.RefreshTokenPurposeWeb,
				AbsoluteExpiryTime:  someday,
				Scopes:              []string{"foo", "bar"},
			}},
			givenBefore: nil,
			wantAfter: []map[string]types.AttributeValue{
				{
					"PK":                  dynamoutils.AttributeValueMemberS("REFRESH#" + firstRefreshToken),
					"SK":                  dynamoutils.AttributeValueMemberS("#REFRESH_SPEC"),
					"Email":               dynamoutils.AttributeValueMemberS("tony@stark.com"),
					"RefreshTokenPurpose": dynamoutils.AttributeValueMemberS("web"),
					"AbsoluteExpiryTime":  dynamoutils.AttributeValueMemberS(someday.Format(time.RFC3339Nano)),
					"Scopes":              dynamoutils.AttributeValueMemberS("foo bar"),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should create another token for the sae user",
			args: args{secondRefreshToken, aclcore.RefreshTokenSpec{
				Email:               "tony@stark.com",
				RefreshTokenPurpose: aclcore.RefreshTokenPurposeWeb,
				AbsoluteExpiryTime:  someday,
			}},
			givenBefore: []map[string]types.AttributeValue{
				{
					"PK":                  dynamoutils.AttributeValueMemberS("REFRESH#" + firstRefreshToken),
					"SK":                  dynamoutils.AttributeValueMemberS("#REFRESH_SPEC"),
					"Email":               dynamoutils.AttributeValueMemberS("tony@stark.com"),
					"RefreshTokenPurpose": dynamoutils.AttributeValueMemberS("web"),
					"AbsoluteExpiryTime":  dynamoutils.AttributeValueMemberS(someday.Format(time.RFC3339Nano)),
					"Scopes":              dynamoutils.AttributeValueMemberS("foo bar"),
				},
			},
			wantAfter: []map[string]types.AttributeValue{
				{
					"PK":                  dynamoutils.AttributeValueMemberS("REFRESH#" + firstRefreshToken),
					"SK":                  dynamoutils.AttributeValueMemberS("#REFRESH_SPEC"),
					"Email":               dynamoutils.AttributeValueMemberS("tony@stark.com"),
					"RefreshTokenPurpose": dynamoutils.AttributeValueMemberS("web"),
					"AbsoluteExpiryTime":  dynamoutils.AttributeValueMemberS(someday.Format(time.RFC3339Nano)),
					"Scopes":              dynamoutils.AttributeValueMemberS("foo bar"),
				},
				{
					"PK":                  dynamoutils.AttributeValueMemberS("REFRESH#" + secondRefreshToken),
					"SK":                  dynamoutils.AttributeValueMemberS("#REFRESH_SPEC"),
					"Email":               dynamoutils.AttributeValueMemberS("tony@stark.com"),
					"RefreshTokenPurpose": dynamoutils.AttributeValueMemberS("web"),
					"AbsoluteExpiryTime":  dynamoutils.AttributeValueMemberS(someday.Format(time.RFC3339Nano)),
					"Scopes":              &types.AttributeValueMemberNULL{Value: true},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should fail to override an existing token",
			args: args{firstRefreshToken, aclcore.RefreshTokenSpec{
				Email:               "peppa@stark.com",
				RefreshTokenPurpose: aclcore.RefreshTokenPurposeWeb,
			}},
			givenBefore: []map[string]types.AttributeValue{
				{
					"PK":                  dynamoutils.AttributeValueMemberS("REFRESH#" + firstRefreshToken),
					"SK":                  dynamoutils.AttributeValueMemberS("#REFRESH_SPEC"),
					"Email":               dynamoutils.AttributeValueMemberS("tony@stark.com"),
					"RefreshTokenPurpose": dynamoutils.AttributeValueMemberS("web"),
					"AbsoluteExpiryTime":  dynamoutils.AttributeValueMemberS(someday.Format(time.RFC3339Nano)),
					"Scopes":              dynamoutils.AttributeValueMemberS("foo bar"),
				},
			},
			wantAfter: []map[string]types.AttributeValue{
				{
					"PK":                  dynamoutils.AttributeValueMemberS("REFRESH#" + firstRefreshToken),
					"SK":                  dynamoutils.AttributeValueMemberS("#REFRESH_SPEC"),
					"Email":               dynamoutils.AttributeValueMemberS("tony@stark.com"),
					"RefreshTokenPurpose": dynamoutils.AttributeValueMemberS("web"),
					"AbsoluteExpiryTime":  dynamoutils.AttributeValueMemberS(someday.Format(time.RFC3339Nano)),
					"Scopes":              dynamoutils.AttributeValueMemberS("foo bar"),
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NotNil(t, err, i)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dyn := dyn.Subtest(t)
			dyn.Must(dyn.WithDbContent(ctx, tt.givenBefore))

			err := r.StoreRefreshToken(tt.args.token, tt.args.spec)
			if tt.wantErr(t, err) && err == nil {
				dyn.MustBool(dyn.EqualContent(ctx, tt.wantAfter))
			}
		})
	}
}
