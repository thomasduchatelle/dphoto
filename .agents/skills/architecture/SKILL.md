---
name: architecture
description: Architecture of DPhoto, describes the different subdomains and their respective components. Required skill to do large scale architecture and designs. Not recommended for simple coding tasks.
---

# Architecture of DPhoto

## Subdomains of the application

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

## How to plan work ?

When asked to plan a task, you need to break down the requirements into stories that don't overlap between domains or layers. One story can only affect one
domain and one layer (deployment, api, web, or CLI).

You need to remind the agent that will develop your task the skill he must load before implementing the change.

Reading the domain model of the subdomain being worked on will give SIGNIFICANT insight of what the application is already doing and how. For example, improving
the loading time of the images will affect the backend of the archive domain. The content of the folder `pkg/archive` is very important
for planning. Another example, changing the order of the albums on the website will affect the frontend of the catalog domain. The content of the folder
`web/src/core/catalog/language/` is very important for planning.
