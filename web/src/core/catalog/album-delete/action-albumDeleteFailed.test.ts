import {albumDeleteFailed} from "./action-albumDeleteFailed";
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
        const action = albumDeleteFailed(errorMsg);
        const resultState = action.reducer(
            {
                ...loadedStateWithTwoAlbums,
                deleteDialog,
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

    it("supports action comparison for testing", () => {
        const action1 = albumDeleteFailed("error");
        const action2 = albumDeleteFailed("error");
        const action3 = albumDeleteFailed("another error");

        expect(action1).toEqual(action2);
        expect(action1).not.toEqual(action3);
        expect([action1]).toContainEqual(action2);
    });
});
