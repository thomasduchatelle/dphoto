import {albumDeleteFailed, reduceAlbumDeleteFailed} from "./action-albumDeleteFailed";
import {DeleteDialogState} from "../language";
import {loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";
import {deleteDialogSelector} from "./selector-deleteDialogSelector";

describe("action:albumDeleteFailed", () => {
    const deleteDialog: DeleteDialogState = {
        deletableAlbums: twoAlbums,
        isLoading: true,
        initialSelectedAlbumId: twoAlbums[0].albumId,
    };

    it("sets the error value in deleteDialog.error when the dialog is open", () => {
        const errorMsg = "Failed to delete album";
        const resultState = reduceAlbumDeleteFailed(
            {
                ...loadedStateWithTwoAlbums,
                deleteDialog,
            },
            albumDeleteFailed(errorMsg)
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
        const resultState = reduceAlbumDeleteFailed(
            loadedStateWithTwoAlbums,
            albumDeleteFailed("Failed to delete album")
        );
        expect(resultState).toBe(loadedStateWithTwoAlbums);
    });

    it("fails to create the action if the error message is empty or blank", () => {
        expect(() => albumDeleteFailed("")).toThrow();
        expect(() => albumDeleteFailed("   ")).toThrow();
    });
});
