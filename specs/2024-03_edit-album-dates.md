# Edit Album Dates Feature

## 1. Feature Summary
The "Edit Dates" feature allows users to update the start and end dates of an album. The feature is accessible via a button on the left menu.

## 2. Ubiquity Language
No new terms or entities are introduced.

## 3. Scenarios

1. **Successful Edit**
   - User navigates to an album they own.
   - User clicks the "Edit Dates" button on the left menu.
   - A dialog appears with the album name, current start and end dates, and an option to input time.
   - User updates the dates and submits.
   - The dialog shows a loading sign.
   - Upon successful API request, the dialog closes, and the album list and media are refreshed.

2. **Failed Edit**
   - User navigates to an album they own.
   - User clicks the "Edit Dates" button on the left menu.
   - A dialog appears with the album name, current start and end dates, and an option to input time.
   - User updates the dates and submits.
   - The dialog shows a loading sign.
   - Upon failed API request, the dialog remains open, displaying the error.

3. **Permission Check**
   - User navigates to an album they don't own.
   - The "Edit Dates" button on the left menu is disabled.

## 4. Technical Context
* The feature reuses the existing dialog and logic from the "Create Dialog".
* The API is used to update the album dates and re-fetch album and media data after a successful update.
* The feature is integrated into the existing album management UI.

### Out of Scope
* Impact on album sharing and media ordering.

## 5. Explorations
None identified.
