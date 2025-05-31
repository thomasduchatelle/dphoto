import {albumFailedToDeleteAction, reduceAlbumFailedToDelete} from "./delete-albumFailedToDeleteAction";
import {DeleteDialogState} from "../catalog-state";
import {loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";
import {selectDeleteAlbumDialog} from "./delete-album-dialog-selector";

describe("AlbumFailedToDeleteAction", () => {
    const deleteDialog: DeleteDialogState = {
        deletableAlbums: twoAlbums,
        isLoading: true,
        initialSelectedAlbumId: twoAlbums[0].albumId,
    };

    it("sets the error value in deleteDialog.error when the dialog is open", () => {
        const errorMsg = "Failed to delete album";
        const resultState = reduceAlbumFailedToDelete(
            {
                ...loadedStateWithTwoAlbums,
                deleteDialog,
            },
            albumFailedToDeleteAction({error: errorMsg})
        );

        expect(selectDeleteAlbumDialog(resultState)).toEqual({
            albums: twoAlbums,
            initialSelectedAlbumId: twoAlbums[0].albumId,
            isOpen: true,
            isLoading: false,
            error: errorMsg,
        });
    });

    it("ignores if the dialog is closed", () => {
        const resultState = reduceAlbumFailedToDelete(
            loadedStateWithTwoAlbums,
            albumFailedToDeleteAction({error: "Failed to delete album"})
        );
        expect(resultState).toBe(loadedStateWithTwoAlbums);
    });

    it("fails to create the action if the error message is empty or blank", () => {
        expect(() => albumFailedToDeleteAction({error: ""})).toThrow();
        expect(() => albumFailedToDeleteAction({error: "   "})).toThrow();
    });
});
