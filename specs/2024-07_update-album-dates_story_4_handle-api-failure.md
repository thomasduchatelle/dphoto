# Story 4: Handle API failure during album date update.

As a user, I want to be informed if saving album dates fails due to a server error, so I can retry or cancel.

## Acceptance Criteria

### AC 4.1: API call to update album dates fails

GIVEN the user has clicked "Save" and the loading indicator is active (AC 3.2)
WHEN the backend API call to update album dates returns an error response (e.g., HTTP 500, network error, specific API error)
THEN the loading indicator should disappear.

### AC 4.2: Display error message on modal

GIVEN the API call has failed (AC 4.1)
WHEN the loading indicator has disappeared
THEN an error message (e.g., "Failed to update album dates. Please try again later." or a more specific message from the API) should be displayed prominently within the "Edit Album Dates" modal.

### AC 4.3: Modal remains open on failure

GIVEN the API call has failed and an error message is displayed (AC 4.2)
WHEN the error occurs
THEN the "Edit Album Dates" modal should remain open, allowing the user to view the error.

### AC 4.4: "Save" button re-enables on failure

GIVEN the API call has failed and an error message is displayed (AC 4.2)
WHEN the error occurs
THEN the "Save" button within the modal should become enabled again, allowing the user to retry the operation or click "Cancel".

## Out of Scope (covered by other stories)

*   Opening the "Edit Album Dates" modal.
*   Client-side validation of date inputs.
*   Successful saving of album dates.
*   Refreshing album or media lists.
