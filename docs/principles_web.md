# Coding Principles

## General Principles

These general principles must be **strictly respected**. In the rare exception where one cannot, a FIXME comment must be added to alert the reviewer.

1. Prioritise exact semantic over correct syntax
2. TDD must be used: no behaviour should be implemented if there is no test to validate it (including but not limited to null checking, field validation, ...)
3. Always use explicit types or interfaces, never use `any`
4. In tests, the results must be asserted as a whole, not properties by properties
5. Do not comment the code to paraphrase what it is doing, nor to explain your change
6. Declared objects must be passed as arguments of a function directly, do not declare a transient variables

## Repository Tree

The file structure is as follows:

* `components/`
    * `catalog-react/` - contains the React Components used to integrate the domain to the other components
* `core/catalog/`
    * `<feature name in dash-case>/` - each feature has a folder containing related actions, thunks, and selectors
    * `language/` - ubiquitous language and definition of the State shared for the domain
    * `common/` - functionalities reused across most features
    * `actions.ts` - where the action interface and the partial reducer are registered
    * `thunks.ts` - where the thunks are registered
    * `adapters/api/` - where the REST API adapter are implemented

## "Action" and "Thunk" Design Pattern

In the UI, the data flow through:

1. **State** - domain model, and ubiquitous language, containing the data used to render the views

2. **Selector** (optional) - function returning a purpose oriented data structure. It must be used to access a complex subset of the _State_ properties for
   composite cases (dialogs or modals)

3. **View** - UI components in TSX (React) rendering the data from the selectors, or from the State. Callbacks (`onClick` and `useEffect`) triggers _Thunks_.

4. **Thunk** - function triggered by a user activity (including browser loading and URL change): often the value of `onClick` and `useEffect`.
   It executes the appropriate requests (through a port) and _dispatches Actions_.

    * **Port** - an interface declared next to the Thunk to abstract specific technologies (REST API, Local store, ...)
    * **Adapter** - a class implementation of the _Port_ for a specific technology (example: Axios)

5. **Action** - interface with a `type: "<action type>"` and a payload ; **dispatching an action is the only way to mutate the _State_**
   The _Reducer_ is a function taking the current state and the action, and returning the mutated state.

A change on the _State_ triggers a refresh on the UI (and back to step 1).

### Actions

The _Action_ naming convention is the **Past Tense Event Naming** and is represented in **camelCase** (example: `albumCreated`)

An _Action_ is defined using the `createAction` function in a single file named
`.../<feature in dash-case>/action-<action name in camelCase>.ts`.

#### Action Definition

* Actions are created using `createAction<StateType, PayloadType?>` from `@light-state`
* The action name is a string literal unique to the action, matching the _Action_ name
* The payload is kept minimum: only what cannot be found on the current state. Examples:
    * Good: ID of the selected object, the rest of the object will be found in the state
    * Good: Value updated from an input field
    * Bad: copy of an object from the state
* The reducer function receives the state and payload directly as parameters

Complete examples:

**Case 1: No payload**

```typescript
// catalog/album-delete/action-loadingStarted.ts
import {createAction} from "src/libs/daction";

export const loadingStarted = createAction<CatalogViewerState>(
    "LoadingStarted",
    (current: CatalogViewerState) => {
        return {
            ...current,
            isLoading: true,
        };
    }
);

export type LoadingStarted = ReturnType<typeof loadingStarted>;
```

**Case 2: Single payload**

```typescript
// catalog/album-delete/action-errorOccurred.ts
import {createAction} from "src/libs/daction";

export const errorOccurred = createAction<CatalogViewerState, string>(
    "ErrorOccurred",
    (current: CatalogViewerState, message: string) => {
        return {
            ...current,
            error: message,
            isLoading: false,
        };
    }
);

export type ErrorOccurred = ReturnType<typeof errorOccurred>;
```

**Case 3: Multiple properties**

