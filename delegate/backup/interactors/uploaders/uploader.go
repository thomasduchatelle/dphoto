package uploaders

import (
	"github.com/thomasduchatelle/dphoto/delegate/backup/backupmodel"
	"github.com/thomasduchatelle/dphoto/delegate/catalog"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"path"
	"strings"
	"sync"
	"time"
)

type Uploader struct {
	timeline      *catalog.Timeline
	timelineLock  sync.Mutex
	catalog       CatalogProxyAdapter
	onlineStorage backupmodel.OnlineStorageAdapter
	owner         string
	postFilter    backupmodel.PostAnalyseFilter // postFilter might be nil (backupmodel.DefaultPostAnalyseFilter)
}

type mediaRecord struct {
	analysedMedia *backupmodel.AnalysedMedia
	createRequest *catalog.CreateMediaRequest
	folderName    string
}

func NewUploader(catalogProxy CatalogProxyAdapter, onlineStorage backupmodel.OnlineStorageAdapter, owner string, postFilter backupmodel.PostAnalyseFilter) (*Uploader, error) {
	albums, err := catalogProxy.FindAllAlbums()
	if err != nil {
		return nil, err
	}

	timeline, err := catalog.NewTimeline(albums)
	if err != nil {
		return nil, err
	}

	return &Uploader{
		timeline:      timeline,
		timelineLock:  sync.Mutex{},
		catalog:       catalogProxy,
		onlineStorage: onlineStorage,
		owner:         owner,
		postFilter:    postFilter,
	}, nil
}

func (u *Uploader) Upload(buffer []*backupmodel.AnalysedMedia, progressChannel chan *backupmodel.ProgressEvent) error {
	defer func() {
		// media must be closed to release local buffer space
		for _, media := range buffer {
			if toClose, ok := media.FoundMedia.(backupmodel.ClosableMedia); ok {
				err := toClose.Close()
				if err != nil {
					log.WithError(err).Warnf("failed to close media %s", toClose)
				}
			}
		}
	}()

	signatures := make([]*catalog.MediaSignature, len(buffer))
	medias := make(map[catalog.MediaSignature]mediaRecord)

	for i, media := range buffer {
		signature := catalog.MediaSignature{
			SignatureSha256: media.Signature.Sha256,
			SignatureSize:   int(media.Signature.Size),
		}

		folderName, created, err := u.findOrCreateAlbum(media.Details.DateTime)
		if err != nil {
			return err
		}
		if created {
			progressChannel <- &backupmodel.ProgressEvent{Type: backupmodel.ProgressEventAlbumCreated, Count: 1, Album: folderName}
		}

		location := catalog.MediaLocation{
			FolderName: folderName,
			Filename:   fmt.Sprintf("%s_%s%s", media.Details.DateTime.Format("2006-01-02_15-04-05"), signature.SignatureSha256[:8], strings.ToLower(path.Ext(media.FoundMedia.MediaPath().Filename))),
		}

		signatures[i] = &signature
		if _, duplicated := medias[signature]; !duplicated {
			medias[signature] = mediaRecord{
				analysedMedia: media,
				folderName:    folderName,
				createRequest: &catalog.CreateMediaRequest{
					Location: location,
					Type:     catalog.MediaType(media.Type),
					Details: catalog.MediaDetails{
						Width:         media.Details.Width,
						Height:        media.Details.Height,
						DateTime:      media.Details.DateTime,
						Orientation:   catalog.MediaOrientation(media.Details.Orientation),
						Make:          media.Details.Make,
						Model:         media.Details.Model,
						GPSLatitude:   media.Details.GPSLatitude,
						GPSLongitude:  media.Details.GPSLongitude,
						Duration:      media.Details.Duration,
						VideoEncoding: media.Details.VideoEncoding,
					},
					Signature: signature,
				},
			}
		}
	}

	err := u.filterKnownMedias(signatures, medias, progressChannel)
	if err != nil {
		return err
	}

	uploaded := make([]catalog.CreateMediaRequest, len(medias))
	index := 0
	for _, media := range medias {
		err = u.doUpload(media.analysedMedia.FoundMedia, &media.createRequest.Location)
		if err != nil {
			return err
		}
		uploaded[index] = *media.createRequest
		index++

		progressChannel <- &backupmodel.ProgressEvent{
			Type:      backupmodel.ProgressEventUploaded,
			Count:     1,
			Size:      media.analysedMedia.FoundMedia.SimpleSignature().Size,
			Album:     media.folderName,
			MediaType: media.analysedMedia.Type,
		}
	}

	return u.catalog.InsertMedias(uploaded)
}

