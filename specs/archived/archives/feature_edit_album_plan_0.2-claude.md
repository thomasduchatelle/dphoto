# Date Edit Dialog - Language Layer Design

## Dialog Properties Interface

Based on the Create Album Dialog pattern and the date edit requirements, the dialog needs these properties:

### Selector Interface

```typescript
interface EditDatesDialogProperties {
    isOpen: boolean
    albumId?: AlbumId
    albumName: string
    currentStart: Dayjs
    currentEnd: Dayjs
    startsAtStartOfTheDay: boolean
    endsAtEndOfTheDay: boolean
    isLoading: boolean
    errorMessage?: string
}
```

### Selector Signature

```typescript
function selectEditDatesDialog(state: CatalogViewerState): EditDatesDialogProperties
```

### Main State Property

```typescript
interface EditDatesDialogState {
    albumId: AlbumId
    albumName: string
    originalStart: Date
    originalEnd: Date
    currentStart: Dayjs
    currentEnd: Dayjs
    startsAtStartOfTheDay: boolean
    endsAtEndOfTheDay: boolean
    isLoading: boolean
    errorMessage?: string
}
```

Add to `CatalogViewerState`:

```typescript
editDatesDialog ? : EditDatesDialogState
```

## User Interactions and Required Actions

### Thunk Analysis Notes:

1. **Open Dialog Thunk**: User clicks "Edit Dates" → needs current album data to populate dialog
2. **Update Date Thunk**: User changes start/end date → needs to update dialog state immediately
3. **Toggle Day Mode Thunk**: User checks/unchecks day mode → needs to adjust date precision
4. **Submit Edit Thunk**: User clicks Save → calls API, handles loading/success/error states
5. **Close Dialog Thunk**: User cancels or completes → clears dialog state

### Required Actions:

1. `editDatesDialogOpened` - Opens dialog with album data
2. `editDatesDialogClosed` - Closes dialog and clears state
3. `editDatesStartDateChanged` - Updates start date in dialog
4. `editDatesEndDateChanged` - Updates end date in dialog
5. `editDatesStartDayModeToggled` - Toggles start date day mode
6. `editDatesEndDayModeToggled` - Toggles end date day mode
7. `editDatesSubmissionStarted` - Sets loading state
8. `editDatesSubmissionSucceeded` - Closes dialog and triggers refresh
9. `editDatesSubmissionFailed` - Shows error message

## Task Breakdown

### Task 1: Create Edit Dates Dialog State and Selector

**Requirements (BDD Style):**

```
GIVEN the system needs to support date editing functionality
WHEN defining the state model for the edit dates dialog
THEN the selector should return dialog properties for component rendering
AND the state should track current date values and UI modes
AND the state should handle loading and error states
```

**Implementation Details:**

- Create file: `web/src/core/catalog/language/edit-dates-dialog-state.ts`
- Export interface: `EditDatesDialogState`
- Export interface: `EditDatesDialogProperties`
- Export function: `selectEditDatesDialog(state: CatalogViewerState): EditDatesDialogProperties`
- Update `CatalogViewerState` in `catalog-state.ts` to include `editDatesDialog?: EditDatesDialogState`

---

### Task 2: Implement Edit Dates Dialog Opened Action

**Requirements (BDD Style):**

```
GIVEN a user wants to edit album dates
WHEN the editDatesDialogOpened action is dispatched with an album
THEN the selector should return isOpen: true
AND the selector should return the album's current start and end dates
AND the selector should return albumId and albumName
AND the selector should return default day mode settings (both true)
AND the selector should return isLoading: false and no error message
```

**Implementation Details:**

- Create file: `web/src/core/catalog/edit-dates/action-editDatesDialogOpened.ts`
- Export interface: `EditDatesDialogOpened`
- Export function: `editDatesDialogOpened(album: Album): EditDatesDialogOpened`
- Export function: `reduceEditDatesDialogOpened(state: CatalogViewerState, action: EditDatesDialogOpened): CatalogViewerState`
- Export function: `editDatesDialogOpenedReducerRegistration(handlers: any)`

**References:**

- Task 1: `EditDatesDialogState` and `selectEditDatesDialog` from `edit-dates-dialog-state.ts`

---

### Task 3: Implement Edit Dates Dialog Closed Action

**Requirements (BDD Style):**

```
GIVEN the edit dates dialog is open
WHEN the editDatesDialogClosed action is dispatched
THEN the selector should return isOpen: false
AND the editDatesDialog property should be undefined in the main state
```

**Implementation Details:**

- Create file: `web/src/core/catalog/edit-dates/action-editDatesDialogClosed.ts`
- Export interface: `EditDatesDialogClosed`
- Export function: `editDatesDialogClosed(): EditDatesDialogClosed`
- Export function: `reduceEditDatesDialogClosed(state: CatalogViewerState, action: EditDatesDialogClosed): CatalogViewerState`
- Export function: `editDatesDialogClosedReducerRegistration(handlers: any)`

**References:**

- Task 1: `selectEditDatesDialog` from `edit-dates-dialog-state.ts`

---

### Task 4: Implement Edit Dates Start Date Changed Action

**Requirements (BDD Style):**

```
GIVEN the edit dates dialog is open
WHEN the editDatesStartDateChanged action is dispatched with a new start date
THEN the selector should return the updated currentStart date
AND all other dialog properties should remain unchanged
```

**Implementation Details:**

