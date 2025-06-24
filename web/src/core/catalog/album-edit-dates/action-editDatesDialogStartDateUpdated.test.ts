import {editDatesDialogStartDateUpdated} from "./action-editDatesDialogStartDateUpdated";
import {CatalogViewerState} from "../language";
import {loadedStateWithTwoAlbums} from "../tests/test-helper-state";
import {DEFAULT_EDIT_DATES_DIALOG_SELECTION, editDatesDialogSelector} from "./selector-editDatesDialogSelector";

describe("action:editDatesDialogStartDateUpdated", () => {
    it("updates the start date when dialog is open", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            editDatesDialog: {
                albumId: {owner: "myself", folderName: "summer-trip"},
                albumName: "Summer Trip",
                startDate: new Date("2023-07-01T00:00:00"),
                endDate: new Date("2023-08-01T00:00:00"),
                isLoading: false,
                startAtDayStart: true,
                endAtDayEnd: true,
            },
        };

        const newStartDate = new Date("2023-07-15T00:00:00");
        const action = editDatesDialogStartDateUpdated(newStartDate);
        const got = action.reducer(state, action);
        const selection = editDatesDialogSelector(got);

        expect(selection).toEqual({
            ...DEFAULT_EDIT_DATES_DIALOG_SELECTION,
            isOpen: true,
            albumName: "Summer Trip",
            startDate: newStartDate,
            endDate: new Date("2023-08-01T00:00:00"),
            startAtDayStart: true,
            endAtDayEnd: true,
        });
    });

    it("does nothing when dialog is closed", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            editDatesDialog: undefined,
        };

        const newStartDate = new Date("2023-07-15T00:00:00");
        const action = editDatesDialogStartDateUpdated(newStartDate);
        const got = action.reducer(state, action);

        expect(got).toBe(state);
    });

    it("clears error when updating start date", () => {
        const stateWithError: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            editDatesDialog: {
                albumId: {owner: "myself", folderName: "summer-trip"},
                albumName: "Summer Trip",
                startDate: new Date("2023-07-01T00:00:00"),
                endDate: new Date("2023-08-01T00:00:00"),
                isLoading: false,
                error: "Previous error message",
                startAtDayStart: true,
                endAtDayEnd: true,
            },
        };

        const newStartDate = new Date("2023-07-15T00:00:00");
        const action = editDatesDialogStartDateUpdated(newStartDate);
        const got = action.reducer(stateWithError, action);
        const selection = editDatesDialogSelector(got);

        expect(selection).toEqual({
            ...DEFAULT_EDIT_DATES_DIALOG_SELECTION,
            isOpen: true,
            albumName: "Summer Trip",
            startDate: newStartDate,
            endDate: new Date("2023-08-01T00:00:00"),
            startAtDayStart: true,
            endAtDayEnd: true,
        });
    });

    it("creates invalid date range when start date is set after end date", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            editDatesDialog: {
                albumId: {owner: "myself", folderName: "summer-trip"},
                albumName: "Summer Trip",
                startDate: new Date("2023-01-10T00:00:00Z"),
                endDate: new Date("2023-01-20T00:00:00Z"),
                isLoading: false,
                startAtDayStart: true,
                endAtDayEnd: true,
            },
        };

        const newStartDate = new Date("2023-01-25T00:00:00Z");
        const action = editDatesDialogStartDateUpdated(newStartDate);
        const got = action.reducer(state, action);
        const selection = editDatesDialogSelector(got);

        expect(selection).toEqual({
            ...DEFAULT_EDIT_DATES_DIALOG_SELECTION,
            isOpen: true,
            albumName: "Summer Trip",
            startDate: newStartDate,
            endDate: new Date("2023-01-20T00:00:00Z"),
            startAtDayStart: true,
            endAtDayEnd: true,
            dateRangeError: "The end date cannot be before the start date",
            isSaveEnabled: false,
        });
    });

    it("corrects invalid date range when start date is set before end date", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            editDatesDialog: {
                albumId: {owner: "myself", folderName: "summer-trip"},
                albumName: "Summer Trip",
                startDate: new Date("2023-01-25T00:00:00Z"),
                endDate: new Date("2023-01-20T00:00:00Z"),
                isLoading: false,
                startAtDayStart: true,
                endAtDayEnd: true,
            },
        };

        const newStartDate = new Date("2023-01-15T00:00:00Z");
        const action = editDatesDialogStartDateUpdated(newStartDate);
        const got = action.reducer(state, action);
        const selection = editDatesDialogSelector(got);

        expect(selection).toEqual({
            ...DEFAULT_EDIT_DATES_DIALOG_SELECTION,
            isOpen: true,
            albumName: "Summer Trip",
            startDate: newStartDate,
            endDate: new Date("2023-01-20T00:00:00Z"),
            startAtDayStart: true,
            endAtDayEnd: true,
        });
    });
});
