import {editDatesDialogEndAtDayEndUpdated} from "./action-editDatesDialogEndAtDayEndUpdated";
import {CatalogViewerState} from "../language";
import {loadedStateWithTwoAlbums} from "../tests/test-helper-state";
import {editDatesDialogSelector} from "./selector-editDatesDialogSelector";

describe("action:editDatesDialogEndAtDayEndUpdated", () => {
    it("unchecks end at day end and keeps current time", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            editDatesDialog: {
                albumId: {owner: "myself", folderName: "summer-trip"},
                albumName: "Summer Trip",
                startDate: new Date("2023-07-01T00:00:00Z"),
                endDate: new Date("2023-08-01T00:00:00Z"),
                startAtDayStart: true,
                endAtDayEnd: true,
                isLoading: false,
            },
        };

        const action = editDatesDialogEndAtDayEndUpdated(false);
        const got = action.reducer(state, action);

        expect(editDatesDialogSelector(got)).toEqual({
            isOpen: true,
            albumName: "Summer Trip",
            startDate: new Date("2023-07-01T00:00:00Z"),
            endDate: new Date("2023-08-01T23:59:00Z"),
            startAtDayStart: true,
            endAtDayEnd: false,
            isLoading: false,
            errorCode: undefined,
        });
    });

    it("checks end at day end and resets time to 23:59", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            editDatesDialog: {
                albumId: {owner: "myself", folderName: "summer-trip"},
                albumName: "Summer Trip",
                startDate: new Date("2023-07-01T00:00:00Z"),
                endDate: new Date("2023-08-01T15:00:00Z"),
                startAtDayStart: true,
                endAtDayEnd: false,
                isLoading: false,
            },
        };

        const action = editDatesDialogEndAtDayEndUpdated(true);
        const got = action.reducer(state, action);

        // End date will have an undefined time: it will be ignored by the update logic, and reset if endAtDayEnd becomes false
        expect(editDatesDialogSelector(got).endAtDayEnd).toEqual(true);
    });

    it("does nothing when dialog is closed", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            editDatesDialog: undefined,
        };

        const action = editDatesDialogEndAtDayEndUpdated(true);
        const got = action.reducer(state, action);

        expect(got).toBe(state);
    });

    it("unchecks end at day end and sets time to 23:59 even if it was different", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            editDatesDialog: {
                albumId: {owner: "myself", folderName: "summer-trip"},
                albumName: "Summer Trip",
                startDate: new Date("2023-07-01T00:00:00Z"),
                endDate: new Date("2023-08-01T10:30:00Z"),
                startAtDayStart: true,
                endAtDayEnd: true,
                isLoading: false,
            },
        };

        const action = editDatesDialogEndAtDayEndUpdated(false);
        const got = action.reducer(state, action);

        expect(editDatesDialogSelector(got)).toEqual({
            isOpen: true,
            albumName: "Summer Trip",
            startDate: new Date("2023-07-01T00:00:00Z"),
            endDate: new Date("2023-08-01T23:59:00Z"),
            startAtDayStart: true,
            endAtDayEnd: false,
            isLoading: false,
            errorCode: undefined,
        });
    });
});
