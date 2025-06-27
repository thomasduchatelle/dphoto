# Migration Guide: Consolidating Dialogs into a Single `CatalogDialog` Type

This document outlines the process for migrating existing dialog states (`DeleteDialogState` and `ShareModal`) into a unified `CatalogDialog` type within the
`CatalogViewerState`. The goal is to ensure that only one dialog can be open at a time, improving state management and UI consistency.

The `EditDatesDialog` has already been migrated as an example. This guide provides instructions for an LLM agent to perform the same migration for the remaining
dialogs.

## Current State (`CatalogViewerState` before this migration)

```typescript
export interface CatalogViewerState {
    // ... other properties
    deleteDialog?: DeleteDialogState;
    shareModal?: ShareModal;
    dialog?: CatalogDialog; // Currently only includes EditDatesDialog
}

export interface DeleteDialogState {
    deletableAlbums: Album[]
    initialSelectedAlbumId?: AlbumId
    isLoading: boolean
    error?: string
}

export interface ShareModal {
    sharedAlbumId: AlbumId
    sharedWith: Sharing[]
    suggestions: UserDetails[]
    error?: ShareError
}

export interface EditDatesDialog {
    type: "EditDatesDialog"
    albumId: AlbumId
    albumName: string
    startDate: Date
    endDate: Date
    startAtDayStart: boolean
    endAtDayEnd: boolean
    isLoading: boolean
    error?: string
}

export type CatalogDialog = EditDatesDialog; // This will be extended
```

## Migration Steps for an LLM Agent

The LLM agent should perform the following steps for each remaining dialog (`DeleteDialogState` and `ShareModal`):

### Step 1: Define the New Dialog Interface

For each dialog, create a new interface that includes a `type` discriminator property. This `type` property will be a string literal unique to that dialog.

**Example for `DeleteDialogState`:**

```typescript
export interface DeleteDialog {
    type: "DeleteDialog"; // New discriminator property
    deletableAlbums: Album[];
    initialSelectedAlbumId?: AlbumId;
    isLoading: boolean;
    error?: string;
}
```

**Example for `ShareDialog` (from `ShareModal`):**

```typescript
export interface ShareDialog {
    type: "ShareDialog"; // New discriminator property
    sharedAlbumId: AlbumId;
    sharedWith: Sharing[];
    suggestions: UserDetails[];
    error?: ShareError;
}
```

These new interfaces should be placed in `web/src/core/catalog/language/catalog-state.ts` alongside `EditDatesDialog`.

### Step 2: Update `CatalogDialog` Type

Modify the `CatalogDialog` type in `web/src/core/catalog/language/catalog-state.ts` to be a union of all dialog types.

**Example:**

```typescript
export type CatalogDialog = EditDatesDialog | DeleteDialog | ShareDialog;
```

### Step 3: Remove Old Dialog Properties from `CatalogViewerState`

Remove the old, individual dialog properties (`deleteDialog` and `shareModal`) from `CatalogViewerState` in `web/src/core/catalog/language/catalog-state.ts`.

**Example:**

```typescript
export interface CatalogViewerState {
    // ... other properties
    // deleteDialog?: DeleteDialogState; // Remove this
    // shareModal?: ShareModal; // Remove this
    dialog?: CatalogDialog; // Keep this
}
```

### Step 4: Update `initialCatalogState`

In `web/src/core/catalog/language/initial-catalog-state.ts`, ensure that `dialog` is initialized to `undefined` (or not present if `undefined` is the default
for optional properties). The old dialog properties should no longer be initialized.

**Example:**

```typescript
export const initialCatalogState = (currentUser: CurrentUserInsight): CatalogViewerState => ({
    // ...
    // deleteDialog: undefined, // Remove this
    // shareModal: undefined, // Remove this
    dialog: undefined, // Ensure this is undefined or omitted
})
```

### Step 5: Update Actions

For each action that previously interacted with `deleteDialog` or `shareModal`, update it to interact with the new `dialog` property and use the `type`
discriminator.

**General pattern for actions that *open* a dialog:**

When an action opens a dialog, it should set `current.dialog` to the new dialog object, including its `type` property.

```typescript
// Before (example for delete dialog)
// return { ...current, deleteDialog: { ... } };

// After
return {
    ...current,
    dialog: {
        type: "DeleteDialog", // Set the discriminator
        // ... other properties of DeleteDialog
    },
};
```

**General pattern for actions that *modify* an open dialog:**

When an action modifies an open dialog, it must first check if `current.dialog` exists and if its `type` matches the expected dialog. If it matches, then update
the properties.

