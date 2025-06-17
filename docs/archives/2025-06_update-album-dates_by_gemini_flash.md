# Feature Requirement Document: Update Album Dates

## 1. Feature Summary

Users will be able to update the start and end dates of albums they own through a dedicated "Edit Album Dates" modal, providing precise control over the album's
temporal boundaries and ensuring data consistency.

## 2. Ubiquity Language

* **"Edit Album Dates" modal:** A dedicated dialog window that appears when a user initiates the date editing process for an album. It contains input fields for
  start and end dates, along with time options.
* **"at the beginning of the day" checkbox:** A UI control associated with the start date input within the "Edit Album Dates" modal. When ticked, it sets the
  time component of the start date to 00:00.
* **"at the end of the day" checkbox:** A UI control associated with the end date input within the "Edit Album Dates" modal. When ticked, it sets the time
  component of the end date to 23:59 (which internally translates to the beginning of the next day for the exclusive end date model).

## 3. Scenarios

### Scenario 1: Successfully Updating an Album's Dates

**Goal:** A user wants to adjust the start and end dates of an album they own.

**Pre-conditions:**

* The user is logged in and viewing the Catalog.
* The user has navigated to an album they own.

**Steps:**

1. The user is viewing the details of an album (e.g., "Summer Vacation 2023").
2. In the left menu, the user observes the "Edit" button, which is enabled because they own this album.
3. The user clicks the "Edit" button.
4. An "Edit Album Dates" modal dialog appears, pre-populated with the current "Start Date" and "End Date" of the album.
    * Next to the "Start Date" field, a checkbox labeled "at the beginning of the day" is displayed. If the current album's start time is 00:00, this checkbox
      is ticked by default.
    * Next to the "End Date" field, a checkbox labeled "at the end of the day" is displayed. If the current album's end time is 23:59, this checkbox is ticked
      by default.
5. The user clicks on the "Start Date" field, and a date picker calendar appears. They select a new start date, e.g., "July 1, 2023".
6. The user clicks on the "End Date" field, and a date picker calendar appears. They select a new end date, e.g., "July 15, 2023".
7. The user decides they want the end date to be precise, so they uncheck the "at the end of the day" checkbox for the End Date. This reveals a time input
   field, which they set to "18:00".
8. The user reviews the selected dates and times, ensuring the start date is strictly before the end date (internally, July 1st 00:00 is before July 15th 18:
   00).
9. The user clicks the "Save" button within the modal.
10. The "Save" button changes to a loading/progress indicator, and the modal becomes unresponsive, indicating an ongoing operation. This loading state persists
    until the album and media lists are fully refreshed.
11. After a short delay, the loading indicator disappears. The modal dialog automatically closes.
12. The main catalog view shows a loading indicator for both the album list and the media list as they refresh.
13. Once refreshed, the album "Summer Vacation 2023" in the album list now displays its updated date range (e.g., "July 1 - July 15, 2023"), and the displayed
    media within the album reflects the new date boundaries.

### Scenario 2: Attempting to Edit an Album Not Owned by the User

**Goal:** A user tries to edit an album they do not own.

**Pre-conditions:**

* The user is logged in and viewing the Catalog.
* The user has navigated to an album they *do not* own (e.g., an album shared with them).

**Steps:**

1. The user is viewing the details of an album (e.g., "Friends' Trip to Paris").
2. In the left menu, the user observes the "Edit" button.
3. Since the user does not own this album, the "Edit" button is visibly disabled (e.g., grayed out).
4. The user attempts to click the disabled "Edit" button.
5. No action occurs; the "Edit Album Dates" modal dialog does not appear, and there is no change in the UI. The user is prevented from initiating the edit
   process.

### Scenario 3: Attempting to Save Invalid Dates (Start Date Not Strictly Before End Date)

**Goal:** A user attempts to save an album with an invalid date range where the start date is not strictly before the end date.

**Pre-conditions:**

* The user is logged in and viewing the Catalog.
* The user has navigated to an album they own.
* The "Edit Album Dates" modal is open.

**Steps:**

1. The user is viewing the "Edit Album Dates" modal for an album they own.
2. The user clicks on the "Start Date" field and selects "July 15, 2023" from the date picker. The "at the beginning of the day" checkbox is ticked.
3. The user then clicks on the "End Date" field and selects "July 1, 2023" from the date picker. The "at the end of the day" checkbox is ticked.
4. Immediately, or upon attempting to click "Save", an error message appears prominently within the modal, for example: "Start Date must be strictly before End
   Date."
5. The "Save" button within the modal becomes disabled (or remains disabled if it was already in that state due to the invalid input), preventing the user from
   submitting the form.
