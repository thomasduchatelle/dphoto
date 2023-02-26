package aclrefreshdynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamotestutils"
	"testing"
	"time"
)

func Test_repository_DeleteRefreshToken(t *testing.T) {
	awsSession, _, table := dynamotestutils.NewDbContext(t)
	r := Must(New(awsSession, table)).(*repository)

	const secretRefreshToken = "1234567890qwertyuiop"

	type args struct {
		token string
	}
	tests := []struct {
		name        string
		args        args
		givenBefore []map[string]*dynamodb.AttributeValue
		wantAfter   []map[string]*dynamodb.AttributeValue
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name: "it should delete an existing token",
			args: args{secretRefreshToken},
			givenBefore: []map[string]*dynamodb.AttributeValue{
				{
					"PK": {S: aws.String("REFRESH#" + secretRefreshToken)},
					"SK": {S: aws.String("#REFRESH_SPEC")},
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
			dynamotestutils.SetContent(t, r.db, r.table, tt.givenBefore)

			err := r.DeleteRefreshToken(tt.args.token)
			if tt.wantErr(t, err, "DeleteRefreshToken() error = %v, wantErr %v", err, tt.wantErr) && err == nil {
				dynamotestutils.AssertAfter(t, r.db, r.table, tt.wantAfter)
			}
		})
	}
}

