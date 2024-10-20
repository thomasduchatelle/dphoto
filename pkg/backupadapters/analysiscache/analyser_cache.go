package analysiscache

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
)

func NewAnalyserCache(localDatabase string) (*AnalyserCache, error) {
	db, err := badger.Open(badger.DefaultOptions(localDatabase).WithLogger(log.StandardLogger()))
	return &AnalyserCache{
		DB: db,
	}, err
}

type AnalyserCache struct {
	DB *badger.DB
}

func (d *AnalyserCache) Close() error {
	log.Debugln("Closing Badger embedded database.")

	err := d.DB.RunValueLogGC(RecommendedGCRatio)
	if err != nil {
		return err
	}

	err = d.DB.Close()
	if errors.Is(err, badger.ErrNoRewrite) {
		return nil
	}
	return err
}

func (d *AnalyserCache) computeKey(found backup.FoundMedia) Key {
	key := Key{
		AbsolutePath: found.MediaPath().Absolute(),
		Size:         found.Size(),
	}
	return key
}

func (d *AnalyserCache) RestoreCache(found backup.FoundMedia) (*backup.AnalysedMedia, bool, error) {
	key := d.computeKey(found)

	payload, err := d.get(key)
	if errors.Is(err, badger.ErrKeyNotFound) {
		return nil, true, nil
	}
	if err != nil {
		return nil, false, errors.Wrapf(err, "failed to restore cache for key=%s", key)
	}

	if d.lastModificationMatchOrIsIgnored(found, payload) {
		return nil, true, nil
	}

	return &backup.AnalysedMedia{
		FoundMedia: found,
		Type:       backup.MediaType(payload.Type),
		Sha256Hash: payload.Sha256Hash,
		Details:    &payload.Details,
	}, false, nil
}

func (d *AnalyserCache) StoreCache(analysedMedia *backup.AnalysedMedia) error {
	key := d.computeKey(analysedMedia.FoundMedia)

	var details backup.MediaDetails
	if analysedMedia.Details != nil {
		details = *analysedMedia.Details
	}

	return d.set(key, Payload{
		LastModification: analysedMedia.FoundMedia.LastModification(),
		Type:             string(analysedMedia.Type),
		Sha256Hash:       analysedMedia.Sha256Hash,
		Details:          details,
	})
}

func (d *AnalyserCache) lastModificationMatchOrIsIgnored(found backup.FoundMedia, payload *Payload) bool {
	return !found.LastModification().IsZero() && !found.LastModification().Equal(payload.LastModification)
}
