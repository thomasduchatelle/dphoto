# Task 2: Thunk - Open Edit Date Dialog

## Introduction
You are implementing the **thunk** `openEditDateDialog` for opening the edit dates dialog.

## Requirements (BDD Style)

```
GIVEN a catalog state with albums loaded
WHEN I call openEditDateDialog thunk with an AlbumId
THEN it dispatches editDatesDialogOpened action with that AlbumId
```

## Implementation Details

**TDD Principle**: Implement the tests first, then implement the code **the simplest and most readable possible**: no behaviour should be implemented if it is not required by one test.

**Folder Structure**: Create the thunk in `web/src/core/catalog/album-edit-dates/` following the naming convention `thunk-openEditDateDialog.ts`.

**Files to Edit**:
- `web/src/core/catalog/thunks.ts` - Register the thunk declaration

**Recommendations**:
- This thunk should be simple and not perform validation (validation is handled by the action)
- No port/adapter is needed for this thunk as it only dispatches an action
- Use the simple case pattern for the thunk factory (Case 1 from principles)

## Interface Specification

```typescript
// Thunk function signature
export function openEditDateDialogThunk(
    dispatch: (action: EditDatesDialogOpened) => void,
    albumId: AlbumId
): void

// Thunk declaration signature
export const openEditDateDialogDeclaration: ThunkDeclaration<
    CatalogViewerState,
    {},
    (albumId: AlbumId) => void,
    CatalogFactoryArgs
>
```

## References

- `web/src/core/catalog/album-edit-dates/action-editDatesDialogOpened.ts` - For the action to dispatch (from Task 1)
- `docs/principles_web.md` - For thunk implementation patterns and testing guidelines
- Existing thunk examples in the codebase for factory pattern implementation
