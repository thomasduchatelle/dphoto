package album

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type UpdateMediaFilter struct {
	AlbumFolderNames map[string]interface{} // AlbumFolderNames is a set of folder names (map value is nil)
	Ranges           []TimeRange            // empty = no restriction
}

func NewUpdateFilter() *UpdateMediaFilter {
	return &UpdateMediaFilter{
		AlbumFolderNames: make(map[string]interface{}),
	}
}

func (m *UpdateMediaFilter) WithAlbum(folderNames ...string) *UpdateMediaFilter {
	for _, name := range folderNames {
		m.AlbumFolderNames[name] = nil
	}
	return m
}

func (m *UpdateMediaFilter) WithinRange(start, end time.Time) *UpdateMediaFilter {
	m.Ranges = append(m.Ranges, TimeRange{
		Start: start,
		End:   end,
	})

	return m
}

func (m *UpdateMediaFilter) String() string {
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

	return fmt.Sprintf("albums in [%s]%s", strings.Join(albums, ", "), rangesString)
}
