# DPhoto - Frontend Developer API Reference

## Overview

This document provides the technical API reference for frontend developers building UI components for DPhoto. It describes the data structures, selectors, and
thunks (callbacks) available to power the visual components.

All code examples use TypeScript. The architecture follows a unidirectional data flow pattern where:

- **Selectors** provide data to display in components (via props)
- **Thunks** are callback functions that handle user interactions (via props)
- Components remain purely presentational and framework-agnostic

## Core Data Structures

### Media

```typescript
enum MediaType {
    IMAGE,
    VIDEO,
    OTHER
}

interface Media {
    id: string;
    type: MediaType;
    time: Date;
    uiRelativePath: string; // Internal link for navigation
    contentPath: string;    // URL to fetch the actual media file
    source: string;         // Original source location
}
```

### MediaWithinADay

Medias are organized by day for display:

```typescript
interface MediaWithinADay {
    day: Date;              // The date for this group
    medias: Media[];        // All medias captured on this day
}
```

### AlbumId

Composite identifier for albums:

```typescript
interface AlbumId {
    owner: string;          // Owner identifier
    folderName: string;     // Technical folder name
}
```

### Album

```typescript
interface Album {
    albumId: AlbumId;
    name: string;                    // Display name
    start: Date;                     // Start date/time of the album
    end: Date;                       // End date/time of the album
    totalCount: number;              // Total number of medias
    temperature: number;             // Medias per day (absolute)
    relativeTemperature: number;     // 0.0-1.0, relative to user's hottest album
    ownedBy?: OwnerDetails;          // Present if owned by someone else
    sharedWith: Sharing[];           // List of users with access
}
```

**Note**: If `ownedBy` is `undefined`, the album is owned by the current user.

### OwnerDetails

```typescript
interface OwnerDetails {
    name: string;           // Owner display name
    users: UserDetails[];   // Users associated with this owner
}
```

### UserDetails

```typescript
interface UserDetails {
    name: string;
    email: string;
    picture?: string;       // Optional profile picture URL
}
```

### Sharing

```typescript
interface Sharing {
    user: UserDetails;      // User who has access to the album
}
```

### AlbumFilterCriterion

```typescript
interface AlbumFilterCriterion {
    owners: string[];       // Empty with selfOwned=false means all albums
    selfOwned?: boolean;    // true = owned by current user
}
```

### AlbumFilterEntry

```typescript
interface AlbumFilterEntry {
    criterion: AlbumFilterCriterion;
    avatars: string[];      // Avatar URLs to display for this filter
    name: string;           // Display name of the filter
}
```

### CurrentUserInsight

```typescript
interface CurrentUserInsight {
    picture?: string;       // User's profile picture
    isOwner: boolean;       // Whether user can create albums
}
```

## Selectors (Data for Display)

Selectors transform the application state into data structures needed by UI components.

### Catalog Viewer Page

Main page displaying albums and medias:

```typescript
interface CatalogViewerPageSelection {
    albumsLoaded: boolean;              // true when albums have loaded
    albums: Album[];                    // Filtered list of albums to display
    displayedAlbum: Album | undefined;  // Currently selected album
    medias: MediaWithinADay[];          // Medias grouped by day
    mediasLoaded: boolean;              // true when medias have loaded
    albumNotFound: boolean;             // true if requested album doesn't exist
    error?: Error;                      // Technical error if occurred
}

// Usage in component props
interface CatalogViewerPageProps {
    data: CatalogViewerPageSelection;
    // ... callbacks
}
```

### Album List Actions

Data for the album list toolbar (filter, create, delete buttons):

```typescript
interface AlbumListActionsProps {
    albumFilter: AlbumFilterEntry;          // Currently active filter
    albumFilterOptions: AlbumFilterEntry[]; // Available filter options
    displayedAlbumIdIsOwned: boolean;       // true if current album is owned by user
    hasAlbumsToDelete: boolean;             // true if user has albums they can delete
    canCreateAlbums: boolean;               // true if user can create albums
}
```

### Create Album Dialog

