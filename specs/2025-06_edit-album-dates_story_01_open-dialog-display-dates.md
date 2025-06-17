# As a user, I can open the "Edit Dates" dialog for an album I own, displaying its current dates and name.

## Acceptance Criteria:
*   **GIVEN** I am on the album view page
*   **AND** I own the current album
*   **WHEN** I click the "Edit Dates" button in the left menu
*   **THEN** A modal dialog titled "Edit Album Dates" appears
*   **AND** The dialog displays the current album name prominently
*   **AND** The dialog displays the current start date of the album in a date picker field
*   **AND** The dialog displays the current end date of the album in a date picker field
*   **AND** Below the start date picker, there is a checkbox labeled "at the start of the day"
*   **AND** Below the end date picker, there is a checkbox labeled "at the end of the day"
*   **AND** If the album's start time is 00:00:00, the "at the start of the day" checkbox is checked by default. Otherwise, it is unchecked.
*   **AND** If the album's end time is 23:59:59, the "at the end of the day" checkbox is checked by default. Otherwise, it is unchecked.
*   **AND** If the "at the start of the day" checkbox is checked, the time input for the start date is disabled and shows "00:00".
*   **AND** If the "at the end of the day" checkbox is checked, the time input for the end date is disabled and shows "23:59".
*   **AND** The dialog contains "Save" and "Cancel" buttons.
*   **AND** Clicking the "Cancel" button closes the dialog without making any changes.

## Out of Scope:
*   Saving the changes
*   Validation of dates
*   Error handling
