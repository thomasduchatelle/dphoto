# DPhoto - Functional Requirements

## Overview

DPhoto is a photo and video management application that allows users to organize their media into albums, share albums with others, and view their collection
through a web interface. The application supports multi-user access with a robust sharing and permission system.

## Core Concepts

### Media

A media is an individual photo or video with the following properties:

- Unique identifier
- Type: image, video, or other
- Capture time (when the photo or video was taken)
- Source location (original file path or location)
- Content path (for accessing the actual media file)

Medias are organized by the day they were captured. When viewing an album, medias are grouped by day and sorted from newest to oldest.

### Albums

An album is a collection of medias grouped by a date range. Each album has:

- A name (user-friendly display name)
- A folder name (technical identifier used in URLs and storage)
- A start date and time
- An end date and time
- An owner (the person or entity who created it)
- A count of total medias in the album
- A temperature metric (number of medias per day, used to indicate activity level)
- A list of users it has been shared with

Albums are automatically populated with all medias whose capture time falls within the album's date range. Albums are sorted by start date, with the most recent
albums appearing first.

The temperature of an album indicates how many photos were taken per day during that period. This helps users identify particularly active periods (vacations,
events, etc.). The temperature is displayed with both an absolute value and a relative indicator compared to the user's most active album.

### Owners

An owner is an entity to which albums and medias belong. An owner can be:

- The current user (shown as "owned by me")
- Another person or family group that has shared albums with the current user

Each owner has:

- A unique identifier
- A name
- A list of associated users (people who can act on behalf of this owner)

### Users

A user is an authenticated person who can access the application. Each user has:

- A name
- An email address
- An optional profile picture

Users can own albums (be owners themselves) or have albums shared with them by others.

### Sharing

Sharing allows album owners to grant view access to other users by email address. When an album is shared:

- The recipient can view all medias in that album
- The recipient can see who else has access to the album
- Only the owner can modify the album (name, dates) or delete it
- Only the owner can share or revoke access to the album

## User Capabilities

### Viewing Albums

Users can see a list of all albums they have access to, which includes:

- Albums they own
- Albums that have been shared with them by others

For each album in the list, users can see:

- The album name
- The date range (start and end dates)
- The total number of medias
- The temperature indicator (how active that period was)
- The owner information (if not owned by the current user)
- Visual indicators of who the album is shared with

Users can filter the album list by:

- Albums owned by the current user only
- Albums owned by a specific other owner
- All albums the user has access to

When a user selects an album, they can view all medias in that album.

### Viewing Medias

When viewing an album's medias, users see:

- All medias in the album organized by day
- The day grouping (each day shows all medias captured on that date)
- Medias within each day sorted from newest to oldest
- The media type (image or video)
- The capture time

Users can click on individual medias to view them in full size or access additional details.

### Creating Albums

Users who are owners can create new albums. To create an album, the user must provide:

- An album name (required)
- A start date (required)
- An end date (required)
- Optional: a custom folder name (if not provided, one is generated from the album name)

For the date range:

- The user can specify whether the start date should begin at the start of the day (00:00:00) or at a specific time
- The user can specify whether the end date should end at the end of the day (23:59:59) or at a specific time

The album name must not be empty. If a custom folder name is provided, it must also not be empty and must follow technical naming constraints.

After creating an album, the user is automatically redirected to view that album's medias.

### Editing Album Names

Users can edit the name of albums they own. When editing an album name, they can change:

- The display name of the album (required)
- The folder name (optional custom identifier)

If the user enables custom folder name editing:

- They must provide a valid folder name
- The folder name must follow technical constraints

If changing the folder name would affect the album's identifier:

- The system handles this change transparently
- The user is redirected to the new album location

Users cannot edit albums they don't own (albums shared with them).

### Editing Album Dates

Users can modify the date range of albums they own. When editing dates, they can change:

- The start date
- The end date
- Whether the start is at the beginning of the day
- Whether the end is at the end of the day

Date validation ensures:

- Both start and end dates must be provided
- The end date must be after the start date
- The date range must be valid

Changing an album's dates will automatically reorganize which medias appear in that album. If the date change would result in medias becoming "orphaned" (not
belonging to any album), the system prevents the change and notifies the user.

After successfully updating dates, the album is refreshed to show the updated media collection.

### Deleting Albums

Users can delete albums they own. When deleting an album:

- The user selects which album to delete from a list of their owned albums
- The system confirms the deletion
- The album and all its metadata are removed
- The actual media files are not deleted (they remain in storage)

Users cannot delete albums they don't own.

After deletion, if the deleted album was currently being viewed, the user is redirected to another available album or to the album list.

### Sharing Albums

Users can share albums they own with other users. The sharing process:

To share an album:

