DPhoto is an application to back up and visualise photos and videos on the cloud (AWS) used either through a website (deployed on the cloud), or through a
command line interface installed on user computers.

This is a mono-repository containing the code of the backend, the CLI, the website, and the deployments. **Its architecture and the design of each component are
EXTREMELY IMPORTANT**. You must take your time, and think carefully, when accomplishing a coding task.

## Project structure

- `pkg/` - core business logic is always implemented here, in **Golang**. Think of the subdomains that a tasks will affect:
    - `pkg/acl/` - access control and permissions management.
    - `pkg/archive/` - long term storage of the photos and videos, with compression and generation of miniatures.
    - `pkg/catalog/` - organisation of the medias into albums and albums management.
    - `pkg/backup/` - analyse medias and load them into the archive, and index then into the catalog.
    - `pkg/**adapters/` - adapters for AWS services: DynamoDB, S3, SNS/SQS. Adapters are the only place 3rd party libraries are allowed in `pkg/`, the subdomain
      must remain pure.
    - `DATA_MODEL.md` - documentation of the indexes and records of the single-table structure in DynamoDB. It must be kept up to date with the data model.
- `cmd/dphoto/` - in **Golang/Cobra**, CLI source to exposes features from `pkg/`, and presentation logic for the terminal.
- `api/lambdas/` - in **Golang/AWS SDK** source of the REST API deployed as AWS Lambdas with an AWS API Gateway (v2 HTTP). Each operation is deployed as one
  lambda handler and is in its own folder.
    - `api/lambdas/common/` contains utilities shared by most handlers (get context from authorizer, ...)
- `web/` - **DEPRECATED! Project will be replaced by web-nextjs** Typescript / React 19 / Waku framework Website built on top of the REST API, deployed as a
  lambda.
    - `web/src/core/catalog/language/` - data structures used across the web application, **very important for context**.
    - `web/src/core/catalog/**/` - other folders are handlers for the operations available on the UI.
    - `web/src/components/` - React components, usually pure.
    - `web/src/pages/` - Waku page-driven navigation built from the components.
- `web-nextjs/` - **Typescript / React 19 / NextJS framework** Website built on top of the REST API, deployed using SST.
    - this is a new project, implemented incrementally to replace `web/`.
    - **NextJS App Router** is used: file structure must respect its principles.
    - Files structure must follow the best practices from NextJS.
- `deployments/cdk/` - AWS CDK infrastructure (TypeScript/Jest), and SST to deploy `web-nextjs`.
- `.github/` - CI definitions, and Agent instructions.
    - `.github/actions/` - customised actions used within this repository only
    - `.github/workflows/job-*.yml` - reusable sub-workflow to build, test, and deploy the application
    - `.github/workflows/workflow-*.yml` - workflows triggered by external events, they call the "job workflow", never replicate their content.
- `internal/` - **Golang**: mocks and utilities that lower the complexity of the CLI but is not part of the domain of the application.
- `Makefile` - comprehensive list of all the commands to build and test the application.

## How to choose the context ?

Important for all agents before planning: **choose which domains the feature affects** - between archive, catalog, or backup - and **read the files from the
domain**. It will give SIGNIFICANT insight of what the application is already doing and how.

For example, improving the loading time of the images will affect the backend of the archive domain. The content of the folder `pkg/archive` is very important
for planning.

Another example, changing the order of the albums on the website will affect the frontend of the catalog domain. The content of the folder
`web/src/core/catalog/language/` is very important for planning.

## How to get a Pull Request accepted ?

As an agent, your primary objective is to fulfil the feature and have a pull request ACCEPTED by the lead developer. To be accepted, it must be conformed with
the priorities:

1. **no data loss** - the medias stored are very valuable and irreplaceable, everything must be done to never lose a single one.
2. **architecture integrity** - each sub-project defined its design principles, its testing strategy, and its coding standard. Any deviation will lead to the
   pull request being rejected.
3. **simplicity** - the resulting code must be simple and easy to read, even if it requires a complex and large changes to implement a feature: we prefer
   refactoring that simplifies the codebase rather than small changes that adds on the complexity.
4. **cost** - this is a pet-project: operating cost must remain low while not requiring any ongoing effort to operate it.
5. **security** - any reasonable efforts and good practices must be made to avoid data leaks

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

**Trust these instructions** - validated against repository. Search only if incomplete/incorrect.