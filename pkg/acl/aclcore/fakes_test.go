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
	Scopes []*aclcore.Scope
}

func NewScopeRepositoryInMemory(scopes ...aclcore.Scope) *ScopeRepositoryInMemory {
	repo := &ScopeRepositoryInMemory{}
	for _, scope := range scopes {
		s := scope
		repo.Scopes = append(repo.Scopes, &s)
	}
	return repo
}

func (s *ScopeRepositoryInMemory) ListScopesByUser(ctx context.Context, email usermodel.UserId, types ...aclcore.ScopeType) ([]*aclcore.Scope, error) {
	var results []*aclcore.Scope
	for _, scope := range s.Scopes {
		if scope.GrantedTo == email && (len(types) == 0 || slices.Contains(types, scope.Type)) {
			results = append(results, scope)
		}
	}
	return results, nil
}

func (s *ScopeRepositoryInMemory) FindScopesById(ids ...aclcore.ScopeId) ([]*aclcore.Scope, error) {
	var results []*aclcore.Scope
	for _, scope := range s.Scopes {
		if slices.Contains(ids, scope.Id()) {
			results = append(results, scope)
		}
	}
	return results, nil
}

func (s *ScopeRepositoryInMemory) ListScopesByOwner(ctx context.Context, owner ownermodel.Owner, types ...aclcore.ScopeType) ([]*aclcore.Scope, error) {
	var results []*aclcore.Scope
	for _, scope := range s.Scopes {
		if scope.ResourceOwner == owner && (len(types) == 0 || slices.Contains(types, scope.Type)) {
			results = append(results, scope)
		}
	}
	return results, nil
}

func (s *ScopeRepositoryInMemory) ListScopesByOwners(ctx context.Context, owners []ownermodel.Owner, types ...aclcore.ScopeType) ([]*aclcore.Scope, error) {
	var results []*aclcore.Scope
	for _, owner := range owners {
		found, err := s.ListScopesByOwner(ctx, owner, types...)
		if err != nil {
			return nil, err
		}
		results = append(results, found...)
	}
	return results, nil
}

func (s *ScopeRepositoryInMemory) SaveIfNewScope(scope aclcore.Scope) error {
	for _, existing := range s.Scopes {
		if existing.Id() == scope.Id() {
			return nil
		}
	}
	saved := scope
	s.Scopes = append(s.Scopes, &saved)
	return nil
}

func (s *ScopeRepositoryInMemory) DeleteScopes(ids ...aclcore.ScopeId) error {
	var kept []*aclcore.Scope
	for _, scope := range s.Scopes {
		if !slices.Contains(ids, scope.Id()) {
			kept = append(kept, scope)
		}
	}
	s.Scopes = kept
	return nil
}

type IdentityRepositoryInMemory struct {
	Identities map[usermodel.UserId]aclcore.Identity
}

func NewIdentityRepositoryInMemory(identities ...aclcore.Identity) *IdentityRepositoryInMemory {
	repo := &IdentityRepositoryInMemory{Identities: make(map[usermodel.UserId]aclcore.Identity)}
	for _, identity := range identities {
		repo.Identities[identity.Email] = identity
	}
	return repo
}

func (r *IdentityRepositoryInMemory) StoreIdentity(identity aclcore.Identity) error {
	r.Identities[identity.Email] = identity
	return nil
}

func (r *IdentityRepositoryInMemory) FindIdentity(email usermodel.UserId) (*aclcore.Identity, error) {
	identity, ok := r.Identities[email]
	if !ok {
		return nil, aclcore.IdentityDetailsNotFoundError
	}
	return &identity, nil
}

func (r *IdentityRepositoryInMemory) FindIdentities(emails []usermodel.UserId) ([]*aclcore.Identity, error) {
	var results []*aclcore.Identity
	for _, email := range emails {
		if identity, ok := r.Identities[email]; ok {
			found := identity
			results = append(results, &found)
		}
	}
	return results, nil
}

type RefreshTokenRepositoryInMemory struct {
	Tokens  map[string]aclcore.RefreshTokenSpec
	NowFunc func() time.Time
}

func NewRefreshTokenRepositoryInMemory() *RefreshTokenRepositoryInMemory {
	return &RefreshTokenRepositoryInMemory{
		Tokens:  make(map[string]aclcore.RefreshTokenSpec),
		NowFunc: time.Now,
	}
}

func (r *RefreshTokenRepositoryInMemory) StoreRefreshToken(token string, spec aclcore.RefreshTokenSpec) error {
	r.Tokens[token] = spec
	return nil
}

func (r *RefreshTokenRepositoryInMemory) FindRefreshToken(token string) (*aclcore.RefreshTokenSpec, error) {
	spec, ok := r.Tokens[token]
	if !ok {
		return nil, aclcore.InvalidRefreshTokenError
	}
	return &spec, nil
}

func (r *RefreshTokenRepositoryInMemory) DeleteRefreshToken(token string) error {
	delete(r.Tokens, token)
	return nil
}

func (r *RefreshTokenRepositoryInMemory) HouseKeepRefreshToken() (int, error) {
	now := r.NowFunc()
	deleted := 0
	for token, spec := range r.Tokens {
		if spec.AbsoluteExpiryTime.Before(now) {
			delete(r.Tokens, token)
			deleted++
		}
	}
	return deleted, nil
}

type AccessTokenGeneratorInMemory struct {
	ExpiryTime   time.Time
	GeneratedFor []usermodel.UserId
}

func NewAccessTokenGeneratorInMemory() *AccessTokenGeneratorInMemory {
	return &AccessTokenGeneratorInMemory{
		ExpiryTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

func (g *AccessTokenGeneratorInMemory) GenerateAccessToken(email usermodel.UserId) (*aclcore.Authentication, error) {
	g.GeneratedFor = append(g.GeneratedFor, email)
	return &aclcore.Authentication{
		AccessToken: "at-" + email.Value(),
		ExpiryTime:  g.ExpiryTime,
		ExpiresIn:   42,
	}, nil
}

type RefreshTokenGeneratorInMemory struct {
	Repository   *RefreshTokenRepositoryInMemory
	GeneratedFor []aclcore.RefreshTokenSpec
}

func NewRefreshTokenGeneratorInMemory() *RefreshTokenGeneratorInMemory {
	return &RefreshTokenGeneratorInMemory{}
}

func (g *RefreshTokenGeneratorInMemory) GenerateRefreshToken(spec aclcore.RefreshTokenSpec) (string, error) {
	g.GeneratedFor = append(g.GeneratedFor, spec)
	token := "rt-" + spec.Email.Value() + "-" + string(spec.RefreshTokenPurpose)
	if g.Repository != nil {
		if err := g.Repository.StoreRefreshToken(token, spec); err != nil {
			return "", err
		}
	}
	return token, nil
}
