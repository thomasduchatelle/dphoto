import {editNameDialogClosed} from "./action-editNameDialogClosed";
import {editNameDialogSelector} from "./selector-editNameDialogSelector";
import {deleteDialogWithOneAlbum, editJanAlbumNameDialog, loadedStateWithTwoAlbums} from "../tests/test-helper-state";

describe('action:editNameDialogClosed', () => {
    it('should close the edit name dialog', () => {
        const state = {
            ...loadedStateWithTwoAlbums,
            dialog: editJanAlbumNameDialog,
        };

        const action = editNameDialogClosed();
        const got = action.reducer(state, action);

        expect(editNameDialogSelector(got).isOpen).toEqual(false);
    });

    it('should not change state when dialog is not open', () => {
        const state = {
            ...loadedStateWithTwoAlbums,
            dialog: deleteDialogWithOneAlbum
        };

        const action = editNameDialogClosed();
        const got = action.reducer(state, action);

        expect(editNameDialogSelector(got).isOpen).toEqual(false);
        expect(got).toBe(state);
    });
});
