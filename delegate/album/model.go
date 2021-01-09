package album

import (
	"fmt"
	"time"
)

type Album struct {
	Name       string
	FolderName string // unique and immutable
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
	const layout = "2006-01-02T03"
	return fmt.Sprintf("[%s-%s] %s (%s)", a.Start.Format(layout), a.End.Format(layout), a.FolderName, a.Name)
}

// Use unique identifier to compare both albums
func (a *Album) IsEqual(other *Album) bool {
	return a.FolderName == other.FolderName
}

type MediaFilter struct {
	Paginated bool
	Page      int
	Size      int

	AlbumFolderNames map[string]interface{}
	Ranges           []TimeRange
}

type MediaUpdate struct {
	FolderName string
}

func NewFilter() *MediaFilter {
	return &MediaFilter{
		AlbumFolderNames: make(map[string]interface{}),
	}
}

func (m *MediaFilter) WithPage(page, size int) *MediaFilter {
	m.Paginated = true
	m.Page = page
	m.Size = size
	return m
}

func (m *MediaFilter) WithAlbum(folderNames ...string) *MediaFilter {
	for _, name := range folderNames {
		m.AlbumFolderNames[name] = nil
	}

	return m
}

func (m *MediaFilter) WithinRange(start, end time.Time) *MediaFilter {
	m.Ranges = append(m.Ranges, TimeRange{
		Start: start,
		End:   end,
	})

	return m
}

func NewTimeRangeFromAlbum(album Album) TimeRange {
	if album.Start.After(album.End) {
		panic("Album must end AFTER its start: " + album.String())
	}

	return TimeRange{
		Start: album.Start,
		End:   album.End,
	}
}

func NewMoveMediaUpdate(album Album) MediaUpdate {
	return MediaUpdate{
		FolderName: album.FolderName,
	}
}
