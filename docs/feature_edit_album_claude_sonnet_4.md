# Album Edit - dates and name

## Requirements

### Feature Summary

Album owners can edit their albums through a dropdown menu offering two options: editing the album name/folderName or editing the start/end dates. Each option
opens a specific dialog with appropriate form fields and handles the corresponding API endpoint.

### Scenarios

Scenario 1: User edits album name only

1. User clicks "Edit Album" button on an owned album
2. User selects "Edit Name" from dropdown
3. Dialog opens showing current name and folderName (disabled)
4. User changes album name, keeps "Force folder name" unchecked
5. User clicks Save, dialog shows loading spinner
6. API updates name, dialog closes, album name updates in UI state

Scenario 2: User edits both name and folderName

1. User clicks "Edit Album" button
2. User selects "Edit Name" from dropdown
3. User changes album name and checks "Force folder name" box
4. FolderName field becomes enabled, user modifies it
5. User clicks Save, dialog shows loading spinner
6. API updates both, dialog closes, user redirects to new album URL

Scenario 3: User edits album dates

1. User clicks "Edit Album" button
2. User selects "Edit Dates" from dropdown
3. Dialog opens with existing date picker logic (reused from create dialog)
4. User modifies start/end dates using date pickers
5. User clicks Save, dialog shows loading spinner
6. API updates dates, dialog closes, both albums list and medias refresh

Scenario 4: API error during edit

1. User attempts any edit operation
2. Dialog shows loading spinner
3. API returns error
4. Dialog remains open showing error message
5. User can retry or cancel

Scenario 5: User attempts to edit non-owned album

1. User views album they don't own
2. Edit Album button is not visible/available

### Technical Context

* Two separate API endpoints: one for name/folderName changes, one for date changes
* Existing date picker components from create dialog can be reused
* albumIsOwnedByCurrentUser() determines if edit option is available
* CatalogViewerState needs updates for name changes and full refresh for date changes
* URL redirection required when folderName changes
* Loading states prevent concurrent operations

### Explorations

* Should we create new dialogs or modify existing create dialog to handle edit mode?
* How to handle the transition/animation when redirecting after folderName change?
* What specific error messages should be shown for different API failure scenarios?
* Should there be any confirmation step before applying changes that affect many medias?