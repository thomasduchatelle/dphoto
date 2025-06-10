# Detailed Catalog Actions Refactoring Plan

## Refined Step-by-Step Implementation

### Step 1: Fix `delete-startDeleteAlbum.ts` (Highest Priority - Simple Case)

Fix the action structure for `delete-startDeleteAlbum.ts` to comply with coding principles.

**Files to modify:**
- `web/src/core/catalog/domain/actions/delete-startDeleteAlbum.ts`
- `web/src/core/catalog/domain/actions/delete-startDeleteAlbum.test.ts`

**Changes required:**

1. **Rename action interface** from `StartDeleteAlbumAction` to `DeleteAlbumStarted` (past tense event naming)
2. **Update type property** from `"StartDeleteAlbum"` to `"deleteAlbumStarted"`
3. **Rename factory function** from `startDeleteAlbumAction()` to `deleteAlbumStartedAction()`
4. **Rename reducer function** from `reduceStartDeleteAlbum()` to `reduceDeleteAlbumStarted()`
5. **Rename registration function** from `startDeleteAlbumReducerRegistration()` to `deleteAlbumStartedReducerRegistration()`
6. **Update handler key** in registration to `"deleteAlbumStarted"`

**Pattern to follow (Case 1: No additional properties):**
```typescript
export interface DeleteAlbumStarted {
    type: "deleteAlbumStarted";
}

export function deleteAlbumStartedAction(): DeleteAlbumStarted {
    return {
        type: "deleteAlbumStarted",
    };
}

export function reduceDeleteAlbumStarted(
    current: CatalogViewerState,
    action: DeleteAlbumStarted
): CatalogViewerState {
    // existing implementation
}

export function deleteAlbumStartedReducerRegistration(handlers: any) {
    handlers["deleteAlbumStarted"] = reduceDeleteAlbumStarted as (
        state: CatalogViewerState,
        action: DeleteAlbumStarted
    ) => CatalogViewerState;
}
```

**Test file updates:**
1. **Update describe block** to `describe("action:deleteAlbumStarted", () => {`
2. **Update all function calls** to use new names
3. **Update interface references** to `DeleteAlbumStarted`
4. **Ensure complete state assertions** (already compliant)

SUPER IMPORTANT:

* Semantic integrity must be the **priority** when tests are edited. Never prioritise the syntax correctness at the cost of the initial intention of the test.
* You must update all reference to a name when you rename it (functions, variables, ...) or when its signature change (function, constructor)
* All tests must continue to pass after these changes.

### Step 2: Fix `sharing-closeSharingModalAction.ts`

Fix the action structure for `sharing-closeSharingModalAction.ts` to comply with coding principles.

**Files to modify:**
- `web/src/core/catalog/domain/actions/sharing-closeSharingModalAction.ts`
- `web/src/core/catalog/domain/actions/sharing-closeSharingModalAction.test.ts`

**Changes required:**

1. **Rename action interface** from `CloseSharingModalAction` to `SharingModalClosed`
2. **Update type property** from `"CloseSharingModalAction"` to `"sharingModalClosed"`
3. **Rename factory function** from `closeSharingModalAction()` to `sharingModalClosedAction()`
4. **Rename reducer function** from `reduceCloseSharingModal()` to `reduceSharingModalClosed()`
5. **Rename registration function** from `closeSharingModalReducerRegistration()` to `sharingModalClosedReducerRegistration()`
6. **Update handler key** in registration to `"sharingModalClosed"`

**Pattern:**
```typescript
export interface SharingModalClosed {
    type: "sharingModalClosed";
}

export function sharingModalClosedAction(): SharingModalClosed {
    return {type: "sharingModalClosed"};
}

export function reduceSharingModalClosed(
    {shareModal, ...rest}: CatalogViewerState,
    _: SharingModalClosed,
): CatalogViewerState {
    return rest;
}

export function sharingModalClosedReducerRegistration(handlers: any) {
    handlers["sharingModalClosed"] = reduceSharingModalClosed as (
        state: CatalogViewerState,
        action: SharingModalClosed
    ) => CatalogViewerState;
}
```

**Test file updates:**
1. **Update describe block** to `describe("action:sharingModalClosed", () => {`
2. **Update all function calls** to use new names
3. **Update interface references** to `SharingModalClosed`

All tests must continue to pass.

### Step 3: Fix `delete-closeDeleteAlbumDialog.ts`

Fix the action structure for `delete-closeDeleteAlbumDialog.ts` to comply with coding principles.

**Files to modify:**
- `web/src/core/catalog/domain/actions/delete-closeDeleteAlbumDialog.ts`
- `web/src/core/catalog/domain/actions/delete-closeDeleteAlbumDialog.test.ts`

**Changes required:**

