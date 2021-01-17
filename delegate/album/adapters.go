package album

var (
	Repository RepositoryPort
)

type RepositoryPort interface {
	FindAllAlbums() ([]*Album, error)
	InsertAlbum(album Album) error
	DeleteEmptyAlbum(folderName string) error
	// return (nil, NotFoundError) when not found
	FindAlbum(folderName string) (*Album, error)
	// update data of matching Album.FolderName
	UpdateAlbum(album Album) error

	// InsertMedias bulk insert medias
	InsertMedias(media []CreateMediaRequest) error
	// UpdateMedias update metadata and mark the media to be moved, the AlbumFolderName is never updated (part of the primary key)
	UpdateMedias(filter *UpdateMediaFilter, newFolderName string) (string, int, error)
	// Find signatures of media already existing
	FindExistingSignatures(signatures []*MediaSignature) ([]*MediaSignature, error)
	// FindMedias is a paginated search of medias within an album
	FindMedias(folderName string, request PageRequest) (*MediaPage, error)
}

type PageRequest struct {
	Size     int64  // defaulted to 50 if not defined
	NextPage string // empty for the first page
}
