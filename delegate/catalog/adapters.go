package catalog

var (
	Repository RepositoryPort
)

type RepositoryPort interface {
	FindAllAlbums() ([]*Album, error)
	InsertAlbum(album Album) error
	DeleteEmptyAlbum(folderName string) error
	// FindAlbum returns (nil, NotFoundError) when not found
	FindAlbum(folderName string) (*Album, error)
	// UpdateAlbum updates data of matching Album.FolderName
	UpdateAlbum(album Album) error
	// CountMedias counts number of media within the album
	CountMedias(folderName string) (int, error)

	// InsertMedias bulks insert medias
	InsertMedias(media []CreateMediaRequest) error
	// FindMedias is a paginated search of medias within an album
	FindMedias(folderName string, request PageRequest) (*MediaPage, error)
	// FindExistingSignatures returns the signatures that are already known
	FindExistingSignatures(signatures []*MediaSignature) ([]*MediaSignature, error)
	// UpdateMedias updates metadata and mark the media to be moved, the AlbumFolderName is never updated (part of the primary key)
	UpdateMedias(filter *UpdateMediaFilter, newFolderName string) (string, int, error)

	FindReadyMoveTransactions() ([]*MoveTransaction, error)
	FindFilesToMove(transactionId, pageToken string) ([]*MovedMedia, string, error)
	UpdateMediasLocation(transactionId string, moves []*MovedMedia) error
}

type PageRequest struct {
	Size     int64  // defaulted to 50 if not defined
	NextPage string // empty for the first page
}

type MoveTransaction struct {
	TransactionId string
	Count         int // Number of medias to be moved as part of this transaction
}