1. **Rename action interface** from `CloseDeleteAlbumDialogAction` to `DeleteAlbumDialogClosed`
2. **Update type property** from `"CloseDeleteAlbumDialog"` to `"deleteAlbumDialogClosed"`
3. **Rename factory function** from `closeDeleteAlbumDialogAction()` to `deleteAlbumDialogClosedAction()`
4. **Rename reducer function** from `reduceCloseDeleteAlbumDialog()` to `reduceDeleteAlbumDialogClosed()`
5. **Rename registration function** from `closeDeleteAlbumDialogReducerRegistration()` to `deleteAlbumDialogClosedReducerRegistration()`
6. **Update handler key** in registration to `"deleteAlbumDialogClosed"`

**Pattern:**
```typescript
export interface DeleteAlbumDialogClosed {
    type: "deleteAlbumDialogClosed";
}

export function deleteAlbumDialogClosedAction(): DeleteAlbumDialogClosed {
    return {type: "deleteAlbumDialogClosed"};
}

export function reduceDeleteAlbumDialogClosed(
    state: CatalogViewerState,
    _action: DeleteAlbumDialogClosed
): CatalogViewerState {
    if (!state.deleteDialog) return state;
    return {
        ...state,
        deleteDialog: undefined,
    };
}

export function deleteAlbumDialogClosedReducerRegistration(handlers: any) {
    handlers["deleteAlbumDialogClosed"] = reduceDeleteAlbumDialogClosed as (
        state: CatalogViewerState,
        action: DeleteAlbumDialogClosed
    ) => CatalogViewerState;
}
```

**Test file updates:**
1. **Update describe block** to `describe("action:deleteAlbumDialogClosed", () => {`
2. **Update all function calls** to use new names
3. **Update interface references** to `DeleteAlbumDialogClosed`

All tests must continue to pass.

### Step 4: Fix `content-noAlbumAvailableAction.ts`

Fix the action structure for `content-noAlbumAvailableAction.ts` to comply with coding principles.

**Files to modify:**
- `web/src/core/catalog/domain/actions/content-noAlbumAvailableAction.ts`
- `web/src/core/catalog/domain/actions/content-noAlbumAvailableAction.test.ts`

**Changes required:**

1. **Rename action interface** from `NoAlbumAvailableAction` to `NoAlbumAvailable`
2. **Update type property** from `"NoAlbumAvailableAction"` to `"noAlbumAvailable"`
3. **Rename factory function** from `noAlbumAvailableAction()` to `noAlbumAvailableAction()` (already correct)
4. **Rename reducer function** from `reduceNoAlbumAvailable()` to `reduceNoAlbumAvailable()` (already correct)
5. **Rename registration function** from `noAlbumAvailableReducerRegistration()` to `noAlbumAvailableReducerRegistration()` (already correct)
6. **Update handler key** in registration to `"noAlbumAvailable"`

**Pattern:**
```typescript
export interface NoAlbumAvailable {
    type: "noAlbumAvailable";
}

export function noAlbumAvailableAction(): NoAlbumAvailable {
    return {type: "noAlbumAvailable"};
}

export function reduceNoAlbumAvailable(
    current: CatalogViewerState,
    action: NoAlbumAvailable
): CatalogViewerState {
    // existing implementation
}

export function noAlbumAvailableReducerRegistration(handlers: any) {
    handlers["noAlbumAvailable"] = reduceNoAlbumAvailable as (
        state: CatalogViewerState,
        action: NoAlbumAvailable
    ) => CatalogViewerState;
}
```

**Test file updates:**
1. **Update describe block** to `describe("action:noAlbumAvailable", () => {`
2. **Update interface references** to `NoAlbumAvailable`

All tests must continue to pass.

### Step 5: Fix `content-startLoadingMediasAction.ts`

Fix the action structure for `content-startLoadingMediasAction.ts` to comply with coding principles.

**Files to modify:**
- `web/src/core/catalog/domain/actions/content-startLoadingMediasAction.ts`
- `web/src/core/catalog/domain/actions/content-startLoadingMediasAction.test.ts`

**Changes required:**

1. **Rename action interface** from `StartLoadingMediasAction` to `MediasLoadingStarted`
2. **Update type property** from `"StartLoadingMediasAction"` to `"mediasLoadingStarted"`
3. **Rename factory function** from `startLoadingMediasAction()` to `mediasLoadingStartedAction()`
4. **Rename reducer function** from `reduceStartLoadingMedias()` to `reduceMediasLoadingStarted()`
5. **Rename registration function** from `startLoadingMediasReducerRegistration()` to `mediasLoadingStartedReducerRegistration()`
6. **Update handler key** in registration to `"mediasLoadingStarted"`

**Pattern to follow (Case 2: Single additional property):**
```typescript
export interface MediasLoadingStarted {
    type: "mediasLoadingStarted";
    albumId: AlbumId;
}

export function mediasLoadingStartedAction(albumId: AlbumId): MediasLoadingStarted {
    return {type: "mediasLoadingStarted", albumId};
}

export function reduceMediasLoadingStarted(
    current: CatalogViewerState,
    action: MediasLoadingStarted
): CatalogViewerState {
    // existing implementation
}

export function mediasLoadingStartedReducerRegistration(handlers: any) {
    handlers["mediasLoadingStarted"] = reduceMediasLoadingStarted as (
        state: CatalogViewerState,
        action: MediasLoadingStarted
    ) => CatalogViewerState;
}
```

