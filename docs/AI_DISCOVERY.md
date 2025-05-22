AI Discovery
=======================================

This file is not related to DPhoto, it is my learning notes on how IA is leveraged on the project.

IA Concepts
---------------------------------------

1. Workflow: _how to integrate IA to each step of the development process (what prompts, what tools) ?_
2. Models: _what models perform better on what type of tasks?_
3. Tools: _how to interact with an IA agent and get to update the code ?_

### Prompts (wip)

#### Testing: Fake

```
You need to use Fake implementation of SharingAPI: an inmemory implementation of the interface. The requests are stored in a property
and its value is asserted in test (1). A property with an error to return is set on the (2) so the fake return a Promise.reject.

You need to use Fake implementation of SharingAPI: an inmemory implementation of the interface. The requests to grantAccessToAlbum are
 stored in a property and its value is asserted in test (1). A property with an error to return is set on the (2) and (3) so the fake
return a Promise.reject.
The requests to sharingAPI.loadUserDetails are not stored. The Fake is a simple in-memory implementation returning what it has on the
relevant property (a list of user details). 
```

Discovery path (and tasks)
---------------------------------------

IA will be used for the following tasks:

1. merging two React State into one more consistent to improve user experience of Sharing feature
    1. refactoring to break down the massive reducer - needs suggestion and implementation
2. migrating to AWS CDK from Terraform and Serverless
    1. aiming for 2 stacks: long term data stores, and WEB overlay
    2. migration paths to be included and executed
    3. new domain to be used
3. migrating to Waku from native React: very light React Framework in beta/exploration
    1. progressive migration - decoupling behaviours from framework before re-integrating them to new framework
    2. parallel deployment - /v2 would get to the new UI
    3. Auth0 - moving to public IDP must be considered
    4. NPM - Yarn or PNPM don't look justified for this project
    5. Visual testing to rethink as the tools used seem discontinued and incompatible with new ones
4. End-to-End testing integrated to Github pipelines to validate CRITICAL path(s)
5. Documentation: C4 models and Screenshot of DPhoto
6. Auto-Synchronisation
    1. Mobile -> S3 Landing
    2. S3 Landing -> DPhoto backup (support deletion and modification)
7. Other features:
    1. deletion of pictures
    2. re-sync times between device

Models notes
-----------------------------------

### aider / gpt-4.1 (OpenAI)

* code is spot on ! Updates I'm making are only for cosmetic ; tests are written following my preferences, the implementation is passing test on the first try
* looks expensive: ~$1 for a session
* response time to watch, I do have to wait before being able to start reviewing the changes and that disengage me
* architect mode wasn't efficient for simple tasks (as expected ?) ; the weak model was confused
