# As a user, I am prevented from setting an invalid date range.

## Acceptance Criteria:

```
GIVEN the "Edit Dates" dialog is open for an album I own
AND the start date is "2023-01-10" and the end date is "2023-01-20"
WHEN I change the start date to "2023-01-25" (making it after the current end date)
THEN an error message is displayed next to the end date field, indicating that the start date cannot be after the end date
AND the "Save" button is disabled

GIVEN the "Edit Dates" dialog is open for an album I own
AND the start date is "2023-01-10" and the end date is "2023-01-20"
WHEN I change the end date to "2023-01-05" (making it before the current start date)
THEN an error message is displayed next to the end date field, indicating that the end date cannot be before the start date
AND the "Save" button is disabled

GIVEN the "Edit Dates" dialog is open for an album I own
AND the start date is "2023-01-25" and the end date is "2023-01-20" (an invalid range, with the Save button disabled)
WHEN I correct the end date to "2023-01-30" (making it valid)
THEN the error message next to the end date field disappears
AND the "Save" button becomes enabled

GIVEN the "Edit Dates" dialog is open for an album I own
AND the start date is "2023-01-25" and the end date is "2023-01-20" (an invalid range, with the Save button disabled)
WHEN I correct the start date to "2023-01-15" (making it valid)
THEN the error message next to the end date field disappears
AND the "Save" button becomes enabled
```

## Out of scope:

* Specific wording of the error message, as long as it clearly indicates the issue.
* Validation of date format (assuming the date picker ensures valid date formats).
* Validation of time inputs (covered in a separate story).