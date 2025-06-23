import {editDatesDialogEndDateUpdated} from "./action-editDatesDialogEndDateUpdated";
import {CatalogViewerState} from "../language";
import {loadedStateWithTwoAlbums} from "../tests/test-helper-state";
import {editDatesDialogSelector} from "./selector-editDatesDialogSelector";

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
});
