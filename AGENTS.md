# DPhoto - Agent Development Guide

This guide provides essential information for AI coding agents working on the DPhoto codebase.

DPhoto is an application to back up photos and videos on the cloud (AWS), and visualise through a website, or using a command line interface installed on user
computers.

This is a mono-repository containing the code of the backend, the CLI, the website, and the deployments. **Its architecture and the design of each component are
EXTREMELY IMPORTANT**. You must take your time, and think carefully, when accomplishing a coding task.

## Project structure

- `pkg/` - core business logic is always implemented here, in **Golang**. Think of the subdomains that a tasks will affect:
    - `pkg/acl/` - access control and permissions management.
    - `pkg/archive/` - long term storage of the photos and videos, with compression and generation of miniatures.
    - `pkg/catalog/` - organisation of the medias into albums and albums management.
    - `pkg/backup/` - analyse local files to upload them into the archive and to index then into the catalog.
    - `pkg/**adapters/` - adapters for AWS services: DynamoDB, S3, SNS/SQS. Adapters are the only place 3rd party libraries are allowed in `pkg/`, the subdomain
      must remain pure.
    - `DATA_MODEL.md` - documentation of the indexes and records of the single-table structure in DynamoDB. It must be kept up to date with the data model.
- `cmd/dphoto/` - in **Golang/Cobra**, CLI source to exposes features from `pkg/`, and presentation logic for the terminal.
- `api/lambdas/` - in **Golang/AWS SDK** source of the REST API deployed as AWS Lambdas with an AWS API Gateway (v2 HTTP). Each operation is deployed as one
  lambda handler and is in its own folder.
    - `api/lambdas/common/` contains utilities shared by most handlers (get context from authorizer, ...)
    - `api/lambdas/authorizer/` contains the API Gateway authorizer logic, running before any REST request is processed
- `web-nextjs/` - **Typescript / React / NextJS framework** Website built on top of the REST API, deployed using SST.
    - this is a new project, implemented incrementally to replace `web/`.
    - **NextJS App Router** is used: file structure must respect its principles.
    - Files structure must follow the best practices from NextJS.
- `deployments/cdk/` - AWS CDK infrastructure (TypeScript/Jest), and SST configuration to deploy `web-nextjs`.
- `.github/` - CI definitions, and Agent instructions.
    - `.github/actions/` - customised actions used within this repository only
    - `.github/workflows/job-*.yml` - reusable sub-workflow to build, test, and deploy the application
    - `.github/workflows/workflow-*.yml` - workflows triggered by external events, they call the "job workflow", never replicate their content.
- `internal/` - **Golang**: mocks and utilities that lower the complexity of the CLI but is not part of the domain of the application.
- `Makefile` - comprehensive list of all the commands to build and test the application.
- `web/` - **DEPRECATED! Project will be replaced by web-nextjs** ; Typescript / React / Waku framework Website built on top of the REST API, deployed as a
  lambda.
    - `web/src/core/catalog/language/` - data structures used across the web application, **very important for context**.
    - `web/src/core/catalog/**/` - other folders are handlers for the operations available on the UI.
    - `web/src/components/` - React components, usually pure.
    - `web/src/pages/` - Waku page-driven navigation built from the components.

## Priorities: How to get a Pull Request accepted ?

As an agent, your primary objective is to fulfil the feature and have a pull request ACCEPTED by the lead developer. To be accepted, it must be conformed with
the priorities:

1. **no data loss** - the medias stored are very valuable and irreplaceable, everything must be done to never lose a single one.
2. **architecture integrity** - each sub-project defined its design principles, its testing strategy, and its coding standard. Any deviation will lead to the
   pull request being rejected.
3. **simplicity** - the resulting code must be simple and easy to read, even if it requires a complex and large changes to implement a feature: we prefer
   refactoring that simplifies the codebase rather than small changes that adds on the complexity.
4. **cost** - this is a pet-project: operating cost must remain low while not requiring any ongoing effort to operate it.
5. **security** - any reasonable efforts and good practices must be made to avoid data leaks

## Architecture: subdomain of the application

DPhoto is designed around 3 core domains and 1 supporting domain, each can have its data storage, API endpoints, business logic, and presentation layer (UI or
CLI), and/or deployment code.

The domains are as follows:

### Catalog

Main features of the **catalog domain** is to organise the medias into albums. An album is a set of medias regrouped by the date they were captured. When albums
are
edited, the medias are re-indexed to be placed in the appropriate album.

The logic of the domain is within the following path:

* `pkg/catalog` - core business logic
* `pkg/catalogadapters` - adapters to access AWS infrastructure
* `pkg/catalogviews` - list of album from the point of view of a user (restricted to what he owns and has been shared with him)
* `pkg/catalogviewsadapters` - adapters to access AWS infrastructure.
* `api/lambdas` - adapters to expose the features of the domain with REST API (shared with other domains)
* `web/src/core/catalog` - WEB UI components to allow user to interact with the domain
* `cmd/dphoto/cmd` - command line interface is the secondary option to interact with the domain. It runs on the client, with direct access to DynamoDB and S3.
* `deployments/cdk/lib/catalog` - AWS infrastructure definitions (CDK)

### Archive

