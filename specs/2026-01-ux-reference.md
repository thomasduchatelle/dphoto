# DPhoto - Frontend Developer API Reference

## Overview

This document provides the TypeScript API reference for frontend developers implementing UI components for DPhoto. The application follows a strict architecture where:

- **UI Components** are purely presentational - they receive data via props and trigger callbacks
- **Thunks** contain all business logic and are passed as callbacks to components
- **Selectors** transform state into the data structures components need
- **State** is the single source of truth, using ubiquitous language from the domain

## Architecture Pattern

```typescript
// Component receives data from selectors and callbacks from thunks
<Component 
  data={selector(state)}
  onAction={thunk}
/>

// Thunk executes business logic and dispatches actions
const thunk = async (dispatch, port, args) => {
  dispatch(actionStarted());
  const result = await port.doOperation(args);
  dispatch(actionCompleted(result));
};

// Action mutates state
const action = (payload) => ({
  type: 'eventOccurred',
  payload,
  reduce: (state) => ({ ...state, ...changes })
});

// Selector transforms state for component consumption
const selector = (state) => ({
  displayData: transform(state.rawData),
  isReady: state.loaded
});
```

## Core Data Types

### Album

```typescript
interface AlbumId {
  owner: string;           // Owner identifier
  folderName: string;      // Album folder name (unique per owner)
}

interface Album {
  albumId: AlbumId;
  name: string;                    // Display name
  start: Date;                     // Start date of album range
  end: Date;                       // End date of album range
  totalCount: number;              // Total number of media items
  temperature: number;             // Media per day metric
  relativeTemperature: number;     // Normalized temperature for visual comparison
  ownedBy?: OwnerDetails;          // Owner info (present when not owned by current user)
  sharedWith: Sharing[];           // Users this album is shared with
}
```

### Media

```typescript
enum MediaType {
  IMAGE,
  VIDEO,
  OTHER
}

type MediaId = string;

interface Media {
  id: MediaId;
  type: MediaType;
  time: Date;                      // Capture date/time
  uiRelativePath: string;          // UI internal link (from album)
  contentPath: string;             // Path to media content
  source: string;                  // Original source location
}

interface MediaWithinADay {
  day: Date;                       // Date grouping
  medias: Media[];                 // Media captured on this day
}
```

### User and Ownership

```typescript
interface CurrentUserInsight {
  picture?: string;                // User's profile picture URL
  isOwner: boolean;                // Whether user can create albums
}

interface UserDetails {
  name: string;
  email: string;
  picture?: string;
}

interface OwnerDetails {
  name: string;
  users: UserDetails[];            // Users associated with this owner
}

interface Sharing {
  user: UserDetails;
}
```

### Album Filtering

```typescript
interface AlbumFilterCriterion {
  owners: string[];                // Owner IDs to filter by (empty = all accessible albums)
  selfOwned?: boolean;             // Filter to only current user's albums
}

interface AlbumFilterEntry {
  criterion: AlbumFilterCriterion;
  avatars: string[];               // Avatar URLs for display
  name: string;                    // Display name for filter option
}
```

### Dialogs

```typescript
interface DateRangeState {
  startDate: Date | null;
  endDate: Date | null;
  startAtDayStart: boolean;        // Include full day start
  endAtDayEnd: boolean;            // Include full day end
}

interface NameEditBase {
  originalFolderName?: string;     // Used to pre-fill custom folder name
  albumName: string;
  customFolderName: string;
  isCustomFolderNameEnabled: boolean;
  nameError: {
    nameError?: string;
    folderNameError?: string;
  };
}

interface CreateDialog extends DateRangeState, NameEditBase {
  type: "CreateDialog";
  isLoading: boolean;
  error?: string;
}

interface EditDatesDialog extends DateRangeState {
  type: "EditDatesDialog";
  albumId: AlbumId;
  albumName: string;
  isLoading: boolean;
  error?: string;
}

interface EditNameDialog extends NameEditBase {
  type: "EditNameDialog";
  albumId: AlbumId;
  technicalError?: string;
  isLoading: boolean;
}

interface DeleteDialog {
  type: "DeleteDialog";
  deletableAlbums: Album[];        // Only albums owned by current user
  initialSelectedAlbumId?: AlbumId;
  isLoading: boolean;
  error?: string;
}

interface ShareDialog {
  type: "ShareDialog";
  sharedAlbumId: AlbumId;
  sharedWith: Sharing[];           // Current users with access
  suggestions: UserDetails[];      // Suggested users to share with
  error?: {
    type: "grant" | "revoke";
    message: string;
    email: string;
  };
}

type CatalogDialog = CreateDialog | EditDatesDialog | DeleteDialog | ShareDialog | EditNameDialog;
```

