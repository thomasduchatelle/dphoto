# As an album owner, I want to see and handle API validation errors for blank folder name when the folder name checkbox is checked.

**Acceptance Criteria:**
```
GIVEN I am an album owner viewing an album
AND the "Edit Name" dialog is open
AND I have checked the "Folder Name" checkbox
AND I leave the folder name field blank
WHEN I click "Save"
THEN the dialog enters a disabled mode with a loading bar
AND the API returns an AlbumNameMandatoryErr
AND an error message appears below the "Folder Name" field indicating it cannot be blank
AND the "Save" button is disabled
AND the dialog exits the disabled mode

GIVEN the error message is shown
WHEN I enter a valid folder name
THEN the error message disappears
AND the "Save" button is enabled
```

**Out of scope:**
- Client-side validation of folder name
- Other API errors
```

specs/2025-06_rename-album_story_06_api_error_folder_name_taken.md
```markdown
<<<<<<< SEARCH