**Test file updates:**
1. **Update describe block** to `describe("action:mediasLoadingStarted", () => {`
2. **Update all function calls** to use new names
3. **Update interface references** to `MediasLoadingStarted`

All tests must continue to pass.

### Step 6: Fix `delete-albumFailedToDeleteAction.ts`

Fix the action structure for `delete-albumFailedToDeleteAction.ts` to comply with coding principles.

**Files to modify:**
- `web/src/core/catalog/domain/actions/delete-albumFailedToDeleteAction.ts`
- `web/src/core/catalog/domain/actions/delete-albumFailedToDeleteAction.test.ts`

**Changes required:**

1. **Rename action interface** from `AlbumFailedToDeleteAction` to `AlbumDeleteFailed`
2. **Update type property** from `"AlbumFailedToDelete"` to `"albumDeleteFailed"`
3. **Rename factory function** from `albumFailedToDeleteAction()` to `albumDeleteFailedAction()`
4. **Rename reducer function** from `reduceAlbumFailedToDelete()` to `reduceAlbumDeleteFailed()`
5. **Rename registration function** from `albumFailedToDeleteReducerRegistration()` to `albumDeleteFailedReducerRegistration()`
6. **Update handler key** in registration to `"albumDeleteFailed"`

**Pattern to follow (Case 2: Single additional property):**
```typescript
export interface AlbumDeleteFailed {
    type: "albumDeleteFailed";
    error: string;
}

export function albumDeleteFailedAction(error: string): AlbumDeleteFailed {
    if (!error || error.trim() === "") {
        throw new Error("AlbumDeleteFailed requires a non-empty error message");
    }
    return {
        error,
        type: "albumDeleteFailed",
    };
}

export function reduceAlbumDeleteFailed(
    current: CatalogViewerState,
    action: AlbumDeleteFailed,
): CatalogViewerState {
    // existing implementation
}

export function albumDeleteFailedReducerRegistration(handlers: any) {
    handlers["albumDeleteFailed"] = reduceAlbumDeleteFailed as (
        state: CatalogViewerState,
        action: AlbumDeleteFailed
    ) => CatalogViewerState;
}
```

**Test file updates:**
1. **Update describe block** to `describe("action:albumDeleteFailed", () => {`
2. **Update all function calls** to use new names and single parameter
3. **Update interface references** to `AlbumDeleteFailed`
4. **Update factory calls** from `albumFailedToDeleteAction({error: "..."})` to `albumDeleteFailedAction("...")`

All tests must continue to pass.

### Step 7: Fix `sharing-removeSharingAction.ts`

Fix the action structure for `sharing-removeSharingAction.ts` to comply with coding principles.

**Files to modify:**
- `web/src/core/catalog/domain/actions/sharing-removeSharingAction.ts`
- `web/src/core/catalog/domain/actions/sharing-removeSharingAction.test.ts`

**Changes required:**

1. **Rename action interface** from `RemoveSharingAction` to `SharingRemoved`
2. **Update type property** from `"RemoveSharingAction"` to `"sharingRemoved"`
3. **Rename factory function** from `removeSharingAction()` to `sharingRemovedAction()`
4. **Rename reducer function** from `reduceRemoveSharing()` to `reduceSharingRemoved()`
5. **Rename registration function** from `removeSharingReducerRegistration()` to `sharingRemovedReducerRegistration()`
6. **Update handler key** in registration to `"sharingRemoved"`
7. **Simplify factory function** to accept only string parameter (remove overload)

**Pattern to follow (Case 2: Single additional property):**
```typescript
export interface SharingRemoved {
    type: "sharingRemoved";
    email: string;
}

export function sharingRemovedAction(email: string): SharingRemoved {
    return {
        email,
        type: "sharingRemoved",
    };
}

export function reduceSharingRemoved(
    current: CatalogViewerState,
    action: SharingRemoved,
): CatalogViewerState {
    // existing implementation
}

export function sharingRemovedReducerRegistration(handlers: any) {
    handlers["sharingRemoved"] = reduceSharingRemoved as (
        state: CatalogViewerState,
        action: SharingRemoved
    ) => CatalogViewerState;
}
```

**Test file updates:**
1. **Update describe block** to `describe("action:sharingRemoved", () => {`
2. **Update all function calls** to use new names
3. **Update interface references** to `SharingRemoved`
4. **Update factory calls** to use string parameter directly

All tests must continue to pass.

### Step 8: Fix `content-mediasLoadedAction.ts`

Fix the action structure for `content-mediasLoadedAction.ts` to comply with coding principles.

