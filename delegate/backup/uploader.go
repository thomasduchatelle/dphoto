package backup

import (
	"duchatelle.io/dphoto/dphoto/backup/model"
	"duchatelle.io/dphoto/dphoto/catalog"
	"fmt"
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
	onlineStorage OnlineStorageAdapter
}

type mediaRecord struct {
	analysedMedia *model.AnalysedMedia
	createRequest *catalog.CreateMediaRequest
}

func NewUploader(catalogProxy CatalogProxyAdapter, onlineStorage OnlineStorageAdapter) (*Uploader, error) {
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
	}, nil
}

func (u *Uploader) Upload(buffer []*model.AnalysedMedia) error {
	log.Infof("Upload %d medias", len(buffer))
	defer func() {
		log.Infof("Closing all medias...")
		// media must be closed to release local buffer space
		for _, media := range buffer {
			if toClose, ok := media.FoundMedia.(ClosableMedia); ok {
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
			SignatureSize:   media.Signature.Size,
		}

		folderName, err := u.findOrCreateAlbum(media.Details.DateTime)
		if err != nil {
			return err
		}

		location := catalog.MediaLocation{
			FolderName: folderName,
			Filename:   fmt.Sprintf("%s_%s%s", media.Details.DateTime.Format("2006-01-02_15-04-05"), signature.SignatureSha256[:8], strings.ToLower(path.Ext(media.FoundMedia.Filename()))),
		}

		signatures[i] = &signature
		if _, duplicated := medias[signature]; !duplicated {
			medias[signature] = mediaRecord{
				analysedMedia: media,
				createRequest: &catalog.CreateMediaRequest{
					Location: location,
					Type:     catalog.MediaType(media.Type),
					Details: catalog.MediaDetails{
						Width:        media.Details.Width,
						Height:       media.Details.Height,
						DateTime:     media.Details.DateTime,
						Orientation:  catalog.MediaOrientation(media.Details.Orientation),
						Make:         media.Details.Make,
						Model:        media.Details.Model,
						GPSLatitude:  media.Details.GPSLatitude,
						GPSLongitude: media.Details.GPSLongitude,
					},
					Signature: signature,
				},
			}
		}
	}

	err := u.filterKnownMedias(signatures, medias)
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
	}

	return u.catalog.InsertMedias(uploaded)
}

func (u *Uploader) filterKnownMedias(signatures []*catalog.MediaSignature, medias map[catalog.MediaSignature]mediaRecord) error {
	knownSignatures, err := u.catalog.FindSignatures(signatures)
	if err != nil {
		return err
	}

	for _, signature := range knownSignatures {
		if m, ok := medias[*signature]; ok {
			log.Debugf("Uploader > skipping duplicate %s", m.analysedMedia.FoundMedia)
		}
		delete(medias, *signature)
	}
	return nil
}

func (u *Uploader) findOrCreateAlbum(mediaTime time.Time) (string, error) {
	u.timelineLock.Lock()
	defer u.timelineLock.Unlock()

	if album, ok := u.timeline.FindAt(mediaTime); ok {
		return album.FolderName, nil
	}

	year := mediaTime.Year()
	quarter := mediaTime.Month() / 4

	createRequest := catalog.CreateAlbum{
		Name:             fmt.Sprintf("Q%d %d", quarter+1, year),
		Start:            time.Date(year, quarter*3+1, 1, 0, 0, 0, 0, time.UTC),
		End:              time.Date(year, (quarter+1)*3+1, 1, 0, 0, 0, 0, time.UTC),
		ForcedFolderName: fmt.Sprintf("%d-Q%d", year, quarter+1),
	}

	log.Infof("Creates new album '%s' to accomodate media at %s", createRequest.ForcedFolderName, mediaTime.Format(time.RFC3339))

	err := u.catalog.Create(createRequest)
	if err != nil {
		return "", err
	}

	u.timeline, err = u.timeline.AppendAlbum(&catalog.Album{
		Name:       createRequest.Name,
		FolderName: createRequest.ForcedFolderName,
		Start:      createRequest.Start,
		End:        createRequest.End,
	})
	return createRequest.ForcedFolderName, err
}

func (u *Uploader) doUpload(media model.FoundMedia, location *catalog.MediaLocation) (err error) {
	log.Debugf("Uploader > Upload media %s", media)
	location.Filename, err = u.onlineStorage.UploadFile(media, location.FolderName, location.Filename)
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
