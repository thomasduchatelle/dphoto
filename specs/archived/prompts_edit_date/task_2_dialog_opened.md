# Task 2: Implement Edit Dates Dialog Opened Action

You are implementing the action that opens the edit dates dialog with a specific album context.

## Requirements (BDD Style)

```
GIVEN the dialog is closed
WHEN I dispatch editDatesDialogOpened with an AlbumId
THEN selectEditDatesDialog returns an open dialog with the appropriate album name

GIVEN the dialog is closed
WHEN I dispatch editDatesDialogOpened with an AlbumId that doesn't exist
THEN selectEditDatesDialog returns a closed dialog
```

## Implementation Details

**Create file:** `web/src/core/catalog/edit-dates/action-editDatesDialogOpened.ts`

**Export the following:**
- Interface: `EditDatesDialogOpened`
- Function: `editDatesDialogOpened(albumId: AlbumId): EditDatesDialogOpened`
- Function: `reduceEditDatesDialogOpened(state: CatalogViewerState, action: EditDatesDialogOpened): CatalogViewerState`
- Function: `editDatesDialogOpenedReducerRegistration(handlers: any)`

**Register action in:** `web/src/core/catalog/catalog-reducer-v2.ts`
- Add to `catalogActions` object
- Add to `CatalogViewerAction` union type
- Add to `reducerRegistrations` array

## Action Interface Specification

```typescript
interface EditDatesDialogOpened {
    type: "editDatesDialogOpened";
    albumId: AlbumId;
}
```

## References

- Task 1: `EditDatesDialogState`, `selectEditDatesDialog` from `edit-dates-dialog-state.ts`
- `web/src/core/catalog/language/catalog-state.ts` - for `CatalogViewerState`, `AlbumId` interfaces
- `web/src/core/catalog/album-delete/action-albumDeleted.ts` - for action pattern examples
- `web/src/core/catalog/tests/test-helper-state.ts` - for test data like `twoAlbums`

## Relevant Principles

> The Action naming convention is the Past Tense Event Naming and is represented in camelCase

> The reducer function that takes two parameters: the current state, and the action interface. Make sure the types are explicit.

> if the interface has a single property on top of type: parameter is that property. Make sure the type is respected.

> TDD must be used: no behaviour should be implemented if there is no test to validate it

## TDD Approach

Write tests first that validate:
1. The action factory creates the correct action object with the provided albumId
2. The reducer sets the editDatesDialog state when album exists in allAlbums
3. The reducer does not set editDatesDialog state when album doesn't exist in allAlbums
4. The selector returns the expected properties after the reducer runs

Use test data from `test-helper-state.ts` like `twoAlbums[0].albumId` for valid album tests and a non-existent albumId for invalid tests.
