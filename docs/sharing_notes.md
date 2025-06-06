Sharing Dialog Re-Design
=======================================

Plan the plan
---------------------------------------

I'm thinking about a different approach based on generating the suggestions when the dialog is open, then keeping them up to date on grant and revoke.

The plan will be layed out as follow:

1. update the openSharingModal and the sharingDialogSelector to implement the following cases:

  - initialises the list of suggestions based on the list of all known users (sharedWith of all albums) minus the user for which the album is already shared
  - create an empty list of suggestions if all known users already have access to this album
  - create an empty list of suggestions if there is no known users
  - suggestions must be sorted by popularity (how many albums are shared to this user), then by name

2. update the granting thunk to implement the following cases:

  - removes the newly granted email from the list of suggestion if it was present, and uses it details if a picture is present and its name is different from
    the email
  - restore sharedWith list and suggestion list when the request failed and the email was already in the suggestion list, dispatch the action
    `sharingModalError` and throw a error
  - restore sharedWith list and suggestion list with the new email when the request failed and the email was not in the suggestion list, dispatch the action
    `sharingModalErr ` and throw an error

3. update the revoke chunk to implement the following cases:

  - removes the revoked email from the sharedWith list and add it to the suggestions list
  - restore the sharedWith and suggestion list if the request failed

---

Ignore the plan you describe earlier, only work with the plan I just wrote.

Provide a series of prompts for a code-generation LLM that will implement each step in a test-driven manner. Prioritize best practices, incremental progress,
and early tes
, ensuring no big jumps in complexity at any stage. Make sure that each prompt builds on the previous prompts, and ends with wiring things together. There
should be no hang
or orphaned code that isn't integrated into a previous step.

Make sure and separate each prompt section. Use markdown. Each prompt should be tagged as text using code tags. The goal is to output prompts, but context, etc
is important as well.

Each prompt should contain at least:

- give an abstract of the change, not a list of things to do
- the list of the files where the changes must be done (including the tests!) and the files necessary of the context list the principles and definition of the
  states
- the list of test cases to implement, the one I wrote

GPT prompts
---------------------------------------

### **Prompt 1: Generate Suggestions When Opening the Share Dialog**

```text
**Abstract:**
Implement logic to generate the list of sharing suggestions when the share dialog is opened. Suggestions are derived from all known users (i.e., users present in the `sharedWit
h` list of all albums), excluding users who already have access to the current album. The suggestions list should be empty if all known users already have access or if there ar
e no known users. Suggestions must be sorted by popularity (number of albums shared with the user), then by name.

**Files to update:**
- `web/src/core/catalog/domain/actions/sharing-openSharingModalAction.ts`
- `web/src/core/catalog/domain/actions/sharing-openSharingModalAction.test.ts`
- `web/src/core/catalog/domain/actions/selector-sharingDialogSelector.ts`
- `web/src/core/catalog/domain/catalog-state.ts` (ensure state shape supports suggestions)
- (readonly) `web/src/core/catalog/domain/tests/test-helper-state.ts`
- (readonly) `web/src/core/catalog/domain/utils-albumIdEquals.ts`

**Principles and context:**
- Follow the state and thunk principles in `docs/principles_actions.md`.
- Use explicit types and full-object assertions in tests.
- Suggestions are a list of `UserDetails`.

**Test cases to implement:**
1. Initializes the list of suggestions based on all known users (from all albums' `sharedWith`), minus users already in the current album's `sharedWith`.
2. Creates an empty list of suggestions if all known users already have access to this album.
3. Creates an empty list of suggestions if there are no known users.
4. Suggestions are sorted by popularity (descending), then by name (ascending).
```

---

### **Prompt 2: Update Granting Thunk to Maintain Suggestions and Handle Errors**

