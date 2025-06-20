# As a user, when I select round times for end dates, they are stored exactly as entered for better predictability.

## Acceptance Criteria:

```
GIVEN the "Edit Dates" dialog is open for an album I own
AND I uncheck the "at the end of the day" checkbox for the end date
WHEN I enter a round time like "16:00" into the end time input field
AND I click the "Save" button
THEN the API request to update the album includes the end date with the exact time entered, for example, "YYYY-MM-DDT16:00:00" (without adding 1 minute)

GIVEN the "Edit Dates" dialog is open for an album I own
AND I uncheck the "at the end of the day" checkbox for the end date
WHEN I enter a round time ending in "30" like "09:30" into the end time input field
AND I click the "Save" button
THEN the API request to update the album includes the end date with the exact time entered, for example, "YYYY-MM-DDT09:30:00" (without adding 1 minute)

GIVEN the "Edit Dates" dialog is open for an album I own
AND I uncheck the "at the end of the day" checkbox for the end date
WHEN I enter a precise time like "11:42" or "14:17" into the end time input field
AND I click the "Save" button
THEN the API request to update the album includes the end date with 1 minute added to make it exclusive, for example, "YYYY-MM-DDT11:43:00" or "YYYY-MM-DDT14:18:00"
```

## Definition:
- **Round times**: Times that end in ":00" (full hours) or ":30" (half hours)
- **Precise times**: Times that end in any other minute value (e.g., ":15", ":42", ":07")

## Rationale:
When users select round times like 16:00 or 09:30, they likely mean "until exactly that time" and expect it to be stored as-is. This provides better predictability and simplifies the conversion logic. Only non-round times (like 11:42) are treated as inclusive and get the +1 minute conversion.

## Out of scope:
- Validation of time input format
- UI indicators showing which times will be stored exactly vs. with +1 minute
- Retroactive changes to existing albums with round end times
