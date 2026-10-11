package catalog

import (
	"github.com/pkg/errors"
)

const MaxCoversPerAlbum = 4

var (
	TooManyCoversErr = errors.Errorf("an album cannot have more than %d covers", MaxCoversPerAlbum)
)

type CoverOrigin string

const (
	CoverOriginRandom       CoverOrigin = "RANDOM"
	CoverOriginCherryPicked CoverOrigin = "CHERRY_PICKED"
)

type Cover struct {
	MediaId  MediaId
	Filename string
	Origin   CoverOrigin
}