```typescript
// Before (example for delete dialog)
// if (!current.deleteDialog) return current;
// return { ...current, deleteDialog: { ...current.deleteDialog, isLoading: true } };

// After
if (!isDeleteDialog(current.dialog)) { // Use the new type predicate function
    return current;
}
return {
    ...current,
    dialog: {
        ...current.dialog, // Spread the existing dialog properties
        isLoading: true,
    },
};
```

**General pattern for actions that *close* a dialog:**

When an action closes a dialog, it should set `current.dialog` to `undefined`. It should also check the `type` discriminator before closing to ensure it's
closing the correct dialog.

```typescript
// Before (example for delete dialog)
// if (!current.deleteDialog) return current;
// return { ...current, deleteDialog: undefined };

// After
if (!isDeleteDialog(current.dialog)) { // Use the new type predicate function
    return current;
}
return {
    ...current,
    dialog: undefined,
};
```

### Step 6: Update Selectors

For each selector that previously accessed `deleteDialog` or `shareModal`, update it to access the new `dialog` property and use the `type` discriminator.

**General pattern for selectors:**

```typescript
// Before (example for delete dialog selector)
// if (!state.deleteDialog) { return DEFAULT_DELETE_DIALOG_SELECTION; }
// return { isOpen: true, ...state.deleteDialog };

// After
if (!isDeleteDialog(state.dialog)) { // Use the new type predicate function
    return DEFAULT_DELETE_DIALOG_SELECTION;
}
return {
    isOpen: true,
    // ... map properties from state.dialog
};
```

### Step 7: Update Thunks

For each thunk that previously accessed `deleteDialog` or `shareModal` via its `selector` property, update the `selector` to access the new `dialog` property
and use the `type` discriminator.

**General pattern for thunk selectors:**

```typescript
// Before (example for delete thunk selector)
// selector: (state: CatalogViewerState) => (state.deleteDialog ? { ... } : undefined),

// After
selector: (state: CatalogViewerState) => {
    const dialog = state.dialog;
    if (!isDeleteDialog(dialog)) { // Use the new type predicate function
        return undefined;
    }
    return {
        // ... map properties from dialog
    };
},
```

### Step 8: Update Tests

All unit tests (`.test.ts` files) that interact with the dialog state must be updated to reflect the new `dialog` structure and the `type` discriminator.

* **Action tests:** When creating a `CatalogViewerState` for testing, ensure the `dialog` property is set correctly with the `type` discriminator.
* **Selector tests:** Ensure assertions correctly check the `isOpen` property and other mapped properties from the selector's output.
* **Thunk tests:** Ensure the initial state passed to the thunk's reducer (if applicable) and any expected dispatched actions correctly reflect the new dialog
  structure.

### Step 9: Create Type Predicate Function for Dialog Access

For each dialog type, create a type predicate function (e.g., `isEditDatesDialog`, `isDeleteDialog`, `isShareDialog`) that takes `CatalogDialog | undefined` as
input and returns `true` if the dialog is of the specific type, narrowing the type.

**Example for `EditDatesDialog` (updated):**

```typescript
// web/src/core/catalog/language/catalog-state.ts
export function isEditDatesDialog(dialog: CatalogDialog | undefined): dialog is EditDatesDialog {
    return dialog?.type === "EditDatesDialog";
}
```

This type predicate function should be used in all actions, selectors, and thunks that need to access properties of a specific dialog type.

### Important Considerations for the LLM Agent:

* **Thoroughness:** The agent must identify *all* occurrences where `deleteDialog` or `shareModal` are used (in actions, selectors, thunks, and tests) and
  update them.
* **Type Safety:** Ensure all changes maintain TypeScript type safety. The `type` discriminator is crucial for this.
* **No New Functionality:** The task is a refactoring. No new business logic should be introduced.
* **Test Coverage:** All existing tests must pass after the migration. If any tests fail, the agent must debug and fix them.
* **File Paths:** Pay close attention to file paths when identifying files to modify.
* **Discriminator Property:** The `type` property is essential for TypeScript to correctly narrow down the union type `CatalogDialog`. It must be present in
  every new dialog interface and checked when accessing dialog-specific properties.
* **Closing Dialogs:** When a dialog is closed, `state.dialog` should be set to `undefined`. This ensures only one dialog can be open at a time.
* **`mediasLoadedFromAlbumId` in `albumDatesUpdated`:** Pay attention to how `mediasLoadedFromAlbumId` is set in `albumDatesUpdated`. It should now correctly
  access `current.dialog.albumId` only if `current.dialog` is of type `EditDatesDialog`.
* **New Type Predicate Function Usage:** Ensure the newly created type predicate functions (e.g., `isEditDatesDialog`) are consistently used across all relevant
  files to reduce code duplication and improve readability.

By following these detailed steps, the LLM agent should be able to successfully migrate the remaining dialogs to the new single-dialog model.
