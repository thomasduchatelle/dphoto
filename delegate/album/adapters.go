package album

var (
	Repository RepositoryPort
)

type RepositoryPort interface {
	FindAll() ([]Album, error)
	Insert(album Album) error
	Delete(name string) error
	Find(name string) (*Album, error)
	Update(name string, album Album) error
	UpdateMedias(filter *MediaFilter, update MediaUpdate) error
}
