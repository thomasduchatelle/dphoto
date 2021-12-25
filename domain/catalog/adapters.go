package catalog

import "github.com/thomasduchatelle/dphoto/domain/catalogmodel"

var (
	Repository RepositoryPort
)

type RepositoryPort interface {
	FindAllAlbums() ([]*catalogmodel.Album, error)
	InsertAlbum(album catalogmodel.Album) error
	DeleteEmptyAlbum(folderName string) error
	// FindAlbum returns (nil, NotFoundError) when not found
	FindAlbum(folderName string) (*catalogmodel.Album, error)
	// UpdateAlbum updates data of matching catalogmodel.Album.FolderName
	UpdateAlbum(album catalogmodel.Album) error
	// CountMedias counts number of media within the album
	CountMedias(folderName string) (int, error)

	// InsertMedias bulks insert medias
	InsertMedias(media []catalogmodel.CreateMediaRequest) error
	// FindMedias is a paginated search of medias within an album, and optionally within a time range
	FindMedias(folderName string, filter catalogmodel.FindMediaFilter) (*catalogmodel.MediaPage, error)
	// FindExistingSignatures returns the signatures that are already known
	FindExistingSignatures(signatures []*catalogmodel.MediaSignature) ([]*catalogmodel.MediaSignature, error)
	// UpdateMedias updates metadata and mark the media to be moved, the AlbumFolderName is never updated (part of the primary key)
	UpdateMedias(filter *catalogmodel.UpdateMediaFilter, newFolderName string) (string, int, error)

	FindReadyMoveTransactions() ([]*catalogmodel.MoveTransaction, error)
	FindFilesToMove(transactionId, pageToken string) ([]*catalogmodel.MovedMedia, string, error)
	UpdateMediasLocation(transactionId string, moves []*catalogmodel.MovedMedia) error
	DeleteEmptyMoveTransaction(transactionId string) error
}
