# As an album owner, I want to open the "Edit Name" dialog pre-filled with the current album name and folder name checkbox unchecked.

**Acceptance Criteria:**
```
GIVEN I am an album owner viewing an album
WHEN I click the "Edit" button and select "Edit Name" from the dropdown
THEN the "Edit Name" dialog opens

GIVEN the dialog is open
WHEN I look at the "Album Name" field
THEN it is pre-filled with the current album display name

GIVEN the dialog is open
WHEN I look at the "Folder Name" checkbox
THEN it is unchecked by default

GIVEN the dialog is open
WHEN the "Folder Name" checkbox is unchecked
THEN the folder name text field is disabled or hidden
```

**Out of scope:**
- Actual renaming functionality
- Validation of fields
- API interactions
```

specs/2025-06_rename-album_story_02_rename_display_name_only.md
```markdown
<<<<<<< SEARCH
