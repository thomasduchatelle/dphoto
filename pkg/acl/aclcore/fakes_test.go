package aclcore_test

import (
	"context"
	"slices"
	"time"

	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type ScopeRepositoryInMemory struct {
	scopes []*aclcore.Scope
}

func NewScopeRepositoryInMemory(scopes ...aclcore.Scope) *ScopeRepositoryInMemory {
	repo := &ScopeRepositoryInMemory{}
	for _, scope := range scopes {
		s := scope
		repo.scopes = append(repo.scopes, &s)
	}
	return repo
}

func (r *ScopeRepositoryInMemory) ListScopesByUser(_ context.Context, email usermodel.UserId, types ...aclcore.ScopeType) ([]*aclcore.Scope, error) {
	var found []*aclcore.Scope
	for _, scope := range r.scopes {
		if scope.GrantedTo == email && slices.Contains(types, scope.Type) {
			found = append(found, scope)
		}
	}
	return found, nil
}

func (r *ScopeRepositoryInMemory) FindScopesById(ids ...aclcore.ScopeId) ([]*aclcore.Scope, error) {
	var found []*aclcore.Scope
	for _, scope := range r.scopes {
		if slices.Contains(ids, scope.Id()) {
			found = append(found, scope)
		}
	}
	return found, nil
}

func (r *ScopeRepositoryInMemory) ListScopesByOwner(_ context.Context, owner ownermodel.Owner, types ...aclcore.ScopeType) ([]*aclcore.Scope, error) {
	var found []*aclcore.Scope
	for _, scope := range r.scopes {
		if scope.ResourceOwner == owner && slices.Contains(types, scope.Type) {
			found = append(found, scope)
		}
	}
	return found, nil
}

func (r *ScopeRepositoryInMemory) ListScopesByOwners(_ context.Context, owners []ownermodel.Owner, types ...aclcore.ScopeType) ([]*aclcore.Scope, error) {
	var found []*aclcore.Scope
	for _, scope := range r.scopes {
		if slices.Contains(owners, scope.ResourceOwner) && slices.Contains(types, scope.Type) {
			found = append(found, scope)
		}
	}
	return found, nil
}

func (r *ScopeRepositoryInMemory) SaveIfNewScope(scope aclcore.Scope) error {
	for _, existing := range r.scopes {
		if existing.Id() == scope.Id() {
			return nil
		}
	}
	s := scope
	r.scopes = append(r.scopes, &s)
	return nil
}

func (r *ScopeRepositoryInMemory) DeleteScopes(ids ...aclcore.ScopeId) error {
	var kept []*aclcore.Scope
	for _, scope := range r.scopes {
		if !slices.Contains(ids, scope.Id()) {
			kept = append(kept, scope)
		}
	}
	r.scopes = kept
	return nil
}

type IdentityRepositoryInMemory struct {
	identities map[usermodel.UserId]aclcore.Identity
}

func NewIdentityRepositoryInMemory(identities ...aclcore.Identity) *IdentityRepositoryInMemory {
	repo := &IdentityRepositoryInMemory{identities: make(map[usermodel.UserId]aclcore.Identity)}
	for _, identity := range identities {
		repo.identities[identity.Email] = identity
	}
	return repo
}

func (r *IdentityRepositoryInMemory) StoreIdentity(identity aclcore.Identity) error {
	r.identities[identity.Email] = identity
	return nil
}

func (r *IdentityRepositoryInMemory) FindIdentity(email usermodel.UserId) (*aclcore.Identity, error) {
	identity, ok := r.identities[email]
	if !ok {
		return nil, aclcore.IdentityDetailsNotFoundError
	}
	return &identity, nil
}

func (r *IdentityRepositoryInMemory) FindIdentities(emails []usermodel.UserId) ([]*aclcore.Identity, error) {
	var found []*aclcore.Identity
	seen := make(map[usermodel.UserId]bool)
	for _, email := range emails {
		if seen[email] {
			continue
		}
		seen[email] = true
		if identity, ok := r.identities[email]; ok {
			i := identity
			found = append(found, &i)
		}
	}
	return found, nil
}

type RefreshTokenRepositoryInMemory struct {
	tokens  map[string]aclcore.RefreshTokenSpec
	NowFunc func() time.Time
}

func NewRefreshTokenRepositoryInMemory() *RefreshTokenRepositoryInMemory {
	return &RefreshTokenRepositoryInMemory{
		tokens:  make(map[string]aclcore.RefreshTokenSpec),
		NowFunc: time.Now,
	}
}

func (r *RefreshTokenRepositoryInMemory) StoreRefreshToken(token string, spec aclcore.RefreshTokenSpec) error {
	r.tokens[token] = spec
	return nil
}

func (r *RefreshTokenRepositoryInMemory) FindRefreshToken(token string) (*aclcore.RefreshTokenSpec, error) {
	spec, ok := r.tokens[token]
	if !ok {
		return nil, aclcore.InvalidRefreshTokenError
	}
	return &spec, nil
}

func (r *RefreshTokenRepositoryInMemory) DeleteRefreshToken(token string) error {
	delete(r.tokens, token)
	return nil
}

func (r *RefreshTokenRepositoryInMemory) HouseKeepRefreshToken() (int, error) {
	now := r.NowFunc()
	deleted := 0
	for token, spec := range r.tokens {
		if spec.AbsoluteExpiryTime.Before(now) {
			delete(r.tokens, token)
			deleted++
		}
	}
	return deleted, nil
}

var accessTokenGeneratorFakeExpiry = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

type AccessTokenGeneratorFake struct {
	GeneratedFor []usermodel.UserId
}

func NewAccessTokenGeneratorFake() *AccessTokenGeneratorFake {
	return &AccessTokenGeneratorFake{}
}

func (f *AccessTokenGeneratorFake) GenerateAccessToken(email usermodel.UserId) (*aclcore.Authentication, error) {
	f.GeneratedFor = append(f.GeneratedFor, email)
	return &aclcore.Authentication{
		AccessToken: "at-" + string(email),
		ExpiryTime:  accessTokenGeneratorFakeExpiry,
		ExpiresIn:   42,
	}, nil
}

type RefreshTokenGeneratorFake struct {
	GeneratedFor []aclcore.RefreshTokenSpec
}

func NewRefreshTokenGeneratorFake() *RefreshTokenGeneratorFake {
	return &RefreshTokenGeneratorFake{}
}

func (f *RefreshTokenGeneratorFake) GenerateRefreshToken(spec aclcore.RefreshTokenSpec) (string, error) {
	f.GeneratedFor = append(f.GeneratedFor, spec)
	return "rt-" + string(spec.Email) + "-" + string(spec.RefreshTokenPurpose), nil
}
