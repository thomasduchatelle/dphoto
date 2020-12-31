package album

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"time"
)

func FindAll() ([]Album, error) {
	return Repository.FindAll()
}

func Create(createAlbum CreateAlbum) error {
	// todo - re-assign medias to new album, if any
	if createAlbum.Name == "" {
		return errors.Errorf("Album name is mandatory")
	}

	if createAlbum.Start == nil || createAlbum.End == nil {
		return errors.Errorf("Start and End times are mandatory")
	}

	album := Album{
		Name:       createAlbum.Name,
		FolderName: createAlbum.ForcedFolderName,
		Start:      *createAlbum.Start,
		End:        *createAlbum.End,
	}

	if album.FolderName == "" {
		album.FolderName = generateAlbumFolder(createAlbum.Name, *createAlbum.Start)
	}

	return Repository.Insert(album)
}

func generateAlbumFolder(name string, start time.Time) string {
	re := regexp.MustCompile("[^a-zA-Z0-9]+")
	return fmt.Sprintf("%s_%s", start.Format("2006-01"), re.ReplaceAllString(name, "_"))
}

func Find(name string) (*Album, error) {
	return Repository.Find(name)
}

func Delete(name string, emptyOnly bool) error {
	// todo - move medias out of the album when emptyOnly == false ; create catch-all album if necessary
	return Repository.Delete(name)
}

// Rename an album, and flag all medias to be moved...
// folderName: optional, force to use a specific name
func Rename(previousName, newName string) error {
	found, err := Repository.Find(previousName)
	if err != nil {
		return err
	}
	if found == nil {
		return errors.Errorf("album '%s' not found", previousName)
	}

	album := Album{
		Name:       newName,
		FolderName: generateAlbumFolder(newName, found.Start),
		Start:      found.Start,
		End:        found.End,
	}

	err = Repository.Insert(album)
	if err != nil {
		return err
	}

	err = Repository.UpdateMedias(NewFilter().withAlbum(previousName), MediaUpdate{
		Album:          album.Name,
		ToMoveToFolder: album.FolderName,
	})
	if err != nil {
		return err
	}

	return Repository.Delete(previousName)
}

func Update(name string, start, end time.Time) error {
	// todo - re-assign medias
	return nil
}
