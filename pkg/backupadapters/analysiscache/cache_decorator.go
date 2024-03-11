package analysiscache

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
)

func NewCacheDecorator(localDatabase string) (backup.AnalyserDecorator, error) {
	db, err := badger.Open(badger.DefaultOptions(localDatabase).WithLogger(log.StandardLogger()))
	return &Decorator{
		db: db,
	}, err
}

type Decorator struct {
	db *badger.DB
}

func (d *Decorator) Decorate(analyseFunc backup.RunnerAnalyserFunc) backup.RunnerAnalyserFunc {
	instance := DecoratorInstance{
		DB:       d.db,
		Delegate: analyseFunc,
	}
	return instance.Analyse
}

func (d *Decorator) Close() error {
	log.Debugln("Closing Badger embedded database.")

	err := d.db.RunValueLogGC(RecommendedGCRatio)
	if err != nil {
		return err
	}

	err = d.db.Close()
	if errors.Is(err, badger.ErrNoRewrite) {
		return nil
	}
	return err
}

type DecoratorInstance struct {
	DB       *badger.DB
	Delegate backup.RunnerAnalyserFunc
}

func (d *DecoratorInstance) Analyse(found backup.FoundMedia, progressChannel chan *backup.ProgressEvent) (*backup.AnalysedMedia, error) {
	key := Key{
		AbsolutePath: found.MediaPath().Absolute(),
		Size:         found.Size(),
	}

	payload, err := d.get(key)
	missedCache := errors.Is(err, badger.ErrKeyNotFound)
	if err != nil && !missedCache {
		return nil, err
	}

	if missedCache || d.lastModificationMatchOrIsIgnored(found, payload) {
		log.Debugf("Runner > missed analysis cache %s", key)
		analysedMedia, err := d.Delegate(found, progressChannel)
		if err == nil {
			var details backup.MediaDetails
			if analysedMedia.Details != nil {
				details = *analysedMedia.Details
			}
			err = d.set(key, Payload{
				LastModification: found.LastModification(),
				Type:             string(analysedMedia.Type),
				Sha256Hash:       analysedMedia.Sha256Hash,
				Details:          details,
			})
		}
		return analysedMedia, err
	}

	progressChannel <- &backup.ProgressEvent{Type: backup.ProgressEventAnalysedFromCache, Count: 1, Size: found.Size()}

	return &backup.AnalysedMedia{
		FoundMedia: found,
		Type:       backup.MediaType(payload.Type),
		Sha256Hash: payload.Sha256Hash,
		Details:    &payload.Details,
	}, nil
}

func (d *DecoratorInstance) lastModificationMatchOrIsIgnored(found backup.FoundMedia, payload *Payload) bool {
	return !found.LastModification().IsZero() && !found.LastModification().Equal(payload.LastModification)
}
