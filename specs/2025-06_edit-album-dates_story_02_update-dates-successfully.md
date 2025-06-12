# As a user, I can update the dates of an album successfully.

## Acceptance Criteria:

```
GIVEN the "Edit Dates" dialog is open for an album I own
WHEN I update the start and end dates in the dialog
AND I click the "Save" button
THEN the dialog shows a loading indicator
AND an API request is sent to update the album with the new dates (e.g., if I set start to "2023-07-10" and end to "2023-07-20", the API request will contain start: 2023-07-10T00:00:00, end: 2023-07-21T00:00:00)
AND upon successful API response, the dialog closes, the albums list is refreshed, and the medias for the current album are refreshed simultaneously.
```

## Out of scope:

* Handling cases where the API request fails (this is covered in a separate story).
* Validation of the dates (e.g., start date before end date) before submission (this is covered in a separate story).
* Specifying time for the start and end dates (this is covered in a separate story).