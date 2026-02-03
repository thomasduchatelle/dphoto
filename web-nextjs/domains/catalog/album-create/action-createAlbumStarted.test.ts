import {CatalogViewerState} from "../language";
import {createAlbumStarted} from "./action-createAlbumStarted";
import {createDialogPrefilledForMar25, deleteDialogWithOneAlbum, loadedStateWithTwoAlbums} from "../tests/test-helper-state";

describe('action:createAlbumStarted', () => {
    it('should set loading status and clear error when dialog is a CreateDialog', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                ...createDialogPrefilledForMar25,
                error: 'Some error',
            }
        };

        const action = createAlbumStarted();
        const newState = action.reducer(state, action);

        expect(newState.dialog).toEqual({
            ...state.dialog,
            isLoading: true,
            error: undefined,
        });
    });

    it('should ignore when dialog is not a CreateDialog', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: deleteDialogWithOneAlbum,
        };

        const action = createAlbumStarted();
        const newState = action.reducer(state, action);

        expect(newState).toBe(state);
    });

    it('should ignore when dialog is closed', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: undefined,
        };

        const action = createAlbumStarted();
        const newState = action.reducer(state, action);

        expect(newState).toBe(state);
    });
});