**Files to modify:**
- `web/src/core/catalog/domain/actions/content-mediasLoadedAction.ts`
- `web/src/core/catalog/domain/actions/content-mediasLoadedAction.test.ts`

**Changes required:**

1. **Rename action interface** from `MediasLoadedAction` to `MediasLoaded`
2. **Update type property** from `"MediasLoadedAction"` to `"mediasLoaded"`
3. **Rename factory function** from `mediasLoadedAction()` to `mediasLoadedAction()` (already correct)
4. **Rename reducer function** from `reduceMediasLoaded()` to `reduceMediasLoaded()` (already correct)
5. **Rename registration function** from `mediasLoadedReducerRegistration()` to `mediasLoadedReducerRegistration()` (already correct)
6. **Update handler key** in registration to `"mediasLoaded"`

**Pattern to follow (Case 3: Multiple additional properties):**
```typescript
export interface MediasLoaded {
    type: "mediasLoaded";
    albumId: AlbumId;
    medias: MediaWithinADay[];
}

export function mediasLoadedAction(props: Omit<MediasLoaded, "type">): MediasLoaded {
    return {
        ...props,
        type: "mediasLoaded",
    };
}

export function reduceMediasLoaded(
    current: CatalogViewerState,
    action: MediasLoaded,
): CatalogViewerState {
    // existing implementation
}

export function mediasLoadedReducerRegistration(handlers: any) {
    handlers["mediasLoaded"] = reduceMediasLoaded as (state: CatalogViewerState, action: MediasLoaded) => CatalogViewerState;
}
```

**Test file updates:**
1. **Update describe block** to `describe("action:mediasLoaded", () => {`
2. **Update interface references** to `MediasLoaded`

All tests must continue to pass.

### Step 9: Fix `content-mediaFailedToLoadAction.ts`

Fix the action structure for `content-mediaFailedToLoadAction.ts` to comply with coding principles.

**Files to modify:**
- `web/src/core/catalog/domain/actions/content-mediaFailedToLoadAction.ts`
- `web/src/core/catalog/domain/actions/content-mediaFailedToLoadAction.test.ts`

**Changes required:**

1. **Rename action interface** from `MediaFailedToLoadAction` to `MediaLoadFailed`
2. **Update type property** from `"MediaFailedToLoadAction"` to `"mediaLoadFailed"`
3. **Rename factory function** from `mediaFailedToLoadAction()` to `mediaLoadFailedAction()`
4. **Rename reducer function** from `reduceMediaFailedToLoad()` to `reduceMediaLoadFailed()`
5. **Rename registration function** from `mediaFailedToLoadReducerRegistration()` to `mediaLoadFailedReducerRegistration()`
6. **Update handler key** in registration to `"mediaLoadFailed"`

**Pattern to follow (Case 3: Multiple additional properties):**
```typescript
export interface MediaLoadFailed {
    type: "mediaLoadFailed";
    albums?: Album[];
    selectedAlbum?: Album;
    error: Error;
}

export function mediaLoadFailedAction(props: Omit<MediaLoadFailed, "type">): MediaLoadFailed {
    return {
        ...props,
        type: "mediaLoadFailed",
    };
}

export function reduceMediaLoadFailed(
    current: CatalogViewerState,
    action: MediaLoadFailed,
): CatalogViewerState {
    // existing implementation
}

export function mediaLoadFailedReducerRegistration(handlers: any) {
    handlers["mediaLoadFailed"] = reduceMediaLoadFailed as (
        state: CatalogViewerState,
        action: MediaLoadFailed
    ) => CatalogViewerState;
}
```

**Test file updates:**
1. **Update describe block** to `describe("action:mediaLoadFailed", () => {`
2. **Update all function calls** to use new names
3. **Update interface references** to `MediaLoadFailed`

All tests must continue to pass.

### Step 10: Fix `sharing-addSharingAction.ts`

Fix the action structure for `sharing-addSharingAction.ts` to comply with coding principles.

**Files to modify:**
- `web/src/core/catalog/domain/actions/sharing-addSharingAction.ts`
- `web/src/core/catalog/domain/actions/sharing-addSharingAction.test.ts`

**Changes required:**

1. **Rename action interface** from `AddSharingAction` to `SharingAdded`
2. **Update type property** from `"AddSharingAction"` to `"sharingAdded"`
3. **Rename factory function** from `addSharingAction()` to `sharingAddedAction()`
4. **Rename reducer function** from `reduceAddSharing()` to `reduceSharingAdded()`
5. **Rename registration function** from `addSharingReducerRegistration()` to `sharingAddedReducerRegistration()`
6. **Update handler key** in registration to `"sharingAdded"`
7. **Simplify factory function** to accept only Sharing parameter (remove overload)

