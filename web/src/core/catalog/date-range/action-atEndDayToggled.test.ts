import {atEndDayToggled} from "./action-atEndDayToggled";
import {CatalogViewerState, CreateDialog, DateRangeState} from "../language";
import {createDialogPrefilledForMar25, deleteDialogWithOneAlbum, loadedStateWithTwoAlbums} from "../tests/test-helper-state";
import {dateRangeIsValid, DateRangeValidation, validateDateRange} from "./date-helper";


describe("action:atEndDayToggled", () => {
    it("set endAtDayEnd=true, and keep endDate Date (time doesn't matter), when at day end is checked ", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: createDialogPrefilledForMar25,
        };

        const action = atEndDayToggled(true);
        const got = action.reducer(state, action);

        expect(got.dialog).toEqual<CreateDialog>({
            ...createDialogPrefilledForMar25,
            endDate: expect.any(Date),
            endAtDayEnd: true,
        });
        expect((got.dialog as CreateDialog)!.endDate!.toISOString().startsWith("2025-03-31")).toBe(true);

        expect(validateDateRange(got.dialog as DateRangeState)).toEqual<DateRangeValidation>(dateRangeIsValid)
    });

    it("uncheck 'at day end' set endAtDayEnd=false, and set the time to 23:59", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...createDialogPrefilledForMar25,
                endAtDayEnd: true,
                endDate: new Date("2025-03-31T10:00:00"),
            },
        };

        const action = atEndDayToggled(false);
        const got = action.reducer(state, action);

        expect(got.dialog).toEqual<CreateDialog>({
            ...createDialogPrefilledForMar25,
            endDate: new Date("2025-03-31T23:59:00.000Z"),
            endAtDayEnd: false,
        });

        expect(validateDateRange(got.dialog as DateRangeState)).toEqual<DateRangeValidation>(dateRangeIsValid)
    });

    it("does nothing when dialog is not a date range dialog", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: undefined,
        };

        const action = atEndDayToggled(true);
        const got = action.reducer(state, action);

        expect(got).toBe(state);
    });

    it("does nothing if there is no end date on the dialog", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...createDialogPrefilledForMar25,
                endDate: null, // No end date
            },
        };

        const action = atEndDayToggled(false);
        const got = action.reducer(state, action);

        expect(got).toBe(state);
    });

    it("does nothing if the dialog is not a DateRangeStatus", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: deleteDialogWithOneAlbum,
        };

        const action = atEndDayToggled(false);
        const got = action.reducer(state, action);

        expect(got).toBe(state);
    });

    it("clears error when making a change", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...createDialogPrefilledForMar25,
                error: "Some error message",
            },
        };

        const action = atEndDayToggled(true);
        const got = action.reducer(state, action);

        expect(got.dialog).toEqual<CreateDialog>({
            ...createDialogPrefilledForMar25,
            error: undefined,
            endDate: expect.any(Date), // Added to match the actual state after the action
            endAtDayEnd: true, // Added to match the actual state after the action
        });
    });
});
