# Album Edit Feature - Implementation Plan

## Data Flow Analysis by Layer

Based on the principles handbook and requirements, here's the breakdown by architectural layers:

### Language Layer (Domain Model)

- Define edit dialog state interfaces
- Add edit-related properties to existing state
- Create utility functions for edit operations

### Action Layer

- Actions for opening/closing edit dialogs
- Actions for updating album name/folderName
- Actions for updating album dates
- Actions for handling edit errors and loading states

### Thunk Layer

- Business logic for editing album names
- Business logic for editing album dates
- Validation and error handling
- API coordination

### Adapter Layer

- REST API calls for name/folderName updates
- REST API calls for date updates

### React Component Layer

- Edit button component
- Dropdown menu for edit options
- Name edit dialog component
- Date edit dialog component

## Task Breakdown and Implementation Order

### Phase 1: Language Layer (Foundation)

#### Task 1: Define Edit Dialog State Interfaces

**Prompt:**

You are implementing the foundational state interfaces for the Album Edit feature. This task establishes the domain model for edit operations.

**Requirements (BDD Style):**

```
GIVEN an album editing feature is needed
WHEN defining the state model for edit dialogs
THEN the state should support both name editing and date editing modes
AND the state should track loading states and errors
AND the state should maintain current values being edited
AND the state should integrate with the existing CatalogViewerState
```

**Implementation Details:**

- Create file: `web/src/core/catalog/language/edit-dialog-state.ts`
- Export interfaces: `EditNameDialogState`, `EditDatesDialogState`, `EditDialogState`
- Export utility functions: `isEditNameDialog()`, `isEditDatesDialog()`
- Update `CatalogViewerState` in `catalog-state.ts` to include `editDialog?: EditDialogState`

**References:**

- `web/src/core/catalog/language/catalog-state.ts` - for existing state patterns and AlbumId interface
- `docs/feature_edit_album_claude_sonnet_4.md` - for understanding edit dialog requirements

**Relevant Principles:**
> The state model is common across the slice of the application, it should not be based on the REST model nor a specific UI component, it should be designed to
> represent the best the business.
>
> Always use explicit types or interfaces, never use `any`
>
> Properties that change together must be within the same state

**TDD Approach:**
Write tests first that validate the interface contracts and utility functions before implementing the interfaces.

---

#### Task 2: Create Edit Dialog Selectors

**Prompt:**

You are implementing selectors for the Album Edit dialog states to provide clean access patterns for React components.

**Requirements (BDD Style):**

```
GIVEN edit dialog states are defined
WHEN components need to access edit dialog data
THEN selectors should provide type-safe access to edit dialog properties
AND selectors should handle undefined states gracefully
AND selectors should extract only the data needed by specific components
```

**Implementation Details:**

- Create file: `web/src/core/catalog/language/edit-dialog-selectors.ts`
- Export functions: `selectEditNameDialog()`, `selectEditDatesDialog()`, `selectIsEditDialogOpen()`, `selectEditDialogType()`
- Each selector takes `CatalogViewerState` and returns the appropriate subset

**References:**

- Previous Task 1: `EditDialogState` interfaces and utility functions
- `web/src/core/catalog/language/catalog-state.ts` - for CatalogViewerState structure

**Relevant Principles:**
> Selector function can be provided as part of the model to extract a set of properties used by a component like a Dialog.
>
> Always use explicit types or interfaces, never use `any`

**TDD Approach:**
Write tests that verify selectors return correct data for various state configurations before implementing the selector logic.

---

### Phase 2: Action Layer (State Management)

#### Task 3: Implement Edit Dialog Open Actions

**Prompt:**

You are implementing actions to open edit dialogs for album name and date editing.

**Requirements (BDD Style):**

```
GIVEN a user wants to edit an album
WHEN they choose to edit the album name
THEN an action should open the name edit dialog with current album data
WHEN they choose to edit album dates  
THEN an action should open the dates edit dialog with current album data
AND the actions should follow the past tense event naming convention
```

**Implementation Details:**

