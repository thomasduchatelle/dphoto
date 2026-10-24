import {CatalogViewerState} from "../language";
import {createAlbumFailed} from "./action-createAlbumFailed";
import {
    createDialogPrefilledForMar25,
    createDialogSelectionPrefilledForMar25,
    deleteDialogWithOneAlbum,
    loadedStateWithTwoAlbums
} from "../tests/test-helper-state";
import {createDialogSelector} from "./selector-createDialogSelector";

describe('action:createAlbumFailed', () => {
    it('should set the general error and clear the loading status when dialog is a CreateDialog', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...createDialogPrefilledForMar25,
                isLoading: true,
            },
        };

        const action = createAlbumFailed({error: "Failed to create album"});
        const newState = action.reducer(state, action);

        expect(createDialogSelector(newState)).toEqual({
            ...createDialogSelectionPrefilledForMar25,
            isLoading: false,
            error: "Failed to create album",
        });
    });

    it('should set the name error and disable submit when a nameError is provided', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...createDialogPrefilledForMar25,
                isLoading: true,
            },
        };

        const action = createAlbumFailed({nameError: {nameError: "An album with this folder name already exists."}});
        const newState = action.reducer(state, action);

        expect(createDialogSelector(newState)).toEqual({
            ...createDialogSelectionPrefilledForMar25,
            isLoading: false,
            canSubmit: false,
            nameError: "An album with this folder name already exists.",
        });
    });

    it('should ignore when dialog is not a CreateDialog', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: deleteDialogWithOneAlbum,
        };

        const action = createAlbumFailed({error: "Failed to create album"});
        const newState = action.reducer(state, action);

        expect(newState).toBe(state);
    });

    it('should ignore when dialog is closed', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: undefined,
        };

        const action = createAlbumFailed({error: "Failed to create album"});
        const newState = action.reducer(state, action);

        expect(newState).toBe(state);
    });
});