```typescript
// catalog/album-delete/action-albumDeleted.ts
import {createAction} from "src/libs/daction";

interface AlbumDeletedPayload {
    albums: Album[];
    redirectTo?: AlbumId;
}

export const albumDeleted = createAction<CatalogViewerState, AlbumDeletedPayload>(
    "AlbumDeleted",
    (current: CatalogViewerState, {albums, redirectTo}: AlbumDeletedPayload) => {
        return {
            ...current,
            allAlbums: albums,
            albums: albums,
            error: undefined,
            albumsLoaded: true,
            deleteDialog: undefined,
        };
    }
);

export type AlbumDeleted = ReturnType<typeof albumDeleted>;
```

#### Action testing

* **Every action MUST have tests associated with it - tested in combination of selector(s) and the state.**
* naming convention of "describe" is `action:<action name>` (example: `action:albumCreated`)
* Typical unit-test will:
    1. initiate a state using the helpers in `web/src/core/catalog/tests/test-helper-state.ts`: only the minimum properties should be set on top of the helpers
    2. execute the reducer
    3. execute the selector: the state is considered private and is not asserted directly
    4. assert the whole result of the selector, not on individual properties

```typescript
// catalog/album-delete/action-albumDeleted.test.ts

import {loadedStateWithTwoAlbums, twoAlbums, marchAlbum} from "../tests/test-helper-state";

describe("action:albumDeleted", () => {
    const deleteDialog = {deletableAlbums: twoAlbums, isLoading: true};

    const loadedStateWithThreeAlbums: CatalogViewerState = {
        ...loadedStateWithTwoAlbums,
        allAlbums: [...twoAlbums, marchAlbum],
        albums: [...twoAlbums, marchAlbum],
    }

    it("closes the dialog and update the lists of albums list like an initial loading", () => {
        const action = albumDeleted({albums: twoAlbums});
        const got = action.reducer(
            {
                ...initialCatalogState(myselfUser),
                deleteDialog,
            },
            action
        );

        expect(listOfAlbumsSelector(got)).toEqual({ // always test the COMPLETE selection as a single assertion, never test each property independently
            loading: false,
            albums: twoAlbums,
            filter: loadedStateWithTwoAlbums.filter,
        });
    });
});
```

### Selectors

* Selectors are returning a called `<name of the selector>Selection` (example:
  `function sharingDialogSelector(state: CatalogViewerState): SharingDialogSelection`)
* Data in the main state is **duplicated ONLY to be edited**, IDs are used as reference.
  Example:
    ```
    export interface State {
        albums: []Album,
        currentAlbum: AlbumId, // ID is used to not duplicate data
    }
    
    export function currentAlbumNameSelector({albums, currentAlbum}: CatalogViewerState): CurrentAlbumNameSelection {
        return { name: albums.find(album => isAlbumIdEquals(album, currentAlbum))?.name }
    }
    ```
* Data is **transformed in the selectors**, not in the UI components.
  Example:
  ```
  export function createDialogSelector({createDialog}: CatalogViewerState): CreateDialogSelection {
      if (!createDialog) return { open: false };

      return { capitalizedName: capitalize(createDialog.name) }
  }
  ```

### Thunks

* Naming convention of thunk is a verb with a complement, example: `createAlbum`, or `sendEmail`
* A thunk implement one and only one responsibility: single-responsibility principle (good example: `renameAlbum`, bad example: `updateAlbum`)
* Thunks are **pure-business logic functions** handling the user commands triggered by the view: they do not use any framework of technology
    * adapters are injected to access lower-level technologies like REST API or Local Storage
* Thunks are stateless
* Thunks dispatch one or several action(s) during the processing. The actions must be as concise as possible to represent the change and contain all the
  data required to update the state.

#### Thunk function

* The thunk function is named after the thunk, suffixed by `Thunk`, example: `createAlbumThunk`
* The thunk function implements the business logic by executing Adapter methods, and dispatching actions to update the state for progress, failure, and/or
  success.
    * Adapters naming conventions is `<ThunkName>Port` (example: `DeleteAlbumPort`)
* The thunk functions first argument is a `dispatch` function accepting the specific action type(s) that this thunk will dispatch. Use the specific
  action type (e.g., `AlbumsLoaded`) rather than the broad union type (`CatalogViewerAction`) to make the thunk's behavior explicit.
* The thunk functions second argument is the dependencies (adapters) the port requires
    * if the thunk has no port, the argument is skipped
* The thunk function last argument(s) are the data, it can be a single object or several arguments

