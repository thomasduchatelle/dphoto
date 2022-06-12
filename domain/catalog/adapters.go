package catalog

var (
	Repository RepositoryPort
)

type RepositoryPort interface {
	FindAllAlbums(owner string) ([]*Album, error)
	InsertAlbum(album Album) error
	DeleteEmptyAlbum(owner string, folderName string) error
	// FindAlbum returns (nil, NotFoundError) when not found
	FindAlbum(owner string, folderName string) (*Album, error)
	// UpdateAlbum updates data of matching Album.FolderName
	UpdateAlbum(album Album) error
	// CountMedias counts number of media within the album
	CountMedias(owner string, folderName string) (int, error)

	// InsertMedias bulks insert medias
	InsertMedias(owner string, media []CreateMediaRequest) error
	// FindMedias is a paginated search of medias within an album, and optionally within a time range
	FindMedias(owner string, folderName string, filter FindMediaFilter) (*MediaPage, error)
	// FindExistingSignatures returns the signatures that are already known
	FindExistingSignatures(owner string, signatures []*MediaSignature) ([]*MediaSignature, error)
	// UpdateMedias updates metadata and mark the media to be moved, the AlbumFolderName is never updated (part of the primary key)
	UpdateMedias(filter *UpdateMediaFilter, newFolderName string) (string, int, error)
	// FindMediaLocations get the list of all locations media might be in (expected only 1 unless media is pending a physical re-location)
	FindMediaLocations(owner string, signature MediaSignature) ([]*MediaLocation, error)

	FindReadyMoveTransactions(owner string) ([]*MoveTransaction, error)
	FindFilesToMove(transactionId, pageToken string) ([]*MovedMedia, string, error)
	UpdateMediasLocation(owner string, transactionId string, moves []*MovedMedia) error
	DeleteEmptyMoveTransaction(transactionId string) error
}
