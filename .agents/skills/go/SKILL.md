---
name: go
description: Golang coding standards and testing style using table-driven testing patterns. Required skill to write code or review on the backend (.go) (`pkg/`), APIs (`api/`), and CLI (`cmd/`).
---

# Golang Coding Standards for DPhoto

## How to write a test

Use the golang idiomatic **table-driven** testing with a slice of test cases. The **canonical
structure of each case is fixed** — do not omit or rename fields on a whim:

1. **`name`** — mandatory, always. Descriptive names starting with "it should" and summarising the
   trigger and the expectation. Examples:
    * `"it should GRANT access to the media owner"`
    * `"it should DENY access to a visitor with no permission"`

2. **`fields`** — the actual fields of the structure under test. Mirror the struct literally.
   **Ignore only if the code under test is a plain function, not a struct.**
    * Do NOT drop `fields` because the tests happen to share the same collaborators — the fields
      must still be declared so each case can override any one of them when it needs to.
    * Use concrete Fake types (not the interface) when `expectX` below needs to read state back
      through a test-only accessor.
    * Every `fields{...}` literal must be **readable in isolation**: a reviewer should
      understand what state each Fake is in by looking only at the case's expression, without
      chasing helper functions that assemble several Fakes at once.

3. **`args`** — the arguments passed to the function under test. **Mandatory even if there is only
   a single argument.** Ignore only if the function takes no argument at all.

4. **`want`** — the first returned value (other than `error`). Ignore if the function returns only
   an error, or if more than one non-error value is returned (use `wantX` instead). Asserted with
   deep equals.

5. **`wantX`** — use when the function returns multiple non-error values; replace `X` with the
   parameter name (e.g. `wantAlbum`, `wantMedias`). Asserted with deep equals.

6. **`wantErr`** — of type `assert.ErrorAssertionFunc`
   (`func(TestingT, error, ...interface{}) bool`). Ignore if the function does not return an
   error.
    * If no error is expected: `assert.NoError`.
    * If a specific error is expected: an anonymous function using `assert.ErrorIs` /
      `assert.ErrorAs` (see the full example).

7. **`expectX`** — use it to assert the state of a fake dependency: an event has been published, or a data has been saved. **Never to assert an interaction with a dependency.**
   * prefer complete state validation: `expectEventsFired []Event`, `expectStoredAlbums []*Album`.
   * use projection if a part of the state is not relevant for the test (already present and not modified for example): `expectStoredAlbumIds []AlbumId`, `expectAlbumDates map[AlbumId]{startDate time.Time ; endDate time.Time}`.
   * never assert if a function has been called or not (bad examples: `expectCalled int`) ; verify the state of the fake instead.
   * never assert an action has been executed (bad example: `expectAlbumDeleted bool`, `expectHasBeenSaved bool`) ; verify a projection of the state instead.
   * never assert if a function has been called with a specific argument (bad example: `expectXBeCalledFor AlbumId`).

**Reminder: `want*` must only be used for argument returned by the function under test. `expect*` must be used to validate the state of a fake.**

### Test-file-level fixtures

Small immutable value fixtures (IDs, owner names, canned certificates, sample records that never
change) live as `var` or `const` at the top of the test file. Each case declares only what is
specific to it.

### Fake instantiation in `fields`

Two goals compete: 

* **explicit, readable case expressions** - what's specific to the test is defined for the test, the rest is shared 
  through constants or helper functions with clear naming: a reviewer should not have to jump into a
  helper to know what a case is doing.
* **no cross-pollution between cases** - each case gets a fresh Fake so mutations
  don't leak)
 
Follow this decision tree, in order:

1. **if the fake is stateless** (no mutable backing state — e.g. a token generator that returns a canned
   value, a resizer that returns a canned image). Share **one test-function-level instance** as a `var`
   with a descriptive name, and reference it by name in every case:

    ```go
    var simpleAccessTokenGenerator = NewAccessTokenGeneratorFake()

    // ...
    fields: fields{AccessTokenGenerator: simpleAccessTokenGenerator}
    ```

2. **if the fake is stateful and the case wants it empty** — use its default constructor **inline** in the
   case, so it's obvious the state is empty:

    ```go
    fields: fields{RefreshTokenRepository: NewRefreshTokenRepositoryInMemory()}
    ```

3. **if the fake is stateful and the case wants a shared "happy path" seeded state** — expose a small
   **named constructor function** that returns a fresh instance every time, and call it inline in
   each case:

    ```go
    func identityRepositoryWithStark() *IdentityRepositoryInMemory {
        return NewIdentityRepositoryInMemory(aclcore.Identity{
            Email: "tony@stark.com", Name: "Tony Stark",
        })
    }

    // ...
    fields: fields{IdentityRepository: identityRepositoryWithStark()}
    ```

   The constructor's **name** documents the seeded state. The **function call** guarantees each
   case gets its own instance — no cross-pollution.

4. **if a specific case needs one-off seeding different from the shared happy path** — build the Fake
   **inline** in that case with its regular constructor, so the specialness is visible right there:

    ```go
    fields: fields{IdentityRepository: NewIdentityRepositoryInMemory(
        aclcore.Identity{Email: "peter@parker.com", Name: "Peter Parker"},
    )}
    ```

### Testing example

```go
func TestCatalogAuthorizer_IsAuthorisedToViewMedia(t *testing.T) {
    userId := usermodel.UserId("user-1")
    owner := ownermodel.Owner("owner-1")
    mediaId := catalog.MediaId("media-1")

    isAnAccessForbiddenError := func(t assert.TestingT, err error, i ...interface{}) bool {
        return assert.ErrorIs(t, err, aclcore.AccessForbiddenError)
    }

    // Named constructor for the seeded happy-path scope repository (fresh instance per case,
    // no cross-pollution between cases).
    scopeRepositoryWithVisitor := func() *ScopeRepositoryInMemory {
        return NewScopeRepositoryInMemory(aclcore.Scope{
            Type:          aclcore.AlbumVisitorScope,
            GrantedTo:     userId,
            ResourceOwner: owner,
            ResourceId:    "media-1",
        })
    }

    type fields struct {
        HasPermissionPort catalogacl.HasPermissionPort
    }
    type args struct {
        currentUser usermodel.CurrentUser
        owner       ownermodel.Owner
        mediaId     catalog.MediaId
    }
    tests := []struct {
        name    string
        fields  fields
        args    args
        wantErr assert.ErrorAssertionFunc
    }{
        {
            name:   "it should GRANT access to the media owner",
            fields: fields{HasPermissionPort: scopeRepositoryWithVisitor()},
            args: args{
                currentUser: usermodel.CurrentUser{UserId: userId, Owner: &owner},
                owner:       owner,
                mediaId:     mediaId,
            },
            wantErr: assert.NoError,
        },
        {
            name:   "it should DENY access to a visitor with no permission",
            fields: fields{HasPermissionPort: NewScopeRepositoryInMemory()}, // empty repository, inline
            args: args{
                currentUser: usermodel.CurrentUser{UserId: userId},
                owner:       owner,
                mediaId:     mediaId,
            },
            wantErr: isAnAccessForbiddenError,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            authorizer := &catalogacl.CatalogAuthorizer{HasPermissionPort: tt.fields.HasPermissionPort}
            tt.wantErr(t, authorizer.IsAuthorisedToViewMedia(ctx, tt.args.currentUser, tt.args.owner, tt.args.mediaId))
        })
    }
}
```
