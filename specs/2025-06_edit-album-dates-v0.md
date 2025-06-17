# Edit Album Dates Feature

## 1. Feature Summary
The "Edit Dates" feature allows users to update the start and end dates of an album. The feature is accessible via a button on the left menu.

## 2. Ubiquity Language
* **Inclusive Start Date**: The start date of an album, inclusive of the time specified.
* **Exclusive End Date**: The end date of an album, exclusive of the time specified (i.e., the album ends just before this date).

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

4. **Date Validation**
   - User navigates to an album they own and clicks the "Edit Dates" button.
   - User updates the start and end dates.
   - If the start date is after the end date, an error label is displayed on the end date field, and the save button is disabled.
   - User corrects the dates to be valid (start date before or equal to end date).
   - The save button becomes enabled.

5. **Date and Time Selection**
   - User navigates to an album they own and clicks the "Edit Dates" button.
   - User can select a start date and optionally a time by unticking "at the start of the day".
   - User can select an end date and optionally a time by unticking "at the end of the day".
   - The selected dates are displayed in a user-friendly format.

## 4. Technical Context
* The feature uses a new dialog for editing album dates.
* The dialog allows users to select dates and optionally times.
* The API is used to update the album dates and re-fetch album and media data after a successful update.
* The feature is integrated into the existing album management UI.
* The UI handles the conversion between the displayed date format and the expected API format (inclusive start date and exclusive end date).

### Out of Scope
* Impact on album sharing and media ordering.

## 5. Explorations
None identified.