Complete example:

```typescript
// catalog/album-delete/thunk-createAlbum.ts

// Port interface: abstracts external dependencies (e.g., REST calls)
export interface CreateAlbumPort {
    createAlbum(request: CreateAlbumRequest): Promise<AlbumId>;

    fetchAlbums(): Promise<Album[]>;
}

export async function createAlbumThunk( // the function is async only when required
    dispatch: (action: AlbumsLoaded) => void, // use 'AlbumsLoaded' as the type implemented by actions raised by thunks in 'core/catalog/thunks'
    createAlbumPort: CreateAlbumPort,
    request: CreateAlbumRequest
): Promise<void> {  // the function returns void or Promise<void> unless explicitely specified in the test cases
    const albumId: AlbumId = await createAlbumPort.createAlbum(request);
    const albums: Album[] = await createAlbumPort.fetchAlbums();
    dispatch(albumsLoaded({albums, redirectTo: albumId}));
    // Note: albumsLoaded() is an action creator that returns an AlbumsLoaded action object
    // AlbumsLoaded is part of the CatalogViewerAction union type
}
```

#### Thunk pre-selector

A thunk have a `payload` and/or access to the current state. It might not require any.

* The Payload are the attributes from the UI components, usually IDs referencing object in the state or new values to set in the state
* The pre selector function argument is the main State, and it returns the properties required by the thunk
* The pre selector function might not be required in which case it returns an empty object
* The pre-selector's output is passed to the factory as `partialState` and gets bound to the thunk function

#### Thunk factory

The factory function wires up dependencies and returns the thunk handler used by views. The data flow is:
**State → Selector extracts needed data → Factory binds it → View provides remaining data → Thunk executes**

**Parameter passing patterns:**

* **Case 1: Simple case** (≤1 adapter, ≤2 total arguments from state+view)
  Pass arguments individually: `dispatch, adapter, stateData, viewData`

* **Case 2: Multiple adapters** (>1 adapter)
  Pass adapters as single object: `dispatch, {adapter1, adapter2}, stateData, viewData`

* **Case 3: Complex case** (≥3 arguments OR optional arguments)
  Merge state and view data into single composite object

* **Case 4: function solely dispatches a single action** (i.e., its logic is `({dispatch}) => (...args) => { dispatch(actionCreator(...args)); }`)
  Use `createSimpleThunkDeclaration` to declare it. Do not add tests.

Complete examples:

```typescript
// catalog/sharing/thunk-grantAlbumSharing.ts
import {ThunkDeclaration} from "src/libs/dthunks";

export interface GrantAlbumSharingAPI {
    grantSharing(albumId: AlbumId, email: string): Promise<void>;
}

export function grantAlbumSharingThunk(
    dispatch: (action: CatalogViewerAction) => void,
    sharingAPI: GrantAlbumSharingAPI,
    albumId: AlbumId | undefined,
    email: string
): Promise<void> {
    // ... implementation
}

export const grantAlbumSharingDeclaration: ThunkDeclaration<
    CatalogViewerState, // this is the global state interface, always this type in 'core/catalog/thunks'
    { albumId?: AlbumId }, // this is the type returned by 'selector'
    (email: string) => Promise<void>, // be specific for the type, this is the business function minus the injected arguments
    CatalogFactoryArgs // this is the type of the argument of the factory, always this type in 'core/catalog/thunks'
> = {
    // Selector: extracts albumId from the share modal state
    selector: ({shareModal}: CatalogViewerState) => ({
        albumId: shareModal?.sharedAlbumId,
    }),

    // Factory: wires up dependencies and returns the thunk
    factory: ({dispatch, app, partialState: {albumId}}) => {
        const sharingAPI: GrantAlbumSharingAPI = new CatalogAPIAdapter(app.axiosInstance, app);
        // Case 1: Simple case - bind arguments individually
        return grantAlbumSharingThunk.bind(null, dispatch, sharingAPI, albumId);
    },
};

// Cases examples:
// Case 1: Simple case
factory: ({dispatch, app, partialState: {albumId}}) => {
    const sharingAPI = new CatalogAPIAdapter(app.axiosInstance, app);
    return grantAlbumSharingThunk.bind(null, dispatch, sharingAPI, albumId);
}

// Case 2: Multiple adapters
factory: ({dispatch, app, partialState}) => {
    const adapters = {sharingAPI: new SharingAPI(), storageAPI: new StorageAPI()};
    return thunkFunction.bind(null, dispatch, adapters, partialState.data);
}

// Case 3: Complex case
factory: ({dispatch, app, partialState: {albumId}}) => {
    const adapter = new SomeAdapter();
    return (viewData) => thunkFunction(dispatch, adapter, {...partialState, ...viewData});
}
```

