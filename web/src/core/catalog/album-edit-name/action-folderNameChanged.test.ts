import {folderNameChanged} from "./action-folderNameChanged";
import {EditNameDialogSelection, editNameDialogSelector} from "./selector-editNameDialogSelector";
import {CatalogViewerState} from "../language";
import {deleteDialogWithOneAlbum, editJanAlbumNameDialog, editJanAlbumNameSelection, loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";

describe('action:folderNameChanged', () => {
    const albumId = twoAlbums[0].albumId;
    const newFolderName = "new-folder-name";

    it('should update folder name in edit name dialog', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...editJanAlbumNameDialog,
                isCustomFolderNameEnabled: true,
            },
        };

        const action = folderNameChanged(newFolderName);
        const got = action.reducer(state, action);

        expect(editNameDialogSelector(got)).toEqual<EditNameDialogSelection>({
            ...editJanAlbumNameSelection,
            customFolderName: newFolderName,
            isCustomFolderNameEnabled: true,
        });
    });

    it('should disable save button when folder name is blank and folder name is enabled', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...editJanAlbumNameDialog,
                isCustomFolderNameEnabled: true,
            },
        };

        const action = folderNameChanged("");
        const got = action.reducer(state, action);

        expect(editNameDialogSelector(got)).toEqual<EditNameDialogSelection>({
            ...editJanAlbumNameSelection,
            customFolderName: "",
            folderNameError: "Folder name is mandatory",
            isCustomFolderNameEnabled: true,
            isSaveEnabled: false,
        });
    });

    it('should not change state when dialog is not EditNameDialog', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: deleteDialogWithOneAlbum,
        };

        const action = folderNameChanged(newFolderName);
        const got = action.reducer(state, action);

        expect(got).toBe(state);
    });

    it('should clear technical error when folder name is changed', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...editJanAlbumNameDialog,
                isCustomFolderNameEnabled: true,
                error: {technicalError: "Some technical error"},
            },
        };

        const action = folderNameChanged("new-folder");
        const got = action.reducer(state, action);

        expect(editNameDialogSelector(got)).toEqual<EditNameDialogSelection>({
            ...editJanAlbumNameSelection,
            customFolderName: "new-folder",
            isCustomFolderNameEnabled: true,
        });
    });
});
