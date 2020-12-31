package album

import (
	"fmt"
	"time"
)

type Album struct {
	Name       string
	FolderName string
	Start      time.Time
	End        time.Time
}

type CreateAlbum struct {
	Name             string
	Start            *time.Time
	End              *time.Time
	ForcedFolderName string
}

func (a *Album) String() string {
	layout := "2006-01-02T03"
	return fmt.Sprintf("%s [%s-%s]", a.Name, a.Start.Format(layout), a.End.Format(layout))
}

// Use unique identifier to compare both albums
func (a *Album) IsEqual(other *Album) bool {
	return a.FolderName == other.FolderName
}

type MediaFilter struct {
	Paginated bool
	Page      int
	Size      int

	AlbumName string
}

type MediaUpdate struct {
	Album          string
	FolderName     string
	ToMoveToFolder string
}

func NewFilter() *MediaFilter {
	return new(MediaFilter)
}

func (m *MediaFilter) withPage(page, size int) *MediaFilter {
	m.Paginated = true
	m.Page = page
	m.Size = size
	return m
}

func (m *MediaFilter) withAlbum(name string) *MediaFilter {
	m.AlbumName = name
	return m
}