- Create file: `web/src/core/catalog/edit-dialog/action-editNameDialogOpened.ts`
- Export interface: `EditNameDialogOpened`
- Export function: `editNameDialogOpened(album: Album): EditNameDialogOpened`
- Export function: `reduceEditNameDialogOpened(state: CatalogViewerState, action: EditNameDialogOpened): CatalogViewerState`
- Export function: `editNameDialogOpenedReducerRegistration(handlers: any)`

- Create file: `web/src/core/catalog/edit-dialog/action-editDatesDialogOpened.ts`
- Export interface: `EditDatesDialogOpened`
- Export function: `editDatesDialogOpened(album: Album): EditDatesDialogOpened`
- Export function: `reduceEditDatesDialogOpened(state: CatalogViewerState, action: EditDatesDialogOpened): CatalogViewerState`
- Export function: `editDatesDialogOpenedReducerRegistration(handlers: any)`

**References:**

- Task 1: `EditDialogState` interfaces from `edit-dialog-state.ts`
- `web/src/core/catalog/language/catalog-state.ts` - for Album interface and CatalogViewerState

**Relevant Principles:**
> The Action naming convention is the Past Tense Event Naming and is represented in camelCase
>
> The reducer function that takes two parameters: the current state, and the action interface. Make sure the types are explicit.
>
> TDD must be used: no behaviour should be implemented if there is no test to validate it

**TDD Approach:**
Write tests that verify the reducers correctly set the editDialog state with the provided album data before implementing the reducer logic.

---

#### Task 4: Implement Edit Dialog Close and Error Actions

**Prompt:**

You are implementing actions to close edit dialogs and handle edit operation errors.

**Requirements (BDD Style):**

```
GIVEN an edit dialog is open
WHEN the user cancels or completes an edit operation
THEN an action should close the edit dialog and clear the state
WHEN an edit operation fails
THEN an action should set the error state while keeping the dialog open
AND the actions should follow the past tense event naming convention
```

**Implementation Details:**

- Create file: `web/src/core/catalog/edit-dialog/action-editDialogClosed.ts`
- Export interface: `EditDialogClosed`
- Export function: `editDialogClosed(): EditDialogClosed`
- Export function: `reduceEditDialogClosed(state: CatalogViewerState, action: EditDialogClosed): CatalogViewerState`
- Export function: `editDialogClosedReducerRegistration(handlers: any)`

- Create file: `web/src/core/catalog/edit-dialog/action-editErrorOccurred.ts`
- Export interface: `EditErrorOccurred`
- Export function: `editErrorOccurred(message: string): EditErrorOccurred`
- Export function: `reduceEditErrorOccurred(state: CatalogViewerState, action: EditErrorOccurred): CatalogViewerState`
- Export function: `editErrorOccurredReducerRegistration(handlers: any)`

**References:**

- Task 1: `EditDialogState` interfaces from `edit-dialog-state.ts`
- `web/src/core/catalog/language/catalog-state.ts` - for CatalogViewerState

**Relevant Principles:**
> if the interface has no property other than type: no parameter
> if the interface has a single property on top of type: parameter is that property
>
> The reducer function returns the updated state, same type as the first parameter

**TDD Approach:**
Write tests that verify the close action clears editDialog state and error action sets error message while preserving dialog data.

---

#### Task 5: Implement Album Update Success Actions

**Prompt:**

You are implementing actions to handle successful album updates for both name/folderName changes and date changes.

**Requirements (BDD Style):**

```
GIVEN an album name edit is successful
WHEN the API returns success
THEN an action should update the album in state and close the dialog
AND if folderName changed, it should include redirect information
GIVEN an album dates edit is successful
WHEN the API returns success
THEN an action should trigger a full refresh of albums and medias
AND close the edit dialog
```

**Implementation Details:**

- Create file: `web/src/core/catalog/edit-dialog/action-albumNameUpdated.ts`
- Export interface: `AlbumNameUpdated`
- Export function: `albumNameUpdated(props: Omit<AlbumNameUpdated, "type">): AlbumNameUpdated`
- Include properties: `updatedAlbum: Album`, `redirectTo?: AlbumId`
- Export function: `reduceAlbumNameUpdated(state: CatalogViewerState, action: AlbumNameUpdated): CatalogViewerState`
- Export function: `albumNameUpdatedReducerRegistration(handlers: any)`