```typescript
interface CreateDialogSelection {
    open: boolean;                      // true when dialog is open
    albumName: string;                  // Current album name input
    customFolderName: string;           // Current custom folder name input
    isCustomFolderNameEnabled: boolean; // true if using custom folder name
    start: Date | null;                 // Selected start date
    end: Date | null;                   // Selected end date
    startsAtStartOfTheDay: boolean;     // true if start is at 00:00:00
    endsAtEndOfTheDay: boolean;         // true if end is at 23:59:59
    isLoading: boolean;                 // true when creation is in progress
    error?: string;                     // Error message if creation failed
    canSubmit: boolean;                 // true if form is valid and can be submitted
    dateRangeError?: string;            // Specific date validation error
    nameError?: string;                 // Album name validation error
    folderNameError?: string;           // Folder name validation error
}
```

### Edit Name Dialog

```typescript
interface EditNameDialogSelection {
    isOpen: boolean;                    // true when dialog is open
    albumName: string;                  // Current album name input
    originalName: string;               // Original album name (for display)
    customFolderName: string;           // Current custom folder name input
    isCustomFolderNameEnabled: boolean; // true if using custom folder name
    technicalError?: string;            // Technical error message
    isLoading: boolean;                 // true when save is in progress
    isSaveEnabled: boolean;             // true if form is valid and can be saved
    nameError?: string;                 // Album name validation error
    folderNameError?: string;           // Folder name validation error
}
```

### Edit Dates Dialog

```typescript
interface EditDatesDialogSelection {
    isOpen: boolean;                    // true when dialog is open
    albumName: string;                  // Album name (for display)
    startDate: Date | null;             // Current start date
    endDate: Date | null;               // Current end date
    startAtDayStart: boolean;           // true if start is at 00:00:00
    endAtDayEnd: boolean;               // true if end is at 23:59:59
    isLoading: boolean;                 // true when save is in progress
    errorCode?: string;                 // Error code if save failed
    dateRangeError?: string;            // Date validation error message
    isSaveEnabled: boolean;             // true if dates are valid and can be saved
}
```

### Delete Album Dialog

```typescript
interface DeleteDialogFrag {
    albums: Album[];                    // Albums user can delete
    initialSelectedAlbumId?: AlbumId;   // Pre-selected album (if any)
    isOpen: boolean;                    // true when dialog is open
    isLoading: boolean;                 // true when deletion is in progress
    error?: string;                     // Error message if deletion failed
}
```

### Share Album Dialog

```typescript
interface ShareError {
    type: "grant" | "revoke";           // Which operation failed
    message: string;                    // Error message
    email: string;                      // Email address involved
}

interface SharingDialogFrag {
    open: boolean;                      // true when dialog is open
    sharedWith: Sharing[];              // Users who have access
    suggestions: UserDetails[];         // Suggested users to share with
    error?: ShareError;                 // Error if sharing operation failed
}
```

### Displayed Album

Helper selector to get information about the currently displayed album:

```typescript
interface DisplayedAlbumSelection {
    displayedAlbumId?: AlbumId;         // ID of displayed album
    displayedAlbumIdIsOwned: boolean;   // true if owned by current user
}
```

## Thunks (User Interaction Callbacks)

Thunks are callback functions provided to components to handle user interactions. All thunks are framework-agnostic and can be called directly from event
handlers.

### Navigation Thunks

#### onPageRefresh

Called when the page loads or refreshes to fetch initial data.

```typescript
(albumId?: AlbumId) => Promise<void>
```

**Parameters:**

- `albumId` (optional): If provided, loads this specific album. If omitted, loads the default album.

**Use cases:**

- Initial page load
- User manually refreshes the page
- Navigating to a specific album via URL

**Example:**

```typescript
interface PageProps {
    onLoad: (albumId?: AlbumId) => Promise<void>;
}

// In component
useEffect(() => {
    props.onLoad(albumIdFromUrl);
}, []);
```

#### onAlbumFilterChange

Called when user changes the album filter.

```typescript
(criterion: AlbumFilterCriterion) => void
```

**Parameters:**

- `criterion`: The new filter criterion to apply

**Use cases:**

- User selects "My Albums" filter
- User selects a specific owner's albums
- User selects "All Albums"

**Example:**

