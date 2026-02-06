# Task 5: Implement Close Edit Dates Dialog Thunk

You are implementing the thunk that handles the business logic for closing the edit dates dialog.

## Requirements (BDD Style)

```
GIVEN any state
WHEN the closeEditDatesDialog thunk is called
THEN it dispatches editDatesDialogClosed
```

## Implementation Details

**Create file:** `web/src/core/catalog/edit-dates/thunk-closeEditDatesDialog.ts`

**Export the following:**
- Function: `closeEditDatesDialogThunk(dispatch: (action: EditDatesDialogClosed) => void): void`
- Constant: `closeEditDatesDialogDeclaration: ThunkDeclaration<CatalogViewerState, {}, () => void, CatalogFactoryArgs>`

**Register thunk in:** `web/src/core/catalog/thunks.ts`
- Add `closeEditDatesDialog: closeEditDatesDialogDeclaration` to the `catalogThunks` object

## Thunk Declaration Pattern

Follow the pattern from existing thunks:
- Use empty object `{}` for selector return type since no state data is needed
- Factory should bind the dispatch function to the thunk
- Use simple parameter binding pattern since there are no parameters beyond dispatch

## References

- Task 3: `EditDatesDialogClosed`, `editDatesDialogClosed` from `action-editDatesDialogClosed.ts`
- `web/src/core/catalog/album-delete/thunk-closeDeleteAlbumDialog.ts` - for simple thunk pattern
- `web/src/core/catalog/thunks.ts` - for registration pattern
- `web/src/core/catalog/tests/test-helper-state.ts` - for test state setup

## Relevant Principles

> A thunk implement one and only one responsibility: single-responsibility principle

> Thunks are pure-business logic functions handling the user commands triggered by the view

> The thunk functions first argument is a dispatch function accepting the specific action interface type(s)

> TDD must be used: no behaviour should be implemented if there is no test to validate it

## TDD Approach

Write tests first that validate:
1. The thunk function dispatches the correct editDatesDialogClosed action
2. The thunk declaration has the correct selector (returns empty object)
3. The thunk declaration factory correctly binds the dispatch function
4. The dispatched action can be verified by checking the dispatched array

The test should be simple since this thunk has no parameters and always dispatches the same action.
