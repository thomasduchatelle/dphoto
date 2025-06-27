import {albumRenamingStarted} from "./action-albumRenamingStarted";
import {editNameDialogSelector} from "./selector-editNameDialogSelector";
import {CatalogViewerState} from "../language";
import {deleteDialogWithOneAlbum, editJanAlbumNameDialog, editJanAlbumNameSelection, loadedStateWithTwoAlbums} from "../tests/test-helper-state";

describe('action:albumRenamingStarted', () => {
    it('should set loading state and clear error in edit name dialog', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...editJanAlbumNameDialog,
                error: {folderNameError: "Previous error"},
            },
        };

        const action = albumRenamingStarted();
        const got = action.reducer(state, action);

        expect(editNameDialogSelector(got)).toEqual({
            ...editJanAlbumNameSelection,
            isLoading: true,
            isSaveEnabled: false,
        });
    });

    it('should not change state when dialog is not EditNameDialog', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: deleteDialogWithOneAlbum,
        };

        const action = albumRenamingStarted();
        const got = action.reducer(state, action);

        expect(got).toBe(state);
    });
});
