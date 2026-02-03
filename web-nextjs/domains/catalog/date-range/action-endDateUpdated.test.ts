import {endDateUpdated} from "./action-endDateUpdated";
import {CatalogViewerState, CreateDialog, EditDatesDialog} from "../language";
import {createDialogPrefilledForMar25, deleteDialogWithOneAlbum, editDatesDialogForJanAlbum, loadedStateWithTwoAlbums} from "../tests/test-helper-state";
import {DateRangeValidation, datesMustBeSetError, validateDateRange} from "./date-helper";

describe("action:endDateUpdated", () => {
    const validDateRange = {
        areDatesValid: true,
        isDateRangeValid: true
    };

    it("updates end date in create dialog and clears date error", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...createDialogPrefilledForMar25,
                endDate: new Date("2025-02-15"),
                error: "AlbumStartAndEndDateMandatoryErr",
            },
        };

        const newEndDate = new Date("2025-03-31");
        const action = endDateUpdated(newEndDate);
        const got = action.reducer(state, action);

        expect(got.dialog).toEqual<CreateDialog>({
            ...createDialogPrefilledForMar25,
            endDate: newEndDate,
            error: undefined,
        });

        expect(validateDateRange(got.dialog as CreateDialog)).toEqual<DateRangeValidation>(validDateRange);
    });

    it("updates end date in edit dialog and clears error", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...editDatesDialogForJanAlbum,
                endDate: new Date("2025-02-15"),
                error: "Some error",
            },
        };

        const newEndDate = new Date("2025-03-31");
        const action = endDateUpdated(newEndDate);
        const got = action.reducer(state, action);

        expect(got.dialog).toEqual<EditDatesDialog>({
            ...editDatesDialogForJanAlbum,
            endDate: newEndDate,
            error: undefined,
        });

        expect(validateDateRange(got.dialog as EditDatesDialog)).toEqual<DateRangeValidation>(validDateRange);
    });

    it("creates invalid date range when end date is before start date", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...createDialogPrefilledForMar25,
                startDate: new Date("2025-02-15"),
            },
        };

        const newEndDate = new Date("2025-02-10");
        const action = endDateUpdated(newEndDate);
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

        const action = endDateUpdated(new Date("2025-03-31"));
        const got = action.reducer(state, action);

        expect(got).toEqual(state);
    });

    it("does nothing when dialog is closed", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: undefined,
        };

        const action = endDateUpdated(new Date("2025-03-31"));
        const got = action.reducer(state, action);

        expect(got).toEqual(state);
    });

    it("makes the range invalid when setting a null value", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...createDialogPrefilledForMar25,
                startDate: new Date("2025-02-01"),
            },
        };

        const action = endDateUpdated(null);
        const got = action.reducer(state, action);

        const validation = validateDateRange(got.dialog as CreateDialog);
        expect(validation).toEqual(datesMustBeSetError);
    });
});
