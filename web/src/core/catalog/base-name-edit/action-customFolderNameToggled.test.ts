import {customFolderNameToggled} from "./action-customFolderNameToggled";
import {baseEditNameSelector, BaseEditNameSelection} from "./selector-baseEditNameSelector";
import {CatalogViewerState} from "../language";
import {deleteDialogWithOneAlbum, editJanAlbumNameDialog, editJanAlbumNameSelection, loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";

describe('action:folderNameEnabledChanged', () => {
    const albumId = twoAlbums[0].albumId;
    it('should enable folder name field and pre-fill with original folder name', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...editJanAlbumNameDialog,
            },
        };

        const action = customFolderNameToggled(true);
        const got = action.reducer(state, action);

        const {isSavable, ...baseSelection} = baseEditNameSelector(got, got.dialog);
        expect(baseSelection).toEqual<BaseEditNameSelection>({
            albumName: editJanAlbumNameSelection.albumName,
            originalName: editJanAlbumNameSelection.originalName,
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

        const {isSavable, ...baseSelection} = baseEditNameSelector(got, got.dialog);
        expect(baseSelection).toEqual<BaseEditNameSelection>({
            albumName: editJanAlbumNameSelection.albumName,
            originalName: editJanAlbumNameSelection.originalName,
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
});
