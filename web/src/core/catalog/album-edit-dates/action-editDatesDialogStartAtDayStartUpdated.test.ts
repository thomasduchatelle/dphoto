import {editDatesDialogStartAtDayStartUpdated} from "./action-editDatesDialogStartAtDayStartUpdated";
import {CatalogViewerState} from "../language";
import {loadedStateWithTwoAlbums} from "../tests/test-helper-state";
import {editDatesDialogSelector} from "./selector-editDatesDialogSelector";

describe("action:editDatesDialogStartAtDayStartUpdated", () => {
    it("unchecks start at day start and keeps current time", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            editDatesDialog: {
                albumId: {owner: "myself", folderName: "summer-trip"},
                albumName: "Summer Trip",
                startDate: new Date("2023-07-01T00:00:00"),
                endDate: new Date("2023-08-01T00:00:00"),
                startAtDayStart: true,
                endAtDayEnd: true,
                isLoading: false,
            },
        };

        const action = editDatesDialogStartAtDayStartUpdated(false);
        const got = action.reducer(state, action);

        expect(editDatesDialogSelector(got)).toEqual({
            isOpen: true,
            albumName: "Summer Trip",
            startDate: new Date("2023-07-01T00:00:00"),
            endDate: new Date("2023-08-01T00:00:00"),
            startAtDayStart: false,
            endAtDayEnd: true,
            isLoading: false,
            errorCode: undefined,
        });
    });

    it("checks start at day start and resets time to 00:00", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            editDatesDialog: {
                albumId: {owner: "myself", folderName: "summer-trip"},
                albumName: "Summer Trip",
                startDate: new Date("2023-07-01T10:30:00"),
                endDate: new Date("2023-08-01T00:00:00"),
                startAtDayStart: false,
                endAtDayEnd: true,
                isLoading: false,
            },
        };

        const action = editDatesDialogStartAtDayStartUpdated(true);
        const got = action.reducer(state, action);

        expect(editDatesDialogSelector(got)).toEqual({
            isOpen: true,
            albumName: "Summer Trip",
            startDate: new Date("2023-07-01T00:00:00"),
            endDate: new Date("2023-08-01T00:00:00"),
            startAtDayStart: true,
            endAtDayEnd: true,
            isLoading: false,
            errorCode: undefined,
        });
    });

    it("does nothing when dialog is closed", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            editDatesDialog: undefined,
        };

        const action = editDatesDialogStartAtDayStartUpdated(false);
        const got = action.reducer(state, action);

        expect(got).toBe(state);
    });
});
