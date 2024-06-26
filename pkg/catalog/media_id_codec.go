package catalog

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"github.com/pkg/errors"
	"strings"
)

type MediaId string

func (m MediaId) Value() string {
	return string(m)
}

// GenerateMediaId generate a unique ID for a media.
func GenerateMediaId(signature MediaSignature) (MediaId, error) {
	idBuffer, err := hex.DecodeString(signature.SignatureSha256)
	buf := make([]byte, 8, 8)
	binary.PutUvarint(buf, uint64(signature.SignatureSize))

	for _, b := range buf {
		if b != 0 {
			idBuffer = append(idBuffer, b)
		}
	}

	return MediaId(strings.ReplaceAll(base64.StdEncoding.EncodeToString(idBuffer), "/", "_")), errors.Wrapf(err, "failed to generate media ID for signature %s", signature)
}

// DecodeMediaId reverse what the GenerateMediaId has done to find original signature.
func DecodeMediaId(encodedId MediaId) (*MediaSignature, error) {
	decoded, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(encodedId.Value(), "_", "/"))
	if err != nil {
		return nil, errors.Wrapf(err, "invalid encoded identifier")
	}
	if len(decoded) < 32 {
		return nil, errors.Errorf("invalid encoded identifier: not long enough")
	}

	size, n := binary.Uvarint(decoded[32:])
	if n <= 0 {
		err = errors.Errorf("size can't be read as a var int")
	}
	return &MediaSignature{
		SignatureSha256: hex.EncodeToString(decoded[:32]),
		SignatureSize:   int(size),
	}, err
}
