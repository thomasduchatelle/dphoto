# Conding Principles: THUNKS

Use the following rules and coding standards to implement thunks in this repository. Follow these rules strictly. If you must break a rule, add a `TODO` comment
explaining why.

## File Structure

A thunk is declared in its own file, which contains:

1. **Business Logic Function or Class**
   - Contains the core logic, calls a Port to interact with the server, and dispatches actions to update state for progress, failure, and success.
   - Use a function in most cases. It takes only the dependencies it uses (e.g., `dispatch`, ports, context) as arguments, in an order that allows use of
     `.bind(null, ...)` in the factory.
   - Use a class if more than one port is used, for readability: pass `dispatch` and ports in the constructor, and pass state context and new values as a
     merged object to the method.
   - If the thunk interacts with external systems, define a Port interface to abstract those dependencies.

2. **Selector**
   - A function taking the `CatalogViewerState` and selecting the context necessary for the thunk implementation to work.

3. **Factory**
   - Returns a function with the selected state context and the properties of `CatalogFactoryArgs` injected (recommended to use `.bind(null, ...)`).

4. **(Optional) Port Interface**
   - Exposes the functions wrapping REST calls, stores, and other technologies. The port interface is instantiated in the factory and injected into the
     business logic.

5. **ThunkDeclaration**
   - Export a `ThunkDeclaration` with the selector and the factory. It is referred to in the `index.ts` file in the `catalogThunks`, and the Port interface (if
     any) is exported.

---

## Examples

### 1. Business Function (with Port Interface if necessary)

```typescript
// Port interface: abstracts external dependencies (e.g., REST calls)
export interface CreateAlbumPort {
    createAlbum(request: CreateAlbumRequest): Promise<AlbumId>;

    fetchAlbums(): Promise<Album[]>;
}

export async function createAlbumThunk( // the function is async only when required
        dispatch: (action: AlbumsLoadedAction) => void, // use 'AlbumsLoadedAction' as the type implemented by a actions raised by thunks in 'core/catalog/thunks'
    createAlbumPort: CreateAlbumPort,
    request: CreateAlbumRequest
): Promise<AlbumId> {  // the function returns null unless explicitely specified in the test cases
    const albumId: AlbumId = await createAlbumPort.createAlbum(request);
    const albums: Album[] = await createAlbumPort.fetchAlbums();
   dispatch(catalogActions.albumsLoadedAction({albums, redirectTo: albumId}));

    return albumId;
}
```

---

### 2. Declaration (Selector and Factory)

```typescript
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

---

### 3. Testing Principles

- Tests are written against the business function, **not** the `ThunkDeclaration`.
- Use **Fakes** (in-memory implementations) for ports instead of mocks, to decouple tests from adapter signatures.
- Assert write requests by inspecting the fake’s state; assert read requests by checking outputs and outcomes.

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

    // Dispatched actions are tested in a sigle assertion of an array (no individually)
    expect(dispatched).toEqual([
        catalogActions.albumsLoadedAction({
            albums: expect.any(Array),
            redirectTo: expect.any(Object)
        })
    ]);
});
```

## General coding principles

- Do not add comments to the functions you create.
- Always use explicit types; never use `any`.
- Inline variables used once as long as it remains readable.
- Prefer asserting whole object in tests (list of dispatched actions, full state, ...), instead of asserting each property individually.