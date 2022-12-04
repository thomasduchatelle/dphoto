package catalog

var (
	repositoryPort RepositoryAdapter
	archivePort    CArchiveAdapter
)

// Init must be called before using this package.
func Init(repositoryAdapter RepositoryAdapter, archive CArchiveAdapter) {
	repositoryPort = repositoryAdapter
	archivePort = archive
}

// RepositoryAdapter brings persistence layer to catalog package
type RepositoryAdapter interface {
	FindAlbumsByOwner(owner string) ([]*Album, error)
	InsertAlbum(album Album) error
	DeleteEmptyAlbum(owner string, folderName string) error
	// FindAlbums only returns found albums
	FindAlbums(ids ...AlbumId) ([]*Album, error)
	// UpdateAlbum updates data of matching Album.FolderName
	UpdateAlbum(album Album) error

	// InsertMedias bulks insert medias
	InsertMedias(owner string, media []CreateMediaRequest) error
	// FindMedias is a paginated search for media with their details
	FindMedias(request *FindMediaRequest) (medias []*MediaMeta, err error)
	// FindMediaIds is a paginated search to only get the media ids
	FindMediaIds(request *FindMediaRequest) (ids []string, err error)
	// FindMediaCurrentAlbum returns the folderName the media is currently in
	FindMediaCurrentAlbum(owner, mediaId string) (folderName string, err error)
	// FindExistingSignatures returns the signatures that are already known
	FindExistingSignatures(owner string, signatures []*MediaSignature) ([]*MediaSignature, error)
	// TransferMedias to a different album, and returns list of moved media ids
	TransferMedias(owner string, mediaIds []string, newFolderName string) error
}

// CArchiveAdapter forward events to archive package
type CArchiveAdapter interface {
	MoveMedias(owner string, ids []string, name string) error
}