func Test_repository_FindRefreshToken(t *testing.T) {
	awsSession, _, table := dynamotestutils.NewDbContext(t)
	r := Must(New(awsSession, table)).(*repository)

	const secretRefreshToken = "1234567890qwertyuiop"

	type args struct {
		token string
	}
	someday := time.Date(2021, 12, 24, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		givenBefore []map[string]*dynamodb.AttributeValue
		args        args
		want        *aclcore.RefreshTokenSpec
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name: "it should find a token that exists",
			args: args{secretRefreshToken},
			givenBefore: []map[string]*dynamodb.AttributeValue{
				{
					"PK":                  {S: aws.String("REFRESH#" + secretRefreshToken)},
					"SK":                  {S: aws.String("#REFRESH_SPEC")},
					"Email":               {S: aws.String("tony@stark.com")},
					"RefreshTokenPurpose": {S: aws.String("web")},
					"AbsoluteExpiryTime":  {S: aws.String(someday.Format(time.RFC3339Nano))},
					"Scopes":              {S: aws.String("foo bar")},
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
			dynamotestutils.SetContent(t, r.db, r.table, tt.givenBefore)

			got, err := r.FindRefreshToken(tt.args.token)
			if tt.wantErr(t, err) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_repository_HouseKeepRefreshToken(t *testing.T) {
	awsSession, _, table := dynamotestutils.NewDbContext(t)
	r := Must(New(awsSession, table)).(*repository)

	const firstRefreshToken = "1234567890qwertyuiop"
	const secondRefreshToken = "poiuytrewq0987654321"
	const thirdRefreshToken = "asdfghjklzxcvbnm"
	someday := time.Date(2021, 12, 24, 0, 0, 0, 0, time.UTC)

	aclcore.TimeFunc = func() time.Time {
		return someday
	}

	tests := []struct {
		name        string
		givenBefore []map[string]*dynamodb.AttributeValue
		wantAfter   []map[string]*dynamodb.AttributeValue
		wantErr     assert.ErrorAssertionFunc
		want        int
	}{
		{
			name: "it should delete any token that already have expired",
			givenBefore: []map[string]*dynamodb.AttributeValue{
				{
					"PK":                 {S: aws.String("REFRESH#" + firstRefreshToken)},
					"SK":                 {S: aws.String("#REFRESH_SPEC")},
					"AbsoluteExpiryTime": {S: aws.String(someday.Add(-1 * time.Hour).Format(time.RFC3339Nano))},
				},
				{
					"PK":                 {S: aws.String("REFRESH#" + secondRefreshToken)},
					"SK":                 {S: aws.String("#REFRESH_SPEC")},
					"AbsoluteExpiryTime": {S: aws.String(someday.Format(time.RFC3339Nano))},
				},
				{
					"PK":                 {S: aws.String("REFRESH#" + thirdRefreshToken)},
					"SK":                 {S: aws.String("#REFRESH_SPEC")},
					"AbsoluteExpiryTime": {S: aws.String(someday.Add(1 * time.Hour).Format(time.RFC3339Nano))},
				},
			},
			wantAfter: []map[string]*dynamodb.AttributeValue{
				{
					"PK":                 {S: aws.String("REFRESH#" + thirdRefreshToken)},
					"SK":                 {S: aws.String("#REFRESH_SPEC")},
					"AbsoluteExpiryTime": {S: aws.String(someday.Add(1 * time.Hour).Format(time.RFC3339Nano))},
				},
			},
			want:    2,
			wantErr: assert.NoError,
		},
		{
			name: "it should not delete anything if nothing has expired",
			givenBefore: []map[string]*dynamodb.AttributeValue{
				{
					"PK":                 {S: aws.String("REFRESH#" + thirdRefreshToken)},
					"SK":                 {S: aws.String("#REFRESH_SPEC")},
					"AbsoluteExpiryTime": {S: aws.String(someday.Add(1 * time.Hour).Format(time.RFC3339Nano))},
				},
			},
			wantAfter: []map[string]*dynamodb.AttributeValue{
				{
					"PK":                 {S: aws.String("REFRESH#" + thirdRefreshToken)},
					"SK":                 {S: aws.String("#REFRESH_SPEC")},
					"AbsoluteExpiryTime": {S: aws.String(someday.Add(1 * time.Hour).Format(time.RFC3339Nano))},
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
			dynamotestutils.SetContent(t, r.db, r.table, tt.givenBefore)

			got, err := r.HouseKeepRefreshToken()
			if tt.wantErr(t, err) {
				dynamotestutils.AssertAfter(t, r.db, r.table, tt.wantAfter)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_repository_StoreRefreshToken(t *testing.T) {
	awsSession, _, table := dynamotestutils.NewDbContext(t)
	r := Must(New(awsSession, table)).(*repository)

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
		givenBefore []map[string]*dynamodb.AttributeValue
		wantAfter   []map[string]*dynamodb.AttributeValue
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
			wantAfter: []map[string]*dynamodb.AttributeValue{
				{
					"PK":                  {S: aws.String("REFRESH#" + firstRefreshToken)},
					"SK":                  {S: aws.String("#REFRESH_SPEC")},
					"Email":               {S: aws.String("tony@stark.com")},
					"RefreshTokenPurpose": {S: aws.String("web")},
					"AbsoluteExpiryTime":  {S: aws.String(someday.Format(time.RFC3339Nano))},
					"Scopes":              {S: aws.String("foo bar")},
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
			givenBefore: []map[string]*dynamodb.AttributeValue{
				{
					"PK":                  {S: aws.String("REFRESH#" + firstRefreshToken)},
					"SK":                  {S: aws.String("#REFRESH_SPEC")},
					"Email":               {S: aws.String("tony@stark.com")},
					"RefreshTokenPurpose": {S: aws.String("web")},
					"AbsoluteExpiryTime":  {S: aws.String(someday.Format(time.RFC3339Nano))},
					"Scopes":              {S: aws.String("foo bar")},
				},
			},
			wantAfter: []map[string]*dynamodb.AttributeValue{
				{
					"PK":                  {S: aws.String("REFRESH#" + firstRefreshToken)},
					"SK":                  {S: aws.String("#REFRESH_SPEC")},
					"Email":               {S: aws.String("tony@stark.com")},
					"RefreshTokenPurpose": {S: aws.String("web")},
					"AbsoluteExpiryTime":  {S: aws.String(someday.Format(time.RFC3339Nano))},
					"Scopes":              {S: aws.String("foo bar")},
				},
				{
					"PK":                  {S: aws.String("REFRESH#" + secondRefreshToken)},
					"SK":                  {S: aws.String("#REFRESH_SPEC")},
					"Email":               {S: aws.String("tony@stark.com")},
					"RefreshTokenPurpose": {S: aws.String("web")},
					"AbsoluteExpiryTime":  {S: aws.String(someday.Format(time.RFC3339Nano))},
					"Scopes":              {NULL: aws.Bool(true)},
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
			givenBefore: []map[string]*dynamodb.AttributeValue{
				{
					"PK":                  {S: aws.String("REFRESH#" + firstRefreshToken)},
					"SK":                  {S: aws.String("#REFRESH_SPEC")},
					"Email":               {S: aws.String("tony@stark.com")},
					"RefreshTokenPurpose": {S: aws.String("web")},
					"AbsoluteExpiryTime":  {S: aws.String(someday.Format(time.RFC3339Nano))},
					"Scopes":              {S: aws.String("foo bar")},
				},
			},
			wantAfter: []map[string]*dynamodb.AttributeValue{
				{
					"PK":                  {S: aws.String("REFRESH#" + firstRefreshToken)},
					"SK":                  {S: aws.String("#REFRESH_SPEC")},
					"Email":               {S: aws.String("tony@stark.com")},
					"RefreshTokenPurpose": {S: aws.String("web")},
					"AbsoluteExpiryTime":  {S: aws.String(someday.Format(time.RFC3339Nano))},
					"Scopes":              {S: aws.String("foo bar")},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NotNil(t, err, i)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dynamotestutils.SetContent(t, r.db, r.table, tt.givenBefore)

			err := r.StoreRefreshToken(tt.args.token, tt.args.spec)
			if tt.wantErr(t, err) {
				dynamotestutils.AssertAfter(t, r.db, r.table, tt.wantAfter)
			}
		})
	}
}
