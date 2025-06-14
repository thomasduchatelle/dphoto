# Task 3: UI Component - Edit Dates Dialog

## Introduction
You are implementing the **React component** `EditDatesDialog` for displaying the edit dates dialog.

## Requirements

**GIVEN** an `EditDatesDialogSelection` with `open: false`
**WHEN** the `EditDatesDialog` component is rendered
**THEN** no dialog is displayed

**GIVEN** an `EditDatesDialogSelection` with `open: true` and album data
**WHEN** the `EditDatesDialog` component is rendered
**THEN** a modal dialog is displayed showing the album name (read-only) and the current start and end dates in a user-friendly format

## Implementation Details

**TDD Principle**: Implement the tests first, then implement the code **the simplest and most readable possible**: no behaviour should be implemented if it is not required by one test.

**Folder Structure**: Create the component in `web/src/components/catalog-react/` following the naming convention `EditDatesDialog.tsx`.

**Files to Edit**:
- Main catalog viewer component to integrate the new dialog

**Recommendations**:
- Use existing modal/dialog components and styling patterns from the codebase
- Display dates in a user-friendly format (consider using dayjs formatting)
- Album name should be read-only text display
- Component should not handle any actions (this story is display-only)
- Use the selector `editDatesDialogSelector` to get dialog data

## Interface Specification

```typescript
// Component props (derived from selector)
interface EditDatesDialogProps {
    dialogSelection: EditDatesDialogSelection
}

// Component signature
export function EditDatesDialog(props: EditDatesDialogProps): JSX.Element
```

## References

- `web/src/core/catalog/edit-dates/action-editDatesDialogOpened.ts` - For the selector `editDatesDialogSelector` (from Task 1)
- `web/src/core/catalog/language/catalog-state.ts` - For the `EditDatesDialogSelection` interface
- Existing dialog components in the codebase for styling and modal patterns
- `docs/principles_web.md` - For React component testing guidelines