**Pattern to follow (Case 2: Single additional property):**
```typescript
export interface SharingAdded {
    type: "sharingAdded";
    sharing: Sharing;
}

export function sharingAddedAction(sharing: Sharing): SharingAdded {
    return {
        sharing,
        type: "sharingAdded",
    };
}

export function reduceSharingAdded(
    current: CatalogViewerState,
    action: SharingAdded,
): CatalogViewerState {
    // existing implementation
}

export function sharingAddedReducerRegistration(handlers: any) {
    handlers["sharingAdded"] = reduceSharingAdded as (state: CatalogViewerState, action: SharingAdded) => CatalogViewerState;
}
```

**Test file updates:**
1. **Update describe block** to `describe("action:sharingAdded", () => {`
2. **Update all function calls** to use new names
3. **Update interface references** to `SharingAdded`
4. **Update factory calls** to use Sharing parameter directly

All tests must continue to pass.

### Step 11: Fix `sharing-openSharingModalAction.ts`

Fix the action structure for `sharing-openSharingModalAction.ts` to comply with coding principles.

**Files to modify:**
- `web/src/core/catalog/domain/actions/sharing-openSharingModalAction.ts`
- `web/src/core/catalog/domain/actions/sharing-openSharingModalAction.test.ts`

**Changes required:**

1. **Rename action interface** from `OpenSharingModalAction` to `SharingModalOpened`
2. **Update type property** from `"OpenSharingModalAction"` to `"sharingModalOpened"`
3. **Rename factory function** from `openSharingModalAction()` to `sharingModalOpenedAction()`
4. **Rename reducer function** from `reduceOpenSharingModal()` to `reduceSharingModalOpened()`
5. **Rename registration function** from `openSharingModalReducerRegistration()` to `sharingModalOpenedReducerRegistration()`
6. **Update handler key** in registration to `"sharingModalOpened"`
7. **Simplify factory function** to accept only AlbumId parameter (remove overload)

**Pattern to follow (Case 2: Single additional property):**
```typescript
export interface SharingModalOpened {
    type: "sharingModalOpened";
    albumId: AlbumId;
}

export function sharingModalOpenedAction(albumId: AlbumId): SharingModalOpened {
    return {
        albumId,
        type: "sharingModalOpened",
    };
}

export function reduceSharingModalOpened(
    current: CatalogViewerState,
    action: SharingModalOpened,
): CatalogViewerState {
    return withOpenShareModal(current, action.albumId);
}

export function sharingModalOpenedReducerRegistration(handlers: any) {
    handlers["sharingModalOpened"] = reduceSharingModalOpened as (
        state: CatalogViewerState,
        action: SharingModalOpened
    ) => CatalogViewerState;
}
```

**Test file updates:**
1. **Update describe block** to `describe("action:sharingModalOpened", () => {`
2. **Update all function calls** to use new names
3. **Update interface references** to `SharingModalOpened`
4. **Update factory calls** to use AlbumId parameter directly

All tests must continue to pass.

### Step 12: Fix `sharing-sharingModalErrorAction.ts`

Fix the action structure for `sharing-sharingModalErrorAction.ts` to comply with coding principles.

**Files to modify:**
- `web/src/core/catalog/domain/actions/sharing-sharingModalErrorAction.ts`
- `web/src/core/catalog/domain/actions/sharing-sharingModalErrorAction.test.ts`

**Changes required:**

1. **Rename action interface** from `SharingModalErrorAction` to `SharingModalErrorOccurred`
2. **Update type property** from `"SharingModalErrorAction"` to `"sharingModalErrorOccurred"`
3. **Rename factory function** from `sharingModalErrorAction()` to `sharingModalErrorOccurredAction()`
4. **Rename reducer function** from `reduceSharingModalError()` to `reduceSharingModalErrorOccurred()`
5. **Rename registration function** from `sharingModalErrorReducerRegistration()` to `sharingModalErrorOccurredReducerRegistration()`
6. **Update handler key** in registration to `"sharingModalErrorOccurred"`
7. **Simplify factory function** to accept only ShareError parameter (remove overload)

**Pattern to follow (Case 2: Single additional property):**
```typescript
export interface SharingModalErrorOccurred {
    type: "sharingModalErrorOccurred";
    error: ShareError;
}

export function sharingModalErrorOccurredAction(error: ShareError): SharingModalErrorOccurred {
    return {
        error,
        type: "sharingModalErrorOccurred",
    };
}

export function reduceSharingModalErrorOccurred(
    current: CatalogViewerState,
    {error}: SharingModalErrorOccurred,
): CatalogViewerState {
    // existing implementation
}

export function sharingModalErrorOccurredReducerRegistration(handlers: any) {
    handlers["sharingModalErrorOccurred"] = reduceSharingModalErrorOccurred as (
        state: CatalogViewerState,
        action: SharingModalErrorOccurred
    ) => CatalogViewerState;
}
```

**Test file updates:**
1. **Update describe block** to `describe("action:sharingModalErrorOccurred", () => {`
2. **Update all function calls** to use new names
3. **Update interface references** to `SharingModalErrorOccurred`
4. **Update factory calls** to use ShareError parameter directly

