package catalog

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

type FindAlbumsByOwnerPort interface {
	FindAlbumsByOwner(ctx context.Context, owner ownermodel.Owner) ([]*Album, error)
}

type FindAlbumsByOwnerFunc func(ctx context.Context, owner ownermodel.Owner) ([]*Album, error)

func (f FindAlbumsByOwnerFunc) FindAlbumsByOwner(ctx context.Context, owner ownermodel.Owner) ([]*Album, error) {
	return f(ctx, owner)
}
