package common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"strings"
	"testing"
)

func TestCodec(t *testing.T) {
	a := assert.New(t)

	signature := catalog.MediaSignature{
		SignatureSha256: "dbd318c1c462aee872f41109a4dfd3048871a03dedd0fe0e757ced57dad6f2d7",
		SignatureSize:   42,
	}

	encoded, err := EncodeMediaId(signature)
	if a.NoError(err) {
		fmt.Printf("encoded > %s\n", encoded)

		a.False(strings.Index(encoded, "/") >= 0, "encoded should not contains a '/'")

		decoded, err := DecodeMediaId(encoded)
		if a.NoError(err) {
			a.Equal(signature.SignatureSha256, decoded.SignatureSha256)
			a.Equal(signature.SignatureSize, decoded.SignatureSize)
		}
	}
}
