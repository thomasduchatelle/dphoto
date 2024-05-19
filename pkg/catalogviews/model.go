package catalogviews

import (
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type VisibleAlbum struct {
	catalog.Album
	MediaCount         int                // MediaCount is the number of medias on the album
	Visitors           []usermodel.UserId // Visitors are the users that can see the album ; only visible to the owner of the album
	OwnedByCurrentUser bool               // OwnedByCurrentUser is set to true when the user is an owner of the album
}

type ListAlbumsFilter struct {
	OnlyDirectlyOwned bool // OnlyDirectlyOwned provides a sub-view where only resources directly owned by user are displayed and accessible
}
