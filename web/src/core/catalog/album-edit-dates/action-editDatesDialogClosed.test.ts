import {editDatesDialogClosed, reduceEditDatesDialogClosed} from "./action-editDatesDialogClosed";
import {CatalogViewerState} from "../language";
import {loadedStateWithTwoAlbums} from "../tests/test-helper-state";

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

        const got = reduceEditDatesDialogClosed(state, editDatesDialogClosed());

        expect(got).toEqual({
            ...state,
            editDatesDialog: undefined,
        });
    });

    it("does nothing when dialog is already closed", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            editDatesDialog: undefined,
        };

        const got = reduceEditDatesDialogClosed(state, editDatesDialogClosed());

        expect(got).toBe(state);
    });
});
