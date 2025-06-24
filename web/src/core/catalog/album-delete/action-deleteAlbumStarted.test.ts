import {CatalogViewerState, DeleteDialog} from "../language";
import {deleteAlbumStarted} from "./action-deleteAlbumStarted";
import {loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";

describe("action:deleteAlbumStarted", () => {
    it("sets the dialog.loading to true and clear the error when dialog is defined and is a DeleteDialog", () => {
        const deleteDialog: DeleteDialog = {
            type: "DeleteDialog",
            deletableAlbums: twoAlbums,
            isLoading: false,
            error: "Some error",
        };
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: deleteDialog,
        };

        const action = deleteAlbumStarted();
        const got = action.reducer(state, action);

        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...deleteDialog,
                isLoading: true,
                error: undefined,
            },
        });
    });

    it("ignores when the dialog is closed", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: undefined,
        };

        const action = deleteAlbumStarted();
        const got = action.reducer(state, action);

        expect(got).toBe(state);
    });

    it("ignores when the dialog is open but not a DeleteDialog", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {type: "EditDatesDialog", albumId: {owner: "o", folderName: "f"}, albumName: "n", startDate: new Date(), endDate: new Date(), startAtDayStart: true, endAtDayEnd: true, isLoading: false},
        };

        const action = deleteAlbumStarted();
        const got = action.reducer(state, action);

        expect(got).toBe(state);
    });
});