- Create file: `web/src/core/catalog/edit-dialog/action-albumDatesUpdated.ts`
- Export interface: `AlbumDatesUpdated`
- Export function: `albumDatesUpdated(): AlbumDatesUpdated`
- Export function: `reduceAlbumDatesUpdated(state: CatalogViewerState, action: AlbumDatesUpdated): CatalogViewerState`
- Export function: `albumDatesUpdatedReducerRegistration(handlers: any)`

**References:**

- Task 1: `EditDialogState` interfaces
- `web/src/core/catalog/language/catalog-state.ts` - for Album, AlbumId, and RedirectToAlbumIdAction pattern
- `web/src/core/catalog/language/utils-albumIdEquals.ts` - for album comparison

**Relevant Principles:**
> if the interface has more properties, it takes a single argument of the type `Omit<ActionName, "type">`
>
> Properties that change together must be within the same state

**TDD Approach:**
Write tests that verify album updates in state arrays and proper redirect/refresh behavior before implementing reducer logic.

---

### Phase 3: Thunk Layer (Business Logic)

#### Task 6: Implement Edit Album Name Thunk

**Prompt:**

You are implementing the business logic for editing album names and folderNames.

**Requirements (BDD Style):**

```
GIVEN a user wants to edit an album name
WHEN they submit the edit name form
THEN the system should call the API to update the album
AND dispatch loading states during the operation
AND on success, update the album in state and close dialog
AND if folderName changed, include redirect information
AND on failure, display error message while keeping dialog open
```

**Implementation Details:**

- Create file: `web/src/core/catalog/edit-dialog/thunk-editAlbumName.ts`
- Export interface: `EditAlbumNamePort`
- Export function: `editAlbumNameThunk(dispatch, editAlbumNamePort, request)`
- Export interface: `EditAlbumNameRequest` with properties: `albumId: AlbumId`, `name: string`, `folderName?: string`
- Export const: `editAlbumNameDeclaration: ThunkDeclaration<...>`

**References:**

- Tasks 3-5: Edit dialog actions for loading, success, error, and close states
- `web/src/core/catalog/language/catalog-state.ts` - for Album and AlbumId interfaces
- `docs/feature_edit_album_claude_sonnet_4.md` - for understanding name/folderName edit requirements

**Relevant Principles:**
> A thunk implement one and only one responsibility: single-responsibility principle
>
> Thunks are pure-business logic functions handling the user commands triggered by the view
>
> The thunk functions first argument is a dispatch function accepting the specific action interface type(s)

**TDD Approach:**
Write tests using fakes that verify the thunk calls the port correctly, dispatches appropriate actions for success/error scenarios, and handles folderName
changes.

---

#### Task 7: Implement Edit Album Dates Thunk

**Prompt:**

You are implementing the business logic for editing album start and end dates.

**Requirements (BDD Style):**

```
GIVEN a user wants to edit album dates
WHEN they submit the edit dates form
THEN the system should call the API to update the album dates
AND dispatch loading states during the operation
AND on success, trigger a full refresh of albums and medias
AND close the edit dialog
AND on failure, display error message while keeping dialog open
```

**Implementation Details:**

- Create file: `web/src/core/catalog/edit-dialog/thunk-editAlbumDates.ts`
- Export interface: `EditAlbumDatesPort`
- Export function: `editAlbumDatesThunk(dispatch, ports, request)`
- Export interface: `EditAlbumDatesRequest` with properties: `albumId: AlbumId`, `start: Date`, `end: Date`
- Export const: `editAlbumDatesDeclaration: ThunkDeclaration<...>`
- Use multiple ports pattern for API and refresh operations

**References:**

- Tasks 3-5: Edit dialog actions
- Existing album loading thunks for refresh pattern
- `web/src/core/catalog/language/catalog-state.ts` - for Album and AlbumId interfaces

