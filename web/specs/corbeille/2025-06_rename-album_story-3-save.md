# As a user, I want to persist the new album name(s) and see relevant feedback for success or API errors.

## Acceptance Criteria

```
GIVEN I have the "Edit Album Name" dialog open for an album (e.g., "My Holiday Pics")
AND I have entered a valid new "Album Name" (e.g., "My Awesome Vacation Photos")
AND the "Change Folder Name" checkbox is unchecked
WHEN I click the "Save" button
THEN the dialog should enter a disabled state (e.g., with a loading spinner or overlay)
AND a request to rename the album should be sent to the API with the new display name
AND upon a successful API response, the dialog should close
AND the album list should be updated to reflect the new display name (e.g., "My Awesome Vacation Photos")
AND I should remain on the current page if the AlbumId did not change

GIVEN I have the "Edit Album Name" dialog open for an album (e.g., "My Holiday Pics")
AND I have entered a valid new "Album Name" (e.g., "Summer Adventures 2024")
AND I have checked the "Change Folder Name" checkbox
AND I have entered a valid new "Folder Name" (e.g., "summer-adventures-2024")
WHEN I click the "Save" button
THEN the dialog should enter a disabled state
AND a request to rename the album should be sent to the API with the new display name and new folder name
AND upon a successful API response, the dialog should close
AND the album list should be updated to reflect the new display name
AND if the AlbumId changed due to the folder name update, I should be redirected to the newly renamed album

GIVEN I have the "Edit Album Name" dialog open
AND I have entered a new "Album Name" (e.g., "My Holiday Pics")
AND I have checked the "Change Folder Name" checkbox
AND I have cleared the "Folder Name" input field (making it invalid client-side)
AND I have somehow bypassed client-side validation (e.g., by disabling JavaScript)
WHEN I click the "Save" button
THEN the dialog should enter a disabled state
AND a request to rename the album should be sent to the API
AND the API responds with an `AlbumNameMandatoryErr` (for folder name)
THEN an error message "Folder name cannot be empty" should appear below the "Folder Name" field
AND the dialog should exit the disabled state, allowing further interaction
AND the "Save" button should be disabled until the error is resolved

GIVEN I have the "Edit Album Name" dialog open
AND I have entered a new "Album Name" (e.g., "My Holiday Pics")
AND I have checked the "Change Folder Name" checkbox
AND I have entered a "Folder Name" that already exists (e.g., "existing-album")
WHEN I click the "Save" button
THEN the dialog should enter a disabled state
AND a request to rename the album should be sent to the API
AND the API responds with an `AlbumFolderNameAlreadyTakenErr`
THEN an error message "Folder name already exists. Please choose a different one." should appear below the "Folder Name" field
AND the dialog should exit the disabled state, allowing further interaction
AND the "Save" button should be disabled until the error is resolved

GIVEN I have the "Edit Album Name" dialog open
AND I have entered valid input for renaming an album
WHEN I click the "Save" button
THEN the dialog should enter a disabled state
AND a request to rename the album should be sent to the API
AND the API responds with an unexpected error (e.g., a 500 Internal Server Error or a generic `CatalogError`)
THEN a generic error message (e.g., "Something went wrong. Please try again.") should be displayed within the dialog, possibly in an alert block
AND the dialog should exit the disabled state, allowing further interaction
AND the "Save" button should be enabled
```

## Out of scope

* The specific implementation details of the API endpoint for renaming albums.
* The exact visual design of the loading state or error messages, beyond their presence and content.
* Error handling for network issues (e.g., no internet connection) which are typically handled by a global error mechanism.
* The behavior of the "Edit" button for albums not owned by the current user (covered by existing functionality).