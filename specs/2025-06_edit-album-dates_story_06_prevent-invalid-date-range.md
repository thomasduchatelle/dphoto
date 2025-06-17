# As a user, I am prevented from setting an invalid date range.

## Acceptance Criteria:
*   **GIVEN** I have opened the "Edit Album Dates" dialog
*   **WHEN** I set the start date to be after the end date (e.g., Start: 2025-01-10, End: 2025-01-05)
*   **THEN** An error message is displayed next to the end date field, indicating an invalid range
*   **AND** The "Save" button is disabled
*   **WHEN** I correct the dates so that the start date is on or before the end date
*   **THEN** The error message is removed
*   **AND** The "Save" button becomes enabled

## Out of Scope:
*   Specific wording of the error message
*   Real-time validation as dates are typed (validation can occur on blur or change)
