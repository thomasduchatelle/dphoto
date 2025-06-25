import {CatalogViewerState} from "../language";
import {createAlbumFailed} from "./action-createAlbumFailed";
import {loadedStateWithTwoAlbums} from "../tests/test-helper-state";

describe('action:createAlbumFailed', () => {
    it('should set the error and clear the loading status when dialog is a CreateDialog', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                type: "CreateDialog",
                name: "Test Album",
                startDate: new Date(),
                endDate: new Date(),
                startAtDayStart: true,
                endAtDayEnd: true,
                forceFolderName: "",
                withCustomFolderName: false,
                isLoading: true,
                error: undefined,
            },
        };
        const errorMessage = "Failed to create album";

        const action = createAlbumFailed(errorMessage);
        const newState = action.reducer(state, action);

        expect(newState.dialog).toEqual({
            ...state.dialog,
            isLoading: false,
            error: errorMessage,
        });
    });

    it('should ignore when dialog is not a CreateDialog', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                type: "DeleteDialog",
                deletableAlbums: [],
                isLoading: false,
            },
        };
        const errorMessage = "Failed to create album";

        const action = createAlbumFailed(errorMessage);
        const newState = action.reducer(state, action);

        expect(newState.dialog).toEqual(state.dialog);
    });

    it('should ignore when dialog is closed', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: undefined,
        };
        const errorMessage = "Failed to create album";

        const action = createAlbumFailed(errorMessage);
        const newState = action.reducer(state, action);

        expect(newState.dialog).toBeUndefined();
    });
});