## Component Props Interfaces

### Main Page

```typescript
interface CatalogViewerPageSelection {
  albumsLoaded: boolean;
  albums: Album[];                 // Filtered album list
  displayedAlbum: Album | undefined;
  medias: MediaWithinADay[];
  mediasLoaded: boolean;
  albumNotFound: boolean;
  error?: Error;
}
```

### Album List Actions

```typescript
interface AlbumListActionsProps {
  albumFilter: AlbumFilterEntry;
  albumFilterOptions: AlbumFilterEntry[];
  displayedAlbumIdIsOwned: boolean;
  hasAlbumsToDelete: boolean;      // Any owned albums exist
  canCreateAlbums: boolean;        // Current user is an owner
}

interface AlbumListActionsCallbacks {
  onAlbumFiltered: (criterion: AlbumFilterCriterion) => void;
  openCreateDialog: () => void;
  openDeleteAlbumDialog: () => void;
  openEditDatesDialog: () => void;
  openEditNameDialog: () => void;
}
```

## Thunks Reference

Thunks are the callbacks provided to UI components. They encapsulate all business logic.

### Navigation Thunks

#### onPageRefresh

Loads albums and/or media when the page loads or refreshes.

```typescript
function onPageRefresh(albumId?: AlbumId): Promise<void>
```

**Usage**: Call on page mount or when navigating to a specific album.

**Behavior**:
- If no albumId: loads default album (most recent)
- If albumId provided: loads that specific album
- Handles loading states and errors
- Updates state with albums and media

#### onAlbumFilterChange

Filters the album list by owner criterion.

```typescript
function onAlbumFilterChange(criterion: AlbumFilterCriterion): void
```

**Usage**: Call when user selects a different filter option.

**Behavior**:
- Filters albums based on criterion
- If current album doesn't match filter, selects first matching album
- Updates filtered album list in state

### Album Creation Thunks

#### openCreateDialog

Opens the create album dialog.

```typescript
function openCreateDialog(): void
```

**Usage**: Call when user clicks the "Create Album" button.

**Behavior**:
- Opens dialog with empty form
- Initializes date range to null
- Sets loading state to false

#### closeCreateDialog

Closes the create album dialog.

```typescript
function closeCreateDialog(): void
```

**Usage**: Call when user cancels or completes creation.

#### submitCreateAlbum

Submits the create album form.

```typescript
function submitCreateAlbum(): Promise<void>
```

**Usage**: Call when user submits the create form.

**Behavior**:
- Validates dates (both must be provided, end after start)
- Creates album via API
- Refreshes album list
- Selects newly created album
- Closes dialog on success
- Shows error message on failure

### Album Edit Thunks

#### openEditDatesDialog

Opens the edit dates dialog for the current album.

```typescript
function openEditDatesDialog(): void
```

**Usage**: Call when user selects "Edit Dates" for an album they own.

**Behavior**:
- Populates dialog with current album dates
- Only works for albums owned by current user

#### closeEditDatesDialog

Closes the edit dates dialog.

```typescript
function closeEditDatesDialog(): void
```

#### updateAlbumDates

Saves the modified album date range.

```typescript
function updateAlbumDates(): Promise<void>
```

**Usage**: Call when user submits the edit dates form.

**Behavior**:
- Validates date range
- Updates album dates via API
- Re-indexes media for the album
- Refreshes media list
- Shows specific error if edit would orphan media
- Closes dialog on success

#### openEditNameDialog

Opens the edit name dialog for the current album.

```typescript
function openEditNameDialog(): void
```

**Usage**: Call when user selects "Edit Name" for an album they own.

**Behavior**:
- Populates dialog with current album name and folder name
- Enables/disables custom folder name based on current settings

#### closeEditNameDialog

Closes the edit name dialog.

```typescript
function closeEditNameDialog(): void
```

#### saveAlbumName

Saves the modified album name and optional folder name.

