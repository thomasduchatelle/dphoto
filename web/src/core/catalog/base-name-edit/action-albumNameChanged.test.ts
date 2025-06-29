import {albumNameChanged} from "./action-albumNameChanged";
import {BaseEditNameSelectionWithSavable, baseEditNameSelector} from "./selector-baseEditNameSelector";
import {CatalogViewerState, EditNameDialog} from "../language";
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

        const selection = baseEditNameSelector(got, got.dialog as EditNameDialog);
        expect(selection).toEqual<BaseEditNameSelectionWithSavable>({
            isSavable: true,
            albumName: newAlbumName,
            originalName: editJanAlbumNameSelection.originalName,
            customFolderName: editJanAlbumNameSelection.customFolderName,
            isCustomFolderNameEnabled: editJanAlbumNameSelection.isCustomFolderNameEnabled,
        });
    });

    it('should not change state when dialog is not EditNameDialog', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: editDatesDialogForJanAlbum,
        };

        const action = albumNameChanged(newAlbumName);
        const got = action.reducer(state, action);

        expect(got.dialog?.type).not.toBe("EditNameDialog");
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

        const selection = baseEditNameSelector(got, got.dialog as EditNameDialog);
        expect(selection).toEqual<BaseEditNameSelectionWithSavable>({
            isSavable: false,
            albumName: "",
            originalName: editJanAlbumNameSelection.originalName,
            customFolderName: editJanAlbumNameSelection.customFolderName,
            isCustomFolderNameEnabled: editJanAlbumNameSelection.isCustomFolderNameEnabled,
            nameError: "Album name is mandatory",
        });
    });

    it('should preserve technical error when album name is changed', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...editJanAlbumNameDialog,
                nameError: {technicalError: "Some technical error", folderNameError: "Some folder error"},
            },
        };

        const action = albumNameChanged("New Name");
        const got = action.reducer(state, action);

        const selection = baseEditNameSelector(got, got.dialog as EditNameDialog);
        expect(selection).toEqual<BaseEditNameSelectionWithSavable>({
            isSavable: false,
            albumName: "New Name",
            originalName: editJanAlbumNameSelection.originalName,
            customFolderName: editJanAlbumNameSelection.customFolderName,
            isCustomFolderNameEnabled: editJanAlbumNameSelection.isCustomFolderNameEnabled,
            technicalError: "Some technical error",
            folderNameError: "Some folder error",
        });
    });
});
