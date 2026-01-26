# DPhoto - Functional Specification for UX Design

## Overview

DPhoto is a personal photo and video management application that allows users to organize their media into albums, view them through a web interface, and share albums with other users.

## Core Concepts

### Users and Ownership

- **User**: A person who authenticates to access the web interface using their Google account
- **Owner**: A role representing the entity that owns media and albums. When a user creates content, they become an owner
- **Current User**: The authenticated user viewing the application, identified by their profile picture

### Media

- **Media**: Individual photos or videos in the system
- Each media has:
    - A capture date/time
    - A type (image, video, or other)
    - A source location
    - A content path for display
- Media are always displayed grouped by day of capture

### Albums

- **Album**: A collection of media grouped by a date range
- Each album has:
    - A name (e.g., "Summer Vacation 2023")
    - A start date and end date (defines which media belong to it)
    - A folder name (technical identifier, can be custom or auto-generated)
    - An owner (who created/owns it)
    - A total count of media
    - A "temperature" metric (media per day) indicating activity level
    - A relative temperature for visual comparison with other albums
- Albums can be:
    - Owned by the current user
    - Shared with the current user by others
- Albums are sorted by date (most recent first)

### Sharing

- **Sharing**: Albums can be shared with other users via their email address
- When shared:
    - The recipient gains read-only access to view the album and its media
    - The recipient can see who owns the album (name and picture)
    - The owner can see who the album is shared with
    - The owner can revoke access at any time

## User Capabilities

### 1. Browsing and Navigation

#### View Albums List

- Users can see a list of all albums they have access to (owned or shared)
- For each album in the list, users can see:
    - Album name
    - Date range (start to end)
    - Number of media items
    - Visual indicator of album "temperature" (how active/dense the album is)
    - Owner information (if not owned by current user)
    - Who the album is shared with (avatars of shared users)
- Albums are displayed in chronological order
- The list shows loading skeletons while data is being fetched

#### Filter Albums

- Users can filter the album list by owner:
    - "My Albums" - shows only albums owned by the current user
    - "All Albums" - shows all accessible albums (owned + shared)
    - Individual owner filters - shows albums from a specific owner
- Filtering updates the album list immediately
- The current filter selection is visually indicated

#### Select and View an Album

- Users can click on an album to view its media
- When an album is selected:
    - The album is highlighted in the list
    - The media from that album are displayed in the main area
    - If the album changes during filtering, the first matching album is automatically selected

#### View Media

- Media are displayed in a grid layout
- Media are grouped by day with date headers
- For each media item, users can see:
    - The thumbnail (image or video preview)
    - Clicking on a media opens it in full view
- The view shows loading indicators while media are being fetched
- If an album is not found, an appropriate message is displayed

### 2. Album Management (Owners Only)

Users who are owners can create and manage their own albums.

#### Create Album

- Users can initiate album creation
- During creation, users provide:
    - Album name (required)
    - Start date (required) with option to include/exclude the day start time
    - End date (required) with option to include/exclude the day end time
    - Optional: Custom folder name (technical identifier)
        - Can be enabled/disabled
        - When disabled, folder name is auto-generated from the album name
        - When enabled, user provides a custom folder name
- Date validation:
    - Both start and end dates must be provided
    - End date must be after start date
    - Users can specify if dates should include full days or specific times
- After creation:
    - The album list is refreshed
    - The newly created album is automatically selected and displayed
- If creation fails, an error message is shown
- The creation dialog shows a loading state while processing

#### Edit Album Dates

- Users can edit the date range of their own albums
- Users can modify:
    - Start date and whether it includes the full day start
    - End date and whether it includes the full day end
- After saving:
    - The album's media are re-indexed based on the new dates
    - The media list is refreshed to reflect the new content
    - Media that no longer fit the date range may be removed from the album
- If editing would orphan media (create media with no album), a specific error is shown
- The dialog shows a loading state while processing
- Users can cancel the operation

#### Edit Album Name

- Users can edit the name of their own albums
- Users can modify:
    - Album display name
    - Custom folder name (optional)
        - Can enable/disable custom folder name
        - When enabled, provide a new folder name
- Name validation:
    - Album name must not be empty
    - Folder name must be valid (when custom naming is enabled)