```typescript
interface FilterSelectorProps {
    onFilterChange: (criterion: AlbumFilterCriterion) => void;
    currentFilter: AlbumFilterEntry;
    options: AlbumFilterEntry[];
}

// In component
const handleFilterClick = (option: AlbumFilterEntry) => {
    props.onFilterChange(option.criterion);
};
```

### Album Creation Thunks

#### openCreateDialog

Opens the create album dialog.

```typescript
() => void
```

**Example:**

```typescript
interface CreateButtonProps {
    onClick: () => void;
    disabled?: boolean;
}
```

#### closeCreateDialog

Closes the create album dialog without creating.

```typescript
() => void
```

#### changeAlbumName

Updates the album name input in the create dialog.

```typescript
(albumName: string) => void
```

**Example:**

```typescript
interface AlbumNameInputProps {
    value: string;
    onChange: (name: string) => void;
    error?: string;
}

// In component
<input
    value = {props.value}
onChange = {(e)
=>
props.onChange(e.target.value)
}
/>
```

#### changeFolderName

Updates the custom folder name input.

```typescript
(folderName: string) => void
```

#### changeFolderNameEnabled

Toggles whether custom folder name is enabled.

```typescript
(isEnabled: boolean) => void
```

#### updateCreateDialogStartDate

Updates the start date in create dialog.

```typescript
(date: Date | null) => void
```

#### updateCreateDialogEndDate

Updates the end date in create dialog.

```typescript
(date: Date | null) => void
```

#### updateCreateDialogStartsAtStartOfTheDay

Toggles whether start time is at beginning of day.

```typescript
(startsAtStart: boolean) => void
```

#### updateCreateDialogEndsAtEndOfTheDay

Toggles whether end time is at end of day.

```typescript
(endsAtEnd: boolean) => void
```

#### submitCreateAlbum

Creates the album with current dialog values.

```typescript
() => Promise<void>
```

**Example:**

```typescript
interface CreateDialogProps {
    data: CreateDialogSelection;
    onSubmit: () => Promise<void>;
    onClose: () => void;
    onNameChange: (name: string) => void;
    onStartDateChange: (date: Date | null) => void;
    onEndDateChange: (date: Date | null) => void;
    // ... other callbacks
}

// In component
const handleSubmit = async () => {
    try {
        await props.onSubmit();
        // Dialog automatically closes on success
    } catch (error) {
        // Error displayed via data.error
    }
};
```

### Album Editing Thunks

#### openEditNameDialog

Opens the edit name dialog.

```typescript
() => void
```

#### closeEditNameDialog

Closes the edit name dialog.

```typescript
() => void
```

#### saveAlbumName

Saves the modified album name.

```typescript
() => Promise<void>
```

**Example:**

```typescript
interface EditNameDialogProps {
    data: EditNameDialogSelection;
    onSave: () => Promise<void>;
    onClose: () => void;
    onNameChange: (name: string) => void;
    onFolderNameChange: (name: string) => void;
    onFolderNameEnabledChange: (enabled: boolean) => void;
}
```

#### openEditDatesDialog

Opens the edit dates dialog.

```typescript
() => void
```

#### closeEditDatesDialog

Closes the edit dates dialog.

```typescript
() => void
```

#### updateEditDatesDialogStartDate

Updates the start date in edit dialog.

```typescript
(startDate: Date | null) => void
```

#### updateEditDatesDialogEndDate

Updates the end date in edit dialog.

```typescript
(endDate: Date | null) => void
```

#### updateEditDatesDialogStartAtDayStart

Toggles whether start is at beginning of day.

```typescript
(startAtDayStart: boolean) => void
```

#### updateEditDatesDialogEndAtDayEnd

Toggles whether end is at end of day.

```typescript
(endAtDayEnd: boolean) => void
```

#### updateAlbumDates

Saves the modified dates.

```typescript
() => Promise<void>
```

**Example:**

```typescript
interface EditDatesDialogProps {
    data: EditDatesDialogSelection;
    onSave: () => Promise<void>;
    onClose: () => void;
    onStartDateChange: (date: Date | null) => void;
    onEndDateChange: (date: Date | null) => void;
    onStartAtDayStartChange: (atStart: boolean) => void;
    onEndAtDayEndChange: (atEnd: boolean) => void;
}
```

