# Feature: Rename Album

## 1. Feature Summary
This feature allows users to update the display name of an album and, optionally, its underlying folder name. This provides greater flexibility in organizing and identifying albums within the application.

## 2. Ubiquity Language
*   **Album Display Name:** The user-facing name of the album, displayed throughout the application.
*   **Album Folder Name:** The underlying system-level name for the album's folder, which can optionally be changed by the user.
*   **Edit Name Dialog:** A new modal dialog that appears when the user chooses to rename an album, containing fields for the new album display name and folder name.

## 3. Scenarios

### Scenario 1: Successfully Renaming an Album (Display Name Only)
1.  The user navigates to an album they own.
2.  The user clicks the "Edit" button.
3.  A dropdown appears, and the user selects "Edit Name".
4.  The "Edit Name" dialog opens, pre-filled with the current album name and the "Folder Name" checkbox unchecked.
5.  The user enters a new "Album Name" (e.g., "My Awesome Vacation Photos").
6.  The user clicks "Save".
7.  The dialog enters a disabled mode with a loading bar.
8.  Upon successful API response, the dialog closes.
9.  The album list is updated, and the album now displays the new name. The user remains on the current page as the AlbumId did not change.

### Scenario 2: Successfully Renaming an Album (Display Name and Folder Name)
1.  The user navigates to an album they own.
2.  The user clicks the "Edit" button.
3.  A dropdown appears, and the user selects "Edit Name".
4.  The "Edit Name" dialog opens, pre-filled with the current album name and the "Folder Name" checkbox unchecked.
5.  The user enters a new "Album Name" (e.g., "Summer Adventures 2024").
6.  The user checks the "Folder Name" checkbox. The folder name text field becomes editable and defaults to the current folder name.
7.  The user enters a new "Folder Name" (e.g., "summer-adventures-2024").
8.  The user clicks "Save".
9.  The dialog enters a disabled mode with a loading bar.
10. Upon successful API response, the dialog closes.
11. The album list is updated. If the AlbumId changed (due to the folder name update), the user is redirected to the newly renamed album.

### Scenario 3: Attempting to Rename an Album with a Blank Album Name
1.  The user navigates to an album they own.
2.  The user clicks the "Edit" button.
3.  A dropdown appears, and the user selects "Edit Name".
4.  The "Edit Name" dialog opens.
5.  The user clears the "Album Name" field.
6.  A client-side validation error message appears below the "Album Name" field, indicating that the field cannot be blank. The "Save" button becomes disabled.
7.  The user enters a valid "Album Name".
8.  The validation error disappears, and the "Save" button becomes enabled.
9.  The user proceeds to save or cancel.

### Scenario 4: API Returns AlbumNameMandatoryErr (Blank Folder Name)
1.  The user attempts to rename an album, having checked the "Folder Name" checkbox and left the folder name field blank.
2.  The user clicks "Save".
3.  The dialog enters a disabled mode with a loading bar.
4.  The API returns an `AlbumNameMandatoryErr`.
5.  An error message appears below the "Folder Name" field, indicating that the folder name cannot be blank. The "Save" button becomes disabled.
6.  The dialog exits the disabled mode.
7.  The user enters a valid "Folder Name".
8.  The error message disappears, and the "Save" button becomes enabled.
9.  The user proceeds to save or cancel.

### Scenario 5: API Returns AlbumFolderNameAlreadyTakenErr
1.  The user attempts to rename an album, having checked the "Folder Name" checkbox and entered a folder name that already exists.
2.  The user clicks "Save".
3.  The dialog enters a disabled mode with a loading bar.
4.  The API returns an `AlbumFolderNameAlreadyTakenErr`.
5.  An error message appears below the "Folder Name" field, indicating that the folder name is already taken. The "Save" button becomes disabled.
6.  The dialog exits the disabled mode.
7.  The user enters a unique "Folder Name".
8.  The error message disappears, and the "Save" button becomes enabled.
9.  The user proceeds to save or cancel.

### Scenario 6: Attempting to Rename an Album (Permission Denied)
1.  The user navigates to an album they *do not* own.
2.  The "Edit" button is disabled (existing functionality).
3.  The user is unable to initiate the rename process.

### Scenario 7: Unknown API Error During Rename
1.  The user initiates an album rename (either display name only or with folder name).
2.  The user clicks "Save".
3.  The dialog enters a disabled mode with a loading bar.
4.  The API returns an unexpected error (not `AlbumNameMandatoryErr` or `AlbumFolderNameAlreadyTakenErr`).
5.  An Alert block appears within the dialog, displaying a generic "Something went wrong" message or the specific error message if available.
6.  The dialog exits the disabled mode, allowing the user to try again or cancel.

## 4. Technical Context

*   **Provided by the platform:**
    *   Existing "Edit" button and its permission handling (disabled for non-owners).
    *   Modal dialog component for "Edit Name" dialog.
    *   Loading bar/spinner component.
    *   Alert block component for displaying generic errors.
    *   Redirection mechanism for `AlbumId` changes.
*   **Provided by supporting APIs:**
    *   A new API endpoint for renaming albums, which accepts the album ID, new display name, and optional new folder name.
    *   This API will return the new `AlbumId` upon success.
    *   This API will return specific error codes (`AlbumNameMandatoryErr`, `AlbumFolderNameAlreadyTakenErr`) for known validation failures.
*   **Out of Scope:**
    *   Changes to the "Edit Dates" dialog or its functionality.
    *   Renaming albums that are not owned by the current user (permissions are handled externally).
    *   Detailed implementation of the API endpoint itself.
    *   Changes to how `AlbumId` is generated or structured.

## 5. Explorations
*   None at this time.

---
**Document Version:** 1.0
**Date:** 27/06/2025
**Git Commit:** a7da888
