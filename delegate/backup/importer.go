package backup

import (
	"duchatelle.io/dphoto/dphoto/backup/interactors"
	"duchatelle.io/dphoto/dphoto/backup/interactors/analyser"
	"duchatelle.io/dphoto/dphoto/backup/model"
	"github.com/pkg/errors"
	"path"
	"regexp"
	"sort"
	"sync"
	"time"
)

type FoundAlbum struct {
	Name       string
	FolderName string // FolderName is the original folder name (Name with date prefix that have been removed)
	Start, End time.Time
}

type scanCompleteListener interface {
	OnScanComplete(total uint)
}

type analyseProgressListener interface {
	OnAnalyseProgress(count, total uint)
}

var (
	datePrefix = regexp.MustCompile("^[0-9]{4}-[01Q][0-9][-_]")
)

// DiscoverAlbumFromSource scan a source to discover albums based on original folder structure
func DiscoverAlbumFromSource(volume model.VolumeToBackup, listeners ...interface{}) ([]*FoundAlbum, error) {
	medias, err := scanMediaSource(volume)
	if err != nil {
		return nil, err
	}

	triggerScanComplete(listeners, len(medias))

	albums := make(map[string]*FoundAlbum)
	for count, found := range medias {
		_, details, err := analyser.ExtractTypeAndDetails(found)
		if err != nil {
			return nil, err
		}

		dir := path.Dir(found.Filename())
		if album, ok := albums[dir]; ok {
			album.pushBoundaries(details.DateTime)
		} else {
			albums[dir] = newFoundAlbum(dir, details.DateTime)
		}

		triggerProgress(listeners, count, len(medias))
	}

	suggestions := make([]*FoundAlbum, len(albums))
	i := 0
	for _, album := range albums {
		suggestions[i] = album
		i++
	}
	sort.Slice(suggestions, func(i, j int) bool {
		if suggestions[i].Start == suggestions[j].Start {
			return suggestions[i].End.Before(suggestions[j].End)
		}

		return suggestions[i].Start.Before(suggestions[j].Start)
	})

	return suggestions, err
}

func triggerScanComplete(listeners []interface{}, total int) {
	for _, l := range listeners {
		if scanCompleteListener, ok := l.(scanCompleteListener); ok {
			scanCompleteListener.OnScanComplete(uint(total))
		}
	}
}

func triggerProgress(listeners []interface{}, count int, total int) {
	for _, l := range listeners {
		if listener, ok := l.(analyseProgressListener); ok {
			listener.OnAnalyseProgress(uint(count+1), uint(total))
		}
	}
}

func scanMediaSource(volume model.VolumeToBackup) ([]model.FoundMedia, error) {
	source, ok := interactors.SourcePorts[volume.Type]
	if !ok {
		return nil, errors.Errorf("No scanner implementation provided for volume type %s", volume.Type)
	}

	lock := sync.Mutex{}
	var medias []model.FoundMedia
	_, _, err := source.FindMediaRecursively(volume, func(media model.FoundMedia) {
		lock.Lock()
		medias = append(medias, media)
		lock.Unlock()
	})
	if err != nil {
		return nil, err
	}
	return medias, nil
}

func newFoundAlbum(albumFullPath string, date time.Time) *FoundAlbum {
	name := path.Base(albumFullPath)
	name = datePrefix.ReplaceAllString(name, "")

	return &FoundAlbum{
		Name:       name,
		FolderName: path.Base(albumFullPath),
		Start:      atStartOfDay(date),
		End:        atStartOfFollowingDay(date),
	}
}
func (a *FoundAlbum) pushBoundaries(date time.Time) {
	if a.Start.After(date) {
		a.Start = atStartOfDay(date)
	}

	if !a.End.After(date) {
		a.End = atStartOfFollowingDay(date)
	}
}

func atStartOfDay(value time.Time) time.Time {
	year, month, day := value.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, value.Location())
}

func atStartOfFollowingDay(value time.Time) time.Time {
	year, month, day := value.Date()
	return time.Date(year, month, day+1, 0, 0, 0, 0, value.Location())
}