**Relevant Principles:**
> Case 2: Multiple adapters - Pass adapters as single object: `dispatch, {adapter1, adapter2}, stateData, viewData`
>
> Thunks dispatch one or several action(s) during the processing

**TDD Approach:**
Write tests using fakes that verify the thunk updates dates via API and triggers proper refresh sequence of albums and medias.

---

### Phase 4: Adapter Layer (External Integration)

#### Task 8: Implement Edit Album API Adapter Methods

**Prompt:**

You are implementing the REST API adapter methods for album editing operations.

**Requirements (BDD Style):**

```
GIVEN the system needs to update album data via REST API
WHEN editing album name/folderName
THEN the adapter should make a PATCH request to update name fields
WHEN editing album dates
THEN the adapter should make a PATCH request to update date fields
AND the adapter should handle authentication and error responses
AND follow existing API adapter patterns
```

**Implementation Details:**

- Extend existing `CatalogAPIAdapter` class in `web/src/core/catalog/adapters/api/CatalogAPIAdapter.ts`
- Implement `EditAlbumNamePort` interface with method: `updateAlbumName(request: EditAlbumNameRequest): Promise<Album>`
- Implement `EditAlbumDatesPort` interface with method: `updateAlbumDates(request: EditAlbumDatesRequest): Promise<void>`
- Use existing authentication and error handling patterns

**References:**

- Task 6: `EditAlbumNamePort` interface and `EditAlbumNameRequest`
- Task 7: `EditAlbumDatesPort` interface and `EditAlbumDatesRequest`
- Existing `CatalogAPIAdapter` for patterns and authentication setup

**Relevant Principles:**
> The Adapter Implementation is the concrete class that implements the Port interface, abstracting a specific technology or external system
>
> Adapters should be tested independently to verify their contract compliance and data transformation logic

**TDD Approach:**
Write tests that verify the adapter makes correct HTTP requests with proper authentication and handles various response scenarios.

---

### Phase 5: React Component Layer (User Interface)

#### Task 9: Implement Edit Album Button Component

**Prompt:**

You are implementing the Edit Album button component that provides access to album editing functionality.

**Requirements (BDD Style):**

```
GIVEN a user is viewing an album they own
WHEN the album is displayed
THEN an "Edit Album" button should be visible
WHEN the user clicks the "Edit Album" button
THEN a dropdown menu should appear with "Edit Name" and "Edit Dates" options
GIVEN a user is viewing an album they don't own
WHEN the album is displayed
THEN the "Edit Album" button should not be visible
```

**Implementation Details:**

- Create file: `web/src/components/catalog-react/EditAlbumButton.tsx`
- Export component: `EditAlbumButton`
- Props interface: `{album: Album}`
- Use `albumIsOwnedByCurrentUser()` to determine visibility
- Implement dropdown with two options that dispatch appropriate dialog open actions

**References:**

- Task 3: `editNameDialogOpened` and `editDatesDialogOpened` actions
- `web/src/core/catalog/language/catalog-state.ts` - for `albumIsOwnedByCurrentUser()` function
- Existing button and dropdown patterns in the codebase

**Relevant Principles:**
> Every data rendered by the view is coming from the State
>
> User activity, including browser loading and URL change, is triggering the appropriate Thunk

**TDD Approach:**
Write tests that verify the button visibility logic and dropdown option dispatching before implementing the component.

---

#### Task 10: Implement Edit Name Dialog Component

**Prompt:**

You are implementing the dialog component for editing album names and folderNames.

**Requirements (BDD Style):**

```
GIVEN the edit name dialog is open
WHEN the dialog is displayed
THEN it should show current album name in an editable field
AND show folderName field (initially disabled)
AND show a "Force folder name" checkbox
WHEN the user checks "Force folder name"
THEN the folderName field should become enabled for editing
WHEN the user submits the form
THEN it should dispatch the edit album name thunk
AND show loading state during submission
WHEN there is an error
THEN it should display the error message while keeping the dialog open
```

**Implementation Details:**

