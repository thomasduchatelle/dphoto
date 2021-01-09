package album

var (
	Repository RepositoryPort
)

type RepositoryPort interface {
	FindAll() ([]Album, error)
	Insert(album Album) error
	DeleteEmpty(folderName string) error
	Find(folderName string) (*Album, error)
	// update data of matching Album.FolderName
	Update(album Album) error

	UpdateMedias(filter *MediaFilter, update MediaUpdate) error
}