- Create file: `web/src/core/catalog/edit-dates/action-editDatesStartDateChanged.ts`
- Export interface: `EditDatesStartDateChanged`
- Export function: `editDatesStartDateChanged(newStart: Dayjs): EditDatesStartDateChanged`
- Export function: `reduceEditDatesStartDateChanged(state: CatalogViewerState, action: EditDatesStartDateChanged): CatalogViewerState`
- Export function: `editDatesStartDateChangedReducerRegistration(handlers: any)`

**References:**

- Task 1: `EditDatesDialogState` and `selectEditDatesDialog` from `edit-dates-dialog-state.ts`

---

### Task 5: Implement Edit Dates End Date Changed Action

**Requirements (BDD Style):**

```
GIVEN the edit dates dialog is open
WHEN the editDatesEndDateChanged action is dispatched with a new end date
THEN the selector should return the updated currentEnd date
AND all other dialog properties should remain unchanged
```

**Implementation Details:**

- Create file: `web/src/core/catalog/edit-dates/action-editDatesEndDateChanged.ts`
- Export interface: `EditDatesEndDateChanged`
- Export function: `editDatesEndDateChanged(newEnd: Dayjs): EditDatesEndDateChanged`
- Export function: `reduceEditDatesEndDateChanged(state: CatalogViewerState, action: EditDatesEndDateChanged): CatalogViewerState`
- Export function: `editDatesEndDateChangedReducerRegistration(handlers: any)`

**References:**

- Task 1: `EditDatesDialogState` and `selectEditDatesDialog` from `edit-dates-dialog-state.ts`

---

### Task 6: Implement Edit Dates Day Mode Toggle Actions

**Requirements (BDD Style):**

```
GIVEN the edit dates dialog is open
WHEN the editDatesStartDayModeToggled action is dispatched with a boolean value
THEN the selector should return the updated startsAtStartOfTheDay value
AND all other dialog properties should remain unchanged

GIVEN the edit dates dialog is open  
WHEN the editDatesEndDayModeToggled action is dispatched with a boolean value
THEN the selector should return the updated endsAtEndOfTheDay value
AND all other dialog properties should remain unchanged
```

**Implementation Details:**

- Create file: `web/src/core/catalog/edit-dates/action-editDatesStartDayModeToggled.ts`
- Export interface: `EditDatesStartDayModeToggled`
- Export function: `editDatesStartDayModeToggled(startsAtStartOfTheDay: boolean): EditDatesStartDayModeToggled`
- Export function: `reduceEditDatesStartDayModeToggled(state: CatalogViewerState, action: EditDatesStartDayModeToggled): CatalogViewerState`
- Export function: `editDatesStartDayModeToggledReducerRegistration(handlers: any)`

- Create file: `web/src/core/catalog/edit-dates/action-editDatesEndDayModeToggled.ts`
- Export interface: `EditDatesEndDayModeToggled`
- Export function: `editDatesEndDayModeToggled(endsAtEndOfTheDay: boolean): EditDatesEndDayModeToggled`
- Export function: `reduceEditDatesEndDayModeToggled(state: CatalogViewerState, action: EditDatesEndDayModeToggled): CatalogViewerState`
- Export function: `editDatesEndDayModeToggledReducerRegistration(handlers: any)`

**References:**

- Task 1: `EditDatesDialogState` and `selectEditDatesDialog` from `edit-dates-dialog-state.ts`

---

### Task 7: Implement Edit Dates Submission Lifecycle Actions

**Requirements (BDD Style):**

```
GIVEN the edit dates dialog is open
WHEN the editDatesSubmissionStarted action is dispatched
THEN the selector should return isLoading: true
AND the selector should return errorMessage: undefined

GIVEN the edit dates dialog is open and loading
WHEN the editDatesSubmissionSucceeded action is dispatched  
THEN the selector should return isOpen: false
AND the editDatesDialog property should be undefined in the main state

GIVEN the edit dates dialog is open and loading
WHEN the editDatesSubmissionFailed action is dispatched with an error message
THEN the selector should return isLoading: false
AND the selector should return the provided error message
AND the dialog should remain open
```

**Implementation Details:**

- Create file: `web/src/core/catalog/edit-dates/action-editDatesSubmissionStarted.ts`
- Export interface: `EditDatesSubmissionStarted`
- Export function: `editDatesSubmissionStarted(): EditDatesSubmissionStarted`
- Export function: `reduceEditDatesSubmissionStarted(state: CatalogViewerState, action: EditDatesSubmissionStarted): CatalogViewerState`
- Export function: `editDatesSubmissionStartedReducerRegistration(handlers: any)`

- Create file: `web/src/core/catalog/edit-dates/action-editDatesSubmissionSucceeded.ts`
- Export interface: `EditDatesSubmissionSucceeded`
- Export function: `editDatesSubmissionSucceeded(): EditDatesSubmissionSucceeded`
- Export function: `reduceEditDatesSubmissionSucceeded(state: CatalogViewerState, action: EditDatesSubmissionSucceeded): CatalogViewerState`
- Export function: `editDatesSubmissionSucceededReducerRegistration(handlers: any)`

- Create file: `web/src/core/catalog/edit-dates/action-editDatesSubmissionFailed.ts`
- Export interface: `EditDatesSubmissionFailed`
- Export function: `editDatesSubmissionFailed(errorMessage: string): EditDatesSubmissionFailed`
- Export function: `reduceEditDatesSubmissionFailed(state: CatalogViewerState, action: EditDatesSubmissionFailed): CatalogViewerState`
- Export function: `editDatesSubmissionFailedReducerRegistration(handlers: any)`

**References:**

- Task 1: `EditDatesDialogState` and `selectEditDatesDialog` from `edit-dates-dialog-state.ts`

> Tokens: 15k sent, 2.4k received. Cost: $0.08 message, $0.26 session.  
