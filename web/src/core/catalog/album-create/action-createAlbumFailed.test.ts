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
    it('should set the error and clear the loading status when dialog is a CreateDialog', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...createDialogPrefilledForMar25,
                isLoading: true,
            },
        };
        const errorMessage = "Failed to create album";

        const action = createAlbumFailed(errorMessage);
        const newState = action.reducer(state, action);

        expect(createDialogSelector(newState)).toEqual({
            ...createDialogSelectionPrefilledForMar25,
            isLoading: false,
            error: errorMessage,
        });
    });

    it('should ignore when dialog is not a CreateDialog', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: deleteDialogWithOneAlbum,
        };
        const errorMessage = "Failed to create album";

        const action = createAlbumFailed(errorMessage);
        const newState = action.reducer(state, action);

        expect(newState).toBe(state);
    });

    it('should ignore when dialog is closed', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: undefined,
        };
        const errorMessage = "Failed to create album";

        const action = createAlbumFailed(errorMessage);
        const newState = action.reducer(state, action);

        expect(newState).toBe(state);
    });
});
