package catalog

import "context"

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
	TransferMediasPort

	FindAlbumsByOwner(ctx context.Context, owner Owner) ([]*Album, error)
	InsertAlbum(ctx context.Context, album Album) error
	DeleteEmptyAlbum(ctx context.Context, id AlbumId) error
	// FindAlbums only returns found albums
	FindAlbumByIds(ctx context.Context, ids ...AlbumId) ([]*Album, error)
	// UpdateAlbum updates data of matching Album.FolderName
	UpdateAlbum(ctx context.Context, album Album) error

	// InsertMedias bulks insert medias
	InsertMedias(ctx context.Context, owner Owner, media []CreateMediaRequest) error
	// FindMedias is a paginated search for media with their details
	FindMedias(ctx context.Context, request *FindMediaRequest) (medias []*MediaMeta, err error)
	// FindMediaIds is a paginated search to only get the media ids
	FindMediaIds(ctx context.Context, request *FindMediaRequest) (ids []MediaId, err error)
	// FindMediaCurrentAlbum returns the folderName the media is currently in
	FindMediaCurrentAlbum(ctx context.Context, owner Owner, mediaId MediaId) (id *AlbumId, err error)
	// FindExistingSignatures returns the signatures that are already known
	FindExistingSignatures(ctx context.Context, owner Owner, signatures []*MediaSignature) ([]*MediaSignature, error)
	// TransferMedias to a different album, and returns list of moved media ids
	TransferMedias(ctx context.Context, owner Owner, mediaIds []MediaId, newFolderName FolderName) error
}

// CArchiveAdapter forward events to archive package
type CArchiveAdapter interface {
	TimelineMutationObserver

	MoveMedias(owner Owner, ids []MediaId, name FolderName) error
}