Case 4 complete example:

```typescript
// web/src/core/catalog/album-create/thunk-createAlbum.ts
import {createAlbum} from "./action-createAlbum";
import {createSimpleThunkDeclaration} from "src/libs/dthunks";

export const createAlbumDeclaration = createSimpleThunkDeclaration(createAlbum);
```

#### Exporting the thunks

Aggregate Thunks into a Single Export per feature (per folder):

* In the `index.ts` file of the target directory (e.g., `web/src/core/catalog/album-create/index.ts`), create a new `const` export named `[featureName]Thunks` (
  e.g., `albumCreateThunks`).
* Add all the thunks related to the feature
* Add JSDoc Documentation to Aggregated Thunks including a brief description and a list of `Expected handler types` for each thunk within the object. This helps
  LLM agents understand the expected function signature when these thunks are used as handlers.

```typescript
// web/src/core/catalog/album-create/index.ts
import {createAlbumDeclaration} from "./thunk-createAlbum"

/**
 * Thunks related to album creation.
 *
 * Expected handler types:
 * - `createAlbum`: `(name: string) => void`
 */
export const albumCreateThunks = {
    createAlbum: createAlbumDeclaration,
};
```

If the `[featureName]Thunks` is new, add it to the mai thunks:

```
// web/src/core/catalog/index.ts
import {albumCreateThunks} from "./album-create";

export const thunks = {
  ...albumCreateThunks,
  // other aggregated thunk from other features
}
```

#### Thunk Testing

* **Every thunk MUST have tests associated with it.**
    * exception is the case 4 that solely dispatches an action where no test should be added
* Tests are written against the business function, **not** the `ThunkDeclaration`.
* Use **Fakes** (in-memory implementations) for ports instead of mocks, to decouple tests from adapter signatures.
    * assert write requests by inspecting the fake's state;
    * assert read requests by checking outputs and outcomes.

Complete example:

```typescript
// Fake implementation reproduce the expected behaviour of the actual implementation

class CreateAlbumPortFake implements CreateAlbumPort {
    albums: Album[] = [];

    async createAlbum(request: CreateAlbumRequest): Promise<AlbumId> {
        // Simulate album creation
        const albumId = {owner: "myself", folderName: request.forcedFolderName};
        this.albums.push({...request, albumId, ...defaultAlbumValues});
        return albumId;
    }

    async fetchAlbums(): Promise<Album[]> {
        return this.albums;
    }
}

it("should store the new Album and dispatch albumsLoadedAction", async () => {
    const fakePort = new CreateAlbumPortFake([existingAlbum]);
    const dispatched: Action<CatalogViewerState, any>[] = [];

    await createAlbumThunk(dispatched.push.bind(dispatched), fakePort, request);

    expect(fakePort.albums).toContainEqual(expect.objectContaining({name: "Album 1"}));

    // Dispatched actions are tested in a sigle assertion of an array (not individually)
    expect(dispatched).toEqual([
        albumsLoaded({
            albums: expect.any(Array),
            redirectTo: expect.any(Object)
        })
    ]);
});
```

## Adapters

Adapters abstract specific technologies or external systems to keep the codebase pure and technology-agnostic. They follow a clear separation between interface
definition and implementation.

### Port Interface

The **Port** is the interface defined by the thunk based on what the thunk requires. It represents the contract that the thunk needs from external dependencies.

* **Naming convention**: Named after the thunk or the function it fulfills, whichever is most readable
    * Examples: `CreateAlbumPort`, `DeleteAlbumPort`, `FetchAlbumsPort`
