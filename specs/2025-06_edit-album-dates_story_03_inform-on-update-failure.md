# As a user, I am informed when updating album dates fails.

## Acceptance Criteria:
*   **GIVEN** I have opened the "Edit Album Dates" dialog for an album I own
*   **AND** I have modified the start and/or end dates (and optionally times)
*   **WHEN** I click the "Save" button
*   **AND** The API request to update the album dates fails
*   **THEN** The loading indicator is dismissed
*   **AND** The dialog remains open
*   **AND** An error message is displayed within the dialog, indicating the failure

## Out of Scope:
*   Specific content of the error message
*   Retry mechanisms
