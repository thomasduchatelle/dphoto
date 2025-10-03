# As a user, I want to rename an album by changing its display name and optionally its folder name, with appropriate client-side validation

## Acceptance Criteria

```
GIVEN I am viewing an album that I own
WHEN I click the "Edit" button and select "Edit Name"
THEN a modal dialog opens with the title "Edit Name"
AND the dialog contains an "Album Name" text field pre-filled with the current album display name
AND the dialog contains a "Folder Name" checkbox that is unchecked by default
AND the dialog contains a "Folder Name" text field that is disabled and empty
AND the dialog contains "Save" and "Cancel" buttons

GIVEN I have the rename dialog open with "My Vacation Photos" pre-filled
WHEN I clear the "Album Name" field completely
THEN a validation error appears below the field saying "Album name cannot be blank"
AND the "Save" button becomes disabled

GIVEN I have the rename dialog open with an empty "Album Name" field showing an error
WHEN I enter "Summer Adventures 2024" in the Album Name field
THEN the validation error disappears and the "Save" button becomes enabled

GIVEN I have the rename dialog open
WHEN I check the "Folder Name" checkbox
THEN the "Folder Name" text field becomes enabled
AND the field is pre-filled with the current folder name

GIVEN the "Folder Name" checkbox is checked and the field is enabled
WHEN I uncheck the "Folder Name" checkbox
THEN the "Folder Name" text field becomes disabled
AND the field value is cleared
AND a placeholder text appears saying "Folder name will be generated from the new name"

GIVEN I have the "Folder Name" checkbox checked with "summer-vacation" pre-filled
WHEN I clear the folder name field completely
THEN a validation error appears below the field saying "Folder name cannot be blank"
AND the "Save" button becomes disabled

GIVEN I have the rename dialog open with valid data in both fields
WHEN I click "Save"
THEN the dialog shows a loading state with a disabled form and loading bar
```

## Out of scope
- API integration and server-side validation (handled in next story)
- Complex validation rules beyond empty field checks
- Folder name format validation or character restrictions
- Real-time validation of folder name uniqueness
