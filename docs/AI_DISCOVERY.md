AI Discovery
=======================================

This file is not related to DPhoto, it is my learning notes on how IA is leveraged on the project.

IA Concepts
---------------------------------------

1. Workflow: _how to integrate IA to each step of the development process (what prompts, what tools) ?_
2. Models: _what models perform better on what type of tasks?_
3. Tools: _how to interact with an IA agent and get to update the code ?_

### Prompts (wip)

#### Prompting the principles

```
I'm editing the main development handbook of the project: `docs/principles_web.md`.

The document purpose is to be read by LLM in order to:

1. decouple the requirement into the expected, and defined in the document, concepts
2. adopt coding style of the project and reduce the number of edit required after LLM propositions
3. save time when developer is prompting LLM

The document must be directive and leave little to interpretation.
It must be clear and useful for LLMs.

---

You will review the document against the requirement above. You'll ask me questions to clarify points that needs clarification. Ask me one question at a time we can document a thorough, concise, and clear documenta
tion.
Our end goal is to update the document I can handoff to developers and LLM when I'm developing features
```

#### Refactoring: re-applying the principles

The principles from `docs/principles_web.md` have been breached when writing the actions of the catalog domain in `web/src/core/catalog/domain/actions`.

Draft a detailed inventory of issues found in the code, group them by unit (an action), prioritise the actions from the one having deep design issues to those
only requiring cosmetic changes. Once you have a solid and prioritised list of units, break it into small steps. Each step must be small enough to be
implemented
safely with strong testing, but big enough to move the project forward. Each step should keep the existing tests as guaranty no regression is introduced, and
remove the redundant one on the next step. Iterate until you feel that the steps are right sized for this project.

From here you should have the foundation to provide a series of prompts for a code-generation LLM that will implement each step. Prioritize best practices, and
clear instructions, and incremental progress, ensuring no big jumps in complexity at any stage. Make sure that each prompt list the files requiring to be
changed, builds on the previous prompts of the same unit, and ends with wiring things together. There should be no hanging or orphaned code that isn't
integrated into a previous step.

Make sure and separate each prompt section. Use markdown. Each prompt should be tagged as text using code tags. The goal is to output prompts, but context, etc
is important as well.

##### Follow up

Break down the steps further so no action is worked at the same time as another action. Start with high priority ones. Make sure the naming convention (past
tense event) is respected and integrate the renaming ask>  in the plan otherwise. Don't keep the tests for the end, each step must include the test required to
validate it. When functions or classes are renamed, reference to them must be updated as well.

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

#### UI: Thunks

I'd like to adapt the Catalog loader to be used as a Thunk-like. I'm am NOT using Redux on the project.

The thunk should be usable in the React component as follow:

 ```
 const thunks = useThunks(catalogThunks);

 useEffect(() => {
   thunks.onPageRefresh(albumId);
 }, [thunks, albumId]);
 ```

