package archivedynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/domain/archive"
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

	repo, err := New(
		session.Must(session.NewSession(&aws.Config{
			CredentialsChainVerboseErrors: aws.Bool(true),
			Endpoint:                      aws.String("http://localhost:4566"),
			Credentials:                   credentials.NewStaticCredentials("localstack", "localstack", ""),
			Region:                        aws.String("eu-west-1"),
		})),
		"dphoto-unittest-"+time.Now().Format("20060102150405.000"),
		true,
	)
	if err != nil {
		panic(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			for _, add := range tt.addArgs {
				err = repo.AddLocation(add.owner, add.id, add.key)
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
