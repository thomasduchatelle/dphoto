package aclcore

type RevokeAccessTokenAdapter interface {
	DeleteRefreshToken(token string) error
}

type Logout struct {
	RevokeAccessTokenAdapter RevokeAccessTokenAdapter
}

func (l *Logout) RevokeSession(refreshToken string) error {
	return l.RevokeAccessTokenAdapter.DeleteRefreshToken(refreshToken)
}
