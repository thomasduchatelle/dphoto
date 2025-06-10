# Delete Album Feature: Implementation Prompts (Port/Adapter Architecture)

---

## 1. Backend: Expose Delete Album as a REST Endpoint

```text
Implement a REST API endpoint to delete an album.

- **Path:** DELETE /api/v1/owners/{owner}/albums/{folderName}
- **Method:** DELETE
- **Request Parameters:**
  - `owner` (string, path): The owner identifier.
  - `folderName` (string, path): The album's folder name.

**Behavior:**
- On success, return HTTP 204 No Content.
- On failure, return HTTP 400 or 409 with a JSON body:
  - `errorType` (string): A code representing the error (e.g., "MEDIA_REALLOCATION_FAILED").
  - `message` (string): A human-readable error message.

**Error Cases:**
- Album not found: 404 Not Found, errorType: "ALBUM_NOT_FOUND"
- Not authorized: 403 Forbidden, errorType: "NOT_AUTHORIZED"
- Media cannot be re-allocated: 409 Conflict, errorType: "MEDIA_REALLOCATION_FAILED"
- Any other error: 500 Internal Server Error, errorType: "INTERNAL_ERROR"

**Tests:**
- Deleting an existing album you own returns 204.
- Deleting a non-existent album returns 404 with correct errorType.
- Deleting an album you do not own returns 403 with correct errorType.
- Deleting an album with un-reallocatable media returns 409 with correct errorType.
- Unexpected errors return 500 with correct errorType.
```

---

## 2. Delete Dialog: Pure Component (Storybook-Driven)

```text
Create a DeleteAlbumDialog component that has two rendering.

First is the initial one, when the Dialog is open and `state.confirmMode === false`:

- a title: "Delete an album"
- a text: "What album do you want to delete ? The medias will be re-assigned to the appropriate album"
- a dropdown with the list of albums: "<name of the album> <from date to date> (<number of medias> media(s))"
- two buttons: "Cancel" to call `onClose` and "Delete" to update the `confirmState === true`

Second is when the `confirmState === true`:

- a title: "Delete an album"
- a text: "Are you sure you want to delete <name of the album> <from date to date> (<number of medias> media(s))"
- two buttons: "Cancel" to restore the state `confirmState = false` and "Delete" to callback `onDelete`

**Props:**
- `albums`: Array of objects `{ id: string, name: string, mediaCount: number }` (only deletable albums)
- `initialSelectedAlbumId`: string (the album to pre-select when dialog opens)
- `isOpen`: boolean (whether the dialog is visible)
- `isLoading`: boolean (show loading indicator)
- `error`: string | null (error message to display, if any)

**Callbacks:**
- `onDelete(albumId: string)`: called when the user clicks the "Confirm" button, passing the selected album id
- `onClose()`: called when the dialog is closed/cancelled

**Behavior:**
- Manages its own internal state for selected album and deletion/confirmation steps.
- Renders a dropdown with album name and media count.
- "Delete" button is enabled when an album is selected.
- "Confirm" button is shown after "Delete" is clicked.
- Shows loading indicator when `isLoading` is true.
- Shows error message if `error` is not null.

**Storybook Stories:**
- Default (open, with albums, no selection)
- With pre-selected album
- Loading state
- Error state
- No albums available
```

---

## 3. State Management: Reducer Actions

### 3.1. Action: Open Delete Dialog

```text
Update the CatalogViewerState to add the properties:

- deleteDialog: object as following, initially undefined
  - deletableAlbums: Album[]
  - initialSelectedAlbumId?: AlbumId
  - isLoading: bool
  - error ?: string
  
And create a selector (a function taking the state in input) that returns:

- albums: Album[]; // maps from deleteDialog.deletableAlbums, or empty
- initialSelectedAlbumId?: AlbumId; // maps from loadingMediasFor if defined, or mediasLoadedFromAlbumId
- isOpen: boolean; // true if deleteDialog is defined
- isLoading: boolean; // maps from deleteDialog.isLoading
- error: string | null; // maps from deleteDialog.error

Define an action and its reducer to open the delete album dialog.

**Action Type:** `OpenDeleteAlbumDialog`

**Reducer Behavior:**
- Sets dialog open state to true.
- Sets the pre-selected album to the one currently loaded (or loading)
- Resets error state

**Test Cases, values from the selector are to be asserted, not directly the state:**
- it results to dialog open, pre-selected album to the one currently loaded, and albums only containing those deletable, when the dialog was not open
- it results to dialog open, pre-selected album to the one being loaded, and albums only containing those deletable, when the dialog was not open and medias are currently being loaded
- it results to dialog open with deletable albums, no error an loading false, when the dialog were already open
```

---

### 3.2. Action: Deletion Success

