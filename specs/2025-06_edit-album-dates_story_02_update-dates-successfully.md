# As a user, I can update the dates of an album successfully.

## Acceptance Criteria:
*   **GIVEN** I have opened the "Edit Album Dates" dialog for an album I own
*   **WHEN** I modify the start and/or end dates (and optionally times) in the dialog
*   **AND** I click the "Save" button
*   **THEN** A loading indicator is displayed within the dialog
*   **AND** An API request is sent to update the album with the new dates
*   **AND** Upon successful API response, the dialog closes
*   **AND** The album list is refreshed to reflect the updated dates
*   **AND** The media displayed for the album is refreshed to reflect the updated dates

## Out of Scope:
*   Specific visual representation of the loading indicator
*   Error handling for failed updates
*   Validation of date inputs