All tests must continue to pass.

### Step 13: Fix `delete-openDeleteAlbumDialog.ts`

Fix the action structure for `delete-openDeleteAlbumDialog.ts` to comply with coding principles.

**Files to modify:**
- `web/src/core/catalog/domain/actions/delete-openDeleteAlbumDialog.ts`
- `web/src/core/catalog/domain/actions/delete-openDeleteAlbumDialog.test.ts`

**Changes required:**

1. **Rename action interface** from `OpenDeleteAlbumDialogAction` to `DeleteAlbumDialogOpened`
2. **Update type property** from `"OpenDeleteAlbumDialog"` to `"deleteAlbumDialogOpened"`
3. **Rename factory function** from `openDeleteAlbumDialogAction()` to `deleteAlbumDialogOpenedAction()`
4. **Rename reducer function** from `reduceOpenDeleteAlbumDialog()` to `reduceDeleteAlbumDialogOpened()`
5. **Rename registration function** from `openDeleteAlbumDialogReducerRegistration()` to `deleteAlbumDialogOpenedReducerRegistration()`
6. **Update handler key** in registration to `"deleteAlbumDialogOpened"`

**Pattern to follow (Case 1: No additional properties):**
```typescript
export interface DeleteAlbumDialogOpened {
    type: "deleteAlbumDialogOpened";
}

export function deleteAlbumDialogOpenedAction(): DeleteAlbumDialogOpened {
    return {
        type: "deleteAlbumDialogOpened",
    };
}

export function reduceDeleteAlbumDialogOpened(
    current: CatalogViewerState,
    action: DeleteAlbumDialogOpened,
): CatalogViewerState {
    // existing implementation
}

export function deleteAlbumDialogOpenedReducerRegistration(handlers: any) {
    handlers["deleteAlbumDialogOpened"] = reduceDeleteAlbumDialogOpened as (
        state: CatalogViewerState,
        action: DeleteAlbumDialogOpened
    ) => CatalogViewerState;
}
```

**Test file updates:**
1. **Update describe block** to `describe("action:deleteAlbumDialogOpened", () => {`
2. **Update all function calls** to use new names
3. **Update interface references** to `DeleteAlbumDialogOpened`

All tests must continue to pass.

### Step 14: Fix `content-albumsFilteredAction.ts`

Fix the action structure for `content-albumsFilteredAction.ts` to comply with coding principles.

**Files to modify:**
- `web/src/core/catalog/domain/actions/content-albumsFilteredAction.ts`
- `web/src/core/catalog/domain/actions/content-albumsFilteredAction.test.ts`

**Changes required:**

1. **Rename action interface** from `AlbumsFilteredAction` to `AlbumsFiltered`
2. **Update type property** from `"AlbumsFilteredAction"` to `"albumsFiltered"`
3. **Rename factory function** from `albumsFilteredAction()` to `albumsFilteredAction()` (already correct)
4. **Rename reducer function** from `reduceAlbumsFiltered()` to `reduceAlbumsFiltered()` (already correct)
5. **Rename registration function** from `albumsFilteredReducerRegistration()` to `albumsFilteredReducerRegistration()` (already correct)
6. **Update handler key** in registration to `"albumsFiltered"`

**Pattern to follow (Case 3: Multiple additional properties):**
```typescript
export interface AlbumsFiltered extends RedirectToAlbumIdAction {
    type: "albumsFiltered";
    criterion: AlbumFilterCriterion;
}

export function albumsFilteredAction(props: Omit<AlbumsFiltered, "type">): AlbumsFiltered {
    return {...props, type: "albumsFiltered"};
}

export function reduceAlbumsFiltered(
    current: CatalogViewerState,
    action: AlbumsFiltered
): CatalogViewerState {
    // existing implementation
}

export function albumsFilteredReducerRegistration(handlers: any) {
    handlers["albumsFiltered"] = reduceAlbumsFiltered as (
        state: CatalogViewerState,
        action: AlbumsFiltered
    ) => CatalogViewerState;
}
```

**Test file updates:**
1. **Update describe block** to `describe("action:albumsFiltered", () => {`
2. **Update interface references** to `AlbumsFiltered`

All tests must continue to pass.

### Step 15: Fix `content-albumsLoadedAction.ts`

Fix the action structure for `content-albumsLoadedAction.ts` to comply with coding principles.

**Files to modify:**
- `web/src/core/catalog/domain/actions/content-albumsLoadedAction.ts`
- `web/src/core/catalog/domain/actions/content-albumsLoadedAction.test.ts`

**Changes required:**

1. **Rename action interface** from `AlbumsLoadedAction` to `AlbumsLoaded`
2. **Update type property** from `"AlbumsLoadedAction"` to `"albumsLoaded"`
3. **Rename factory function** from `albumsLoadedAction()` to `albumsLoadedAction()` (already correct)
4. **Rename reducer function** from `reduceAlbumsLoaded()` to `reduceAlbumsLoaded()` (already correct)
5. **Rename registration function** from `albumsLoadedReducerRegistration()` to `albumsLoadedReducerRegistration()` (already correct)
6. **Update handler key** in registration to `"albumsLoaded"`
7. **Remove alternative factory signature** that accepts `Album[]` directly