* **Location**: Defined in the same file as the thunk that uses it
* **Purpose**: Abstracts external dependencies and makes thunks testable

Complete example:

```typescript
// thunks/album-createAlbum.ts
export interface CreateAlbumPort {
    createAlbum(request: CreateAlbumRequest): Promise<AlbumId>;

    fetchAlbums(): Promise<Album[]>;
}
```

### Adapter Implementation

The **Adapter Implementation** is the concrete class that implements the Port interface, abstracting a specific technology or external system.

* **Naming convention**: Named after the technology or external system it abstracts
    * Examples: `AxiosCatalogRestApi`, `LocalStorageAdapter`, `S3FileAdapter`
* **Location**: Typically in `adapters/` directory, organized by technology or domain
* **Purpose**: Handles the actual communication with external systems (REST APIs, databases, file systems, etc.)
* **Testing**: Adapters should be tested independently to verify their contract compliance and data transformation logic

Complete example:

```typescript
// adapters/api/CatalogAPIAdapter.ts
export class CatalogAPIAdapter implements CreateAlbumPort, DeleteAlbumPort {
    constructor(
        private readonly authenticatedAxios: AxiosInstance,
        private readonly accessTokenHolder: AccessTokenHolder,
    ) {
    }

    async createAlbum(request: CreateAlbumRequest): Promise<AlbumId> {
        const response = await this.authenticatedAxios.post('/api/v1/albums', request);
        return response.data;
    }

    async fetchAlbums(): Promise<Album[]> {
        const response = await this.authenticatedAxios.get('/api/v1/albums');
        return response.data;
    }
}
```

### App Factory

The **App** is a helper class that provides instances of adapters, centralizing dependency creation and configuration.

* **Purpose**: Centralize adapter instantiation and dependency injection
* **Usage**: Used in thunk factories to get properly configured adapter instances

Complete example:

```typescript
// Factory usage in thunk declaration
factory: ({dispatch, app, partialState}) => {
    const catalogAPI: CreateAlbumPort = app.catalogAdapter(); // App provides the adapter instance
    return createAlbumThunk.bind(null, dispatch, catalogAPI, partialState.data);
}
```

## Testing Strategy

Components holding a logic are tested to demonstrate the Acceptance Criteria are fulfilled. All the tests must follow these principles:

1. **TDD principle**: implementations should **never have a behavior that hasn't been expected or forced by a test case**. Without an appropriate test, code
   must remain extremely simple, even if it means it is wrong.

2. **Robust Tests**: tests must be robust to refactoring
    * test actions and selectors together: the state structure can be changed without affecting the tests
    * use pre-defined constants of the states and the selections in a known situation: adding new properties is done on these constants and do not affect the
      tests
    * use fake implementations: method signature of the Ports can change, the fakes are updated, but it does not affect the tests

3. **Unit test first**: the code is structured so most of the acceptance criteria can be validated on a unit of code, and never depends on the integration of
   several layers
    * actions and selectors: jest unit tests
    * thunks: jest unit tests
    * UI components: StoryBook with the component set in each of the relevant situations (examples: "default", "loading", "with a technical error", "success")
    * react hooks: jest + testing-library

4. **Behaviour tests** (also called integration): validate consistency when a unit might affect another unit ; use one when you identify a case that requires a
   sequence of events to be played to validate the system behaviour is user-friendly and predictable.
    * example 1: when a dialog is closed and re-open, the data shown in it should reflect what have been saved previously ; if the data is not reloaded from the
      server, it
      should be tested as behaviour test
    * example 2: when a dialog is open, it should close another dialog ; run the sequence of actions to confirm it behave as expected

5. **Acceptance tests**: validate the integration of the domain and its adapters ; they are limited to the strick minimum viable end to end features.
    * Domain focused acceptance test: instantiate the domain with the adapters implementation and a REST mock server (but without the React libraries). End to
      end scenario will be played by calling the domain functions and asserting the results through selectors.
    * Complete acceptance test: integrate as much as possible of code, interact with the `testing-library` (simulated browser) and backed with a mock rest
      server

The unit tests must always be implemented alongside the new or updated code.

Behaviour tests and Acceptance tests can be suggested as consideration for next steps, with justification that it cannot be done at unit level.