func (u *Uploader) filterKnownMedias(signatures []*catalog.MediaSignature, medias map[catalog.MediaSignature]mediaRecord, progressChannel chan *backupmodel.ProgressEvent) error {
	knownSignatures, err := u.catalog.FindSignatures(signatures)
	if err != nil {
		return err
	}

	filteredOutCount := uint(0)
	filteredOutSize := uint(0)

	if u.postFilter != nil {
		for sig, record := range medias {
			if !u.postFilter.AcceptAnalysedMedia(record.analysedMedia, record.folderName) {
				log.Debugf("Uploader > skipping filtered out media %s", record.analysedMedia.FoundMedia)
				delete(medias, sig)
				filteredOutCount += 1
				filteredOutSize += record.analysedMedia.Signature.Size
			}
		}
	}

	for _, signature := range knownSignatures {
		if m, ok := medias[*signature]; ok {
			log.Debugf("Uploader > skipping duplicate %s", m.analysedMedia.FoundMedia)
			delete(medias, *signature)
			filteredOutCount += 1
			filteredOutSize += uint(signature.SignatureSize)
		}
	}

	progressChannel <- &backupmodel.ProgressEvent{Type: backupmodel.ProgressEventSkippedAfterAnalyse, Count: filteredOutCount, Size: filteredOutSize}
	return nil
}

func (u *Uploader) findOrCreateAlbum(mediaTime time.Time) (string, bool, error) {
	u.timelineLock.Lock()
	defer u.timelineLock.Unlock()

	if album, ok := u.timeline.FindAt(mediaTime); ok {
		return album.FolderName, false, nil
	}

	year := mediaTime.Year()
	quarter := (mediaTime.Month() - 1) / 3

	createRequest := catalog.CreateAlbum{
		Name:             fmt.Sprintf("Q%d %d", quarter+1, year),
		Start:            time.Date(year, quarter*3+1, 1, 0, 0, 0, 0, time.UTC),
		End:              time.Date(year, (quarter+1)*3+1, 1, 0, 0, 0, 0, time.UTC),
		ForcedFolderName: fmt.Sprintf("/%d-Q%d", year, quarter+1),
	}

	log.Infof("Creates new album '%s' to accomodate media at %s", createRequest.ForcedFolderName, mediaTime.Format(time.RFC3339))

	err := u.catalog.Create(createRequest)
	if err != nil {
		return "", false, errors.Wrapf(err, "failed to create album containing %s [%s]", mediaTime.Format(time.RFC3339), createRequest.String())
	}

	u.timeline, err = u.timeline.AppendAlbum(&catalog.Album{
		Name:       createRequest.Name,
		FolderName: createRequest.ForcedFolderName,
		Start:      createRequest.Start,
		End:        createRequest.End,
	})
	return createRequest.ForcedFolderName, true, err
}

func (u *Uploader) doUpload(media backupmodel.FoundMedia, location *catalog.MediaLocation) (err error) {
	log.Debugf("Uploader > Upload media %s", media)
	location.Filename, err = u.onlineStorage.UploadFile(u.owner, media, location.FolderName, location.Filename)
	return
}

type CatalogProxyAdapter interface {
	FindAllAlbums() ([]*catalog.Album, error)
	InsertMedias(medias []catalog.CreateMediaRequest) error
	Create(createRequest catalog.CreateAlbum) error
	FindSignatures(signatures []*catalog.MediaSignature) ([]*catalog.MediaSignature, error)
}

type CatalogProxy struct{}

func (c CatalogProxy) FindAllAlbums() ([]*catalog.Album, error) {
	return catalog.FindAllAlbums()
}

func (c CatalogProxy) InsertMedias(medias []catalog.CreateMediaRequest) error {
	return catalog.InsertMedias(medias)
}

func (c CatalogProxy) Create(createRequest catalog.CreateAlbum) error {
	return catalog.Create(createRequest)
}

func (c CatalogProxy) FindSignatures(signatures []*catalog.MediaSignature) ([]*catalog.MediaSignature, error) {
	return catalog.FindSignatures(signatures)
}
