# Coding Principles

## General Principles

These general principles must be **strictly respected**. In the rare exception where one cannot, a FIXME comment must be added to alert the reviewer.

1. Prioritise exact semantic over correct syntax
2. TDD must be used: no behaviour should be implemented if there is no test to validate it (including but not limited to null checking, field validation, ...)
3. Always use explicit types or interfaces, never use `any`
4. In tests, the results must be asserted as a whole, not properties by properties
5. Do not comment the code to paraphrase what it is doing
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

An _Action_ is defined by an interface, a factory, and a reducer function, all in a single file named
`.../<feature in dash-case>/action-<action name in camelCase>.ts`.

#### Action interface

* The interface is the name of the _Action_ in PascalCase (ex: `AlbumCreated`)
* It defines the payload of the action, plus the property `type`
* The `type` property is a string literal unique to the action, its value is always the same as the _Action_ name
* The payload is kept minimum: only what cannot be found on the current state. Examples:
  * Good: ID of the selected object, the rest of the object will be found in the state
  * Good: Value updated from an input field
  * Bad: copy of an object from the state

#### Action Factory

* The function is named after the action, in **camelCase** (example: `albumCreated`)
* It returns the action interface
* The parameters are either:
  * if the interface has no property other than `type`: no parameter
    * if the interface has a single property on top of `type`: parameter is that property. Make sure the type is respected.
    * if the interface has more properties, it takes a single argument of the type `Omit<AlbumCreated, "type">`

Complete examples:

**Case 1: No additional properties**

```typescript
// catalog/album-delete/action-loadingStarted.ts
export interface LoadingStarted {
    type: "LoadingStarted";
}

export function loadingStarted(): LoadingStarted {
    return {
        type: "LoadingStarted",
    };
}
```

**Case 2: Single additional property**

```typescript
// catalog/album-delete/album-errorOccurred.ts
export interface ErrorOccurred {
    type: "ErrorOccurred";
    message: string;
}

export function errorOccurred(message: string): ErrorOccurred {
    return {
        message,
        type: "ErrorOccurred",
    };
}
```

**Case 3: Multiple additional properties**

```typescript
// catalog/album-delete/action-albumDeleted.ts
export interface AlbumDeleted {
  type: "AlbumDeleted";
  albums: Album[];
  redirectTo?: AlbumId;
}

export function albumDeleted(props: Omit<AlbumDeletedAction, "type">): AlbumDeleted {
    return {
        ...props,
        type: "AlbumDeleted",
    };
}
```

#### Reducer function

* The reducer function is named after the action, prefixed by "reduce" (example: `reduceAlbumCreated`)
* The reducer function that takes two parameters: the current state, and the action interface. Make sure the types are explicits.
* It returns the updated state, same type as the first parameter.

Complete example:

```typescript
// catalog/album-delete/action-albumDeleted.ts

// current and returns are always a `CatalogViewerState` type
export function reduceAlbumDeleted(
    current: CatalogViewerState,
    {deletedAlbumId}: AlbumDeletedAction,
): CatalogViewerState {
    return {
        ...current,
        allAlbums: current.allAlbums.filter(album => !albumIdEquals(deletedAlbumId.albumId, album.albumId)),
        albums: current.albums.filter(album => !albumIdEquals(deletedAlbumId.albumId, album.albumId)),
        error: undefined,
        albumsLoaded: true,
        deleteDialog: undefined,
    };
}
```

#### Reducer Registration

A function that registers the reducer function in the handlers map, keyed by the action's `type`.

Then, the action needs to be registered in the file `catalog-reducer-v2.ts` :

1. add the creator function to `catalogActions`
2. export the Action Interface
3. adds the action to `CatalogViewerAction` list
4. adds the Reducer Registration to `reducerRegistrations`

Complete example:

```typescript
// catalog/album-delete/action-albumDeleted.ts
export function albumDeletedReducerRegistration(handlers: any) {
    handlers["AlbumDeleted"] = reduceAlbumDeleted as (
        state: CatalogViewerState,
        action: AlbumDeletedAction
    ) => CatalogViewerState;
}
```

#### Action testing