### Album Deletion Thunks

#### openDeleteAlbumDialog

Opens the delete album dialog.

```typescript
() => void
```

#### closeDeleteAlbumDialog

Closes the delete album dialog.

```typescript
() => void
```

#### deleteAlbum

Deletes the specified album.

```typescript
(albumIdToDelete: AlbumId) => Promise<void>
```

**Parameters:**

- `albumIdToDelete`: The ID of the album to delete

**Example:**

```typescript
interface DeleteDialogProps {
    data: DeleteDialogFrag;
    onDelete: (albumId: AlbumId) => Promise<void>;
    onClose: () => void;
}

// In component
const handleDelete = async () => {
    if (selectedAlbum) {
        await props.onDelete(selectedAlbum.albumId);
        // Dialog automatically closes on success
    }
};
```

### Sharing Thunks

#### openSharingModal

Opens the sharing dialog for a specific album.

```typescript
(album: Album) => void
```

**Parameters:**

- `album`: The album to manage sharing for

**Example:**

```typescript
interface AlbumActionsProps {
    album: Album;
    onShare: (album: Album) => void;
}

// In component
<button onClick = {()
=>
props.onShare(props.album)
}>
Share
< /button>
```

#### closeSharingModal

Closes the sharing dialog.

```typescript
() => void
```

#### grantAlbumAccess

Grants access to a user by email.

```typescript
(email: string) => Promise<void>
```

**Parameters:**

- `email`: Email address of user to grant access to

**Example:**

```typescript
interface ShareDialogProps {
    data: SharingDialogFrag;
    onGrantAccess: (email: string) => Promise<void>;
    onRevokeAccess: (email: string) => Promise<void>;
    onClose: () => void;
}

// In component
const handleShare = async (email: string) => {
    try {
        await props.onGrantAccess(email);
        // Success - email added to sharedWith list
    } catch (error) {
        // Error displayed via data.error
    }
};
```

#### revokeAlbumAccess

Revokes access from a user.

```typescript
(email: string) => Promise<void>
```

**Parameters:**

- `email`: Email address of user to revoke access from

**Example:**

```typescript
// In component
const handleRevoke = async (sharing: Sharing) => {
    await props.onRevokeAccess(sharing.user.email);
    // User removed from sharedWith list
};
```

## Common Component Patterns

### Album List Component

```typescript
interface AlbumsListProps {
    albums: Album[];
    loaded: boolean;
    selectedAlbumId?: AlbumId;
    openSharingModal: (albumId: AlbumId) => void;
    onAlbumClick: (albumId: AlbumId) => void;
}
```

### Media List Component

```typescript
interface MediasListProps {
    medias: MediaWithinADay[];
    loaded: boolean;
    displayedAlbum?: Album;
}
```

### Date Range Picker Component

```typescript
interface DateRangePickerProps {
    startDate: Date | null;
    endDate: Date | null;
    startAtDayStart: boolean;
    endAtDayEnd: boolean;
    onStartDateChange: (date: Date | null) => void;
    onEndDateChange: (date: Date | null) => void;
    onStartAtDayStartChange: (atStart: boolean) => void;
    onEndAtDayEndChange: (atEnd: boolean) => void;
    dateRangeError?: string;
}
```

### Folder Name Input Component

```typescript
interface FolderNameInputProps {
    folderName: string;
    isEnabled: boolean;
    onFolderNameChange: (name: string) => void;
    onEnabledChange: (enabled: boolean) => void;
    error?: string;
}
```

## Error Handling

### CatalogError

Typed error from API operations:

```typescript
class CatalogError extends Error {
    code: string;           // Error code for categorization
    message: string;        // Human-readable message
}

function isCatalogError(err: any): err is CatalogError;
```

**Common error codes:**

- `"AlbumStartAndEndDateMandatoryErr"` - Dates are required
- `"OrphanedMediasErr"` - Date change would orphan medias

### Error Display Pattern

```typescript
interface DialogWithErrorProps {
    error?: string;         // Display generic error
    nameError?: string;     // Display field-specific error
    // ... other error fields
}
```

