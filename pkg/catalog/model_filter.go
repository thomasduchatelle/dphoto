package catalog

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// FindMediaRequest is a filter that is applied to find medias within a time range.
type FindMediaRequest struct { // TODO FindMediaRequest should be deprecated and replaced by MediaSelector
	Owner            Owner
	AlbumFolderNames map[FolderName]interface{} // AlbumFolderNames is a set of folder names (map value is nil)
	Ranges           []TimeRange                // Ranges is optional, if empty no restriction will be applied
}

func NewFindMediaRequest(owner Owner) *FindMediaRequest {
	return &FindMediaRequest{
		Owner:            owner,
		AlbumFolderNames: make(map[FolderName]interface{}),
	}
}

func (m *FindMediaRequest) WithAlbum(folderNames ...FolderName) *FindMediaRequest {
	for _, name := range folderNames {
		m.AlbumFolderNames[name] = nil
	}
	return m
}

func (m *FindMediaRequest) WithinRange(start, end time.Time) *FindMediaRequest {
	if start.IsZero() && end.IsZero() {
		return m
	}

	actualEnd := end
	if actualEnd.IsZero() {
		actualEnd = time.Date(2200, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	m.Ranges = append(m.Ranges, TimeRange{
		Start: start,
		End:   actualEnd,
	})

	return m
}

func (m *FindMediaRequest) String() string {
	albums := make([]string, 0, len(m.AlbumFolderNames))
	for name, _ := range m.AlbumFolderNames {
		albums = append(albums, name.String())
	}

	sort.Slice(albums, func(i, j int) bool {
		return albums[i] < albums[j]
	})

	rangesString := ""
	if len(m.Ranges) > 0 {
		var ranges []string
		for _, r := range m.Ranges {
			ranges = append(ranges, r.String())
		}

		rangesString = fmt.Sprintf(" and ranges in [%s}", strings.Join(ranges, ", "))
	}

	return fmt.Sprintf("albums in %s/{%s}%s", m.Owner, strings.Join(albums, ", "), rangesString)
}
