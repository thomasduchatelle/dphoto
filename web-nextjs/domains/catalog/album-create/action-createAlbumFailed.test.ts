import {CatalogError, CatalogViewerState} from "../language";
import {createAlbumFailed} from "./action-createAlbumFailed";
import {
    createDialogPrefilledForMar25,
    createDialogSelectionPrefilledForMar25,
    deleteDialogWithOneAlbum,
    loadedStateWithTwoAlbums
} from "../tests/test-helper-state";
import {createDialogSelector} from "./selector-createDialogSelector";

const nameConflictError = new CatalogError("AlbumFolderNameAlreadyTakenErr", "An album with the same name or folder name already exists.");

describe('action:createAlbumFailed', () => {
    it('should set the technical error and clear the loading status when dialog is a CreateDialog', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...createDialogPrefilledForMar25,
                isLoading: true,
            },
        };

        const action = createAlbumFailed({message: "Failed to create album"});
        const newState = action.reducer(state, action);

        expect(createDialogSelector(newState)).toEqual({
            ...createDialogSelectionPrefilledForMar25,
            isLoading: false,
            error: "Failed to create album",
        });
    });

    it('should set the name error when AlbumFolderNameAlreadyTakenErr and the custom folder name is disabled', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...createDialogPrefilledForMar25,
                isLoading: true,
                isCustomFolderNameEnabled: false,
            },
        };

        const action = createAlbumFailed(nameConflictError);
        const newState = action.reducer(state, action);

        expect(createDialogSelector(newState)).toEqual({
            ...createDialogSelectionPrefilledForMar25,
            isLoading: false,
            canSubmit: false,
            nameError: nameConflictError.message,
        });
    });

    it('should set the folder name error when AlbumFolderNameAlreadyTakenErr and the custom folder name is enabled', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...createDialogPrefilledForMar25,
                isLoading: true,
                isCustomFolderNameEnabled: true,
                customFolderName: "2025-03_march",
            },
        };

        const action = createAlbumFailed(nameConflictError);
        const newState = action.reducer(state, action);

        expect(createDialogSelector(newState)).toEqual({
            ...createDialogSelectionPrefilledForMar25,
            isLoading: false,
            canSubmit: false,
            isCustomFolderNameEnabled: true,
            customFolderName: "2025-03_march",
            folderNameError: nameConflictError.message,
        });
    });

    it('should ignore when dialog is not a CreateDialog', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: deleteDialogWithOneAlbum,
        };

        const action = createAlbumFailed({message: "Failed to create album"});
        const newState = action.reducer(state, action);

        expect(newState).toBe(state);
    });

    it('should ignore when dialog is closed', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: undefined,
        };

        const action = createAlbumFailed({message: "Failed to create album"});
        const newState = action.reducer(state, action);

        expect(newState).toBe(state);
    });
});
