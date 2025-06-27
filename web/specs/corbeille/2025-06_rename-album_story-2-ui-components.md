# As a user, I want to rename an album by changing its display name and optionally its folder name, with appropriate client-side validation.

## Acceptance Criteria

```
GIVEN I have the "Edit Album Name" dialog open for an album (e.g., "My Holiday Pics")
AND the "Change Folder Name" checkbox is unchecked
WHEN I type a valid name (e.g. "My Awesome Vacation Photos") into the "Album Name" input field
THEN the "Save" button should remain enabled

GIVEN I have the "Edit Album Name" dialog open for an album (e.g., "My Holiday Pics")
AND the "Change Folder Name" checkbox is unchecked
WHEN I clear the "Album Name" input field
THEN a validation error message "Album name cannot be empty" should appear below the "Album Name" field
AND the "Save" button should be disabled

GIVEN I have the "Edit Album Name" dialog open for an album (e.g., "My Holiday Pics")
AND the "Change Folder Name" checkbox is unchecked
AND the "Album Name" input field is empty, and the "Save" button is disabled
WHEN I type "My Awesome Vacation Photos" into the "Album Name" input field
THEN the validation error message should disappear
AND the "Save" button should become enabled

GIVEN I have the "Edit Album Name" dialog open for an album (e.g., "My Holiday Pics")
WHEN I check the "Change Folder Name" checkbox
THEN a new input field labeled "Folder Name" should appear below the checkbox
AND the "Folder Name" input field should be pre-filled with the current folder name (e.g., "my-holiday-pics")
AND the "Save" button should remain enabled

GIVEN I have the "Edit Album Name" dialog open for an album (e.g., "My Holiday Pics")
AND the "Change Folder Name" checkbox is checked
WHEN I clear the "Folder Name" input field
THEN a validation error message "Folder name cannot be empty" should appear below the "Folder Name" field
AND the "Save" button should be disabled

GIVEN I have the "Edit Album Name" dialog open for an album (e.g., "My Holiday Pics")
AND the "Change Folder Name" checkbox is checked
AND the "Folder Name" input field is empty, and the "Save" button is disabled
WHEN I type "my-awesome-vacation-photos" into the "Folder Name" input field
THEN the validation error message should disappear
AND the "Save" button should become enabled

GIVEN I have the "Edit Album Name" dialog open for an album (e.g., "My Holiday Pics")
AND the "Change Folder Name" checkbox is checked
WHEN I type "My New Folder Name With Spaces" into the "Folder Name" input field
THEN a validation error message "Folder name can only contain lowercase letters, numbers, and hyphens" should appear below the "Folder Name" field
AND the "Save" button should be disabled

GIVEN I have the "Edit Album Name" dialog open for an album (e.g., "My Holiday Pics")
AND the "Change Folder Name" checkbox is checked
AND the "Folder Name" input field contains invalid characters, and the "Save" button is disabled
WHEN I type "my-new-folder-name" into the "Folder Name" input field
THEN the validation error message should disappear
AND the "Save" button should become enabled
```

## Out of scope

* Server-side validation of the album name or folder name.
* The actual API call to rename the album.
* Handling of specific API errors (e.g., folder name already taken), which will be covered in the next story.
* The visual appearance of the validation messages beyond their presence and content.