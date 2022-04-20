package catalogmodel

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// FindMediaFilter is a filter that is applied to find medias within a time range.
type FindMediaFilter struct {
	PageRequest PageRequest // PageRequest size will use a default if too high or not set (0)
	TimeRange   TimeRange   // TimeRange is optional
}

// UpdateMediaFilter is used internally to filter medias based on album and date range
type UpdateMediaFilter struct {
	Owner            string
	AlbumFolderNames map[string]interface{} // AlbumFolderNames is a set of folder names (map value is nil)
	Ranges           []TimeRange            // empty = no restriction
}

func NewUpdateFilter(owner string) *UpdateMediaFilter {
	return &UpdateMediaFilter{
		Owner:            owner,
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

	return fmt.Sprintf("albums in %s/{%s}%s", m.Owner, strings.Join(albums, ", "), rangesString)
}
