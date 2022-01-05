package oauth

import (
	"github.com/thomasduchatelle/dphoto/domain/oauthmodel"
	"time"
)

var (
	Config         oauthmodel.Config
	UserRepository oauthmodel.UserRepository
	Now            = time.Now
)
