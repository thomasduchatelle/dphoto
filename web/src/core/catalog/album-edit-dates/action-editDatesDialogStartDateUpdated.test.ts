import {editDatesDialogStartDateUpdated} from "./action-editDatesDialogStartDateUpdated";
import {CatalogViewerState} from "../language";
import {loadedStateWithTwoAlbums} from "../tests/test-helper-state";

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

        expect(got.editDatesDialog?.startDate).toEqual(newStartDate);
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

        expect(got.editDatesDialog?.startDate).toEqual(newStartDate);
        expect(got.editDatesDialog?.error).toBeUndefined();
    });
});
