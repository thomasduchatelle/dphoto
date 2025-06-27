# As a user, I want to open and close the "Edit Name" dialog for an album I own.

## Acceptance Criteria

```
GIVEN I am on the catalog viewer page
AND I am logged in as a user who owns an album (e.g., "My Holiday Pics")
WHEN I click on the "Edit" button associated with "My Holiday Pics"
THEN a dropdown menu should appear
AND the dropdown menu should contain an option labeled "Edit Name"

GIVEN I am on the catalog viewer page
AND I am logged in as a user who owns an album (e.g., "My Holiday Pics")
WHEN I click on the "Edit" button associated with "My Holiday Pics"
AND I select "Edit Name" from the dropdown menu
THEN a modal dialog titled "Edit the name of the album <current album name>" should appear
AND the dialog should display a "Cancel" button

GIVEN I have the "Edit Album Name" dialog open
WHEN I click the "Cancel" button
THEN the "Edit Album Name" dialog should close
AND I should be returned to the catalog viewer page
```

## Out of scope

* The behavior of the "Edit" button for albums not owned by the current user (this is existing functionality).
* Any other content in this dialog is covered on the story 2 (text fields, ...)
* The specific styling or exact layout of the dialog, beyond the presence of the specified elements.