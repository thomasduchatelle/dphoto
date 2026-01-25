# Task 6: Implement Edit Dates Dialog Component

You are implementing the React component that renders the edit dates dialog and handles user interactions.

## Requirements (BDD Style)

```
GIVEN the edit dates dialog is closed
WHEN the component renders
THEN the dialog is not visible

GIVEN the edit dates dialog is open with an album
WHEN the component renders
THEN the dialog displays the album name
AND the dialog has a close button that calls the close thunk when clicked
```

## Implementation Details

**Create file:** `web/src/components/catalog-react/EditDatesDialog.tsx`

**Export the following:**
- Component: `EditDatesDialog`

**Component specifications:**
- Use `selectEditDatesDialog` selector to get dialog state
- Conditionally render dialog based on `isOpen` property
- Display the album name when dialog is open
- Include a close button that calls the `closeEditDatesDialog` thunk
- Follow existing dialog patterns in the codebase

## References

- Task 1: `selectEditDatesDialog` from `edit-dates-dialog-state.ts`
- Task 5: `closeEditDatesDialog` thunk
- Existing dialog components in the codebase for UI patterns
- `web/src/core/catalog/tests/test-helper-state.ts` - for test data

## Relevant Principles

> Every data rendered by the view is coming from the State

> a Selector should be used to access a specific subset of the State properties for composite cases (dialogs or modals)

> User activity is triggering the appropriate Thunk (through onClick and useEffect)

> TDD must be used: no behaviour should be implemented if there is no test to validate it

## TDD Approach

Write tests first that validate:
1. The component renders nothing when dialog is closed
2. The component renders the dialog when isOpen is true
3. The component displays the correct album name
4. Clicking the close button calls the closeEditDatesDialog thunk
5. The component uses the selector correctly to get state

Use React Testing Library for component tests and mock the thunk calls to verify they are called correctly.
