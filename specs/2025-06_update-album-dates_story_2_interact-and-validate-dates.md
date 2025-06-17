# Story 2: Interact with date/time inputs and client-side validation.

As a user, I want to precisely set the start and end dates and times for an album, and be guided by immediate validation feedback.

## Acceptance Criteria

### AC 2.1: Select dates using date picker

GIVEN the "Edit Album Dates" modal is open
WHEN the user interacts with the "Start Date" or "End Date" input fields
THEN a date picker component should appear, allowing the user to select a new date.

### AC 2.2: Toggle "at the beginning/end of the day" checkboxes and show/hide time inputs

GIVEN the "Edit Album Dates" modal is open
WHEN the user ticks or unticks the "at the beginning of the day" checkbox for the Start Date
THEN:
    - If ticked, the time component of the Start Date should be set to 00:00, and any associated time input field should be hidden.
    - If unticked, a time input field (e.g., HH:MM) should become visible, allowing the user to specify a precise time.
AND
WHEN the user ticks or unticks the "at the end of the day" checkbox for the End Date
THEN:
    - If ticked, the time component of the End Date should be set to 23:59, and any associated time input field should be hidden.
    - If unticked, a time input field (e.g., HH:MM) should become visible, allowing the user to specify a precise time.

### AC 2.3: Enter specific times

GIVEN the "Edit Album Dates" modal is open and a time input field is visible for either Start or End Date
WHEN the user enters a time into the time input field
THEN the time component of the corresponding date should be updated to the entered value.

### AC 2.4: Client-side validation: Start Date strictly before End Date

GIVEN the "Edit Album Dates" modal is open and the user has selected dates and times
WHEN the calculated Start Date (including time) is not strictly before the calculated End Date (including time, where 23:59 for End Date means the beginning of the next day internally)
THEN:
    - An error message (e.g., "Start Date must be strictly before End Date.") should be displayed prominently within the modal.
    - The "Save" button within the modal should be disabled, preventing submission.

### AC 2.5: Validation error disappears and Save button enables when dates become valid

GIVEN the "Edit Album Dates" modal is open and an invalid date range is currently selected (triggering AC 2.4)
WHEN the user adjusts the Start Date or End Date (or their times) such that the Start Date is now strictly before the End Date
THEN:
    - The error message displayed in AC 2.4 should disappear.
    - The "Save" button within the modal should become enabled.

## Out of Scope (covered by other stories)

*   Opening the "Edit Album Dates" modal.
*   Saving changes to the backend.
*   Handling API success or failure.
*   Refreshing album or media lists after save.
