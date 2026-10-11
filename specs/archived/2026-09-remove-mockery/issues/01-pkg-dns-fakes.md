# 01 — `pkg/dns` Fakes

Status: done
Layer: `pkg/dns`
Depends on: —

## Description

Replace the mockery-generated `mocks.CertificateManager` and `mocks.CertificateAuthority` in
`pkg/dns/certificate_manager_test.go` with in-memory Fakes.

See `../spec.md` for the feature principles (one Fake per real backing store; verify by state;
no "was-called" assertions).

## Files touched

- `pkg/dns/certificate_manager_test.go` (rewrite the test bodies)
- `pkg/dns/fakes_test.go` (new)

## Fakes to introduce

Create `pkg/dns/fakes_test.go` (package `dns_test`) with:

- **`CertificateManagerInMemory`** — implements `dns.CertificateManager` (or whatever interface the
  package exports; today it is a global `dns.CertificateManager` variable of an anonymous interface,
  check `pkg/dns/adapters.go`). Backing state:
  - `Certificates map[string]dnsdomain.ExistingCertificate` keyed by domain
  - `InstalledContent map[string]dnsdomain.CompleteCertificate` keyed by arn (`""` for a fresh install)
  - `SSMEnsured map[string]bool` keyed by arn
  - Methods:
    - `FindCertificate(ctx, domain)` returns from `Certificates` or `dnsdomain.CertificateNotFoundError`
    - `InstallCertificate(ctx, arn, cert)` writes into `InstalledContent`
    - `EnsureSSMParameter(ctx, arn)` marks `SSMEnsured[arn] = true`
  - Test-only accessors (only if there is no read method already): `IsSSMEnsured(arn) bool`,
    `Installed(arn) (dnsdomain.CompleteCertificate, bool)`.
- **`CertificateAuthorityInMemory`** — implements `dns.CertificateAuthority`. Returns a canned
  `dnsdomain.CompleteCertificate` from a `NextCertificate` field. Records requests in
  `Requested []CertificateRequest{Email, Domain string}`. Test-only accessor: `RequestedFor() []CertificateRequest`.

## Test conversion

The three cases in `TestRenewCertificate` become state assertions:

- Case "not create a new certificate if one already exists": seed `CertificateManagerInMemory.Certificates`
  with a non-expired cert. After `dns.RenewCertificate(...)`, assert `IsSSMEnsured("arn::132456")` is
  `true` and `CertificateAuthorityInMemory.RequestedFor()` is empty (no CA call).
- Case "create a new certificate if the existing one has or is about to expire": seed with an expiring
  cert. After the call, assert `Installed("arn::132456")` returns the canned cert and
  `RequestedFor()` contains one entry with the right `(email, domain)`.
- Case "create a new certificate if none were there": seed nothing (natural
  `CertificateNotFoundError`). After the call, assert `Installed("")` contains the canned cert and
  `RequestedFor()` contains one entry.

## Failure-injection tests

None in this story.

## Success criteria

- `pkg/dns/certificate_manager_test.go` no longer imports `github.com/thomasduchatelle/dphoto/internal/mocks`
  or `github.com/stretchr/testify/mock`.
- `go test ./pkg/dns/...` is green.
- No other file in the package regressed.
