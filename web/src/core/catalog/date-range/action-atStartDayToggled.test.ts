import {atStartDayToggled} from "./action-atStartDayToggled";
import {CatalogViewerState, CreateDialog} from "../language";
import {createDialogPrefilledForMar25, loadedStateWithTwoAlbums} from "../tests/test-helper-state";
import {DateRangeValidation, validateDateRange} from "./date-helper";

describe("action:atStartDayToggled", () => {
    const validDateRange = {
        areDatesValid: true,
        isDateRangeValid: true
    };

    it("set startAtDayStart=true, and keep startDate Date (time doesn't matter), when at day start is checked ", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...createDialogPrefilledForMar25,
                startDate: new Date("2025-02-01T10:00:00"),
                startAtDayStart: false,
            },
        };

        const action = atStartDayToggled(true);
        const got = action.reducer(state, action);

        expect(got.dialog).toEqual<CreateDialog>({
            ...createDialogPrefilledForMar25,
            startDate: expect.any(Date),
            startAtDayStart: true,
        });
        expect((got.dialog as CreateDialog)!.startDate!.toISOString().startsWith("2025-02-01")).toBe(true);

        expect(validateDateRange(got.dialog as CreateDialog)).toEqual<DateRangeValidation>(validDateRange)
    });

    it("uncheck 'at day start' set startAtDayStart=false, and set the time to 00:00", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...createDialogPrefilledForMar25,
                startDate: new Date("2025-02-01T10:00:00"),
                startAtDayStart: true,
            },
        };

        const action = atStartDayToggled(false);
        const got = action.reducer(state, action);

        expect(got.dialog).toEqual<CreateDialog>({
            ...createDialogPrefilledForMar25,
            startDate: new Date("2025-02-01T00:00:00.000Z"),
            startAtDayStart: false,
        });

        expect(validateDateRange(got.dialog as CreateDialog)).toEqual<DateRangeValidation>(validDateRange)
    });

    it("clears the errors when updated", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...createDialogPrefilledForMar25,
                error: "Some error"
            },
        };

        const action = atStartDayToggled(false);
        const got = action.reducer(state, action);

        expect(got.dialog).toEqual<CreateDialog>({
            ...createDialogPrefilledForMar25,
            startAtDayStart: false,
            error: undefined,
        });

        expect(validateDateRange(got.dialog as CreateDialog)).toEqual<DateRangeValidation>(validDateRange)
    });

    it("does nothing when dialog is not a date range dialog", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: undefined,
        };

        const action = atStartDayToggled(true);
        const got = action.reducer(state, action);

        expect(got).toBe(state);
    });

    it("does nothing if there is no start date on the dialog", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...createDialogPrefilledForMar25,
                startDate: null, // No start date
            },
        };

        const action = atStartDayToggled(false);
        const got = action.reducer(state, action);

        expect(got).toBe(state);
    });
});
