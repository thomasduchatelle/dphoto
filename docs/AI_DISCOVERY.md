AI Discovery
=======================================

This file is not related to DPhoto, it is my learning notes on how IA is leveraged on the project.

IA Concepts
---------------------------------------

1. Workflow: _how to integrate IA to each step of the development process (what prompts, what tools) ?_
2. Models: _what models perform better on what type of tasks?_
3. Tools: _how to interact with an IA agent and get to update the code ?_

LLM Tradeoffs
---------------------------------------

Starting point is an **idea** -> target is **accepted code**. Adding the feature on existing code is orders of magnitude more complex: the LLM cannot see the
whole context.

* **Effort** vs Autonomy: how much effort is expected by the Supervisor ? ... to break down the tasks and bring engineering expertise.
* **Costs**: how much is the budget by feature to spend on AI models ?
* **Timeline**: how quickly features should be shipped ? ... onboarding has a ramping costs in time

With unlimited budget, Claude Sonnet 4 + Aider would be my preference so far. But if budget is restricted, over 3 cents per request is getting high too quick.
So investigation to have next:

1. how much can 3 cents do with Claude Sonnet 4?
    1. if it goes from idea to accepted code, it would be worth it !
    2. if it doesn't, would a cheap/free model to collect requirements, and one cheap to do the finish on the code could be an approach ?
2. is there cheaper alternative from Claude Sonnet 4 that would give consistent good results ? (need to benchmark the instructions)
3. how much can do an opensource / cheap model ? how much effort would be required to break down the dev cycle into chunks that can be covered by a free/cheap
   LLM ?

Prompts (wip)
---------------------------------------

#### A planning that only an expensive model can perform

##### With `docs/principles_web.md`:

```
Summary:
----------------------------------------------------------------------
Model                                    Cost         Time
----------------------------------------------------------------------
openrouter/x-ai/grok-3-beta                                  $   0.08000       0:39     4/5 - EditDatesDialogState does not have the dates | tests good.
openrouter/anthropic/claude-sonnet-4                         $   0.08000       0:52     5/5 - SPOT ON: State + Tests + File structure (including the selector interfacee in the selector file) + index !
openrouter/google/gemini-2.5-pro-preview                     $   0.08000       1:47     5/5 - State + Tests + File structure
gpt-4.1                                                      $   0.07000       0:25     3/5 - EditDatesDialogState does not have the dates, selector not in its file (and action badly registered)
o3                                                           $   0.05000       0:37     1/5 - EditDatesDialogState does not have the dates | tests does test the reducer
openrouter/google/gemini-2.5-pro-preview-05-06               $   0.05000       0:57     3/5 - Too many tests
gpt-4.1-mini                                                 $   0.03000       0:28     ZER - no actions, no selector, no tests
openrouter/google/gemini-2.5-flash-preview-05-20             $   0.00370       0:23     4/5 - EditDatesDialogState does not have the dates | Tests + File structure are good
openrouter/deepseek/deepseek-r1-0528:free                    $   free          2:30     1/5 - NO TESTS
openrouter/deepseek/deepseek-chat-v3-0324:free               $   free          1:14     3/5 - Testing intermediatary state props by props instead of testing through the selector
openrouter/meta-llama/llama-4-maverick:free                  $   free          0:15     ZER - Gives advises but cannot code.
----------------------------------------------------------------------
```

##### Without the `webprinciples_web.md`:

```
------------------------------------------------------------------------------------------
Model                                                        Cost         Time
------------------------------------------------------------------------------------------
openrouter/anthropic/claude-sonnet-4                         $    0.1200       1:10     4/5 - EditDatesDialogState does not have the dates ; weird test names and props by props testing
openrouter/google/gemini-2.5-pro-preview                     $    0.0800       2:17     4/5 - too many and complex tests
openrouter/x-ai/grok-3-beta                                  $    0.0600       0:33     5/5 - tests props by props
gpt-4.1                                                      $    0.0300       0:23     5/5 - tests props by props
openrouter/mistralai/mistral-large-2411                      $    0.0300       0:43     2/5 - EditDatesDialogState does not have the dates ; no tests 
o3                                                           $    0.0200       0:13     aborted because of missing file
openrouter/deepseek/deepseek-r1                              $    0.0088       2:03     2/5 - no tests
openrouter/qwen/qwen3-235b-a22b                              $    0.0076       3:10     failed to provide valid output
openrouter/google/gemini-2.5-flash-preview-05-20             $    0.0026       0:17     2/5 - no tests
openrouter/qwen/qwen3-30b-a3b:free                                  free       2:52     ZER - no reducer, no tests
------------------------------------------------------------------------------------------
Total                                                        $    0.3590      13:46
```

Conclusions:

1. The principles is not absolutely required at this phase but is ironing the style (especially testing)
2. two next steps are possible:
    1. get more from `claude-sonnet-4` (and other mid-range expensive) -> skip this step and get more done in one go.
    2. break down more to use with `gemini-2.5-flash-preview-05-20`: seems the prompt is not to far away to get it perfect (or is peer programming accepted)
    3. find what the free / opensource LLM are conformable with - this is not the right level

#### Refinement / Story mapping

```
------------------------------------------------------------------------------------------
Model                                                        Cost         Time
------------------------------------------------------------------------------------------
openrouter/anthropic/claude-sonnet-4                         $    0.0500       0:41     4/5 - 5 stories covering everything (except close), a bit technical is places (action and state)
openrouter/x-ai/grok-3-beta                                  $    0.0400       0:22     4/5 - 5 stories covering everything (except close)
openrouter/google/gemini-2.5-pro-preview                     $    0.0400       1:10     4/5 - 5 stories incl closing modal, output messed up in the files (with aider)  
o3                                                           $    0.0300       0:54     3/5 - 7 stories, BDD is very technical
gpt-4.1                                                      $    0.0200       0:18     NEED TO RERUN - didn't write the files and asked for review (6 stories)
openrouter/mistralai/mistral-large-2411                      $    0.0200       1:04     1/5 - 7 stories with no AC and no details leaving too much room for interpretation (misswrite the files) 
openrouter/deepseek/deepseek-r1                              $    0.0073       0:47     2/5 - 7 stories, a bit confused with power user and system
gpt-4.1-mini                                                 $    0.0045       0:26     3/5 - 6 stories, the refresh and save are 2 stories
openrouter/meta-llama/llama-4-maverick                       $    0.0030       1:05     3/5 - 6 stories, lack of details with complex behaviour
openrouter/qwen/qwen3-235b-a22b                              $    0.0016       0:42     4/5 - 5 stories, cover everything with good level of details but didn't write it into files
openrouter/google/gemini-2.5-flash-preview-05-20             $    0.0009       0:07     NEED TO RERUN: only 1 story
openrouter/qwen/qwen3-30b-a3b:free                                  free       2:04     2/5 - 7 stories with fuzy boundaries
------------------------------------------------------------------------------------------
Total                                                        $    0.2173       9:44
```

Need to update the system prompt:

* no technical details (things that are not seen by user)
* be explicit, do not leave anything for interpretation
* write into the files, no review, respect the requested format
* closing case must be added to the original requirements, Gemini pro found it

