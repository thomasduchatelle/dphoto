# As a user, I can specify time for the start and end dates.

## Acceptance Criteria:
*   **GIVEN** I have opened the "Edit Album Dates" dialog for an album I own
*   **WHEN** I uncheck the "at the start of the day" checkbox for the start date
*   **THEN** The time input field for the start date becomes enabled
*   **AND** I can input a specific time (HH:MM) for the start date
*   **WHEN** I uncheck the "at the end of the day" checkbox for the end date
*   **THEN** The time input field for the end date becomes enabled
*   **AND** I can input a specific time (HH:MM) for the end date
*   **WHEN** I re-check the "at the start of the day" checkbox
*   **THEN** The time input for the start date is disabled and resets to "00:00"
*   **WHEN** I re-check the "at the end of the day" checkbox
*   **THEN** The time input for the end date is disabled and resets to "23:59"
*   **WHEN** I save the album dates with specified times
*   **THEN** The start date is sent to the API as an inclusive `Date` object (e.g., `2025-01-15T10:00:00`)
*   **AND** The end date is sent to the API as an exclusive `Date` object (e.g., if user inputs `2025-01-20T15:00:00`, it's sent as `2025-01-20T15:00:01`)

## Examples:
*   **Scenario: Album initially has non-midnight times**
    *   **GIVEN** An album has `start: 2024-03-10T09:30:00` and `end: 2024-03-15T18:00:00`
    *   **WHEN** I open the "Edit Album Dates" dialog
    *   **THEN** The start date picker shows "2024-03-10" and the time input shows "09:30"
    *   **AND** The "at the start of the day" checkbox is unchecked
    *   **AND** The end date picker shows "2024-03-15" and the time input shows "18:00"
    *   **AND** The "at the end of the day" checkbox is unchecked
*   **Scenario: Converting exclusive end date to display**
    *   **GIVEN** An album has `start: 2025-01-01T00:00:00` and `end: 2025-02-01T00:00:00` (representing Jan 2025)
    *   **WHEN** I open the "Edit Album Dates" dialog
    *   **THEN** The start date picker shows "2025-01-01" and the "at the start of the day" checkbox is checked
    *   **AND** The end date picker shows "2025-01-31" and the "at the end of the day" checkbox is checked
*   **Scenario: User selects date and time, then re-checks "at the start of the day"**
    *   **GIVEN** I have unchecked "at the start of the day" and set the start date to "2025-01-15" and time to "10:00"
    *   **WHEN** I re-check the "at the start of the day" checkbox
    *   **THEN** The time input for the start date is disabled and shows "00:00"
    *   **AND** The internal state for the start date is `2025-01-15T00:00:00`
*   **Scenario: User selects date and time, then re-checks "at the end of the day"**
    *   **GIVEN** I have unchecked "at the end of the day" and set the end date to "2025-01-20" and time to "15:00"
    *   **WHEN** I re-check the "at the end of the day" checkbox
    *   **THEN** The time input for the end date is disabled and shows "23:59"
    *   **AND** The internal state for the end date, before conversion to exclusive, is `2025-01-20T23:59:59`
