# Album Edit - dates and name

## Requirements

```
Idea: A user should be able to edit an existing Album. He can either edit the name (and the folder name), or the dates (start and end date).
```

### Feature Summary

Allow album owners to edit an existing album’s name (and folderName) or its start/end dates via dedicated dialogs accessible from the left menu. The UI will
enforce ownership
restrictions and provide clear feedback during the editing process.

### Ubiquity Language

• Edit Album Name/Folder: The action of changing an album’s display name and its internal folderName.
• Edit Album Dates: The action of changing an album’s start and end dates.

### Scenarios

#### Scenario 1 - Edit Album Name and/or FolderName

1. The user navigates to the left menu and selects an album they own.
2. The user clicks the "Edit Name" button above the album list.
3. The "Edit Album Name" dialog opens, displaying the current album name and folderName.
4. The user edits the album name. By default, folderName is auto-generated and disabled.
5. The user may untick "Auto-generate folder name" to manually edit the folderName.
6. The user submits the changes.
7. The dialog shows a loading indicator while the API call is in progress.
8. If the update is successful:
    * If only the name changed, the dialog closes and the album name updates in the UI.
    * If the folderName changed, the dialog closes and the user is redirected to the new album URL.
9. If there is an error, the dialog remains open and displays the error message.

#### Scenario 2 - Edit Album Dates

1. The user navigates to the left menu and selects an album they own.
2. The user clicks the "Edit Dates" button above the album list.
3. The "Edit Album Dates" dialog opens, displaying the current start and end dates.
4. The user edits one or both dates (with basic frontend validation, e.g., start date must be before end date).
5. The user submits the changes.
6. The dialog shows a loading indicator while the API call is in progress.
7. If the update is successful, the dialog closes and both the album list and media list are refreshed in the UI.
8. If there is an error, the dialog remains open and displays the error message.

#### Scenario 3 - Ownership Restriction

1. The user navigates to the left menu and selects an album they do not own.
2. The "Edit Name" and "Edit Dates" buttons above the album list are disabled.
3. When the user hovers over a disabled button, a tooltip appears: "Only the album owner can edit this album."

#### Scenario 4 - Error Handling During Edit

1. The user opens either the "Edit Album Name" or "Edit Album Dates" dialog and makes changes.
2. The user submits the changes.
3. The dialog shows a loading indicator while waiting for the API response.
4. The API returns an error (e.g., invalid folderName, date out of range).
5. The dialog remains open and displays the error message.
6. The user can correct the input and resubmit, or cancel the operation.

#### Scenario 5 - Cancel Edit Operation

1. The user opens either the "Edit Album Name" or "Edit Album Dates" dialog.
2. The user decides not to make any changes and clicks the "Cancel" button.
3. The dialog closes with no changes made to the album.

### Technical Context

* The left menu already contains icon buttons for album actions (delete, create).
* API provides two endpoints: one for updating name/folderName, one for updating dates.
* Only album owners can edit albums.
* Album state and UI updates are managed in React state.
* Existing validation logic from the create dialog can be reused as utility functions.
* AlbumId and ownership logic are defined in the catalog-state domain.
* Redirect logic and album/media refresh are already supported in the platform.

### Explorations

* What are the exact backend error messages and codes for failed updates, and how should they be mapped to user-friendly UI errors?
* Are there any edge cases where folderName changes could break links or references elsewhere in the app?
* Should there be a debounce or enforced delay between sequential API calls if both name/folderName and dates are edited in one session?
* Are there any additional audit or logging requirements for album edits?
* Should the dialogs support keyboard navigation and accessibility features (e.g., ARIA labels, focus management)?