#### Design and planning v2

    ## Structural Improvements

    **1. Add a Quick Context Section at the top:**
    ```markdown
    **Context**:
    - Architecture: [brief description of layers - actions/thunks/selectors/components]
    - Codebase conventions: [key naming patterns, folder structure]
    - Story scope: [what's in/out of scope]
    ```
    
    **2. Streamline the Design Phase:**
    Instead of asking for "layers impacted" separately, use a template:
    ```markdown
    1. **Design Phase** - For the story "[story text]", identify:
       * State changes needed (interfaces/types to add/modify)
       * Actions required (name + payload schema)
       * Thunks required (name + signature)
       * Selectors required (name + return type)
       * UI components required (name + purpose)
       * Data flow: [brief description]
    ```
    
    **3. Add Design Constraints upfront:**
    ```markdown
    **Design Constraints**:
    - Dialog state must include `open: boolean` property
    - Use Dayjs for date types
    - Read-only data should come from selectors, not state duplication
    - [other architectural rules]
    ```
    
    ## Content Improvements
    
    **4. Provide BDD template immediately:**
    ```markdown
    **Task Format Template**:
    GIVEN [initial state description]
    WHEN [action/thunk name with payload]
    THEN [selector return value or UI behavior]
    ```
    
    **5. Add explicit scope boundaries:**
    ```markdown
    **Story Scope**: Opening and displaying dialog only
    **Out of Scope**: Closing dialog, editing functionality, validation
    ```
    
    **6. Pre-define prompt structure:**
    ```markdown
    **Final Prompt Structure** (for reference):
    - Introduction: Component type and name
    - Requirements: BDD statements
    - Implementation: Folder, TDD, files to edit
    - Interface: Signatures and types
    - References: Related tasks/files
    ```
    
    This would reduce the back-and-forth iterations and get to the task breakdown faster while maintaining quality.

#### Design and planning

**Generic prompt in input - generic plan in output**. None of the output is usable. The analysis is far too much generic. The tasks are wide open and unprecise.
Let's try to fix that:

```
/read-only web/src/pages/authenticated/albums/CreateAlbumDialog
/ask

Let's start refining the language + selectors + actions and replace the steps. And we're going to start by the Date Edit.

Start by defining what properties the Dialog would require to function (I added the create dialog as an example).

These properties are what the selector interface must returns. You can define now:
* the interface of the properties (give it an explicit name)
* the signature of the selector (give it an explicit name)
* the property that is added to the main state (give it a name, and define its interface)

From there, think about the thunks (any user interaction) to find how they will mutate with the state (re-read the scenarios). These are actions. List them, define what in
put parameters they will requires.
Write your notes about the thunks, we will use them on the next phase of our deisgn and planning.

Once you have the list, create one task per actions. Each task will have the requirements written using BDD looking like:

    GIVEN <description of the initial state>
    WHEN <name of the action dispatched and description of its payload>
    THEN <description of what will return the selector>

Also add the interfaces and signatures you defined early on the first task: it needs to create them. Then reference in which file they are for nteh following actions (one
per step).
```

The analysis was good but the tasks are not. They are "write this class", "write this function" type in a pretty weird BDD style.

And the scope was far too wide with an integration far too late. I'll try to fix it with:

```
The analysis is good, but I'm not convienced with the task break down: they don't look like independently implementable and testable. They look like a "write this code" type of
 tasks which lead to difficult integration and dead code.

We're going to reduce the scope into a vertical slice: "As a user, I can open and close the Edit Date Dialog that displays the name of the album".
(note to myself: this type of break down need to happen much earlier and need to be added in the prompt)

Then we're going to work more collaboratively:

1. **design**: define the following
  * properties interface used by the dialog (`EditDatesDialogProperties` minus what's not required anymore due to reducing the scope)
  * main state (same comment: remove what's not necessary anymore)
  * list of thunks as you've done it, but define the data it requires. Differentiate the one that are "new" and the ones coming from the state
  * list of actions as you've done it, but define the payload as well
2. **collaboration** ask me a feedback, you might iterate several time before moving to the next step
3. **task breakdown**, each task must be a unit of work:
  * GOOD example (for the action): "GIVEN the dialog is closed WHEN I dispatch the `editDatesDialogOpened` with an AlbumId THEN `selectEditDatesDialog` returns an open dialog with the appropriate name"
    it's good because it is describing the behaviour that the thunk that will use it expect, and the feature only exposes what is meant to be exposed: the action.
  * GOOD example (for the UI component): "WHEN the dialog is open THEN it displays the details (name) of the album"
    it's good because it describe one state the dialog can take and is described with user oriented languae
  * BAD example: "GIVEN the system needs to support date editing functionality WHEN defining the state model for the edit dates dialog THEN the selector should return dialog properties for component rendering"
    it's bad because it is not a behaviour expected by the application running, it's a behaviour the developer should have. But we're not programing the developer, we're progra
ming the application.
4. after writing some (2-3), ask me feedback and we can move to the next until it's complete
```

