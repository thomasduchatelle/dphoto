# Task 3: Implement Edit Dates Dialog Closed Action

You are implementing the action that closes the edit dates dialog and clears its state.

## Requirements (BDD Style)

```
GIVEN the dialog is open
WHEN I dispatch editDatesDialogClosed
THEN selectEditDatesDialog returns a closed dialog with empty album name
```

## Implementation Details

**Create file:** `web/src/core/catalog/edit-dates/action-editDatesDialogClosed.ts`

**Export the following:**
- Interface: `EditDatesDialogClosed`
- Function: `editDatesDialogClosed(): EditDatesDialogClosed`
- Function: `reduceEditDatesDialogClosed(state: CatalogViewerState, action: EditDatesDialogClosed): CatalogViewerState`
- Function: `editDatesDialogClosedReducerRegistration(handlers: any)`

**Register action in:** `web/src/core/catalog/catalog-reducer-v2.ts`
- Add to `catalogActions` object
- Add to `CatalogViewerAction` union type
- Add to `reducerRegistrations` array

## Action Interface Specification

```typescript
interface EditDatesDialogClosed {
    type: "editDatesDialogClosed";
}
```

## References

- Task 1: `selectEditDatesDialog`, `closedEditDatesDialogProperties` from `edit-dates-dialog-state.ts`
- `web/src/core/catalog/language/catalog-state.ts` - for `CatalogViewerState` interface
- `web/src/core/catalog/album-delete/action-albumDeleted.ts` - for action pattern examples
- `web/src/core/catalog/tests/test-helper-state.ts` - for test state setup

## Relevant Principles

> The Action naming convention is the Past Tense Event Naming and is represented in camelCase

> if the interface has no property other than type: no parameter

> The reducer function returns the updated state, same type as the first parameter.

> TDD must be used: no behaviour should be implemented if there is no test to validate it

## TDD Approach

Write tests first that validate:
1. The action factory creates the correct action object with only the type property
2. The reducer clears the editDatesDialog state (sets it to undefined)
3. The selector returns the closed dialog properties after the reducer runs
4. The reducer works correctly whether the dialog was previously open or already closed

Use test data from `test-helper-state.ts` to create initial states with open dialogs for testing the close behavior.
