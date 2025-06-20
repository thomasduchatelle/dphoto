# As a user, I can specify time for the start and end dates.

## Acceptance Criteria:

```
GIVEN the "Edit Dates" dialog is open for an album I own
AND the "at the start of the day" checkbox for the start date is checked
AND the "at the end of the day" checkbox for the end date is checked
WHEN I uncheck the "at the start of the day" checkbox
THEN a time input field appears next to the start date field
AND the time input field defaults to "00:00"
WHEN I enter "10:30" into the start time input field
AND I click the "Save" button
THEN the API request to update the album includes the start date with the specified time, for example, "YYYY-MM-DDT10:30:00"

GIVEN the "Edit Dates" dialog is open for an album I own
AND the "at the start of the day" checkbox for the start date is checked
AND the "at the end of the day" checkbox for the end date is checked
WHEN I uncheck the "at the end of the day" checkbox
THEN a time input field appears next to the end date field
AND the time input field defaults to "23:59"
WHEN I enter "15:00" into the end time input field
AND I click the "Save" button
THEN the API request to update the album includes the end date with the specified time, converted to be exclusive by adding one minute, for example, if the user selected "YYYY-MM-DD 15:00", the API request will contain "YYYY-MM-DDT15:01:00"

GIVEN the "Edit Dates" dialog is open for an album I own
AND the "at the start of the day" checkbox for the start date is unchecked
AND a specific time is entered in the start time input field (e.g., "10:30")
WHEN I check the "at the start of the day" checkbox
THEN the time input field for the start date disappears
AND the internal time for the start date reverts to "00:00"
AND if I click the "Save" button, the API request will contain "YYYY-MM-DDT00:00:00" for the start date

GIVEN the "Edit Dates" dialog is open for an album I own
AND the "at the end of the day" checkbox for the end date is unchecked
AND a specific time is entered in the end time input field (e.g., "15:00")
WHEN I check the "at the end of the day" checkbox
THEN the time input field for the end date disappears
AND the internal time for the end date reverts to "23:59"
AND if I click the "Save" button, the API request will contain "YYYY-MM-DD+1T00:00:00" for the exclusive end date (converting 23:59 end-of-day to next day 00:00:00)

GIVEN I open the "Edit Dates" dialog for an album that has a start date of "2023-01-01T10:00:00" and an end date of "2023-01-05T14:30:00"
THEN the "at the start of the day" checkbox for the start date is unchecked
AND the start time input field displays "10:00"
AND the "at the end of the day" checkbox for the end date is unchecked
AND the end time input field displays "14:30"
```

## Out of scope:

* Validation of time input format (e.g., ensuring it's a valid time).
* Interaction with the date picker itself (only the time input is in scope here).
