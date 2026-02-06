import {albumNameChanged} from "./action-albumNameChanged";
import {BaseEditNameSelectionWithSavable, baseEditNameSelector} from "./selector-baseEditNameSelector";
import {CatalogViewerState, EditNameDialog} from "../language";
import {editDatesDialogForJanAlbum, editJanAlbumNameDialog, loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";

describe('action:albumNameChanged', () => {
    const newAlbumName = "Updated Album Name";

    const baseEditSelectionForJan: BaseEditNameSelectionWithSavable = {
        albumName: twoAlbums[0].name,
        customFolderName: "",
        isCustomFolderNameEnabled: false,
        nameError: undefined,
        folderNameError: undefined,
        isSavable: true,
    }

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
            ...baseEditSelectionForJan,
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
            ...baseEditSelectionForJan,
            isSavable: false,
            albumName: "",
            nameError: "Album name is mandatory",
        });
    });

    it('should clear name error but preserve other errors when album name is changed', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...editJanAlbumNameDialog,
                technicalError: "Some technical error",
                nameError: {folderNameError: "Some folder error", nameError: "Some name error"},
            },
        };

        const action = albumNameChanged("New Name");
        const got = action.reducer(state, action);

        const selection = baseEditNameSelector(got, got.dialog as EditNameDialog);
        expect(selection).toEqual<BaseEditNameSelectionWithSavable>({
            ...baseEditSelectionForJan,
            isSavable: false,
            albumName: "New Name",
            folderNameError: "Some folder error",
        });
    });
});
