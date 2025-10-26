# As a user, I want to open and close the "Edit Name" dialog for an album I own

## Acceptance Criteria

```
GIVEN I am viewing an album that I own
WHEN I click the "Edit" button
THEN I see a dropdown menu that includes an "Edit Name" option

GIVEN I am viewing an album that I do not own
WHEN I look at the album interface
THEN the "Edit" button is disabled and I cannot access any edit options

GIVEN I am viewing an album that I own
WHEN I click the "Edit" button and select "Edit Name"
THEN a modal dialog opens with the title "Edit Name"
AND the dialog contains a "Cancel" button

GIVEN the rename dialog is open
WHEN I click "Cancel" or press Escape or click outside the dialog
THEN the dialog closes without making any changes
```

## Out of scope

- Form fields, form validation, and save functionality (handled in subsequent stories)
- Changes to existing edit options like "Edit Dates"
- Permission logic for album ownership (already handled by existing platform)
- Styling or positioning of the dropdown menu beyond adding the new option
