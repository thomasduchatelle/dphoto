import {startDateUpdated} from "./action-startDateUpdated";
import {CatalogViewerState, CreateDialog, EditDatesDialog} from "../language";
import {createDialogPrefilledForMar25, deleteDialogWithOneAlbum, editDatesDialogForJanAlbum, loadedStateWithTwoAlbums} from "../tests/test-helper-state";
import {dateRangeIsValid, DateRangeValidation, datesMustBeSetError, validateDateRange} from "./date-helper";

describe("action:startDateUpdated", () => {
    it("updates start date in create dialog and clears date error", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...createDialogPrefilledForMar25,
                error: "AlbumStartAndEndDateMandatoryErr",
            },
        };

        const newStartDate = new Date("2025-02-15");
        const action = startDateUpdated(newStartDate);
        const got = action.reducer(state, action);

        expect(got.dialog).toEqual<CreateDialog>({
            ...createDialogPrefilledForMar25,
            startDate: newStartDate,
            error: undefined,
        });

        expect(validateDateRange(got.dialog as CreateDialog)).toEqual<DateRangeValidation>(dateRangeIsValid);
    });

    it("updates start date in edit dialog and clears error", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...editDatesDialogForJanAlbum,
                error: "Some error",
            },
        };

        const newStartDate = new Date("2025-01-15");
        const action = startDateUpdated(newStartDate);
        const got = action.reducer(state, action);

        expect(got.dialog).toEqual<EditDatesDialog>({
            ...editDatesDialogForJanAlbum,
            startDate: newStartDate,
            error: undefined,
        });

        expect(validateDateRange(got.dialog as EditDatesDialog)).toEqual<DateRangeValidation>(dateRangeIsValid);
    });

    it("creates invalid date range when start date is after end date", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...createDialogPrefilledForMar25,
                endDate: new Date("2025-02-15"),
            },
        };

        const newStartDate = new Date("2025-02-20");
        const action = startDateUpdated(newStartDate);
        const got = action.reducer(state, action);

        const validation = validateDateRange(got.dialog as CreateDialog);
        expect(validation).toEqual({
            areDatesValid: true,
            isDateRangeValid: false,
            dateRangeError: "The end date cannot be before the start date"
        });
    });

    it("does nothing when dialog is not a date range dialog", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: deleteDialogWithOneAlbum,
        };

        const action = startDateUpdated(new Date("2025-02-15"));
        const got = action.reducer(state, action);

        expect(got).toEqual(state);
    });

    it("does nothing when dialog is closed", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: undefined,
        };

        const action = startDateUpdated(new Date("2025-02-15"));
        const got = action.reducer(state, action);

        expect(got).toEqual(state);
    });

    it("makes the range invalid when setting a null value", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...createDialogPrefilledForMar25,
                endDate: new Date("2025-03-31"),
            },
        };

        const action = startDateUpdated(null);
        const got = action.reducer(state, action);

        const validation = validateDateRange(got.dialog as CreateDialog);
        expect(validation).toEqual(datesMustBeSetError);
    });
});