```typescript
function saveAlbumName(): Promise<void>
```

**Usage**: Call when user submits the edit name form.

**Behavior**:
- Validates album name (required)
- Validates folder name (if custom naming enabled)
- Renames album via API
- May update URL if folder name changed
- Refreshes album list
- Navigates to new album ID if folder name changed
- Shows validation errors or technical errors
- Closes dialog on success

### Album Deletion Thunks

#### openDeleteAlbumDialog

Opens the delete album dialog.

```typescript
function openDeleteAlbumDialog(): void
```

**Usage**: Call when user clicks the "Delete Album" button.

**Behavior**:
- Shows list of deletable albums (owned by current user only)
- Pre-selects current album if owned by user

#### closeDeleteAlbumDialog

Closes the delete album dialog.

```typescript
function closeDeleteAlbumDialog(): void
```

#### deleteAlbum

Deletes the specified album.

```typescript
function deleteAlbum(albumIdToDelete: AlbumId): Promise<void>
```

**Usage**: Call when user confirms deletion.

**Behavior**:
- Deletes album via API
- Refreshes album list
- If deleted album was displayed, loads default album
- Shows error if deletion fails
- Closes dialog on success

### Sharing Thunks

#### openSharingModal

Opens the sharing dialog for an album.

```typescript
function openSharingModal(albumId: AlbumId): void
```

**Usage**: Call when user clicks sharing indicators on an album they own.

**Behavior**:
- Loads current sharing information
- Shows users album is shared with
- Provides suggestions for new users to share with

#### closeSharingModal

Closes the sharing dialog.

```typescript
function closeSharingModal(): void
```

#### grantAlbumAccess

Grants access to an album for a user by email.

```typescript
function grantAlbumAccess(email: string): Promise<void>
```

**Usage**: Call when user enters an email and clicks "Share".

**Behavior**:
- Validates email format
- Grants access via API
- Loads user details (name, picture)
- Updates shared users list in dialog
- Shows error if granting fails
- User remains in shared list even if profile loading fails

#### revokeAlbumAccess

Revokes access to an album for a user.

```typescript
function revokeAlbumAccess(email: string): Promise<void>
```

**Usage**: Call when user clicks "Revoke" for a shared user.

**Behavior**:
- Immediately removes user from shared list
- Revokes access via API
- Shows error if revocation fails (user remains removed in UI)

### Date Range Edit Thunks

These thunks are used within create and edit dialogs.

#### updateDateRangeStartDate

```typescript
function updateDateRangeStartDate(date: Date | null): void
```

**Usage**: Call when user changes the start date picker.

#### updateDateRangeEndDate

```typescript
function updateDateRangeEndDate(date: Date | null): void
```

**Usage**: Call when user changes the end date picker.

#### updateDateRangeStartAtDayStart

```typescript
function updateDateRangeStartAtDayStart(value: boolean): void
```

**Usage**: Call when user toggles "start at day beginning" checkbox.

#### updateDateRangeEndAtDayEnd

```typescript
function updateDateRangeEndAtDayEnd(value: boolean): void
```

**Usage**: Call when user toggles "end at day end" checkbox.

### Name Edit Thunks

These thunks are used within create and edit name dialogs.

#### changeAlbumName

```typescript
function changeAlbumName(name: string): void
```

**Usage**: Call when user types in the album name field.

**Behavior**:
- Updates album name in dialog state
- Triggers validation
- Updates error messages if invalid

#### changeFolderName

```typescript
function changeFolderName(folderName: string): void
```

**Usage**: Call when user types in the custom folder name field.

**Behavior**:
- Updates folder name in dialog state
- Triggers validation
- Updates error messages if invalid

#### changeFolderNameEnabled

```typescript
function changeFolderNameEnabled(enabled: boolean): void
```

**Usage**: Call when user toggles custom folder name checkbox.

**Behavior**:
- Enables/disables custom folder name input
- Clears folder name validation errors

## Selectors Reference

Selectors transform state into component-ready data structures.

### catalogViewerPageSelector

```typescript
function catalogViewerPageSelector(state: CatalogViewerState): CatalogViewerPageSelection
```

**Returns**:
```typescript
{
  albumsLoaded: boolean;
  albums: Album[];           // Filtered list
  displayedAlbum: Album | undefined;
  medias: MediaWithinADay[];
  mediasLoaded: boolean;
  albumNotFound: boolean;
  error?: Error;
}
```