Claude suggested to write the prompts as defined much earlier from the design then BDD. I asked to write them into files one at a time. Collaboration is
definitively a WIN to bring feedback earlier. The prompts looks OK - I'll try to get them implemented ...

Proposed new prompt:

```
We need to implement a new feature using vertical slices and collaborative design. Let's start with a minimal viable slice.

**Feature Scope**: "As a user, I can open and close the Edit Date Dialog that displays the name of the album"

**Process**:
1. **Design Phase** - Define the minimal state model:
   - Dialog properties interface (what the React component needs)
   - Main state structure (what gets stored)
   - User interactions → actions mapping (with payloads)
   - Thunks analysis (new data vs state data requirements)

2. **Collaboration** - Present the design and ask for feedback before proceeding

3. **Task Breakdown** - Create independently implementable tasks where each task describes **application behavior**, not developer tasks:
   - ✅ GOOD: "GIVEN dialog is closed WHEN I dispatch editDatesDialogOpened with AlbumId THEN selectEditDatesDialog returns open dialog with album name"
   - ❌ BAD: "GIVEN system needs dialog support WHEN defining state model THEN selector should return properties"

   Each task must be:
   - A unit of work (1 action OR 1 thunk OR 1 component)
   - Independently testable
   - Described in BDD format focusing on runtime behavior

**Constraints**:
- Follow the architectural principles from the handbook
- Use TDD approach
- Actions/thunks must be registered when developed
- Tasks should build incrementally (state → actions → thunks → components)

**Deliverable**: Markdown files with detailed prompts for each task, including BDD requirements, implementation details, references, and TDD guidance.

Start with the design phase for the minimal slice.
```

##### Task implementation

Try 1: Claude + task spec. Disastrous. No test. No action. No reducer. Context is corrupted and irrecuperable.

Try 2: Claude + task spec + (readonly) `docs/principles_web.md`. Better. Got test. Still not the right folders. Still no action or reducer.

#### Workflow building

... after getting the required document, if any correction have be requested, prompt to improve the prompt:

```
I'm pleased with this version of the requirements.

You remember the first Prompt I wrote (I'll add it back below) ? What would you change to lead on this final version, rather than the first version ?

Initial prompt:
```

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

#### Refactoring: renaming a concept

Pre-requisite: add the files where the concept is declared, and where it is used... Then:

Rename `sharingRemoved` into `albumAccessRevoked`, and `revokeAlbumSharing` into `revokeAlbumAccess`, in all types, interface, functions, tests, file names, ...
Keep the same case convention as it was. Provide the shell script to rename the files afterward. Provide another script using grep to idenitify any other
file that could use the old names so they can be updated on the next step.

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

1. ~~merging two React State into one more consistent to improve user experience of Sharing feature~~
    1. ~~refactoring to break down the massive reducer - needs suggestion and implementation~~
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

Do NOT use to reorganise the files !!

I found that asking to update the files, and then providing the shell script to move them working well. Asking for another script (with grep) to identify the
missed file worked as well.

### Aider / openrouter/meta-llama/llama-4-maverick

Used for generating the requirement documentation. ($0.0042)

Really basic prompt. Need to try to get the the level of details I got to with gemini-2.5-flash. Document is minimalist but complete.

### Aider / openrouter/google/gemini-2.5-flash-preview-05-20:thinking

Used for generating the requirement documentation. ($0.02)

He struggled to follow the instructions: asked two questions at a time, and keep asking me to give it the scenario steps instead of generating it !!

Result is good, lots of details with examples. Maybe too much and could be confusing (list of state management libs !). Eager to get things implemented and
ask (too much) details on the API contract.