- Create file: `web/src/components/catalog-react/EditNameDialog.tsx`
- Export component: `EditNameDialog`
- Use `selectEditNameDialog()` selector for state
- Form fields: name (always enabled), folderName (conditionally enabled), force checkbox
- Dispatch `editAlbumNameThunk` on form submission
- Handle loading and error states

**References:**

- Task 2: `selectEditNameDialog()` selector
- Task 6: `editAlbumNameThunk` and `EditAlbumNameRequest`
- Task 4: Error handling actions
- Existing dialog patterns in the codebase

**Relevant Principles:**
> a Selector should be used to access a specific subset of the State properties for composite cases (dialogs or modals)
>
> User activity is triggering the appropriate Thunk (through onClick and useEffect)

**TDD Approach:**
Write tests that verify form behavior, conditional field enabling, and proper thunk dispatching before implementing the component.

---

#### Task 11: Implement Edit Dates Dialog Component

**Prompt:**

You are implementing the dialog component for editing album start and end dates.

**Requirements (BDD Style):**

```
GIVEN the edit dates dialog is open
WHEN the dialog is displayed
THEN it should show current start and end dates in date picker fields
AND reuse existing date picker components from create dialog
WHEN the user submits the form
THEN it should dispatch the edit album dates thunk
AND show loading state during submission
WHEN there is an error
THEN it should display the error message while keeping the dialog open
```

**Implementation Details:**

- Create file: `web/src/components/catalog-react/EditDatesDialog.tsx`
- Export component: `EditDatesDialog`
- Use `selectEditDatesDialog()` selector for state
- Reuse existing date picker components and logic
- Dispatch `editAlbumDatesThunk` on form submission
- Handle loading and error states

**References:**

- Task 2: `selectEditDatesDialog()` selector
- Task 7: `editAlbumDatesThunk` and `EditAlbumDatesRequest`
- Task 4: Error handling actions
- Existing create dialog date picker components

**Relevant Principles:**
> Existing date picker components from create dialog can be reused
>
> a Selector should be used to access a specific subset of the State properties for composite cases (dialogs or modals)

**TDD Approach:**
Write tests that verify date picker integration, form submission, and error handling before implementing the component.

---

#### Task 12: Integrate Edit Components into Main Catalog View

**Prompt:**

You are integrating the edit album functionality into the main catalog view and ensuring proper registration of all actions and thunks.

**Requirements (BDD Style):**

```
GIVEN all edit components and logic are implemented
WHEN integrating into the main catalog view
THEN the EditAlbumButton should be displayed for owned albums
AND edit dialogs should render when their respective states are set
AND all actions should be registered in the catalog reducer
AND all thunks should be registered in the catalog thunks
AND URL redirection should work when folderName changes
```

**Implementation Details:**

- Update main catalog view component to include `EditAlbumButton`
- Add conditional rendering for `EditNameDialog` and `EditDatesDialog`
- Register all edit actions in `catalog-reducer-v2.ts`
- Register edit thunks in `catalog/thunks.ts`
- Update exports in `catalog/language/index.ts`
- Handle redirect behavior for folderName changes

**References:**

- Tasks 9-11: Edit components
- Tasks 3-5: All edit actions for registration
- Tasks 6-7: Edit thunks for registration
- Existing registration patterns in catalog reducer and thunks

**Relevant Principles:**
> Then, the action needs to be registered in the file catalog-reducer-v2.ts
>
> URL redirection required when folderName changes

**TDD Approach:**
Write integration tests that verify the complete edit workflow from button click through successful update and proper state management.

> Tokens: 6.7k sent, 4.5k received. Cost: $0.09 message, $0.09 session.  
> web/src/core/catalog/adapters/api/CatalogAPIAdapter.ts  
> Add file to the chat? (Y)es/(N)o/(D)on't ask again [Yes]: y

Looking at the implementation plan and the `CatalogAPIAdapter.ts` file you've added, I can see we have a good foundation. To provide more complete and accurate
prompts, it would be helpful to add these additional files:

## Files to Add for Better Context:

### 1. Action and Reducer Patterns

