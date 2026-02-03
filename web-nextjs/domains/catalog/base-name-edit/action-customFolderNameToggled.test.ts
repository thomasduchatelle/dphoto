import {customFolderNameToggled} from "./action-customFolderNameToggled";
import {BaseEditNameSelectionWithSavable, baseEditNameSelector} from "./selector-baseEditNameSelector";
import {CatalogViewerState, EditNameDialog} from "../language";
import {deleteDialogWithOneAlbum, editJanAlbumNameDialog, loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";

describe('action:folderNameEnabledChanged', () => {
    const albumId = twoAlbums[0].albumId;
    const baseEditSelectionForJan: BaseEditNameSelectionWithSavable = {
        albumName: twoAlbums[0].name,
        customFolderName: "",
        isCustomFolderNameEnabled: false,
        nameError: undefined,
        folderNameError: undefined,
        isSavable: true,
    }

    it('should enable folder name field and pre-fill with original folder name', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...editJanAlbumNameDialog,
            },
        };

        const action = customFolderNameToggled(true);
        const got = action.reducer(state, action);

        const selection = baseEditNameSelector(got, got.dialog as EditNameDialog);
        expect(selection).toEqual<BaseEditNameSelectionWithSavable>({
            ...baseEditSelectionForJan,
            customFolderName: albumId.folderName,
            isCustomFolderNameEnabled: true,
        });
    });

    it('should disable folder name field and clear folder name', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...editJanAlbumNameDialog,
                customFolderName: "some-folder-name",
                isCustomFolderNameEnabled: true,
            },
        };

        const action = customFolderNameToggled(false);
        const got = action.reducer(state, action);

        const selection = baseEditNameSelector(got, got.dialog as EditNameDialog);
        expect(selection).toEqual<BaseEditNameSelectionWithSavable>({
            ...baseEditSelectionForJan,
            customFolderName: "",
            isCustomFolderNameEnabled: false,
        });
    });

    it('should not change state when dialog is not EditNameDialog', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: deleteDialogWithOneAlbum,
        };

        const action = customFolderNameToggled(true);
        const got = action.reducer(state, action);

        expect(got.dialog?.type).not.toBe("EditNameDialog");
    });

    it('should be without error but not be savable if custom is switch of without default value for folder name', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...editJanAlbumNameDialog,
                originalFolderName: undefined,
            },
        };

        const action = customFolderNameToggled(true);
        const got = action.reducer(state, action);

        const selection = baseEditNameSelector(got, got.dialog as EditNameDialog);
        expect(selection).toEqual<BaseEditNameSelectionWithSavable>({
            ...baseEditSelectionForJan,
            customFolderName: "",
            isCustomFolderNameEnabled: true,
            isSavable: false,
        });
    });
});
