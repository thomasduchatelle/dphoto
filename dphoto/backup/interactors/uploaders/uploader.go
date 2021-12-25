package uploaders

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"github.com/thomasduchatelle/dphoto/domain/catalogmodel"
	"github.com/thomasduchatelle/dphoto/dphoto/backup/backupmodel"
	"path"
	"strings"
	"sync"
	"time"
)

type Uploader struct {
	timeline       *catalog.Timeline
	timelineLock   sync.Mutex
	signatures     map[catalogmodel.MediaSignature]interface{}
	signaturesLock sync.Mutex
	catalog        CatalogProxyAdapter
	onlineStorage  backupmodel.OnlineStorageAdapter
	owner          string
	postFilter     backupmodel.PostAnalyseFilter // postFilter might be nil (backupmodel.DefaultPostAnalyseFilter)
}

type mediaRecord struct {
	analysedMedia *backupmodel.AnalysedMedia
	createRequest *catalogmodel.CreateMediaRequest
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
		timeline:       timeline,
		timelineLock:   sync.Mutex{},
		signatures:     make(map[catalogmodel.MediaSignature]interface{}),
		signaturesLock: sync.Mutex{},
		catalog:        catalogProxy,
		onlineStorage:  onlineStorage,
		owner:          owner,
		postFilter:     postFilter,
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

	var signatures []*catalogmodel.MediaSignature
	medias := make(map[catalogmodel.MediaSignature]mediaRecord)

	for _, media := range buffer {
		signature := catalogmodel.MediaSignature{
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

		location := catalogmodel.MediaLocation{
			FolderName: folderName,
			Filename:   fmt.Sprintf("%s_%s%s", media.Details.DateTime.Format("2006-01-02_15-04-05"), signature.SignatureSha256[:8], strings.ToLower(path.Ext(media.FoundMedia.MediaPath().Filename))),
		}

		if _, duplicated := medias[signature]; !duplicated {
			signatures = append(signatures, &signature)
			medias[signature] = mediaRecord{
				analysedMedia: media,
				folderName:    folderName,
				createRequest: &catalogmodel.CreateMediaRequest{
					Location: location,
					Type:     catalogmodel.MediaType(media.Type),
					Details: catalogmodel.MediaDetails{
						Width:         media.Details.Width,
						Height:        media.Details.Height,
						DateTime:      media.Details.DateTime,
						Orientation:   catalogmodel.MediaOrientation(media.Details.Orientation),
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
		return errors.Wrapf(err, "failed to find signatures %+v", signatures)
	}

	uploaded := make([]catalogmodel.CreateMediaRequest, len(medias))
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

	return errors.Wrapf(u.catalog.InsertMedias(uploaded), "failed to insert photos in catalog: ")
}

func (u *Uploader) filterKnownMedias(signatures []*catalogmodel.MediaSignature, medias map[catalogmodel.MediaSignature]mediaRecord, progressChannel chan *backupmodel.ProgressEvent) error {
	filteredOutCount := uint(0)
	filteredOutSize := uint(0)

	u.signaturesLock.Lock()
	for sig, record := range medias {
		if _, duplicated := u.signatures[sig]; duplicated {
			log.Debugf("Uploader > skipping media already backed up %s", record.analysedMedia.FoundMedia)
			delete(medias, sig)
			filteredOutCount += 1
			filteredOutSize += record.analysedMedia.Signature.Size
		}

		u.signatures[sig] = nil
	}
	u.signaturesLock.Unlock()

	knownSignatures, err := u.catalog.FindSignatures(signatures)
	if err != nil {
		return err
	}

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

	createRequest := catalogmodel.CreateAlbum{
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

	u.timeline, err = u.timeline.AppendAlbum(&catalogmodel.Album{
		Name:       createRequest.Name,
		FolderName: createRequest.ForcedFolderName,
		Start:      createRequest.Start,
		End:        createRequest.End,
	})
	return createRequest.ForcedFolderName, true, err
}

func (u *Uploader) doUpload(media backupmodel.FoundMedia, location *catalogmodel.MediaLocation) (err error) {
	log.Debugf("Uploader > Upload media %s", media)
	location.Filename, err = u.onlineStorage.UploadFile(u.owner, media, location.FolderName, location.Filename)
	return
}

type CatalogProxyAdapter interface {
	FindAllAlbums() ([]*catalogmodel.Album, error)
	InsertMedias(medias []catalogmodel.CreateMediaRequest) error
	Create(createRequest catalogmodel.CreateAlbum) error
	FindSignatures(signatures []*catalogmodel.MediaSignature) ([]*catalogmodel.MediaSignature, error)
}

type CatalogProxy struct{}

func (c CatalogProxy) FindAllAlbums() ([]*catalogmodel.Album, error) {
	return catalog.FindAllAlbums()
}

func (c CatalogProxy) InsertMedias(medias []catalogmodel.CreateMediaRequest) error {
	return catalog.InsertMedias(medias)
}

func (c CatalogProxy) Create(createRequest catalogmodel.CreateAlbum) error {
	return catalog.Create(createRequest)
}

func (c CatalogProxy) FindSignatures(signatures []*catalogmodel.MediaSignature) ([]*catalogmodel.MediaSignature, error) {
	return catalog.FindSignatures(signatures)
}
