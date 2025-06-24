import {editDatesDialogEndDateUpdated} from "./action-editDatesDialogEndDateUpdated";
import {CatalogViewerState} from "../language";
import {loadedStateWithTwoAlbums} from "../tests/test-helper-state";
import {DEFAULT_EDIT_DATES_DIALOG_SELECTION, editDatesDialogSelector} from "./selector-editDatesDialogSelector";

describe("action:editDatesDialogEndDateUpdated", () => {
    it("updates the end date when dialog is open", () => {
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

        const newEndDate = new Date("2023-07-25T00:00:00");
        const action = editDatesDialogEndDateUpdated(newEndDate);
        const got = action.reducer(state, action);

        expect(editDatesDialogSelector(got)).toEqual({
            isOpen: true,
            albumName: "Summer Trip",
            startDate: new Date("2023-07-01T00:00:00"),
            endDate: newEndDate,
            startAtDayStart: true,
            endAtDayEnd: true,
            isLoading: false,
            errorCode: undefined,
            dateRangeError: undefined,
            isSaveEnabled: true,
        });
    });

    it("does nothing when dialog is closed", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            editDatesDialog: undefined,
        };

        const newEndDate = new Date("2023-07-25T00:00:00");
        const action = editDatesDialogEndDateUpdated(newEndDate);
        const got = action.reducer(state, action);

        expect(got).toBe(state);
    });

    it("creates invalid date range when end date is set before start date", () => {
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

        const newEndDate = new Date("2023-01-05T00:00:00Z");
        const action = editDatesDialogEndDateUpdated(newEndDate);
        const got = action.reducer(state, action);
        const selection = editDatesDialogSelector(got);

        expect(selection).toEqual({
            ...DEFAULT_EDIT_DATES_DIALOG_SELECTION,
            isOpen: true,
            albumName: "Summer Trip",
            startDate: new Date("2023-01-10T00:00:00Z"),
            endDate: newEndDate,
            startAtDayStart: true,
            endAtDayEnd: true,
            dateRangeError: "The end date cannot be before the start date",
            isSaveEnabled: false,
        });
    });

    it("corrects invalid date range when end date is set after start date", () => {
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

        const newEndDate = new Date("2023-01-30T00:00:00Z");
        const action = editDatesDialogEndDateUpdated(newEndDate);
        const got = action.reducer(state, action);
        const selection = editDatesDialogSelector(got);

        expect(selection).toEqual({
            ...DEFAULT_EDIT_DATES_DIALOG_SELECTION,
            isOpen: true,
            albumName: "Summer Trip",
            startDate: new Date("2023-01-25T00:00:00Z"),
            endDate: newEndDate,
            startAtDayStart: true,
            endAtDayEnd: true,
        });
    });
});
