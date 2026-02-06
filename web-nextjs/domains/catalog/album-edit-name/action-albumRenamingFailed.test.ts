import {albumRenamingFailed} from "./action-albumRenamingFailed";
import {editNameDialogSelector} from "./selector-editNameDialogSelector";
import {CatalogError, CatalogViewerState} from "../language";
import {deleteDialogWithOneAlbum, editJanAlbumNameDialog, editJanAlbumNameSelection, loadedStateWithTwoAlbums} from "../tests/test-helper-state";

const nameConflictError = new CatalogError("AlbumFolderNameAlreadyTakenErr", "Folder name already taken");

describe('action:albumRenamingFailed', () => {
    const errorMessage = "Something went wrong";

    it('should set technical error and stop loading in edit name dialog', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...editJanAlbumNameDialog,
                isLoading: true,
            },
        };

        const action = albumRenamingFailed({message: errorMessage});
        const got = action.reducer(state, action);

        expect(editNameDialogSelector(got)).toEqual({
            ...editJanAlbumNameSelection,
            isLoading: false,
            technicalError: errorMessage,
            isSaveEnabled: true,
        });
    });

    it('should set name error when AlbumFolderNameAlreadyTakenErr and custom folder name is disabled', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...editJanAlbumNameDialog,
                isLoading: true,
                isCustomFolderNameEnabled: false,
            },
        };

        const action = albumRenamingFailed(nameConflictError);
        const got = action.reducer(state, action);

        expect(editNameDialogSelector(got)).toEqual({
            ...editJanAlbumNameSelection,
            isLoading: false,
            nameError: nameConflictError.message,
            isSaveEnabled: false,
        });
    });

    it('should set folder name error when AlbumFolderNameAlreadyTakenErr and custom folder name is enabled', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...editJanAlbumNameDialog,
                isLoading: true,
                isCustomFolderNameEnabled: true,
            },
        };

        const action = albumRenamingFailed(nameConflictError);
        const got = action.reducer(state, action);

        expect(editNameDialogSelector(got)).toEqual({
            ...editJanAlbumNameSelection,
            isLoading: false,
            isCustomFolderNameEnabled: true,
            folderNameError: nameConflictError.message,
            isSaveEnabled: false,
        });
    });

    it('should set name error when AlbumNameMandatoryErr', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...editJanAlbumNameDialog,
                isLoading: true,
            },
        };

        const action = albumRenamingFailed({code: "AlbumNameMandatoryErr", message: "Album name is mandatory"});
        const got = action.reducer(state, action);

        expect(editNameDialogSelector(got)).toEqual({
            ...editJanAlbumNameSelection,
            isLoading: false,
            nameError: "Album name is mandatory",
            isSaveEnabled: false,
        });
    });

    it('should not change state when dialog is not EditNameDialog', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: deleteDialogWithOneAlbum,
        };

        const action = albumRenamingFailed({message: errorMessage});
        const got = action.reducer(state, action);

        expect(got).toBe(state);
    });
});
