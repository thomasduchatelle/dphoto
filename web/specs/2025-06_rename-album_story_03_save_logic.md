# As a user, I want to persist the new album name(s) and see relevant feedback for success or API errors

## Acceptance Criteria

```
GIVEN I have the rename dialog open with "My Old Album" as the current name
AND the "Folder Name" checkbox is unchecked
WHEN I change the "Album Name" field to "My Awesome Vacation Photos" and click "Save"
THEN the API is called with the new display name only
AND when the API responds successfully, the dialog closes
AND the album list updates to show "My Awesome Vacation Photos"
AND if the returned AlbumId differs from the current one, I am redirected to the new album URL

GIVEN I have both "Album Name" set to "Summer Adventures 2024" and "Folder Name" checked with value "summer-adventures-2024"
WHEN I click "Save"
THEN the API is called with both the new display name and new folder name
AND when the API responds successfully with a new AlbumId, the dialog closes
AND if the returned AlbumId differs from the current one, I am redirected to the new album URL


GIVEN I have the rename dialog open with a folder name that already exists
WHEN I click "Save" and the API returns "AlbumFolderNameAlreadyTakenErr"
THEN the loading state ends
AND an error message appears below the "Folder Name" field saying "This name is already taken"
AND the "Save" button becomes disabled

GIVEN I have the rename dialog open with valid data
WHEN I click "Save" and the API returns an unexpected error
THEN the loading state ends
AND an Alert block appears within the dialog showing the message from the API error (or "Something went wrong. Please try again." if no message)
AND the form remains enabled so I can try again or cancel

GIVEN I have the rename dialog open with a field error showing below the "Folder Name" field
WHEN I modify the "Folder Name" field
THEN that field's error message disappears and the "Save" button becomes enabled if all fields are valid

GIVEN I have the rename dialog open with an Alert block showing a generic error
WHEN I modify any field in the form
THEN the Alert block disappears
```

## Out of scope
- Retry mechanisms or automatic error recovery
- Logging or reporting of errors to external systems
- Network connectivity error handling
- Complex redirection logic beyond basic URL updates
