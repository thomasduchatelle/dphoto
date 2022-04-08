![GitHub tag](https://img.shields.io/github/tag/thomasduchatelle/dphoto?include_prereleases=&sort=semver&color=blue)

DPhoto
====================================

Backup photo to your private AWS Cloud, and share them with family and friends. Core features:

| Feature           | Version  | Description                                                                                                        |
|-------------------|----------|--------------------------------------------------------------------------------------------------------------------|
| Backup medias     | v1.0     | Backup photos and videos from USB Sticks, Flash drives, and Camera (USB) when they're plugged                      |
| Organise by album | v1.0     | Photos and videos are organised by album based on their creation date. Each album is a directory in S3.            |
| Migration script  | v1.0     | Medias already uploaded in S3 are re-ordered by albums and de-duplicated (with interactive command line interface) |
| Web viewer        | *in dev* | See photos by albums, tag them, search by tag or dates                                                             |
| Media sharing     | -        | Albums can be shared and contributed by several users                                                              |

Getting Started
------------------------------------

Install 'dphoto' command line interface and configure it using the following:

    go install github.com/thomasduchatelle/dphoto/dphoto@latest
    dphoto configure

Then use command `backup` to upload media in a directory, or scan to interactively organise the albums.

    dphoto scan /x/y/z
    dphoto backup /x/y/z

Contribute
------------------------------------

Components:

* `infra-data`: terraform project to create required infrastructure on AWS for the CLI to work. Project won't be re-usable in a different context without overriding backend and some other defaults.
* `domain`: core domain model and business logic from Hexagonal Architecture. This domain is integrated and used from both CLI and app's APIs
* [DPhoto CLI](./dphoto/README.md): installed on the end-user computer, backup photos and videos using command line interface
* [APP](./app/README.md): deployed on top of `infra-data`, contains the viewer UI, and APIs for the UI and the CLI

### Install development environment

Required tools:

* Infra:
  * terraform: `brew install tfenv`, [Makefile](./Makefile)
  * Serverless Framework: `npm install -g serverless`
  * AWS CLI: `brew install awscli`
* Languages & build tools:
  * `make`
  * GoLang: `brew install golang`
  * Yarn: `brew install yarn`
* Docker and Docker Compose

Setup the environment:

    make install all

    # Run tests & build (all sub-projects)
    make

### Releasing process

To release a new version, without CI:

1. verify everything is checked-in (`git status`)
2. push on `main`
3. **infra-data** - verify and approve the plan on [https://app.terraform.io/app/dphoto/workspaces](https://app.terraform.io/app/dphoto/workspaces)
4. **One-timer pre-requisite**: create an SSL certificate for the domain to provision the API Gateway. Certificate is then automatically re-newed
   ```
   go run ./app/letsencrypt/ignition -domain <domain> -email <email> -env dev
   # to generate manually a SSL certificate:
   go install github.com/go-acme/lego/v4/cmd/lego@latest
   ```

5. **One-timer pre-requisite**: create a SSM parameter /dphoto/{Serverless Stage}/googleLogin/clientId with the client if from https://console.developers.google.com/apis/credentials
6. **APP** - to deploy dev version, run `cd app && make deploy`
7. **dphoto CLI** - create a git tag: `git tag dphoto/v1.x.y && git push --tags`

### Tech debt

1. infra-data is using AWS provider v3.*, should be upgraded to v4.*
2. 