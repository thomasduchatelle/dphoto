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

#### UI: action patterns

The reducer has become a very large switch case and cannot be maintained anymore. I want it to be breakdown following principles.

Each Action is placed in its own file 'action-<name of the action>.ts' with its associated test. The action name is always in camelCase, starting with lower
case, except the interface defining it that starts with an upper case.

The action file contains:

* the interface defining the action with a 'type' and other properties, they should be copied without changes
* the reducer fragment, a function taking 2 parameters: the previous state, and action (of the type of interface), and returning the new state
* the action function: named after the type of the action, it takes as parameters the action interface (except 'type' property) and returns an object
  implementing the Action interface

You need to make an exact copy the implementation of the reducer fragment from the existing 'catalogReducerFunction', and copy the tests relevant to this action
from catalog-reducer.test.ts. Make sure the parameters passed to the reducer fragment are the same one as on the original test. Use the action function to
create the action, but the other params must be exactly the same. The result of the reducer fragment must be exactly the same as it was defined on the original
test.

Do not add comments on the functions you're creating.

Register the reducer fragment in web/src/core/catalog/domain/catalog-reducer.ts.

Update the file web/src/core/catalog/domain/catalog-index.ts to export:

* all action interfaces
* an "catalogActions" object with each action function as property
* the catalog reducer which is a conventional reducer function: parameters are current state and an action of teh type of one of the supported Action type.

#### UI: Reducer mechanics

A function 'createReducer' will be created and used for the catalogReducer. It takes an object with one property per action supported: property name is the name
of the action, and the value is the related reducer fragment.

You will not delete existing code, only copy the code to the code. You'll proceed one action at a time following the following list. Once the changes for one
action are completed and I approved them, I'll tell you to move to the next one.

The actions list is:

1. AlbumsAndMediasLoadedAction
2. AlbumsLoadedAction
3. MediaFailedToLoadAction
4. NoAlbumAvailableAction
5. StartLoadingMediasAction
6. AlbumsFilteredAction
7. MediasLoadedAction
8. OpenSharingModalAction
9. AddSharingAction
10. RemoveSharingAction
11. CloseSharingModalAction
12. SharingModalErrorAction

##### Focus on the index and the creation of the reducer

Create a function 'createReducer' in web/src/core/catalog/domain/catalog-reducer-v2.ts. It takes an object with one property per action supported: property name
is the name of the action, and the value is the related reducer function. It returns a conventional reducer function that take 2 parameters (the current state
and the action that must implement one of teh supported action interfaces) and returns the new state.

The two actions supported so far are AlbumsAndMediasLoadedAction and AlbumsLoadedAction.

In web/src/core/catalog/domain/catalog-index.ts, exports:

* all the supported action interfaces
* a "catalogActions" object with each function that creates an action instance set as a property
* the catalog reducer which is built from the 'createReducer' function with the list of the supported actions.

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
* failed to execute because context was too big (~ 12 files)
* looks good for refactoring but the chat must be cleared on very regular basis, and the number of files very limited ; I used `--map-tokens 1024` to reduce the
  number of token used by the map.

### aider / gpt-4.1-mini (OpenAI)

* (architect: gpt-4.1) aider couldn't apply the changes because the model failed to respect the format ; or some files were changed and it didn't find an exact
  match.
* failed to respect the format when used directly - I wonder if it wasn't because of refactoring with read-only files that it tries to change (?)
* used as weak model for summarisation and git commit looks working well
* 
