import {albumNameChanged} from "./action-albumNameChanged";
import {editNameDialogSelector} from "./selector-editNameDialogSelector";
import {CatalogViewerState} from "../language";
import {editDatesDialogForJanAlbum, editJanAlbumNameDialog, editJanAlbumNameSelection, loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";

describe('action:albumNameChanged', () => {
    const albumId = twoAlbums[0].albumId;
    const newAlbumName = "Updated Album Name";

    it('should update album name in edit name dialog', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...editJanAlbumNameDialog,
            },
        };

        const action = albumNameChanged(newAlbumName);
        const got = action.reducer(state, action);

        expect(editNameDialogSelector(got)).toEqual({
            ...editJanAlbumNameSelection,
            albumName: newAlbumName,
        });
    });

    it('should not change state when dialog is not EditNameDialog', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: editDatesDialogForJanAlbum,
        };

        const action = albumNameChanged(newAlbumName);
        const got = action.reducer(state, action);

        expect(got).toBe(state);
    });

    it('should disable save button when album name is blank', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...editJanAlbumNameDialog,
            },
        };

        const action = albumNameChanged("");
        const got = action.reducer(state, action);

        expect(editNameDialogSelector(got)).toEqual({
            ...editJanAlbumNameSelection,
            albumName: "",
            nameError: "Album name is mandatory",
            isSaveEnabled: false,
        });
    });

    it('should clear technical error when album name is changed', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...editJanAlbumNameDialog,
                error: {technicalError: "Some technical error"},
            },
        };

        const action = albumNameChanged("New Name");
        const got = action.reducer(state, action);

        expect(editNameDialogSelector(got)).toEqual({
            ...editJanAlbumNameSelection,
            albumName: "New Name",
        });
    });
});
