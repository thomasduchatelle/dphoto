import {customFolderNameChanged} from "./action-customFolderNameChanged";
import {baseEditNameSelector, BaseEditNameSelection} from "./selector-baseEditNameSelector";
import {CatalogViewerState, EditNameDialog} from "../language";
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

        const action = customFolderNameChanged(newFolderName);
        const got = action.reducer(state, action);

        const selection = baseEditNameSelector(got, got.dialog as EditNameDialog);
        expect(selection).toEqual<BaseEditNameSelectionWithSavable>({
            isSavable: true,
            albumName: editJanAlbumNameSelection.albumName,
            originalName: editJanAlbumNameSelection.originalName,
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

        const action = customFolderNameChanged("");
        const got = action.reducer(state, action);

        const selection = baseEditNameSelector(got, got.dialog as EditNameDialog);
        expect(selection).toEqual<BaseEditNameSelectionWithSavable>({
            isSavable: false,
            albumName: editJanAlbumNameSelection.albumName,
            originalName: editJanAlbumNameSelection.originalName,
            customFolderName: "",
            folderNameError: "Folder name is mandatory",
            isCustomFolderNameEnabled: true,
        });
    });

    it('should not change state when dialog is not EditNameDialog', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: deleteDialogWithOneAlbum,
        };

        const action = customFolderNameChanged(newFolderName);
        const got = action.reducer(state, action);

        expect(got.dialog?.type).not.toBe("EditNameDialog");
    });

    it('should clear technical error when folder name is changed', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...editJanAlbumNameDialog,
                isCustomFolderNameEnabled: true,
                nameError: {technicalError: "Some technical error"},
            },
        };

        const action = customFolderNameChanged("new-folder");
        const got = action.reducer(state, action);

        const selection = baseEditNameSelector(got, got.dialog as EditNameDialog);
        expect(selection).toEqual<BaseEditNameSelectionWithSavable>({
            isSavable: true,
            albumName: editJanAlbumNameSelection.albumName,
            originalName: editJanAlbumNameSelection.originalName,
            customFolderName: "new-folder",
            isCustomFolderNameEnabled: true,
        });
    });
});
