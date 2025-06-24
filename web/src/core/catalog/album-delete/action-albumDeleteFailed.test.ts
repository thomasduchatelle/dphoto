import {albumDeleteFailed} from "./action-albumDeleteFailed";
import {DeleteDialog} from "../language";
import {loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";
import {deleteDialogSelector} from "./selector-deleteDialogSelector";
import {CatalogViewerState} from "../language";

describe("action:albumDeleteFailed", () => {
    const deleteDialog: DeleteDialog = {
        type: "DeleteDialog",
        deletableAlbums: twoAlbums,
        isLoading: true,
        initialSelectedAlbumId: twoAlbums[0].albumId,
    };

    it("sets the error value in dialog.error when the dialog is open and is a DeleteDialog", () => {
        const errorMsg = "Failed to delete album";
        const action = albumDeleteFailed(errorMsg);
        const resultState = action.reducer(
            {
                ...loadedStateWithTwoAlbums,
                dialog: deleteDialog,
            },
            action
        );

        expect(deleteDialogSelector(resultState)).toEqual({
            albums: twoAlbums,
            initialSelectedAlbumId: twoAlbums[0].albumId,
            isOpen: true,
            isLoading: false,
            error: errorMsg,
        });
    });

    it("ignores if the dialog is closed", () => {
        const action = albumDeleteFailed("Failed to delete album");
        const resultState = action.reducer(
            loadedStateWithTwoAlbums,
            action
        );
        expect(resultState).toBe(loadedStateWithTwoAlbums);
    });

    it("ignores if the dialog is open but not a DeleteDialog", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {type: "EditDatesDialog", albumId: {owner: "o", folderName: "f"}, albumName: "n", startDate: new Date(), endDate: new Date(), startAtDayStart: true, endAtDayEnd: true, isLoading: false},
        };
        const action = albumDeleteFailed("Failed to delete album");
        const resultState = action.reducer(state, action);
        expect(resultState).toBe(state);
    });

    it("fails to create the action if the error message is empty or blank", () => {
        expect(() => {
            const action = albumDeleteFailed("")
            action.reducer(loadedStateWithTwoAlbums, action);
        }).toThrow();
        expect(() => {
            const action = albumDeleteFailed("   ")
            action.reducer(loadedStateWithTwoAlbums, action);
        }).toThrow();
    });
});
