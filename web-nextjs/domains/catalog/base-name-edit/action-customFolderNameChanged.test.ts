import {customFolderNameChanged} from "./action-customFolderNameChanged";
import {BaseEditNameSelectionWithSavable, baseEditNameSelector} from "./selector-baseEditNameSelector";
import {CatalogViewerState, EditNameDialog} from "../language";
import {deleteDialogWithOneAlbum, editJanAlbumNameDialog, loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";

describe('action:folderNameChanged', () => {
    const newFolderName = "new-folder-name";

    const baseEditSelectionForJan: BaseEditNameSelectionWithSavable = {
        albumName: twoAlbums[0].name,
        customFolderName: "",
        isCustomFolderNameEnabled: false,
        nameError: undefined,
        folderNameError: undefined,
        isSavable: true,
    }

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
            ...baseEditSelectionForJan,
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
            ...baseEditSelectionForJan,
            isSavable: false,
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

    it('should clear folderName error bu preserve other errors when folder name is changed', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...editJanAlbumNameDialog,
                isCustomFolderNameEnabled: true,
                technicalError: "Some technical error",
                nameError: {nameError: "Some name error", folderNameError: "Some folder name error"},
            },
        };

        const action = customFolderNameChanged("new-folder");
        const got = action.reducer(state, action);

        const selection = baseEditNameSelector(got, got.dialog as EditNameDialog);
        expect(selection).toEqual<BaseEditNameSelectionWithSavable>({
            ...baseEditSelectionForJan,
            isSavable: false,
            customFolderName: "new-folder",
            isCustomFolderNameEnabled: true,
            nameError: "Some name error",
        });
    });
});