- After saving:
    - The album is renamed
    - If the folder name changed, the URL may update
    - The album list is refreshed
- If renaming fails, appropriate error messages are shown
- The dialog shows a loading state while processing

#### Delete Album

- Users can delete their own albums
- During deletion:
    - Users select which album to delete from a dropdown list
    - Only albums owned by the current user are available for deletion
    - A confirmation is required before deletion
- After deletion:
    - The album is permanently removed
    - The album list is refreshed
    - If the deleted album was currently displayed, another album is selected
- If deletion fails, an error message is shown
- The dialog shows a loading state while processing

### 3. Sharing Management (Owners Only)

#### View Album Sharing

- Users can see who their albums are shared with
- The sharing information is visible in the album list (avatars)
- Users can click on sharing indicators to open the sharing management dialog

#### Open Sharing Dialog

- Users can open a sharing dialog for albums they own
- The dialog displays:
    - List of users the album is currently shared with (email, name, picture)
    - Suggestions for users to share with
    - An input field to enter an email address

#### Grant Access

- Users can share an album by entering an email address
- The system:
    - Validates the email
    - Grants read access to the specified user
    - Loads the user's profile information (name, picture)
    - Updates the shared users list in the dialog
    - Shows the new user in the album list's shared indicators
- If granting access fails:
    - An error message is shown
    - The email address is indicated in the error
    - Users can retry

#### Revoke Access

- Users can revoke sharing access from users
- After revoking:
    - The user is immediately removed from the shared users list
    - The change is reflected in the album list
- If revoking fails:
    - An error message is shown
    - The user is informed to try again

### 4. Error Handling

Users are informed of errors through:

- **Album Not Found**: When navigating to a non-existent album
- **Loading Errors**: When albums or media fail to load
- **Creation/Edit Errors**:
    - Invalid date ranges
    - Orphaned media errors when editing dates
    - Invalid names or folder names
    - Technical errors with descriptive messages
- **Sharing Errors**:
    - Failed to grant access (e.g., invalid email)
    - Failed to revoke access
    - Network or permission errors
- **Deletion Errors**: When an album cannot be deleted

### 5. User Interface States

Throughout the application, users see:

- **Loading States**: Skeleton screens or loading indicators while data is fetched
- **Empty States**: Informative messages when no albums exist
- **Disabled States**: Buttons and actions are disabled when not applicable (e.g., edit buttons disabled for albums not owned by the user)
- **Active States**: Visual feedback for selected items, active filters, and current operations

## Access Control Summary

| Capability | Owner of Album | Viewer (Shared Access) |
|-----------|----------------|------------------------|
| View album list | ✓ | ✓ |
| View media | ✓ | ✓ |
| Filter albums | ✓ | ✓ |
| Create album | ✓ | ✗ |
| Edit album dates | ✓ (own albums) | ✗ |
| Edit album name | ✓ (own albums) | ✗ |
| Delete album | ✓ (own albums) | ✗ |
| Share album | ✓ (own albums) | ✗ |
| Revoke sharing | ✓ (own albums) | ✗ |

## Key User Flows

### First-Time User Flow
1. User logs in with Google account
2. If no albums exist, user sees an empty state message
3. User is informed to use the command-line interface to upload photos
4. Once photos are uploaded via CLI, albums appear automatically

### Album Creation Flow
1. User clicks "Create Album" button
2. User enters album name
3. User selects start and end dates
4. User optionally enables custom folder name and provides it
5. User submits the form
6. System creates the album
7. Album list refreshes and new album is selected

### Album Viewing Flow
1. User sees list of albums
2. User optionally filters by owner
3. User clicks on an album
4. Media from that album loads and displays in day-grouped format
5. User can click individual media to view full-size

### Sharing Flow
1. User clicks on sharing indicator for their album
2. Sharing dialog opens
3. User enters email address of person to share with
4. System grants access and displays the new user
5. Shared user can now see the album in their album list

### Album Editing Flow
1. User selects an album they own
2. User clicks edit button (dates or name)
3. Dialog opens with current values
4. User modifies values
5. User saves changes
6. System updates the album
7. Album list and media view refresh with updated information
