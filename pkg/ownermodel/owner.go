package ownermodel

import "github.com/pkg/errors"

var (
	EmptyOwnerError = errors.New("owner is mandatory and must be not empty")
)

// Owner is a non-empty ID
type Owner string

func (o Owner) IsValid() error {
	if o == "" {
		return EmptyOwnerError
	}

	return nil
}

func (o Owner) Value() string {
	return string(o)
}

func (o Owner) String() string {
	return string(o)
}
