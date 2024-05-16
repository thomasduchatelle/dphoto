package usermodel

import (
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"strings"
)

var (
	InvalidUserEmailError = errors.New("user id must be non-empty and cannot start or end with spaces")
)

type UserId string

func (u UserId) Value() string {
	return string(u)
}

func (u UserId) IsValid() error {
	value := string(u)
	if value == "" || strings.Trim(value, " ") != value || strings.Trim(value, " ") != value {
		return InvalidUserEmailError
	}

	return nil
}

type CurrentUser struct {
	UserId UserId
	Owner  *ownermodel.Owner // Owner is the identifier used to store medias of the user. It might be nil.
}

func NewUserId(value string) UserId {
	return UserId(strings.ToLower(strings.Trim(value, " ")))
}