- The user opens the sharing interface for a specific album
- The user sees who the album is currently shared with
- The user sees a list of suggested users (people they've shared with before)
- The user enters an email address of someone to share with
- The system grants access to that user

The shared user receives access immediately and can view the album in their album list.

To revoke access:

- The user views the list of people the album is shared with
- The user selects a person to revoke access from
- The system removes that person's access
- The person can no longer see the album

Users can only manage sharing for albums they own. For albums shared with them, they can see who has access but cannot modify the sharing settings.

### Filtering and Navigation

Users can filter their album view by owner:

- View only albums they own
- View albums from a specific other owner
- View all albums they have access to

The filter affects which albums appear in the album list. When changing filters:

- If the currently viewed album matches the new filter, it remains selected
- If the currently viewed album doesn't match the new filter, the system selects the first album that does match (if any)

Users can refresh the page at any time, and the application will:

- Load the list of available albums
- Load medias for the currently selected album
- Preserve the user's navigation context

### Loading States and Errors

Throughout the application, users experience:

Loading states when:

- Albums are being fetched from the server
- Medias are being loaded for an album
- An operation is in progress (creating, deleting, sharing, etc.)

Error states when:

- An album cannot be found
- Medias fail to load
- A sharing operation fails (invalid email, network error, etc.)
- An album operation fails (deletion, rename, date update, etc.)
- Date validation fails
- Name validation fails

Each error provides appropriate feedback to help users understand what went wrong and how to proceed.

### User Interface Elements

The application uses dialogs (modal windows) for operations:

- Create Album dialog: for creating a new album with name and date range
- Edit Name dialog: for changing an album's name and folder name
- Edit Dates dialog: for modifying an album's date range
- Delete Album dialog: for confirming album deletion
- Share dialog: for managing who has access to an album

Only one dialog can be open at a time. Users can close dialogs to return to the main view.

The application provides appropriate action buttons based on context:

- Create button: visible only for users who are owners
- Delete button: visible only when the user owns at least one album
- Edit buttons: visible only for albums the user owns
- Share button: visible for albums, with different behaviors for owned vs. shared albums

### User Profile

Users can see their own profile information:

- Profile picture (if available)
- Whether they are an owner (can create albums)

This information helps users understand their permissions and capabilities within the application.

## Business Rules

### Ownership Rules

- Only owners can create albums
- Only the album owner can modify or delete an album
- Only the album owner can share or revoke access to an album

### Date Range Rules

- Albums must have both start and end dates
- End date must be after start date
- Dates can optionally be set to start/end of day
- Changing dates cannot orphan medias (leave medias not belonging to any album)

### Naming Rules

- Album names cannot be empty
- Custom folder names, if enabled, cannot be empty
- Folder names must meet technical constraints for URL and file system compatibility

### Sharing Rules

- Sharing is done by email address
- Users can share with multiple people
- Each person with access can be individually revoked
- Shared users can view but not modify albums

### Filtering Rules

- Filters affect the album list display
- Filtering does not affect album ownership or sharing
- When no albums match a filter, the user sees an empty list

### Navigation Rules

- Selecting an album loads its medias
- After creating an album, the user views that album
- After deleting the current album, the user is redirected to another album
- Page refresh preserves context when possible

## Data Organization

### Album Filtering

The application provides predefined filter options:

- "My Albums" - shows only albums owned by the current user
- Owner-specific filters - one for each owner who has shared albums with the user
- "All Albums" - shows all accessible albums

Each filter option displays:

- A name describing what it filters
- Avatar images representing the owner(s) in that filter

### Media Grouping

Medias are always displayed grouped by day:

- Each group shows the date
- Within each date, medias are sorted newest to oldest
- Days are presented in reverse chronological order (most recent first)

This grouping helps users find medias from specific days and understand the temporal distribution of their photos.

### Album Temperature

The temperature metric helps users identify significant periods:

- High temperature = many photos per day (active period)
- Low temperature = few photos per day (quiet period)
- Relative temperature shows comparison to the user's most active album

This visualization helps users quickly identify vacation albums, event albums, or daily life periods.

## System Behavior

### Automatic Refresh After Changes

After operations that modify data, the system automatically refreshes:

- After creating an album: fetch updated album list and load new album's medias
- After deleting an album: fetch updated album list and load appropriate album
- After editing album dates: fetch updated album list and reload current album's medias
- After sharing/revoking access: update the sharing information display

This ensures users always see current, accurate data.

### Error Recovery

When operations fail:

- User-initiated operations remain interruptible (dialogs stay open)
- Error messages provide specific information about what failed
- Users can retry operations after addressing errors
- The application state remains consistent even after errors

### Performance Considerations

The application optimizes loading:

- Albums are fetched once and cached during a session
- Medias are only loaded for the currently selected album
- Switching between already-loaded albums can use cached data
- Filter changes don't require reloading data, just reorganizing the view

## Not Included (Out of Scope)

This document describes viewing and organizing existing medias. The following are handled by other parts of the system:

- Uploading new medias (handled by the CLI application)
- Processing and compressing medias (handled by the backend archive system)
- Authenticating users (handled by the security layer)
- Storing and retrieving media files (handled by the backend storage system)
- Managing owner groups and relationships (handled by the access control system)
