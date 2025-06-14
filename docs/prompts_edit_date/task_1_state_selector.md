# Task 1: Create Edit Dates Dialog State and Selector

You are implementing the foundational state interfaces and selector for the Edit Dates Dialog feature. This task establishes the domain model that will be used by actions, thunks, and React components.

## Requirements (BDD Style)

```
GIVEN the dialog is closed
WHEN selectEditDatesDialog is called
THEN it returns isOpen: false and albumName: ""

GIVEN the dialog is open with an album that exists in allAlbums
WHEN selectEditDatesDialog is called
THEN it returns isOpen: true and the correct album name

GIVEN the dialog is open with an album that doesn't exist in allAlbums
WHEN selectEditDatesDialog is called
THEN it returns isOpen: false and albumName: ""
```

## Implementation Details

**Create file:** `web/src/core/catalog/language/edit-dates-dialog-state.ts`

**Export the following:**
- Interface: `EditDatesDialogState`
- Interface: `EditDatesDialogProperties` 
- Constant: `closedEditDatesDialogProperties: EditDatesDialogProperties`
- Function: `selectEditDatesDialog(state: CatalogViewerState): EditDatesDialogProperties`

**Update file:** `web/src/core/catalog/language/catalog-state.ts`
- Add `editDatesDialog?: EditDatesDialogState` to the `CatalogViewerState` interface

**Export from:** `web/src/core/catalog/language/index.ts`
- Add exports for the new interfaces and selector

## Interface Specifications

```typescript
interface EditDatesDialogState {
    albumId: AlbumId
}

interface EditDatesDialogProperties {
    isOpen: boolean
    albumName: string
}
```

## References

- `web/src/core/catalog/language/catalog-state.ts` - for `CatalogViewerState`, `AlbumId`, and `Album` interfaces
- `web/src/core/catalog/language/utils-albumIdEquals.ts` - for album comparison utilities
- `web/src/core/catalog/tests/test-helper-state.ts` - for test data like `twoAlbums` and `myselfUser`

## Relevant Principles

> The state model is common across the slice of the application, it should not be based on the REST model nor a specific UI component, it should be designed to represent the best the business.

> Always use explicit types or interfaces, never use `any`

> Selector function can be provided as part of the model to extract a set of properties used by a component like a Dialog.

> Properties that change together must be within the same state

## TDD Approach

Write tests first that validate:
1. The selector returns closed state when `editDatesDialog` is undefined
2. The selector finds the correct album name when dialog is open with valid albumId
3. The selector returns closed state when dialog is open but albumId doesn't exist in allAlbums
4. The constant `closedEditDatesDialogProperties` has the correct default values

Use test data from `test-helper-state.ts` like `twoAlbums` for testing album lookup scenarios.