```text
Define the action and its reducer for successful album deletion.

**Action Type:** `AlbumDeleted`
**Payload:** `{ albums: Album[], redirectTo?: AlbumId }`

**Reducer test cases:**
- it closes the dialog and refresh the list of albums based on albumFilter
  (note: closing the dialog means deleteDialog becomes undefined)
```

---

### 3.3. Action: Deletion Failure

```text
Define the action and its reducer for failed album deletion.

**Action Type:** `AlbumFailedToDelete`
**Payload:** `{ error: string }`

**Test cases:**
- set the error value in deleteDialog.error when the dialog is open
  (dialog open means deleteDialog is defined)
- ignores if the dialog is closed
- fails to create the action if the error message is empty.
```

---

### 3.4. Action: Close Delete Dialog

```text
Define the action and the reducer action to close/cancel the delete album dialog.

**Action Type:** `CloseDeleteAlbumDialog`
**Payload:** none

**Reducer Behavior:**
- Sets dialog open state to false.
- Clears any error state.

**Test Cases:**
- closes the dialog when it was open, no matter its state
  (closing the dialog means to set deleteDialog to undefined ; use a state with an error)
- ignores when the dialog is already closed

Create the files for the code and the related tests. Register the new action and reducer in reducer-v2.ts. In the test, always write assertions 
on the whole state (not on each property), use the already defined test cases, and keep the code concise by not creating variable when inline calls
are possible and readable.
```

---

## 4. Thunks: Backend Interaction Interfaces

### 4.1. Thunk: Open Delete Dialog

```text
Create a thunk and its tests to open the delete album dialog.

**Context**
- _thunk arguments_: none
- _selected state_: none
- _port_: none

**Test Case**
- triggers the action `openDeleteAlbum`, always
```

---

### 4.2. Thunk: Close Delete Dialog

```text
Define a thunk and its tests to close/cancel the delete album dialog.

**Context**
- _thunk arguments_: none
- _selected state_: none
- _port_: none

**Test Case**
- triggers the action `closeDeleteAlbumDialog`, always
```

---

### 4.3. Thunk: Delete Album

```text
Define a thunk and its test to delete an album.

**Context**
- _thunk arguments_: the album id to delete
- _selected state_: currently selected album (`loadingMediasFor ?? mediasLoadedFromAlbumId`)
- _port_: an interface with two methods: one to delete the album from its albumId, one to find all albums (no arguments)

**Test Case**
- dispatches 'startDeleteAlbum' action, then the 'albumDeleted' action with freshly reloaded albums lists and and undefined redirectTo when the deletion succeed for an album different from the one selected
- dispatches 'startDeleteAlbum' action, then the 'albumDeleted' action with the refreshed list of album and redirect to the first album when the deletion succeed for the selected album 
- dispatches 'startDeleteAlbum' action, then 'albumFailedToDelete' action when the deletion failed
- dispatches 'startDeleteAlbum' action, then 'albumFailedToDelete' action when the reloading of albums failed

**Code context**:

    export interface AlbumDeletedAction extends RedirectToAlbumIdAction {
        type: "AlbumDeleted";
        albums: Album[]; // the list returned by the port
        redirectTo?: AlbumId;
    }
    export function albumDeletedAction(props: Omit<AlbumDeletedAction, "type">): AlbumDeletedAction {
        return {
            ...props,
            type: "AlbumDeleted",
        };
    }
    
    export interface AlbumFailedToDeleteAction {
        type: "AlbumFailedToDelete";
        error: string; // the port return an error code which must be converted into a human readable message. 
                       // "OrphanedMediasError" -> "The album cannot be deleted, create an album to which all medias can be re-allocated"
    }
    
    export function albumFailedToDeleteAction(props: Omit<AlbumFailedToDeleteAction, "type">): AlbumFailedToDeleteAction {
        if (!props.error || props.error.trim() === "") {
            throw new Error("AlbumFailedToDeleteAction requires a non-empty error message");
        }
        return {
            ...props,
            type: "AlbumFailedToDelete",
        };
    }
```

---

## 5. Adapters: Backend API Implementation

```text
Implement an adapter that fulfills the backend interaction interface for album deletion.

**Interface:**
- `deleteAlbum(albumId: string): Promise<void>`
  - Makes a DELETE request to `/api/v1/owners/{owner}/albums/{folderName}`.
  - Resolves on HTTP 204.
  - Rejects with `{ errorType: string, message: string }` on error.

**Test Cases:**
- HTTP 204: resolves.
- HTTP 404: rejects with errorType "ALBUM_NOT_FOUND".
- HTTP 403: rejects with errorType "NOT_AUTHORIZED".
- HTTP 409: rejects with errorType "MEDIA_REALLOCATION_FAILED".
- HTTP 500: rejects with errorType "INTERNAL_ERROR".
- Network error: rejects with errorType "NETWORK_ERROR".
```

---