- `web/src/core/catalog/album-delete/action-albumDeleted.ts` - to see the complete action pattern implementation
- `web/src/core/catalog/catalog-reducer-v2.ts` - to understand how actions are registered

### 2. Thunk Patterns

- `web/src/core/catalog/album-create/thunk-createAlbum.ts` - to see thunk implementation patterns
- `web/src/core/catalog/thunks.ts` - to understand thunk registration

### 3. Component Patterns

- Any existing dialog component (like create album dialog) - to see dialog implementation patterns
- Main catalog view component - to understand integration patterns

### 4. Test Patterns

- `web/src/core/catalog/tests/test-helper-state.ts` - for test data patterns
- Any existing action test file - to see testing patterns

### 5. State Management

- `web/src/core/catalog/navigation/index.ts` - to see if there are additional state patterns

These files would help me refine the prompts to:

- Match exact naming conventions used in the codebase
- Follow the precise action/reducer registration patterns
- Use consistent thunk declaration patterns
- Align with existing component and testing approaches
- Ensure proper integration with the existing catalog system

Would you like to add some of these files so I can provide more accurate and detailed prompts?

> Tokens: 13k sent, 330 received. Cost: $0.04 message, $0.13 session.  
> web/src/core/catalog/album-delete/action-albumDeleted.ts  
> Add file to the chat? (Y)es/(N)o/(A)ll/(S)kip all/(D)on't ask again [Yes]: y  
> web/src/core/catalog/navigation/index.ts  
> Add file to the chat? (Y)es/(N)o/(A)ll/(S)kip all/(D)on't ask again [Yes]: y  
> web/src/core/catalog/tests/test-helper-state.ts  
> Add file to the chat? (Y)es/(N)o/(A)ll/(S)kip all/(D)on't ask again [Yes]: y  
> web/src/core/catalog/thunks.ts  
> Add file to the chat? (Y)es/(N)o/(A)ll/(S)kip all/(D)on't ask again [Yes]: y

Perfect! With these additional files, I have a much clearer picture of the patterns used in the codebase. Looking at the files, I can see:

1. **Action patterns** from `action-albumDeleted.ts` - shows the exact structure for interfaces, factory functions, reducers, and registration
2. **Thunk registration** from `thunks.ts` - shows how thunks are exported and registered
3. **Test patterns** from `test-helper-state.ts` - shows the test data structures and naming conventions
4. **Navigation exports** - shows how features export their functionality

## Additional Files That Would Be Helpful:

### 1. Complete Thunk Implementation Example

- A complete thunk file (like `web/src/core/catalog/album-create/thunk-createAlbum.ts` or `web/src/core/catalog/album-delete/thunk-deleteAlbum.ts`) to see:
    - The exact `ThunkDeclaration` pattern
    - Port interface definitions
    - Factory implementation with selector and binding
    - How multiple adapters are handled

### 2. Catalog Reducer Registration

- `web/src/core/catalog/catalog-reducer-v2.ts` to see:
    - How actions are registered in the `catalogActions` object
    - How reducer registrations are added
    - The exact pattern for the `CatalogViewerAction` union type

### 3. Component Pattern (Optional but Helpful)

- Any existing dialog component to understand:
    - How dialogs integrate with the state
    - Component structure and patterns
    - How thunks are called from components

The thunk file would be particularly valuable since I can see from `thunks.ts` that there are declarations like `createAlbumDeclaration` and
`deleteAlbumDeclaration`, but I need to see the actual implementation pattern to write accurate prompts.

Would you like to add one of the thunk implementation files and the catalog reducer file?

> Tokens: 15k sent, 381 received. Cost: $0.05 message, $0.18 session.

#### /drop   web/src/core/catalog/adapters/api/CatalogAPIAdapter.ts web/src/core/catalog/tests/test-helper-state.ts web/src/core/catalog/thunks.ts

> Removed web/src/core/catalog/adapters/api/CatalogAPIAdapter.ts from the chat  
> Removed web/src/core/catalog/tests/test-helper-state.ts from the chat  
> Removed web/src/core/catalog/thunks.ts from the chat

#### /edit  
