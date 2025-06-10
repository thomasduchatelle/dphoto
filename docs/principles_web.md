# Coding Principles

## General Principles

These general principles must be **strictly respected**. In the rare exception where one cannot, a FIXME comment must be added to alert the reviewer.

1. Always use explicit types or interfaces, never use `any`
2. In tests, the results must be asserted as a whole, not properties by properties
3. Do not comment the code to paraphrase what it is doing
4. Declared objects must be passed as arguments of a function directly, do not declare a transient variables

## "Action" and "Thunk" Principles

Data cycle in the UI is as follows:

1. Every data rendered by the view is coming from the _State_
    1. a _Selector_ should be used to access a specific subset of the _State_ properties for composite cases (dialogs or modals)
2. User activity, including browser loading and URL change, is triggering the appropriate _Thunk_ (through `onClick` and `useEffect`)
3. The _Thunk_ executes the appropriated request(s) on the server using an _Adapter_, and _dispatches_ _Action(s)_ containing the new data
4. Dispatched _Actions_ are merged into the _State_ by the _Reducer_
5. A change on the _State_ triggers a refresh on the UI (and back to step 1)

### State and Selectors

#### When / Where to create a state ?

1. Properties displayed in several place should be present in a **single state**, never replicated they are passed down as properties / hooks
2. Properties that **change together must be within the same state**
3. States should be **as small as possible** _without infringing the previous rule_
4. States should be as low as possible in the document tree _without infringing previous rules_

#### State is a read-only domain model

The state model is common across the slice of the application, it should not be based on the REST model nor a specific UI component, it should be designed to
represent the best the business.

The state model is **never updated directly**. Any change is done by dispatching an _Action_ which is then reduce by the _Reducer_ function.

Selector function can the provided as part of the model to extract a set of properties used by a component like a Dialog.

### Action

The _Action_ naming convention is the **Past Tense Event Naming** and is represented in **camelCase** (example: `albumCreated`)

An _Action_ is defined by the following parts in a single file named `<feature name>-<action name>.ts` where "feature name" is a lower-case short name used to
regroup related actions, and action name is in camelCase.

#### Action interface

* The interface is the name of the _Action_, starting with an upper case (ex: `AlbumCreated`)
* It defines the shape of the action object, including a `type` property and any other required properties.
* The `type` property is a string literal unique to the action, its value is always the same as the _Action_ name

Complete example:

```typescript
// actions/album-albumDeleted.ts
export interface AlbumDeletedAction {
    type: "AlbumDeleted";
    albums: Album[];
    redirectTo?: AlbumId;
}
```

#### Action Factory

* The function is named after the action with "Action" suffix, in **camelCase** (example: `albumCreatedAction`)
* It returns the action interface
* The parameters are either:
    * if the interface has only no property other than `type`: no parameter
    * if the interface has a single property on top of `type`: parameter is that property. Make sure the type is respected.
    * if the interface has more properties, it takes a single argument of the type `Omit<AlbumCreated, "type"`

Complete example:

```typescript
// actions/album-albumDeleted.ts
export function albumDeletedAction(props: Omit<AlbumDeletedAction, "type">): AlbumDeletedAction {
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
// actions/album-albumDeleted.ts

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
// actions/album-albumDeleted.ts
export function albumDeletedReducerRegistration(handlers: any) {
    handlers["AlbumDeleted"] = reduceAlbumDeleted as (
        state: CatalogViewerState,
        action: AlbumDeletedAction
    ) => CatalogViewerState;
}
```

#### Action testing

* naming convention of "describe" is `action:<action name>` (example: `action:albumCreated`)
* types predefined in test helper must be used where possible (`web/src/core/catalog/domain/tests/test-helper-state.ts`)
* assertions should be done on the result of selectors, and not directly the state
* assertions must be done on the whole result or whole state, not on individual properties

```typescript
// actions/album-albumDeleted.test.ts

import {loadedStateWithTwoAlbums, twoAlbums, marchAlbum} from "tests/test-helper-state";

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

### Thunks

* Naming convention of thunk is an verb with a complement, example: `createAlbum`, or `sendEmail`
* A thunk implement one and only one responsibility: single-responsibility principle
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
* The thunk functions first argument is a `dispatch` function accepting the type of actions the thunk will generate
* The thunk functions second argument is the dependencies (adapters) the port requires
    * if the thunk has no port, the argument is skipped
* The thunk function last argument(s) are the data, it can be a single objects or several arguments

Complete example:

```typescript
// thunks/album-createAlbum.ts

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

    return albumId;
}
```

#### Thunk pre-selector

The thunk might requires new data and data from the state. Data from the state is extracted using the Pre Selector function.

* Pre selector function argument is the main State, and it returns what the thunk requires
* Pre selector function might not be required in which case it returns an empty object

#### Thunk factory

* The factory function arguments are settled: `app` (which contains the adapters) and `dispacth`.
* The factory function returns another function used as handler in the views:
    * handler function arguments are the new data
* Best implementation of the factory function is to return `<thunk function>.bind(null, dispatch, adapters, ...)`

a function returning another function with the selected state context and the properties of `CatalogFactoryArgs` injected (recommended to use
`.bind(null, ...)`).

* Thunks use a function in most cases. It takes only the dependencies it uses (e.g., `dispatch`, ports, context) as arguments, in an order that allows use of
  `.bind(null, ...)` in the factory.
* Use a class if more than one port is used, for readability: pass `dispatch` and ports in the constructor, and pass state context and new values as a
  merged object to the method.
* If the thunk interacts with external systems, define a Port interface to abstract those dependencies.

Complete example of pre-selector and thunk factory:

```typescript
// thunks/album-createAlbum.ts

import type {ThunkDeclaration} from "../../thunk-engine";
import type {CatalogFactoryArgs} from "./catalog-factory-args";
import {CatalogFactory} from "../catalog-factories";
import {DPhotoApplication} from "../../application";

export const createAlbumDeclaration: ThunkDeclaration<
    CatalogViewerState, // this is the global state interface, always this type in 'core/catalog/thunks'
    {}, // this is the type returned by 'selector' 
    (request: CreateAlbumRequest) => Promise<AlbumId>, // be specific for the type, this is the business function minus the injected arguments
    CatalogFactoryArgs // this is the type of the argument of the factory, always this type in 'core/catalog/thunks'
> = {
    // Selector: extracts context from state (none needed here)
    selector: (_state: CatalogViewerState) => ({}),

    // Factory: wires up dependencies and returns the thunk
    factory: ({dispatch, app}) => {
        const restAdapter = new CatalogFactory(app as DPhotoApplication).restAdapter();
        // Bind dispatch, port, and optionally the partial state returned by the selector
        // returns a function that takes only the request
        return createAlbumThunk.bind(null, dispatch, restAdapter);
    },
};

// for reference, CatalogFactoryArgs is defined as follow:
export interface CatalogFactoryArgs {
    app: DPhotoApplication
    dispatch: (action: CatalogViewerAction) => void
}
```

### Testing

* Tests are written against the business function, **not** the `ThunkDeclaration`.
* Use **Fakes** (in-memory implementations) for ports instead of mocks, to decouple tests from adapter signatures.
    * assert write requests by inspecting the fake’s state;
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