6. The user realizes their mistake. They click on the "End Date" field again and select "July 31, 2023".
7. As soon as the dates are valid (July 15, 2023, is strictly before July 31, 2023), the error message disappears, and the "Save" button becomes enabled again.
8. The user can now proceed to click "Save" (leading to Scenario 1's successful flow or Scenario 4's API failure flow).

### Scenario 4: API Call to Update Album Dates Fails

**Goal:** A user attempts to update album dates, but the server-side operation fails.

**Pre-conditions:**

* The user is logged in and viewing the Catalog.
* The user has navigated to an album they own.
* The "Edit Album Dates" modal is open with valid dates selected.

**Steps:**

1. The user is viewing the "Edit Album Dates" modal for an album they own.
2. The user has selected valid "Start Date" and "End Date" values (e.g., "July 1, 2023" and "July 31, 2023").
3. The user clicks the "Save" button within the modal.
4. The "Save" button changes to a loading/progress indicator, and the modal becomes unresponsive, indicating an ongoing operation.
5. After a short delay (simulating a network error or server-side issue), the loading indicator disappears.
6. An error message is displayed prominently within the modal, for example: "Failed to update album dates. Please try again later." or a more specific error if
   provided by the API (e.g., "Album not found on server").
7. The modal remains open, allowing the user to either correct any potential issues (though in this case, it's a server error) or click "Cancel" to close the
   dialog without saving.
8. The "Save" button becomes enabled again, allowing the user to retry the operation if they wish.

### Scenario 5: User Explicitly Sets a Time for Both Start and End Dates

**Goal:** A user wants to define precise start and end times for an album, rather than using the default "beginning of day" or "end of day" values.

**Pre-conditions:**

* The user is logged in and viewing the Catalog.
* The user has navigated to an album they own.
* The "Edit Album Dates" modal is open.

**Steps:**

1. The user is viewing the "Edit Album Dates" modal for an album they own. The modal displays the current "Start Date" and "End Date".
2. **For the Start Date:**
    * The "at the beginning of the day" checkbox is displayed next to the date input.
    * If it's currently ticked (meaning the album's start time is 00:00), the user *unticks* it.
    * If it's currently unticked (meaning the album already has a specific start time), the user leaves it unticked.
    * Upon unticking (or if already unticked), a time input field (e.g., `HH:MM`) becomes visible next to the date.
    * The user enters a specific time, for example, "09:00".
3. **For the End Date:**
    * The "at the end of the day" checkbox is displayed next to the date input.
    * If it's currently ticked (meaning the album's end time is 23:59), the user *unticks* it.
    * If it's currently unticked (meaning the album already has a specific end time), the user leaves it unticked.
    * Upon unticking (or if already unticked), a time input field (e.g., `HH:MM`) becomes visible next to the date.
    * The user enters a specific time, for example, "17:30".
4. The user ensures that the selected start date and time are strictly before the selected end date and time. If not, the validation error (as in Scenario 3)
   would appear.
5. The user clicks the "Save" button.
6. The modal displays a loading indicator.
7. Upon successful completion, the modal closes, and the album and media lists refresh, reflecting the new precise date and time range for the album.

## 4. Technical Context

* **Provided by the platform/supporting APIs:**
    * Existing UI framework for rendering the left menu, album list, and media display.
    * Mechanism to identify the "current album" being viewed.
    * API for retrieving album details, including ownership information (`album.ownedBy` field in `Album` interface).
    * API for updating album dates (requires `AlbumId`, new start date, and new end date as a single operation; returns empty success on success).
    * API for fetching updated album and media lists.
    * User authentication and authorization services to determine current user and their ownership.
    * The internal data model for album dates stores the start date as inclusive and the end date as exclusive (e.g., an end date of "July 1st 00:00" internally
      represents "June 30th 23:59" for UI display purposes when "at the end of the day" is ticked).
* **Out of Scope:**
    * Detailed API contract specifications (exact endpoints, request/response schemas, specific error codes beyond general failure).
    * Internal implementation details of UI components (e.g., how the date picker is built, specific state management library choices like Redux, Zustand,
      etc.).
    * Performance optimizations beyond the basic refresh mechanism described.

## 5. List of Explorations

* What are the precise API endpoints and payload structures required for updating album dates and for triggering the refresh of album and media lists?
* How is album ownership currently determined and communicated to the frontend, and what is the exact mechanism used to enable/disable UI elements based on this
  information (e.g., is `albumIsOwnedByCurrentUser` sufficient, or are there other checks)?
* What is the existing data fetching and state management strategy for albums and media, and how can the post-save refresh integrate seamlessly and efficiently
  with it to ensure a smooth user experience?
