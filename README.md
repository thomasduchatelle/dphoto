![Licence](https://img.shields.io/github/license/thomasduchatelle/dphoto)
![CLI Version](https://img.shields.io/github/tag/thomasduchatelle/dphoto?include_prereleases=&sort=semver&color=007bff)
[![live version](https://img.shields.io/badge/dynamic/json?label=live+version&query=%24.version&url=https%3A%2F%2Fdphoto.duchatelle.net%2Fapi%2Fv1%2Fversion&color=dc3545)](https://dphoto.duchatelle.net)
[![dev version](https://img.shields.io/badge/dynamic/json?label=dev+version&query=%24.version&url=https%3A%2F%2Fdphoto-dev.duchatelle.net%2Fapi%2Fv1%2Fversion&color=28a745)](https://dphoto-dev.duchatelle.net)

[comment]: <> (Generate badges: https://michaelcurrin.github.io/badge-generator/#/generic or https://shields.io/)

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

Bootstrap an environment with built-in command (one-of pre-requisite):

    go run ./app/infra-bootstrap -domain <domain> -email <email> -env dev -google-client-id <id>

To release a new version:

1. make changes on a feature branch and bump the CLI version:
   ```
   ./ci/pre-release.sh 1.5.0
   ```

2. create a pull request to `develop`, review the terraform plan and tests then merge -> it will deploy to [https://dphoto-dev.duchatelle.net](https://dphoto-dev.duchatelle.net)
3. create a pull request `develop -> main`, review the terraform plan then merge -> it will deploy to [https://dphoto.duchatelle.net](https://dphoto.duchatelle.net) and create a tag for the CLI
4. (optional) update local versions of `dphoto` by running
   ```
   go install github.com/thomasduchatelle/dphoto/dphoto@latest
   ```
   
5. to avoid confusion, next development iteration can be started by running `./ci/pre-release.sh 1.6.0-alpha`.

### Tech debt

1. infra-data is using AWS provider v3.x, should be upgraded to v4.x
2. go cli is using AWS SDK 1.x and should use 2.x