## Helper Functions

### Album Ownership Check

```typescript
function albumIsOwnedByCurrentUser(album: Album): boolean
```

Returns `true` if `album.ownedBy` is `undefined`.

### Album Filter Matching

```typescript
function albumMatchCriterion(
    criterion: AlbumFilterCriterion
): (album: Album) => boolean
```

Returns a predicate function to test if an album matches the criterion.

### Album ID Equality

```typescript
function albumIdEquals(a?: AlbumId, b?: AlbumId): boolean
```

Compares two album IDs for equality. Returns `false` if either is `undefined`.

## Integration Example

Complete example of a catalog viewer page component:

```typescript
interface CatalogViewerProps {
    // Data from selectors
    page: CatalogViewerPageSelection;
    albumActions: AlbumListActionsProps;
    createDialog: CreateDialogSelection;
    editNameDialog: EditNameDialogSelection;
    editDatesDialog: EditDatesDialogSelection;
    deleteDialog: DeleteDialogFrag;
    shareDialog: SharingDialogFrag;

    // Navigation thunks
    onPageLoad: (albumId?: AlbumId) => Promise<void>;
    onAlbumFilterChange: (criterion: AlbumFilterCriterion) => void;

    // Create album thunks
    onOpenCreateDialog: () => void;
    onCloseCreateDialog: () => void;
    onCreateAlbumNameChange: (name: string) => void;
    onCreateFolderNameChange: (name: string) => void;
    onCreateFolderNameEnabledChange: (enabled: boolean) => void;
    onCreateStartDateChange: (date: Date | null) => void;
    onCreateEndDateChange: (date: Date | null) => void;
    onCreateStartAtDayStartChange: (atStart: boolean) => void;
    onCreateEndAtDayEndChange: (atEnd: boolean) => void;
    onSubmitCreateAlbum: () => Promise<void>;

    // Edit name thunks
    onOpenEditNameDialog: () => void;
    onCloseEditNameDialog: () => void;
    onEditAlbumNameChange: (name: string) => void;
    onEditFolderNameChange: (name: string) => void;
    onEditFolderNameEnabledChange: (enabled: boolean) => void;
    onSaveAlbumName: () => Promise<void>;

    // Edit dates thunks
    onOpenEditDatesDialog: () => void;
    onCloseEditDatesDialog: () => void;
    onEditStartDateChange: (date: Date | null) => void;
    onEditEndDateChange: (date: Date | null) => void;
    onEditStartAtDayStartChange: (atStart: boolean) => void;
    onEditEndAtDayEndChange: (atEnd: boolean) => void;
    onUpdateAlbumDates: () => Promise<void>;

    // Delete album thunks
    onOpenDeleteDialog: () => void;
    onCloseDeleteDialog: () => void;
    onDeleteAlbum: (albumId: AlbumId) => Promise<void>;

    // Sharing thunks
    onOpenShareDialog: (album: Album) => void;
    onCloseShareDialog: () => void;
    onGrantAccess: (email: string) => Promise<void>;
    onRevokeAccess: (email: string) => Promise<void>;
}
```

## Notes for Frontend Developers

1. **Pure Components**: All components should be pure (presentational only). Business logic is in thunks and selectors.

2. **TypeScript Types**: Use the exact types defined here. Do not use `any`.

3. **Async Operations**: Thunks that return `Promise<void>` will:
    - Update loading states automatically
    - Update error states if the operation fails
    - Close dialogs and refresh data on success

4. **State Updates**: Never manipulate state directly. Always use the provided thunks.

5. **Validation**: Validation errors are provided via selectors (e.g., `nameError`, `dateRangeError`). Display them near the relevant inputs.

6. **Loading States**: Use `isLoading` and `loaded` flags to show skeletons, spinners, or disabled states.

7. **Optional Values**: Always handle `undefined` values (e.g., `displayedAlbum?: Album`).

8. **Date Handling**: Dates are provided as `Date` objects. Use appropriate date picker components.

9. **Dialog Pattern**: Only one dialog can be open at a time. The `open` or `isOpen` flag controls visibility.

10. **Error Recovery**: Failed operations keep dialogs open with error messages, allowing users to correct and retry.
