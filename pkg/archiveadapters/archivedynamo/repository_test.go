package archivedynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/archive"
	"testing"
	"time"
)

const owner = "ironman"

func TestShouldAddAndFindLocations(t *testing.T) {
	type addArgs struct {
		owner string
		id    string
		key   string
	}
	type findArgs struct {
		owner string
		id    string
	}
	tests := []struct {
		name     string
		addArgs  []addArgs
		findArgs findArgs
		want     string
		wantErr  error
	}{
		{
			name:     "it should not find a key for a non-existing location",
			addArgs:  nil,
			findArgs: findArgs{owner, "media-1"},
			want:     "",
			wantErr:  archive.NotFoundError,
		},
		{
			name:     "it should not find a key even if a media exists for a different owner",
			addArgs:  []addArgs{{owner, "media-2", "avengers/media-2.jpg"}},
			findArgs: findArgs{"captain", "media-2"},
			want:     "",
			wantErr:  archive.NotFoundError,
		},
		{
			name:     "it should store a location and find it",
			addArgs:  []addArgs{{owner, "media-3", "avengers/media-3.jpg"}},
			findArgs: findArgs{owner, "media-3"},
			want:     "avengers/media-3.jpg",
			wantErr:  nil,
		},
		{
			name:     "it should override a location and find th last version of it",
			addArgs:  []addArgs{{owner, "media-4", "avengers/media-4.jpg"}, {owner, "media-4", "thanos/media-4.jpg"}},
			findArgs: findArgs{owner, "media-4"},
			want:     "thanos/media-4.jpg",
			wantErr:  nil,
		},
	}

	repo := Must(New(
		session.Must(session.NewSession(awsConfig())),
		"dphoto-unittest-archive-"+time.Now().Format("20060102150405.000"),
		true,
	)).(*repository)
	defer repo.db.DeleteTable(&dynamodb.DeleteTableInput{TableName: &repo.table})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			for _, add := range tt.addArgs {
				err := repo.AddLocation(add.owner, add.id, add.key)
				if !a.NoError(err, tt.name) {
					a.FailNow(err.Error())
				}
			}

			gotId, err := repo.FindById(tt.findArgs.owner, tt.findArgs.id)
			if tt.wantErr == nil && a.NoError(err, tt.name) {
				a.Equal(tt.want, gotId, tt.name)
			} else if tt.wantErr != nil {
				a.Equal(tt.wantErr, err, tt.name)
			}
		})
	}
}

func TestUpdateLocations(t *testing.T) {
	tests := []struct {
		name    string
		updates []map[string]string
		ids     []string
		want    map[string]string
	}{
		{
			name: "it should creates non-existing locations",
			updates: []map[string]string{
				{
					"id-01": "key-01",
					"id-02": "key-02",
				},
			},
			ids: []string{"id-01", "id-02"},
			want: map[string]string{
				"id-01": "key-01",
				"id-02": "key-02",
			},
		},
		{
			name: "it should creates non-existing locations, then update some",
			updates: []map[string]string{
				{
					"id-11": "key-11",
					"id-12": "key-12",
				},
				{
					"id-12": "key-12",
					"id-13": "key-13",
				},
			},
			ids: []string{"id-11", "id-12", "id-13"},
			want: map[string]string{
				"id-11": "key-11",
				"id-12": "key-12",
				"id-13": "key-13",
			},
		},
	}

	repo := Must(New(
		session.Must(session.NewSession(awsConfig())),
		"dphoto-unittest-archive-location"+time.Now().Format("20060102150405.000"),
		true,
	)).(*repository)
	defer repo.db.DeleteTable(&dynamodb.DeleteTableInput{TableName: &repo.table})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner := fmt.Sprintf("owner-%s", time.Now().Format("20060102150405.000"))
			for _, update := range tt.updates {
				err := repo.UpdateLocations(owner, update)
				if !assert.NoError(t, err, tt.name) {
					assert.FailNow(t, err.Error())
				}
			}

			got, err := repo.FindByIds(owner, tt.ids)
			if assert.NoError(t, err, tt.name) {
				assert.Equal(t, tt.want, got, tt.name)
			}
		})
	}
}

func awsConfig() *aws.Config {
	return &aws.Config{
		CredentialsChainVerboseErrors: aws.Bool(true),
		Endpoint:                      aws.String("http://localhost:4566"),
		Credentials:                   credentials.NewStaticCredentials("localstack", "localstack", ""),
		Region:                        aws.String("eu-west-1"),
	}
}

func TestFindIdsFromKeyPrefix(t *testing.T) {
	tests := []struct {
		name          string
		keyPrefix     string
		withLocations map[string]string
		want          map[string]string
		wantErr       assert.ErrorAssertionFunc
	}{
		{
			name:      "it should not return anything if table is empty",
			keyPrefix: "ironman/album-01",
			wantErr:   assert.NoError,
		},
		{
			name:      "it should not return anything if location is for a different owner",
			keyPrefix: "ironman/album-01",
			withLocations: map[string]string{
				"img1.jpg": "thor/album-01/image-01",
				"img2.jpg": "blackwindow/album-01/image-02",
			},
			wantErr: assert.NoError,
		},
		{
			name:      "it should not return anything if location is for a different album",
			keyPrefix: "ironman/album-02",
			withLocations: map[string]string{
				"img1.jpg": "ironman/album-01/image-01",
				"img2.jpg": "ironman/album-03/image-02",
			},
			wantErr: assert.NoError,
		},
		{
			name:      "it should not return anything if an album starts with the same values",
			keyPrefix: "ironman/album-02",
			withLocations: map[string]string{
				"img1.jpg": "ironman/album-020/image-01",
				"img2.jpg": "ironman/album-021/image-02",
			},
			wantErr: assert.NoError,
		},
		{
			name:      "it should only return what on the folder",
			keyPrefix: "ironman/album-02",
			withLocations: map[string]string{
				"img1.jpg": "ironman/album-01/image-01",
				"img2.jpg": "ironman/album-02/image-02",
				"img3.jpg": "ironman/album-02/image-03",
				"img4.jpg": "ironman/album-03/image-04",
			},
			want: map[string]string{
				"img2.jpg": "ironman/album-02/image-02",
				"img3.jpg": "ironman/album-02/image-03",
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := Must(New(
				session.Must(session.NewSession(awsConfig())),
				"dphoto-unittest-archive-key-prefix-"+time.Now().Format("20060102150405.000"),
				true,
			)).(*repository)
			defer repo.db.DeleteTable(&dynamodb.DeleteTableInput{
				TableName: &repo.table,
			})

			err := repo.UpdateLocations(owner, tt.withLocations)
			if !assert.NoError(t, err) {
				assert.FailNow(t, err.Error())
			}

			got, err := repo.FindIdsFromKeyPrefix(tt.keyPrefix)
			if !tt.wantErr(t, err, fmt.Sprintf("FindIdsFromKeyPrefix(%v)", tt.keyPrefix)) {
				return
			}
			assert.Equalf(t, tt.want, got, "FindIdsFromKeyPrefix(%v)", tt.keyPrefix)
		})
	}
}
