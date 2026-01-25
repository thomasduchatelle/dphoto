# Coding Principles: ACTIONS

This document describes the principles and standards for implementing Actions in this repository. Follow these rules strictly. If you must break a rule, add a
`TODO` comment explaining why.

## File Structure

Each Action is placed in its own file named `action-<actionName>.ts`, where `<actionName>` is the camelCase name of the action (starting with a lowercase
letter). The interface defining the action starts with an uppercase letter. Each action file must have an associated test file.

The action file contains:

1. **Action Interface**
   - Defines the shape of the action object, including a `type` property and any other required properties.
   - The `type` property is a string literal unique to the action.
   - The interface name starts with an uppercase letter and matches the action's intent.

2. **Action Creator Function**
   - Named after the action (camelCase).
   - Takes as parameters all properties of the action interface except `type`.
   - Returns an object implementing the action interface, with the correct `type` property.

3. **Reducer Fragment**
   - A function that takes two parameters: the previous state and the action (typed as the action interface).
   - Returns the new state.

4. **Reducer Registration**
   - A function that registers the reducer fragment in the handlers map, keyed by the action's `type`.

The action must have a test-file attached:

- Use the action creator to create actions in tests.
- All other parameters and expected results must remain exactly as in the original tests.

Then, the action needs to be registered in the file `catalog-reducer-v2.ts` :

1. add the creator function to `catalogActions`
2. export the Action Interface
3. adds the action to `CatalogViewerAction` list
4. adds the Reducer Registration to `reducerRegistrations`

---

## Example

Suppose you want to implement an action for deleting an album.

### 1. Action Interface

```typescript
export interface AlbumDeletedAction {
    type: "AlbumDeleted";
    albums: Album[];
    redirectTo?: AlbumId;
}
```

### 2. Action Creator

```typescript
export function albumDeletedAction(props: Omit<AlbumDeletedAction, "type">): AlbumDeletedAction {
    return {
        ...props,
        type: "AlbumDeleted",
    };
}
```

### 3. Reducer Fragment

```typescript
// current and returns are always a `CatalogViewerState` type
export function reduceAlbumDeleted(
    current: CatalogViewerState,
    action: AlbumDeletedAction,
): CatalogViewerState {
    let {albumFilterOptions, albumFilter, albums} = refreshFilters(current.currentUser, current.albumFilter, action.albums);

    if (
        action.redirectTo &&
        !albums.some(album => albumIdEquals(album.albumId, action.redirectTo))
    ) {
        albumFilter =
            albumFilterOptions.find(option =>
                albumFilterAreCriterionEqual(option.criterion, ALL_ALBUMS_FILTER_CRITERION)
            ) ?? DEFAULT_ALBUM_FILTER_ENTRY;
        albums = action.albums
    }

    return {
        ...current,
        albumFilterOptions,
        albumFilter,
        allAlbums: action.albums,
        albums: albums,
        error: undefined,
        albumsLoaded: true,
        deleteDialog: undefined,
    };
}
```

### 4. Reducer Registration

```typescript
export function albumDeletedReducerRegistration(handlers: any) {
    handlers["AlbumDeleted"] = reduceAlbumDeleted as (
        state: CatalogViewerState,
        action: AlbumDeletedAction
    ) => CatalogViewerState;
}
```

### 5. Test File

```
describe("reduceAlbumDeleted", () => {
 const deleteDialog = {deletableAlbums: twoAlbums, isLoading: true};

 const marchAlbum = {
     albumId: {owner: "myself", folderName: "mar-25"},
     name: "March 25",
     start: new Date(2025, 2, 1),
     end: new Date(2025, 2, 31),
     totalCount: 0,
     temperature: 0,
     relativeTemperature: 0,
     sharedWith: []
 };

 const loadedStateWithThreeAlbums: CatalogViewerState = {
     ...loadedStateWithTwoAlbums, // use the loadedStateWithTwoAlbums defined in web/src/core/catalog/domain/tests/test-helper-state.ts as a base state
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
```

## General Coding Principles

- Do not add comments to the functions you create.
- Always use explicit types; never use `any`.
- Inline variables used once as long as it remains readable.
- Prefer asserting whole object in tests (list of dispatched actions, full state, ...), instead of asserting each property individually.

