package m2ts

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestM2TSReader(t *testing.T) {
	a := assert.New(t)

	reader, err := os.Open("../../../test_resources/scan/00000.MTS")
	_, err = ReadM2TSDetails(reader)
	a.NoError(err)
}