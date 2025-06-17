# Story 1: View "Edit" button and open date edit modal for owned albums.

As a user, I want to see an "Edit" button for albums I own, and be able to open a modal to edit their dates.

## Acceptance Criteria

### AC 1.1: "Edit" button enabled for owned albums

GIVEN the user is viewing an album
WHEN the `albumIsOwnedByCurrentUser` selector is called for the current album and returns `true`
THEN the "Edit" button in the left menu should be enabled.

### AC 1.2: "Edit" button disabled for unowned albums

GIVEN the user is viewing an album
WHEN the `albumIsOwnedByCurrentUser` selector is called for the current album and returns `false`
THEN the "Edit" button in the left menu should be disabled.

### AC 1.3: Open "Edit Album Dates" modal

GIVEN the user is viewing an album they own and the "Edit" button is enabled
WHEN the user clicks the "Edit" button
THEN an "Edit Album Dates" modal dialog should appear, pre-populated with the current album's `start` and `end` dates.

### AC 1.4: "at the beginning/end of the day" checkboxes default state

GIVEN the "Edit Album Dates" modal is open
WHEN the modal is initialized with the current album's dates
THEN:
    - The "at the beginning of the day" checkbox for the Start Date should be ticked if the current album's start time is 00:00.
    - The "at the end of the day" checkbox for the End Date should be ticked if the current album's end time is 23:59.

## Out of Scope (covered by other stories)

*   Actual modification of dates or times within the modal.
*   Client-side validation of date inputs.
*   Saving changes to the backend.
*   Handling API success or failure.
*   Refreshing album or media lists after save.
