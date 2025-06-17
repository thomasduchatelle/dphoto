# Task 4: Implement Open Edit Dates Dialog Thunk

You are implementing the thunk that handles the business logic for opening the edit dates dialog.

## Requirements (BDD Style)

```
GIVEN any state
WHEN the openEditDatesDialog thunk is called with an album ID
THEN it dispatches editDatesDialogOpened with that album ID
```

## Implementation Details

**Create file:** `web/src/core/catalog/edit-dates/thunk-openEditDatesDialog.ts`

**Export the following:**
- Function: `openEditDatesDialogThunk(dispatch: (action: EditDatesDialogOpened) => void, albumId: AlbumId): void`
- Constant: `openEditDatesDialogDeclaration: ThunkDeclaration<CatalogViewerState, {}, (albumId: AlbumId) => void, CatalogFactoryArgs>`

**Register thunk in:** `web/src/core/catalog/thunks.ts`
- Add `openEditDatesDialog: openEditDatesDialogDeclaration` to the `catalogThunks` object

## Thunk Declaration Pattern

Follow the pattern from existing thunks:
- Use empty object `{}` for selector return type since no state data is needed
- Factory should bind the dispatch function to the thunk
- Use simple parameter binding pattern since there's only one adapter (dispatch) and one parameter

## References

- Task 2: `EditDatesDialogOpened`, `editDatesDialogOpened` from `action-editDatesDialogOpened.ts`
- `web/src/core/catalog/album-delete/thunk-openDeleteAlbumDialog.ts` - for simple thunk pattern
- `web/src/core/catalog/thunks.ts` - for registration pattern
- `web/src/core/catalog/tests/test-helper-state.ts` - for test data

## Relevant Principles

> A thunk implement one and only one responsibility: single-responsibility principle

> Thunks are pure-business logic functions handling the user commands triggered by the view

> The thunk functions first argument is a dispatch function accepting the specific action interface type(s)

> TDD must be used: no behaviour should be implemented if there is no test to validate it

## TDD Approach

Write tests first that validate:
1. The thunk function dispatches the correct action with the provided albumId
2. The thunk declaration has the correct selector (returns empty object)
3. The thunk declaration factory correctly binds the dispatch function
4. The dispatched action can be verified by checking the dispatched array

Use test data from `test-helper-state.ts` like `twoAlbums[0].albumId` for testing with valid album IDs.
