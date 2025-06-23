import {editDatesDialogClosed} from "./action-editDatesDialogClosed";
import {CatalogViewerState} from "../language";
import {loadedStateWithTwoAlbums} from "../tests/test-helper-state";
import {DEFAULT_EDIT_DATES_DIALOG_SELECTION, editDatesDialogSelector} from "./selector-editDatesDialogSelector";

describe("action:editDatesDialogClosed", () => {
    it("closes the dialog by setting editDatesDialog to undefined", () => {
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

        const action = editDatesDialogClosed();
        const got = action.reducer(state, action);
        const dialogSelection = editDatesDialogSelector(got);

        expect(dialogSelection).toEqual(DEFAULT_EDIT_DATES_DIALOG_SELECTION);
    });

    it("does nothing when dialog is already closed", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            editDatesDialog: undefined,
        };

        const action = editDatesDialogClosed();
        const got = action.reducer(state, action);
        const dialogSelection = editDatesDialogSelector(got);

        expect(dialogSelection).toEqual(DEFAULT_EDIT_DATES_DIALOG_SELECTION);
        expect(got).toBe(state);
    });
});
