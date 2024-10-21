package analysiscache

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v4"
	"github.com/pkg/errors"
)

const RecommendedGCRatio = 0.5

func (d *AnalyserCache) get(key Key) (*Payload, error) {
	var (
		payload Payload
	)

	err := d.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key.SerialisedKey())
		if err != nil {
			return errors.Wrapf(err, "failed to find key [%s]", key)
		}

		jsonPayload, err := item.ValueCopy(nil)
		if err != nil {
			return errors.Wrapf(err, "failed to read jsonPayload from cache [%s]", key)
		}

		err = json.Unmarshal(jsonPayload, &payload)
		return errors.Wrapf(err, "'%s' jsonPayload couldn't be deserialized [%s]", key, jsonPayload)
	})
	return &payload, err
}

func (d *AnalyserCache) set(key Key, payload Payload) error {
	return d.DB.Update(func(txn *badger.Txn) error {
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return errors.Wrapf(err, "failed to serialise in JSON [%+v]", payload)
		}

		return txn.Set(key.SerialisedKey(), jsonPayload)
	})
}
