![Licence](https://img.shields.io/github/license/thomasduchatelle/dphoto)
![CLI Version](https://img.shields.io/github/tag/thomasduchatelle/dphoto?include_prereleases=&sort=semver&color=007bff)
[![live version](https://img.shields.io/badge/dynamic/json?label=live+version&query=%24.version&url=https%3A%2F%2Fdphoto.duchatelle.net%2Fapi%2Fv1%2Fversion&color=dc3545)](https://dphoto.duchatelle.net)
[![dev version](https://img.shields.io/badge/dynamic/json?label=dev+version&query=%24.version&url=https%3A%2F%2Fdphoto-dev.duchatelle.net%2Fapi%2Fv1%2Fversion&color=28a745)](https://dphoto-dev.duchatelle.net)

[comment]: <> (Generate badges: https://michaelcurrin.github.io/badge-generator/#/generic or https://shields.io/)

DPhoto
====================================

Backup photo to your private AWS Cloud, and share them with family and friends. Core features:

| Feature              | Version | Description                                                                                                             |
|----------------------|---------|-------------------------------------------------------------------------------------------------------------------------|
| Backup medias        | v0.1    | Backup photos and videos from USB Sticks, Flash drives, and Camera (USB) when they're plugged                           |
| Organise by album    | v0.1    | Photos and videos are organised by album based on their creation date. Each album is a directory in S3.                 |
| Migration script     | v0.1    | Medias already uploaded in S3 are re-ordered by albums and de-duplicated (with interactive command line interface)      |
| Web viewer           | v0.5    | See photos by albums, tag them, search by tag or dates                                                                  |
| Faster UI            | v1.0    | Pre-generate miniatures, supports browser caching (miniature + medium quality) and backend cache medium-quality images. |
| Album sharing        | v1.3    | Albums can be shared to other users                                                                                     |
| Upgraded sharing     | -       | Pre-selected emails and globals views, merge owned and shared-with-me albums                                            |
| Upgraded album views | -       | Albums are listed by most recent medias, medias can be listed by dates (without being grouped by albums                 |
| Upgraded backup      | -       | Backup directly to the cloud (WebDAV), cleaning/filtering possible after backup                                         |
| Photo Frame          | -       | Generate a random pick of media for a photo frame                                                                       |
| Album contribution   | -       | Several users can contribute to the same album (families, friend group)                                                 |
| API Driven           | -       | Provide enhanced API to tip the business on the backend side: CLI authenticated by Google Token (vs AWS Credentials)    |
| Android App          | -       | Minimalist app to synchronise local medias to DPhoto.                                                                   |
| Tagging              | -       | Adding tags to medias to find them later, share them, or print them                                                     |
| Google Photo         | -       | Support Google Photo to push images and video, or import from Google Photo.                                             |

Getting Started
------------------------------------

Install 'dphoto' command line interface and configure it using the following:

    go install github.com/thomasduchatelle/dphoto/cmd/...@latest
    dphoto configure

Then use command `backup` to upload media in a directory, or scan to interactively organise the albums.

    dphoto scan /x/y/z
    dphoto backup /x/y/z

Contribute
------------------------------------

Components:

* `deployments/infra-data`: terraform project to create required infrastructure on AWS for the CLI to work. Project
  won't be re-usable in a different context without overriding backend and some other defaults.
* `pkg`: core domain model and business logic from Hexagonal Architecture. This domain is used from both CLI and app's
  APIs
* [DPhoto CLI](cmd/dphoto/README.md): installed on the end-user computer, backup photos and videos using command line
  interface
* [APP](deployments/sls/README.md): deployed on top of `infra-data`, contains the APP (API deployed on AWS lambdas, and
  WEB)

Note: the repository follow https://github.com/golang-standards/project-layout structure convention.

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
    * `brew install docker docker-compose`

Setup the environment:

    make setup all

    # Run tests & build (all sub-projects)
    make

### Releasing process

Bootstrap an environment with built-in command (one-of pre-requisite):

    go run ./tools/infra-bootstrap -domain <domain> -email <email> -env dev -google-client-id <id>

To release a new version:

1. make changes on a feature branch and bump the CLI version:
   ```
   ./scripts/pre-release.sh 1.5.0
   ```

2. create a pull request to `next`, review the terraform plan and tests then merge -> it will deploy
   to [https://dphoto-dev.duchatelle.net](https://dphoto-dev.duchatelle.net)
3. create a pull request `next -> main`, review the terraform plan then merge -> it will deploy
   to [https://dphoto.duchatelle.net](https://dphoto.duchatelle.net) and create a tag for the CLI
4. (optional) update local versions of `dphoto` by running
   ```
   go install github.com/thomasduchatelle/dphoto/cmd/...@latest
   ```

5. to avoid confusion, next development iteration can be started by running `./ci/pre-release.sh 1.6.0-alpha`.

### AWS Support

DPhoto only supports AWS to be deployed as a serverless application using AWS Gateway, Lambdas, DynamoDB, and S3.

DynamoDB is a single table documented in [README.md](pkg/awssupport/appdynamodb/README.md).

### Required Upgrades

1. go cli is using AWS SDK 1.x and should use 2.x
2. React -> NextJS