**Usage**: Use for main page component props.

### albumListActionsSelector

```typescript
function albumListActionsSelector(state: CatalogViewerState): AlbumListActionsProps
```

**Returns**:
```typescript
{
  albumFilter: AlbumFilterEntry;
  albumFilterOptions: AlbumFilterEntry[];
  displayedAlbumIdIsOwned: boolean;
  hasAlbumsToDelete: boolean;
  canCreateAlbums: boolean;
}
```

**Usage**: Use for album action buttons component props.

### Dialog Selectors

Each dialog has a selector that returns the dialog state or undefined:

```typescript
function selectCreateDialog(state: CatalogViewerState): CreateDialog | undefined
function selectEditDatesDialog(state: CatalogViewerState): EditDatesDialog | undefined
function selectEditNameDialog(state: CatalogViewerState): EditNameDialog | undefined
function selectDeleteDialog(state: CatalogViewerState): DeleteDialog | undefined
function selectShareDialog(state: CatalogViewerState): ShareDialog | undefined
```

**Usage**: Use to determine if dialog should be open and populate dialog props.

## Port Interfaces (for reference)

Ports are interfaces that abstract external operations (API calls). Frontend developers don't implement these directly, but they help understand what operations are available.

### FetchAlbumsAndMediasPort

```typescript
interface FetchAlbumsAndMediasPort {
  fetchAlbums(): Promise<Album[]>;
  fetchMedias(albumId: AlbumId): Promise<Media[]>;
}
```

### CreateAlbumPort

```typescript
interface CreateAlbumPort {
  createAlbum(request: {
    name: string;
    start: Date;
    end: Date;
    forcedFolderName: string;  // Empty string for auto-generated
  }): Promise<AlbumId>;
  fetchAlbums(): Promise<Album[]>;
}
```

### UpdateAlbumDatesPort

```typescript
interface UpdateAlbumDatesPort {
  updateAlbumDates(albumId: AlbumId, startDate: Date, endDate: Date): Promise<void>;
  fetchAlbums(): Promise<Album[]>;
  fetchMedias(albumId: AlbumId): Promise<Media[]>;
}
```

### SaveAlbumNamePort

```typescript
interface SaveAlbumNamePort {
  renameAlbum(albumId: AlbumId, newName: string, newFolderName?: string): Promise<AlbumId>;
}
```

### DeleteAlbumPort

```typescript
interface DeleteAlbumPort {
  deleteAlbum(albumId: AlbumId): Promise<void>;
  fetchAlbums(): Promise<Album[]>;
  fetchMedias(albumId: AlbumId): Promise<Media[]>;
}
```

### GrantAlbumAccessAPI

```typescript
interface GrantAlbumAccessAPI {
  grantAccessToAlbum(albumId: AlbumId, email: string): Promise<void>;
  loadUserDetails(email: string): Promise<UserDetails>;
}
```

### RevokeAlbumAccessAPI

```typescript
interface RevokeAlbumAccessAPI {
  revokeSharingAlbum(albumId: AlbumId, email: string): Promise<void>;
}
```

## Error Types

```typescript
interface CatalogError {
  code: string;
  message: string;
}

// Common error codes:
// - "AlbumStartAndEndDateMandatoryErr" - Missing required dates
// - "OrphanedMediasErr" - Date edit would leave media without an album
// - Folder name validation errors
// - Album name validation errors
```

## Helper Functions

### Album Helpers

```typescript
function albumIdEquals(a?: AlbumId, b?: AlbumId): boolean
```

Compares two album IDs for equality.

```typescript
function albumIsOwnedByCurrentUser(album: Album): boolean
```

Returns true if album is owned by current user (ownedBy is undefined).

```typescript
function albumMatchCriterion(criterion: AlbumFilterCriterion): (album: Album) => boolean
```

Returns a predicate function to filter albums by criterion.

## Component Implementation Guidelines

### Pure Components

Components should be pure and presentational:

```typescript
// ✓ Good
export function AlbumCard({ 
  album, 
  selected, 
  onClick 
}: {
  album: Album;
  selected: boolean;
  onClick: (albumId: AlbumId) => void;
}) {
  return (
    <Card onClick={() => onClick(album.albumId)}>
      <Typography>{album.name}</Typography>
      <Typography>{album.totalCount} media</Typography>
    </Card>
  );
}

// ✗ Bad - component contains logic
export function AlbumCard({ album }: { album: Album }) {
  const [selected, setSelected] = useState(false);
  const handleClick = async () => {
    await fetchMedias(album.albumId);  // Business logic in component
    setSelected(true);
  };
  // ...
}
```

### Using Thunks

```typescript
// In parent component or page
function AlbumsPage() {
  const state = useCatalogState();
  const thunks = useCatalogThunks();
  
  const pageData = catalogViewerPageSelector(state);
  const actionsData = albumListActionsSelector(state);
  
  return (
    <MediasPage
      {...pageData}
      albumListActionsProps={{
        ...actionsData,
        onAlbumFiltered: thunks.onAlbumFilterChange,
        openCreateDialog: thunks.openCreateDialog,
        openDeleteAlbumDialog: thunks.openDeleteAlbumDialog,
        openEditDatesDialog: thunks.openEditDatesDialog,
        openEditNameDialog: thunks.openEditNameDialog,
      }}
      openSharingModal={thunks.openSharingModal}
    />
  );
}
```

### Dialog Pattern

```typescript
function CreateAlbumDialogContainer() {
  const state = useCatalogState();
  const thunks = useCatalogThunks();
  
  const dialog = selectCreateDialog(state);
  
  if (!dialog) return null;
  
  return (
    <CreateAlbumDialog
      open={true}
      albumName={dialog.albumName}
      startDate={dialog.startDate}
      endDate={dialog.endDate}
      customFolderName={dialog.customFolderName}
      isCustomFolderNameEnabled={dialog.isCustomFolderNameEnabled}
      startAtDayStart={dialog.startAtDayStart}
      endAtDayEnd={dialog.endAtDayEnd}
      isLoading={dialog.isLoading}
      error={dialog.error}
      nameError={dialog.nameError}
      onClose={thunks.closeCreateDialog}
      onSubmit={thunks.submitCreateAlbum}
      onAlbumNameChange={thunks.changeAlbumName}
      onFolderNameChange={thunks.changeFolderName}
      onFolderNameEnabledChange={thunks.changeFolderNameEnabled}
      onStartDateChange={thunks.updateDateRangeStartDate}
      onEndDateChange={thunks.updateDateRangeEndDate}
      onStartAtDayStartChange={thunks.updateDateRangeStartAtDayStart}
      onEndAtDayEndChange={thunks.updateDateRangeEndAtDayEnd}
    />
  );
}
```

## Testing Components

### Visual Testing with Ladle

Components should have visual tests (stories) covering each state:

```typescript
// create-album-dialog.stories.tsx
import { action } from "@ladle/react";
import { CreateAlbumDialog } from "./CreateAlbumDialog";

export const Default = () => (
  <CreateAlbumDialog
    open={true}
    albumName=""
    startDate={null}
    endDate={null}
    customFolderName=""
    isCustomFolderNameEnabled={false}
    startAtDayStart={true}
    endAtDayEnd={true}
    isLoading={false}
    nameError={{}}
    onClose={action("onClose")}
    onSubmit={action("onSubmit")}
    onAlbumNameChange={action("onAlbumNameChange")}
    // ... other callbacks
  />
);

export const Loading = () => (
  <CreateAlbumDialog
    // ... same props
    isLoading={true}
  />
);

export const WithError = () => (
  <CreateAlbumDialog
    // ... same props
    error="Failed to create album"
  />
);

export const WithValidationErrors = () => (
  <CreateAlbumDialog
    // ... same props
    nameError={{
      nameError: "Album name is required",
      folderNameError: "Invalid folder name"
    }}
  />
);
```

## State Management Summary

The state follows this flow:

1. **User Interaction** → Component calls thunk callback
2. **Thunk Execution** → Performs business logic, calls API via port
3. **Action Dispatch** → Thunk dispatches action(s)
4. **State Update** → Action's reduce function updates state
5. **Selector Evaluation** → Selector transforms updated state
6. **Component Re-render** → Component receives new props from selector

This architecture ensures:
- Components remain pure and testable
- Business logic is centralized in thunks
- State mutations are explicit and traceable
- The UI is always in sync with state
