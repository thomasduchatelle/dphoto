# As a user, I am informed when updating album dates fails.

## Acceptance Criteria:

```
GIVEN the "Edit Dates" dialog is open for an album I own
WHEN I update the start and end dates in the dialog
AND I click the "Save" button
AND the API request to update the album dates fails
THEN the dialog remains open
AND an error message is displayed in an Alert area within the dialog, indicating the failure
AND the loading indicator is no longer displayed
AND the albums list is not refreshed
AND the medias for the current album are not refreshed

GIVEN the "Edit Dates" dialog is open with an error message displayed from a previous failed save attempt
WHEN I change one of the dates (start or end date)
THEN the error message in the Alert area is cleared

GIVEN the "Edit Dates" dialog is open with an error message displayed from a previous failed save attempt
WHEN I click the "Save" button again
THEN the error message in the Alert area is cleared
AND the dialog shows a loading indicator
```

## Out of scope:

* Specific error messages for different failure types (e.g., network error vs. server error). A generic error message is sufficient.
* Automatic retry mechanisms.
* Validation of the dates (e.g., start date before end date) before submission (this is covered in a separate story).
* Specifying time for the start and end dates (this is covered in a separate story).
