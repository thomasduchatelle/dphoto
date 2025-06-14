# Task 1: Action - Edit Dates Dialog Opened

## Introduction
You are implementing the **action** `editDatesDialogOpened` and its associated **selector** `selectEditDatesDialog` for the edit dates dialog feature.

## Requirements (BDD Style)

```
GIVEN the catalog state contains albums in allAlbums
WHEN I dispatch editDatesDialogOpened with an AlbumId that exists in the state
THEN the selector editDatesDialogSelector returns an EditDatesDialogSelection with open: true, the album name, and the dates converted to Dayjs format

GIVEN the catalog state contains albums in allAlbums
WHEN I dispatch editDatesDialogOpened with an AlbumId that does not exist in the state
THEN the selector editDatesDialogSelector returns an EditDatesDialogSelection with open: false

GIVEN a catalog state with no dialog open
WHEN no action is dispatched
THEN the selector editDatesDialogSelector returns the closedEditDatesDialogProperties constant
```

## Implementation Details

**TDD Principle**: Implement the tests first, then implement the code **the simplest and most readable possible**: no behaviour should be implemented if it is not required by one test.

**Folder Structure**: Create the action in `web/src/core/catalog/edit-dates/` following the naming convention `action-editDatesDialogOpened.ts`.

**Files to Edit**:
- `web/src/core/catalog/language/catalog-state.ts` - Add the new interfaces and constant
- `web/src/core/catalog/actions.ts` - Register the action
- `web/src/core/catalog/catalog-reducer-v2.ts` - Register the reducer

**Recommendations**:
- Use `dayjs()` to convert Date objects to Dayjs format
- Handle the case where album is not found by returning a closed dialog state
- The selector should always return a valid `EditDatesDialogSelection`, never undefined

## Interface Specification

```typescript
// In catalog-state.ts
export interface EditDatesDialogState {
    albumId: AlbumId
    startDate: Dayjs
    endDate: Dayjs
}

export interface EditDatesDialogSelection {
    open: boolean
    albumName: string
    startDate: Dayjs
    endDate: Dayjs
}

export const closedEditDatesDialogProperties: EditDatesDialogSelection = {
    open: false,
    albumName: "",
    startDate: dayjs(),
    endDate: dayjs()
}

// Action interface
export interface EditDatesDialogOpened {
    type: "EditDatesDialogOpened"
    albumId: AlbumId
}

// Action factory signature
export function editDatesDialogOpened(albumId: AlbumId): EditDatesDialogOpened

// Reducer signature
export function reduceEditDatesDialogOpened(
    current: CatalogViewerState,
    action: EditDatesDialogOpened
): CatalogViewerState

// Selector signature
export function editDatesDialogSelector(state: CatalogViewerState): EditDatesDialogSelection
```

## References

- `web/src/core/catalog/language/catalog-state.ts` - For existing state interfaces and patterns
- `web/src/core/catalog/language/utils-albumIdEquals.ts` - For album comparison logic
- `docs/principles_web.md` - For action implementation patterns and testing guidelines
