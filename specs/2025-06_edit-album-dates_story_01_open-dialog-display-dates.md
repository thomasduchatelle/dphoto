# As a user, I can open the "Edit Dates" dialog for an album I own, displaying its current dates and name.

## Acceptance Criteria:

```
GIVEN I am viewing an album that I own, for example, an album named "Summer Trip" with a start date of 2023-07-01T00:00:00 and an end date of 2023-08-01T00:00:00
WHEN I click the "Edit Dates" button between "Create" and "Delete" buttons above the list of albums
THEN the "Edit Dates" dialog is displayed
AND the dialog title shows the name of the current album, for example, "Summer Trip"
AND the dialog displays the current start date of the album, for example, "2023-07-01"
AND the "at the start of the day" checkbox for the start date is checked
AND the dialog displays the current end date of the album, for example, "2023-07-31" (derived from the exclusive end date 2023-08-01T00:00:00)
AND the "at the end of the day" checkbox for the end date is checked

GIVEN I am viewing an album that I own, for example, an album named "Winter Holidays" with a start date of 2024-12-20T00:00:00 and an end date of 2025-01-05T00:00:00
WHEN I click the "Edit Dates" button between "Create" and "Delete" buttons above the list of albums
THEN the "Edit Dates" dialog is displayed
AND the dialog title shows the name of the current album, for example, "Winter Holidays"
AND the dialog displays the current start date of the album, for example, "2024-12-20"
AND the "at the start of the day" checkbox for the start date is checked
AND the dialog displays the current end date of the album, for example, "2025-01-04" (derived from the exclusive end date 2025-01-05T00:00:00)
AND the "at the end of the day" checkbox for the end date is checked

GIVEN the "Edit Dates" dialog is open
WHEN I click the "Cancel" button
THEN the "Edit Dates" dialog is closed
AND no changes are applied to the album
```

## Out of scope:

* Saving any changes made in the dialog.
* Validation of the dates entered by the user.
* Handling cases where the user does not own the album (this is covered in a separate story).
* Displaying or allowing input for specific times when the "at the start/end of the day" checkboxes are unchecked (this is covered in a separate story).