```text
**Abstract:**
Update the granting thunk so that when access is granted, the granted email is removed from the suggestions list (if present), and the user's details are used if available (e.g., picture, name). If the grant request fails, restore the previous `sharedWith` and suggestions lists. If the email was already in suggestions, restore it; if not, add it back with the new details. Always dispatch the `sharingModalError` action and throw an error on failure.

**Files to update:**
- `web/src/core/catalog/thunks/share-grantAlbumSharing.ts` (thunk logic)
- `web/src/core/catalog/thunks/share-grantAlbumSharing.test.ts` (add/extend tests for new logic)
- `web/src/core/catalog/domain/catalog-state.ts`
- (readonly) `web/src/core/catalog/domain/actions/sharing-addSharingAction.ts`
- (readonly) `web/src/core/catalog/domain/actions/sharing-sharingModalErrorAction.ts`
- (readonly) `web/src/core/catalog/domain/tests/test-helper-state.ts`
- (readonly) `web/src/core/catalog/domain/utils-albumIdEquals.ts`

**Principles and context:**
- Follow the thunk and action principles in `docs/principles_thunks.md`.
- Use explicit types and full-object assertions in tests.
- Ensure no orphaned state or actions.

**Test cases to implement:**
1. Removes the newly granted email from the suggestions list if it was present, using its details if available.
2. Restores `sharedWith` and suggestions lists when the request fails and the email was already in the suggestions list; dispatches `sharingModalError` and throws.
3. Restores `sharedWith` and suggestions lists with the new email when the request fails and the email was not in the suggestions list; dispatches `sharingModalError` and throws.
```

---

### **Prompt 3: Update Revoke Thunk to Maintain Suggestions and Handle Errors**

```text
**Abstract:**
Update the revoke thunk so that when access is revoked, the revoked email is removed from the `sharedWith` list and added to the suggestions list. If the revoke request fails,
restore the previous `sharedWith` and suggestions lists.

**Files to update:**
- `web/src/core/catalog/thunks/share-revokeAlbumSharing.ts` (thunk logic)
- `web/src/core/catalog/domain/actions/` (any actions for updating suggestions/sharedWith)
- `web/src/core/catalog/domain/catalog-state.ts` (ensure state shape supports this logic)
- Tests:
  - `web/src/core/catalog/thunks/share-revokeAlbumSharing.test.ts` (add/extend tests for new logic)
  - Any test helpers or mock data files as needed

**Principles and context:**
- Follow the thunk and action principles in `docs/principles_thunks.md` and `docs/principles_actions.md`.
- Use explicit types and full-object assertions in tests.
- Ensure no orphaned state or actions.

**Test cases to implement:**
1. Removes the revoked email from the `sharedWith` list and adds it to the suggestions list.
2. Restores the `sharedWith` and suggestions lists if the request fails.
```

---

### **Prompt 4: Wire Everything Together and Ensure Integration**

```text
**Abstract:**
Integrate all changes so that the `ShareDialog` receives the correct `sharedWith`, `suggestions`, and error props from the selector, and that all actions and thunks update the
state as expected. Ensure all tests pass and that the dialog behaves correctly in all scenarios.

**Files to update:**
- `web/src/core/catalog/domain/actions/selector-sharingDialogSelector.ts` (finalize selector logic)
- `web/src/core/catalog/domain/catalog-state.ts` (finalize state shape)
- `web/src/pages/authenticated/albums/ShareDialog/index.tsx` (ensure props are correctly wired)
- All relevant test files (ensure all test cases are covered and pass)

**Principles and context:**
- Follow the integration and test coverage principles in `docs/principles_thunks.md` and `docs/principles_actions.md`.
- Use explicit types and full-object assertions in tests.
- Ensure no orphaned or unused code.

**Test cases to implement:**
- The dialog receives and displays the correct suggestions, sharedWith, and error props in all scenarios.
- All state transitions and error handling work as expected.
- All previous test cases for suggestions, grant, and revoke are covered and pass.
```

(Draft) GPT initial Proposition:
---------------------------------------

**1. Use Chips for Granted Users (like Suggestions)**

- **How it works:**
  Display each user who has access as a Chip, similar to the suggestion chips above.
  Each chip shows the user's avatar, name, and email (in tooltip or subtitle), and includes a "remove" (revoke) icon (e.g., a close or trash icon).
- **Benefits:**
    - Consistent with the suggestion UI.
    - Compact, visually appealing, and easy to scan.
    - Revoking access is a single click.

---

**2. List with Avatars and Actions**

- **How it works:**
  Use a vertical list (e.g., MUI List or Stack) where each item shows the user's avatar, name, and email, with a revoke button (icon) on the right.
- **Benefits:**
    - Still modern and compact.
    - Works well if you want to show more details per user (e.g., when access was granted).

---

**3. Chip Group with Section Title**

- **How it works:**
  Group all granted users' chips under a section title like "Users with access", visually separating them from the suggestion chips.
- **Benefits:**
    - Clear distinction between suggestions and already-granted users.
    - Keeps the UI clean and organized.

---

**Summary:**
The **chip-based approach** (option 1 or 3) is the most consistent and modern, matching your new suggestions UI. It also makes revoking access more intuitive
and visually appealing than a table. If you need to show
extra info, a list with avatars (option 2) is a good compromise.