# Story 3: Successfully save album dates and refresh UI.

As a user, I want to save the updated album dates and see the catalog refresh to reflect my changes.

## Acceptance Criteria

### AC 3.1: Initiate save operation

GIVEN the "Edit Album Dates" modal is open with valid dates and times selected (AC 2.5 is met)
WHEN the user clicks the "Save" button
THEN a request to update the album dates should be initiated to the backend API.

### AC 3.2: Show loading indicator on modal during API call

GIVEN the user has clicked "Save" (AC 3.1)
WHEN the API call to update album dates is in progress
THEN:
    - The "Save" button should change to a loading/progress indicator.
    - The modal should become unresponsive (e.g., inputs disabled, overlay active), indicating an ongoing operation.

### AC 3.3: Close modal on successful API response

GIVEN the API call to update album dates returns a successful response (empty success)
WHEN the loading indicator from AC 3.2 is active
THEN the "Edit Album Dates" modal dialog should automatically close.

### AC 3.4: Show loading indicators for album/media lists during refresh

GIVEN the "Edit Album Dates" modal has closed after a successful save (AC 3.3)
WHEN the album and media lists in the main catalog view are being refreshed
THEN loading indicators should be displayed for both the album list and the media list.

### AC 3.5: Refresh album and media lists to reflect new dates

GIVEN the album and media lists are refreshing (AC 3.4)
WHEN the refresh operation completes successfully
THEN:
    - The loading indicators should disappear.
    - The album in the album list should display its updated date range.
    - The displayed media within the album should reflect the new date boundaries.

## Out of Scope (covered by other stories)

*   Opening the "Edit Album Dates" modal.
*   Client-side validation of date inputs.
*   Handling API errors during save.
