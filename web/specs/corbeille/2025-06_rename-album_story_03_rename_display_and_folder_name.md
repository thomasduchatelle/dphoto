# As an album owner, I want to rename both the album display name and folder name, and be redirected if the AlbumId changes.

**Acceptance Criteria:**
```
GIVEN I am an album owner viewing an album
AND the "Edit Name" dialog is open with the current album name pre-filled
WHEN I enter a new album display name
AND I check the "Folder Name" checkbox
AND I enter a new folder name
AND I click "Save"
THEN the dialog enters a disabled mode with a loading bar
AND upon successful API response, the dialog closes
AND the album list updates to show the new album display name
AND if the AlbumId changed due to folder name update, I am redirected to the new album page
```

**Out of scope:**
- Validation errors
- API error handling
```

specs/2025-06_rename-album_story_04_client_side_validation_blank_album_name.md
```markdown
<<<<<<< SEARCH
