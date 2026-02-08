package aclcore

// TokenDecoder validates and decodes access tokens into Claims
type TokenDecoder interface {
	Decode(accessToken string) (Claims, error)
}