`thunks` must be stable (same reference when refreshed.

And the implementation of the thunk should be a function (or a method):

 ```
 export function async onPageRefresh({albumId, allAlbums, ...otherPropsFromState}: OnPageRefreshProps) {
         const medias = await fetchMediaPort.fetchMedias(albumId);
 ^Idispatch(catalogActions.mediaLoaded({albumId, medias})
 }
 ```

Where `fetchMediaPort` and `dispatch` are somehow injected. The function could be in a class and they become `this.fetchMediaPort`
and `this.dispatch`.

Can you give several options on how to implement the middleware necessary to get this behaviour, and your recommendation ?

A special attention should be placed on how the injection would work (through the implementation of `useThunks` and what is exported as `catalogThunks`.

##### Refinement

Let's explore the option 1.

First level would be the parameters that are not changing once the component (context) is mounted: dispatch and API adapters. Then the state to only have the
required parameters exposed (albumId here).

I would think to use a class with the context parameters set as constructor arguments:

 ```
 // catalog/thunks/thunks-onPageRefresh.ts

 export interface MediaLoaderPort {
   findMedias(albumId: AlbumId): Promise<Media[]>
 }

 export interface OnPageRefreshArgs {
   albumId: AlbumId
   allAlbums: Album[]
   // ...
 }

 export class OnPageRefresh {
   constructor(private dispatch: any, // todo find the appropriate type
               private mediaLoaderPort: MediaLoaderPort,
               ) {}

   const onPageRefresh = async ({albumId, allAlbums, ...others}: OnPageRefreshArgs => {
     // ...
   }
 }
 ```

Then I would have a function indicating what parameters are coming from the state:

```
 // catalog/thunks/thunks-onPageRefresh.ts

 export const onPageRefreshSelector = ({allAlbums}: CatalogViewerState): Omit<OnPageRefreshArgs, "albumId"> => ({
   allAlbums,
 })
```

The naive way to put that together is:

```
 // catalog-react/thunks.ts

 export const useCatalogThunks = () => {
   const {dispatch, state} = useContext(CatalogContext)

   return {
     onPageRefresh: (albumId: AlbumId) => new OnPageRefresh(dispatch, catalogAPIAdapter).onPageRefresh({...onPageRefreshSelector(state), albumId}),
   }
 }
```

But I would need to make the `onPageRefresh` stable when anything other than the selected partial state is changing. How would you suggest to do it ?

##### Drafts

The second level is to be state aware so a function only taking the relevant parameters is exposed (only take the new parameters):

 ```
 // catalog/thunks/index.ts

 export const thunkFactories = {
   onPageRefresh: (callback, selectedState) => {
     return (albumId: AlbumId) => callback({...selectedState, albumId})
   }
 }
 ```

And it is used with a hook:

 ```
 // catalog-react/thunks.ts

 export const useCatalogThunks = () => {
   const {dispatch, state} = useContext(CatalogContext)

   const partialState = onPageRefreshSelector(state)
   const onPageRefresh = useCallback( (albumId: AlbumId) => {
     const callback = new OnPageRefresh(dispatch, catalogAPIAdapter).onPageRefresh
     return callback({...partialState, albumId});
   }, [dispatch, catalogAPIAdapter, ...Object.values(partialState)])

   return {
     onPageRefresh,
   }
 }
 ```

#### UI: Thunk (v2)

Callback used on onClick and similar props on the view are called **Thunk**. A Thunk is characterized by:

* it is a function with any number of arguments which returns either void or Promise<void>
* its arguments are values used to mutate the state, not values that are already in the state, and neither context object (App, credentials, ...)
* the thunk is stable against component refresh

A thunk is declared in its own file which contains:

* business logic function or class: logic of the thunk that call a Port to interact with the server, and dispatch actions to update the state of the progress,
  failure, and success.
    * a function is used in most cases, it takes its dependencies as argument (dispatch, ports, context, ...) in an order allowing to use `.bind(null, ...)` in
      the factory ; only used dependencies are in the arguments
    * a class is used when more than one port is used, for readability: dispatch function and ports are passed on the constructor, state context and new values
      mutating the state are passed a argument to the method as a merged object.
* selector: a function taking the `CalalogViewerState` and selecting the context necessary for the thunk implementation to work
* factory: return a function that have the selected state context and the properties of `CatalogFactoryArgs` injected (recommended to use `.bind(null, ...)`)
* (optional) a Port interface which exposes the functions wrapping REST calls, stores, and other technologies ; the port interface is instantiated in the
  factory and injected into the business logic

A `ThunkDeclaration` is exported with the selector and the factory. It is referred in the index.ts file in the `catalogThunks` and the Port interface, if any,
are exported.

The tests are written against the business logic function or class ; not the `ThunkDeclaration`. Mocks are not used as they coupled the test with the signature
of adapter methods. We use instead Fake: a simple in-memory implementation behaving the same way the abstracted system is expected to. The write requests are
asserted by reading properties of the fake-implementation while the read request are not (only output and outcomes are asserted).

### Summary of previous discussion

Automatically enable it:

```
aider --model gpt-4.1 --max-chat-history-tokens 110000
```

Prompt:

```
+*Briefly* summarize this partial conversation about programming.
+Include less detail about older parts and more detail about the most recent messages.
+Start a new paragraph every time the topic changes!
+
+This is only part of a longer conversation so *DO NOT* conclude the summary with language like "Finally, ...". Because the conversation continues after the summary.
+The summary *MUST* include the function names, libraries, packages that are being discussed.
+The summary *MUST* include the filenames that are being referenced by the assistant inside the ```...``` fenced code blocks!
+The summaries *MUST NOT* include ```...``` fenced code blocks!
+
+Phrase the summary with the USER in first person, telling the ASSISTANT about the conversation.
+Write *as* the user.
+The user should refer to the assistant as *you*.
+Start the summary with "I asked you...".
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
* failed to execute because context was too big (~ 12 files)
* looks good for refactoring but the chat must be cleared on very regular basis, and the number of files very limited ; I used `--map-tokens 1024` to reduce the
  number of token used by the map.

### aider / gpt-4.1-mini (OpenAI)

* (architect: gpt-4.1) aider couldn't apply the changes because the model failed to respect the format ; or some files were changed and it didn't find an exact
  match.
* failed to respect the format when used directly - I wonder if it wasn't because of refactoring with read-only files that it tries to change (?)
* used as weak model for summarisation and git commit looks working well

### OpenAI - recommendation

As of the 10th of June 2025, OpenAI sent an email to recommend the use of gpt-4o for programming purpose. They dropped the price by 80%, matching the price of
4.1.

### Aider / claude sonnet 4

    aider --model openrouter/anthropic/claude-sonnet-4

