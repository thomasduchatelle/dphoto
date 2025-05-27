# Delete Album Feature Specification

## Overview

Enable users to delete any album they own from the UI. The deletion is definitive and immediate, with no undo. The backend already supports album deletion and will return errors if media cannot be re-allocated. The UI should provide a smooth, clear, and responsive experience across all device types.

---

## Functional Requirements

1. **Entry Point**
   - The delete album feature is accessible from the album list view (left menu).
   - A placeholder button in the UI will open the delete dialog.

2. **Delete Dialog**
   - When opened, the dialog pre-selects the currently selected album.
   - The user can change the selection to any other deletable album (i.e., albums owned by the user).
   - The dropdown in the dialog lists only albums the user owns and can delete.
   - Each album in the dropdown displays its name and the number of media it contains.
   - The dialog contains a "Delete" button and a "Confirm" button (two-step confirmation).

3. **Deletion Process**
   - On clicking "Delete" and then "Confirm," the UI sends a request to the backend to delete the selected album.
   - The dialog shows a loading indicator during the deletion process and while the album list is reloading.
   - After successful deletion:
     - The dialog closes automatically.
     - The album list is refreshed to reflect the new state and media distribution.
     - If the deleted album was being viewed, the user is redirected to the first available album.
   - If deletion fails:
     - The dialog displays a detailed, user-friendly error message based on the error type returned by the backend.

4. **User Experience**
   - The feature is available and usable on all device types (desktop, tablet, mobile).
   - No special UI/UX constraints beyond general usability.
   - No additional warnings or reminders are required, but a confirmation step is present.

5. **Permissions & Security**
   - Only albums owned by the current user are deletable.
   - Standard authentication is required; no additional permissions or re-authentication needed.

6. **Data Handling**
   - The backend endpoint follows REST conventions:  
     `DELETE /api/v1/owners/{owner}/albums/{folderName}`
   - The backend returns a detailed error type (string code) and message if deletion fails.

7. **Persistence & Logging**
   - No need to log or track deletions for audit purposes.

---

## Error Handling

- If the backend returns an error (e.g., media cannot be re-allocated), the UI displays a detailed, user-friendly message based on the error type.
- The dialog remains open on error, allowing the user to retry or cancel.

---

## Architecture & API

- **Frontend**
  - Add a delete album dialog component.
  - Integrate with the album list view and ensure only owned albums are selectable.
  - Use the existing API infrastructure to call the backend endpoint.
  - Handle loading, success, and error states as described.

- **Backend**
  - Use the existing album deletion logic.
  - Ensure error responses include a string error type and a detailed message.

---

## Testing Plan

**Unit Tests**
- Dialog opens with the correct album pre-selected.
- Only owned albums appear in the dropdown.
- Album name and media count are displayed in the dropdown.
- Loading indicator appears during deletion and album list refresh.
- Dialog closes automatically on success.
- Error messages are displayed on failure.

**Integration Tests**
- Deleting an album updates the album list and redirects if necessary.
- Attempting to delete a non-owned album is not possible via the UI.
- Backend error codes are correctly mapped to user-friendly messages.

**Manual/UX Tests**
- Test on desktop, tablet, and mobile for usability.
- Confirm that the dialog cannot be opened outside the album list view.
- Confirm that the dialog always defaults to the currently selected album.

---

## Summary Table

| Requirement                | Details                                                                 |
|----------------------------|-------------------------------------------------------------------------|
| Entry Point                | Album list view (left menu)                                             |
| Dialog Pre-selection       | Currently selected album                                                |
| Album Selection            | Dropdown with only owned albums, shows name & media count               |
| Confirmation               | Two-step: Delete, then Confirm                                          |
| Loading Feedback           | Shown during deletion and album list refresh                            |
| Success Behavior           | Dialog closes, album list refreshes, redirect if needed                 |
| Error Handling             | Detailed, user-friendly error messages                                  |
| API Endpoint               | `DELETE /api/v1/owners/{owner}/albums/{folderName}`                     |
| Permissions                | Only owned albums, standard authentication                              |
| Undo/Trash                 | Not required                                                            |
| Logging/Audit              | Not required                                                            |
| Device Support             | All (desktop, tablet, mobile)                                           |
| Testing                    | Unit, integration, and manual/UX as described above                     |

---

This specification should provide all the details a developer needs to implement the delete album feature as discussed. If you need wireframes or further breakdowns, let me know!