**Pattern to follow (Case 3: Multiple additional properties):**
```typescript
export interface AlbumsLoaded extends RedirectToAlbumIdAction {
    type: "albumsLoaded";
    albums: Album[];
}

export function albumsLoadedAction(props: Omit<AlbumsLoaded, "type">): AlbumsLoaded {
    return {
        ...props,
        type: "albumsLoaded",
    };
}

export function reduceAlbumsLoaded(
    current: CatalogViewerState,
    action: AlbumsLoaded,
): CatalogViewerState {
    // existing implementation
}

export function albumsLoadedReducerRegistration(handlers: any) {
    handlers["albumsLoaded"] = reduceAlbumsLoaded as (
        state: CatalogViewerState,
        action: AlbumsLoaded
    ) => CatalogViewerState;
}
```

**Test file updates:**
1. **Update describe block** to `describe("action:albumsLoaded", () => {`
2. **Update interface references** to `AlbumsLoaded`
3. **Convert factory calls** from `albumsLoadedAction(twoAlbums)` to `albumsLoadedAction({albums: twoAlbums})`

All tests must continue to pass.

### Step 16: Fix `content-albumsAndMediasLoadedAction.ts`

Fix the action structure for `content-albumsAndMediasLoadedAction.ts` to comply with coding principles.

**Files to modify:**
- `web/src/core/catalog/domain/actions/content-albumsAndMediasLoadedAction.ts`
- `web/src/core/catalog/domain/actions/content-albumsAndMediasLoadedAction.test.ts`

**Changes required:**

1. **Rename action interface** from `AlbumsAndMediasLoadedAction` to `AlbumsAndMediasLoaded`
2. **Update type property** from `"AlbumsAndMediasLoadedAction"` to `"albumsAndMediasLoaded"`
3. **Rename factory function** from `albumsAndMediasLoadedAction()` to `albumsAndMediasLoadedAction()` (already correct)
4. **Rename reducer function** from `reduceAlbumsAndMediasLoaded()` to `reduceAlbumsAndMediasLoaded()` (already correct)
5. **Rename registration function** from `albumsAndMediasLoadedReducerRegistration()` to `albumsAndMediasLoadedReducerRegistration()` (already correct)
6. **Update handler key** in registration to `"albumsAndMediasLoaded"`

**Pattern to follow (Case 3: Multiple additional properties):**
```typescript
export interface AlbumsAndMediasLoaded extends RedirectToAlbumIdAction {
    type: "albumsAndMediasLoaded";
    albums: Album[];
    medias: MediaWithinADay[];
    selectedAlbum?: Album;
}

export function albumsAndMediasLoadedAction(props: Omit<AlbumsAndMediasLoaded, "type">): AlbumsAndMediasLoaded {
    return {
        ...props,
        type: "albumsAndMediasLoaded",
    };
}

export function reduceAlbumsAndMediasLoaded(
    current: CatalogViewerState,
    action: AlbumsAndMediasLoaded
): CatalogViewerState {
    // existing implementation
}

export function albumsAndMediasLoadedReducerRegistration(handlers: any) {
    handlers["albumsAndMediasLoaded"] = reduceAlbumsAndMediasLoaded as (
        state: CatalogViewerState,
        action: AlbumsAndMediasLoaded
    ) => CatalogViewerState;
}
```

**Test file updates:**
1. **Update describe block** to `describe("action:albumsAndMediasLoaded", () => {`
2. **Update interface references** to `AlbumsAndMediasLoaded`

All tests must continue to pass.

### Step 17: Fix `delete-albumDeletedAction.ts`

Fix the action structure for `delete-albumDeletedAction.ts` to comply with coding principles.

**Files to modify:**
- `web/src/core/catalog/domain/actions/delete-albumDeletedAction.ts`
- `web/src/core/catalog/domain/actions/delete-albumDeletedAction.test.ts`

**Changes required:**

1. **Rename action interface** from `AlbumDeletedAction` to `AlbumDeleted`
2. **Update type property** from `"AlbumDeleted"` to `"albumDeleted"`
3. **Rename factory function** from `albumDeletedAction()` to `albumDeletedAction()` (already correct)
4. **Rename reducer function** from `reduceAlbumDeleted()` to `reduceAlbumDeleted()` (already correct)
5. **Rename registration function** from `albumDeletedReducerRegistration()` to `albumDeletedReducerRegistration()` (already correct)
6. **Update handler key** in registration to `"albumDeleted"`

