package catalog

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// FindMediaRequest is a filter that is applied to find medias within a time range.
type FindMediaRequest struct {
	Owner            string
	AlbumFolderNames map[string]interface{} // AlbumFolderNames is a set of folder names (map value is nil)
	Ranges           []TimeRange            // Ranges is optional, if empty no restriction will be applied
}

func NewFindMediaRequest(owner string) *FindMediaRequest {
	return &FindMediaRequest{
		Owner:            owner,
		AlbumFolderNames: make(map[string]interface{}),
	}
}

func (m *FindMediaRequest) WithAlbum(folderNames ...string) *FindMediaRequest {
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
		albums = append(albums, name)
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
