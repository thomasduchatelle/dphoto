package catalogviewsdynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/catalogviews"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type AlbumViewRepository struct {
	Client    *dynamodb.Client
	TableName string
}

func (a *AlbumViewRepository) InsertAlbumSize(ctx context.Context, albumSizes []catalogviews.AlbumSize) error {
	return nil
}

// TODO GetAvailabilitiesByUser shouldn't return a slice of AlbumSize: the AlbumSize.User field is not used

func (a *AlbumViewRepository) GetAvailabilitiesByUser(ctx context.Context, user usermodel.UserId) ([]catalogviews.AlbumSize, error) {
	//TODO implement me
	panic("implement me")
}