**Pattern to follow (Case 3: Multiple additional properties):**
```typescript
export interface AlbumDeleted extends RedirectToAlbumIdAction {
    type: "albumDeleted";
    albums: Album[];
    redirectTo?: AlbumId;
}

export function albumDeletedAction(props: Omit<AlbumDeleted, "type">): AlbumDeleted {
    return {
        ...props,
        type: "albumDeleted",
    };
}

export function reduceAlbumDeleted(
    current: CatalogViewerState,
    action: AlbumDeleted,
): CatalogViewerState {
    // existing implementation
}

export function albumDeletedReducerRegistration(handlers: any) {
    handlers["albumDeleted"] = reduceAlbumDeleted as (
        state: CatalogViewerState,
        action: AlbumDeleted
    ) => CatalogViewerState;
}
```

**Test file updates:**
1. **Update describe block** to `describe("action:albumDeleted", () => {`
2. **Update interface references** to `AlbumDeleted`

All tests must continue to pass.

### Step 18: Update Index File Integration

Update the main index file to integrate all the refactored actions.

**Files to modify:**
- `web/src/core/catalog/domain/actions/index.ts`

**Changes required:**

1. **Update all import statements** to use new interface names (remove "Action" suffixes)
2. **Update `catalogActions` object** to use new factory function names
3. **Update `CatalogViewerAction` union type** to use new interface names
4. **Update `reducerRegistrations` array** to use new registration function names
5. **Update all export statements** to use new names

**Pattern:**

```typescript
// Update imports
import {SharingAdded, sharingAddedAction, sharingAddedReducerRegistration} from "./sharing-addSharingAction";
import {AlbumsAndMediasLoaded, albumsAndMediasLoadedAction, albumsAndMediasLoadedReducerRegistration} from "./content-albumsAndMediasLoadedAction";
// ... all other imports with new names

// Update catalogActions
export const catalogActions = {
    sharingModalOpenedAction,
    sharingAddedAction,
    sharingRemovedAction,
    sharingModalClosedAction,
    sharingModalErrorOccurredAction,
    albumsAndMediasLoadedAction,
    albumsLoadedAction,
    albumDeletedAction,
    mediaLoadFailedAction,
    noAlbumAvailableAction,
    mediasLoadingStartedAction,
    albumsFilteredAction,
    mediasLoadedAction,
    deleteAlbumDialogOpenedAction,
    albumDeleteFailedAction,
    deleteAlbumDialogClosedAction,
    deleteAlbumStartedAction,
};

// Update union type
export type CatalogViewerAction =
    SharingAdded
    | AlbumsAndMediasLoaded
    | AlbumsFiltered
    | AlbumsLoaded
    | AlbumDeleted
    | SharingModalClosed
    | MediaLoadFailed
    | MediasLoaded
    | NoAlbumAvailable
    | SharingModalOpened
    | SharingRemoved
    | SharingModalErrorOccurred
    | MediasLoadingStarted
    | DeleteAlbumDialogOpened
    | AlbumDeleteFailed
    | DeleteAlbumDialogClosed
    | DeleteAlbumStarted

// Update registrations
const reducerRegistrations = [
    sharingAddedReducerRegistration,
    albumsAndMediasLoadedReducerRegistration,
    albumsLoadedReducerRegistration,
    albumDeletedReducerRegistration,
    mediaLoadFailedReducerRegistration,
    noAlbumAvailableReducerRegistration,
    mediasLoadingStartedReducerRegistration,
    albumsFilteredReducerRegistration,
    sharingModalOpenedReducerRegistration,
    sharingRemovedReducerRegistration,
    sharingModalClosedReducerRegistration,
    sharingModalErrorOccurredReducerRegistration,
    mediasLoadedReducerRegistration,
    deleteAlbumDialogOpenedReducerRegistration,
    albumDeleteFailedReducerRegistration,
    deleteAlbumDialogClosedReducerRegistration,
    deleteAlbumStartedReducerRegistration,
];
```

Ensure all existing functionality continues to work and all tests pass.

### Step 19: Final Integration and Verification

Complete the integration by verifying all references are updated and the system works correctly.

**Files to check and potentially modify:**
- Any files that import from `web/src/core/catalog/domain/actions/`
- Thunk files that dispatch these actions
- Component files that use these actions

**Verification steps:**

1. **Run all tests** to ensure no regressions
2. **Check TypeScript compilation** for any remaining errors
3. **Verify action dispatching** works correctly in the application
4. **Check for unused imports/exports**
5. **Verify all action names follow past tense event naming**

**Expected outcomes:**
- All action interfaces follow `PascalCase` naming without "Action" suffix
- All action types follow `camelCase` past tense event naming
- All factory functions follow `actionNameAction()` pattern
- All reducer functions follow `reduceActionName()` pattern
- All registration functions follow `actionNameReducerRegistration()` pattern
- All tests use `action:actionName` describe blocks
- All tests pass without modification to business logic
- No TypeScript compilation errors
- Action dispatching continues to work correctly

This completes the comprehensive refactoring to align with the coding principles.
