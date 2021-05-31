package backup

import (
	"duchatelle.io/dphoto/dphoto/backup/interactors"
	"duchatelle.io/dphoto/dphoto/backup/interactors/analyser"
	"duchatelle.io/dphoto/dphoto/backup/model"
	"fmt"
	"github.com/pkg/errors"
	"path"
	"sort"
	"time"
)

type FoundAlbum struct {
	Name       string
	Start, End time.Time
}

type scanCompleteListener interface {
	OnScanComplete(total uint)
}

type analyseProgressListener interface {
	OnAnalyseProgress(count, total uint)
}

// DiscoverAlbumFromSource scan a source to discover albums based on original folder structure
func DiscoverAlbumFromSource(volume model.VolumeToBackup, listeners ...interface{}) ([]*FoundAlbum, error) {
	medias, err := scanMediaSource(volume)
	if err != nil {
		return nil, err
	}

	triggerScanComplete(listeners, len(medias))

	albums := make(map[string]*FoundAlbum)
	for count, found := range medias {
		analysed, err := analyser.AnalyseMedia(found)
		if err != nil {
			return nil, err
		}

		dir := path.Dir(analysed.FoundMedia.Filename())
		if album, ok := albums[dir]; ok {
			album.pushBoundaries(analysed.Details.DateTime)
		} else {
			albums[dir] = newFoundAlbum(dir, analysed.Details.DateTime)
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

	var medias []model.FoundMedia
	_, _, err := source.FindMediaRecursively(volume, func(media model.FoundMedia) {
		if media == nil {
			fmt.Printf("THIS IS GOING TO BUG\n\nPREVIOUS: %s\n\n", medias[len(medias)-1])
		}
		medias = append(medias, media)
	})
	if err != nil {
		return nil, err
	}
	return medias, nil
}

func newFoundAlbum(albumFullPath string, date time.Time) *FoundAlbum {
	return &FoundAlbum{
		Name:  path.Base(albumFullPath),
		Start: date,
		End:   date,
	}
}
func (a *FoundAlbum) pushBoundaries(date time.Time) {
	if a.Start.After(date) {
		a.Start = atStartOfDay(a.Start)
	}

	if a.End.Before(date) {
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