* **Every action MUST have tests associated with it - tested in combination of selector(s) and the state.**
* naming convention of "describe" is `action:<action name>` (example: `action:albumCreated`)
* types predefined in test helper must be used where possible (`web/src/core/catalog/tests/test-helper-state.ts`)
* assertions should be done on the result of selectors, and not directly the state
* assertions must be done on the whole result or whole state, not on individual properties

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
        const got = reduceAlbumDeleted(
            {
                ...initialCatalogState(myselfUser),
                deleteDialog,
            },
            albumDeletedAction({albums: twoAlbums})
        );

        expect(got).toEqual({ // always test the COMPLETE state as a single assertion, never test each property independently
            ...loadedStateWithTwoAlbums,
            medias: [],
            mediasLoaded: false,
            mediasLoadedFromAlbumId: undefined,
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
* The thunk functions first argument is a `dispatch` function accepting the specific action interface type(s) that this thunk will dispatch. Use the specific
  action interface (e.g., `AlbumsLoadedAction`) rather than the broad union type (`CatalogViewerAction`) to make the thunk's behavior explicit.
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
    dispatch: (action: AlbumsLoadedAction) => void, // use 'AlbumsLoadedAction' as the type implemented by a actions raised by thunks in 'core/catalog/thunks'
    createAlbumPort: CreateAlbumPort,
    request: CreateAlbumRequest
): Promise<void> {  // the function returns void or Promise<void> unless explicitely specified in the test cases
    const albumId: AlbumId = await createAlbumPort.createAlbum(request);
    const albums: Album[] = await createAlbumPort.fetchAlbums();
    dispatch(catalogActions.albumsLoadedAction({albums, redirectTo: albumId}));
    // Note: catalogActions.albumsLoadedAction() is an action factory that returns an AlbumsLoadedAction object
    // AlbumsLoadedAction is part of the CatalogViewerAction union type
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

Complete example:

```typescript
// catalog/sharing/thunk-grantAlbumSharing.ts

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

#### Thunk Testing

* **Every thunk MUST have tests associated with it.**
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
    const dispatched: any[] = [];

    await createAlbumThunk(dispatched.push.bind(dispatched), fakePort, request);

    expect(fakePort.albums).toContainEqual(expect.objectContaining({name: "Album 1"}));

    // Dispatched actions are tested in a sigle assertion of an array (not individually)
    expect(dispatched).toEqual([
        catalogActions.albumsLoadedAction({
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

The testing strategy follows these principles:

* **Test structure does not match code structure exactly:**
  * **Action Unit tests**: State + single action + selector are tested together to fulfill a requirement
  * **Behavior tests**: Sequence of several different actions tested when risk of collision between actions is identified
  * **Thunk unit tests**: Thunk function tested independently
  * **Adapter unit tests**: Adapters tested independently
  * **Acceptance tests**: Application tested as early as possible (without browser) to as far as possible (without actual API backend, using wiremock or
    equivalent)
  * **End-to-end tests**: Integration validation through one or two critical paths that must never fail

* **TDD principle**: Implementation should **never** have behavior that hasn't been expected or forced by a test case. Without an appropriate test, code must
  remain extremely simple, even if it means it is wrong.

* **Test as low as possible**: Everything should be tested as unit tests when possible. Higher-level tests (behavior, acceptance, e2e) only cover what couldn't
  be tested at the unit level. The goal is robust tests
  that provide high confidence that refactoring hasn't broken anything, not 100% code coverage.

### Test Selection Criteria

**Action Unit tests** - Use when:

- Testing a single user action and its immediate state change
- Validating reducer logic and selector output
- The action operates independently of other actions

**Behavior tests** - Use when:

- Multiple actions modify the same state properties
- Actions have dependencies or ordering requirements
- Risk of state corruption when actions are combined
- Complex workflows spanning multiple user interactions

**Thunk unit tests** - Use for:

- All thunk functions (mandatory)
- Business logic validation
- Error handling scenarios
- Adapter interaction verification

**Adapter unit tests** - Use for:

- All adapters (mandatory)
- Data transformation logic
- External API contract validation

**Acceptance tests** - Use for:

- Critical user journeys
- Integration between major components
- Regression prevention for key features

**End-to-end tests** - Use for:

- Authentication flow
- One primary happy path per major feature
- Critical business processes that must never break