Main feature of the **archive domain** is to store long term all the medias, and provide WEB-friendly compressed versions of the images to optimise the WEB
interface rendering time.

The logic of the domain is within the following paths:

* `pkg/archive` - core business logic
* `pkg/archiveadapters` - adapters to access AWS infrastructure
* `api/lambdas` - adapters to expose the features of the domain with REST API and ASYNC jobs (shared with other domains)
* `deployments/cdk/lib/archive` - AWS infrastructure definitions (CDK)

There is no UI component for this domain.

### Backup

Main feature of the **backup domain** is to scan a folder for medias (images and videos), and if they are not already in the catalog upload them to the *
*archive** and index them into the **catalog**.

The logic of the domain is within the following paths:

* `pkg/backup` - core business logic
* `pkg/backupadapters` - adapters to access AWS infrastructure
* `cmd/dphoto/cmd` - command line interface is the primary and only option to interact with the domain. It runs on the client with the medias the backup on the
  local drive.

There is no UI nor API for this domain.

### ACL (Access Control List)

ACL is a supporting domain, it's main feature is to control what a user can access:

* a _user_ is a person who can authenticate on the WEB UI (there is no concept of a user on the CLI)
* a _owner_ is a role, and also the concept to which all other resources are attached: _medias_ and _albums_.
* _permissions_ are added to the user to allow him to access and change resources: _owner_ or _album_.

The logic of the domain is within the following paths:

* `pkg/acl` - core business logic.
* `pkg/catalogviews` - optimised access to the catalog after the ACL rules have been applied.
* `pkg/catalogviewsadapters` - adapters to access AWS infrastructure
* `api/lambdas/authorizer` - main usage of the ACL domain to authorise each REST request
* `deployments/cdk/lib/access` - AWS infrastructure definitions (CDK)

## Tech-stack and coding standards

DPhoto is contained in a mono-repository, and is built on top of the following technologies:

* **Golang**: main language of the application, used for all core logic running on the backend, or on the client (CLI).
    * Directories: `pkg/`, `cmd/dphoto/`, `api/`
    * Coding standards: `.github/instructions/go.instructions.md`
* **AWS**, deployed using **CDK / Typescript** and **SST / Typescript**: cloud infrastructure, prioritising on serverless technology and optimized for cost.
    * Directories: `deployments/cdk`
    * Coding standards: `.github/instructions/cdk.instructions.md`
* **NextJS / Typescript (React)**
    * Directories: `web-nextjs/`
    * Coding standards: `.github/instructions/nextjs.instructions.md`
* **Github Action**: build and deploy the application:
    * Directories: `.github/`
    * Coding standards: follow existing patter, and then good practices.

**Every agent must strictly follow the code conventions for each tool and language.**

## How to plan work ?

As a planner, you need to break down the requirements into stories that doesn't overlap between domains or layers. One story can only affect one domain, and one
layer (deployment, api, web, or CLI).

You need to remind the developer (agent) the coding standard file he must read before implementing the change.

Reading the domain model of the subdomain being worked on will give SIGNIFICANT insight of what the application is already doing and how. For example, improving
the loading time of the images will affect the backend of the archive domain. The content of the folder `pkg/archive` is very important
for planning. Another example, changing the order of the albums on the website will affect the frontend of the catalog domain. The content of the folder
`web/src/core/catalog/language/` is very important for planning.

## How to build and test ?

### Golang - `pkg` and `cmd/dphoto/`

**Always run `make setup-go` before executing the tests !**

Then,

```shell
# run all tests
go test ./...
```

Warning: `pkg/acl/jwks` test fails without internet (accounts.google.com) - EXPECTED, ignore it

### Golang - `api/lambdas/`

```shell
# Run all tests
cd api/lambdas && go test ./...

# Build the code
make build-api
```

### Typescript - `web-nextjs/`

**Always run from the web folder `cd web-nextjs`, and always run `npm install` before executing other commands !**

```shell
npm run test          # run unit tests only (~5s)
npm run test:visual   # run visual tests (~30s)
npm run laddle        # run Laddle to take screenshots of the component on :61000
```

### Typescript - `web/`

**Always run from the web folder `cd web`, and always run `npm install` before executing other commands !**

```shell
npx vitest run        # run unit tests only (~17s)
npx playwright test   # run the visual tests 
npm run build         # build the application for deployment

npm run ladle         # Component viewer on :61000, useful to take screenshots of the changes on the UI
```

### Typescript - `deployments/cdk/`

**Always run from the cdk folder `cd deployments/cdk`, and always run `npm install` before executing other commands !**

```shell
npm test              # run all unit tests
npm run synth:test    # verify the CDK template can be built using stub data
```

## Checklist before raising a pull-request

Before requesting a code review, you must ensure:

1. **coding standards have been strictly followed**: changes are conformed with the architecture and designs.
2. **the resulting code is simple and cannot be improved**: think of clean code principles with no excessive comments (NO comment paraphrasing the code!)
3. **conform with the testing strategy**: each project must adhere to the strict testing strategy that guaranty the robustness of the tests with a low coupling
   with the code.
4. the code can be built and is immediately shippable to production.
5. the tests are passing.

---

**Trust these instructions** - validated against the repository. Search only if missing information.