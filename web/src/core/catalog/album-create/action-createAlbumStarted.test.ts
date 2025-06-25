import {CatalogViewerState} from "../language";
import {createAlbumStarted} from "./action-createAlbumStarted";
import {loadedStateWithTwoAlbums} from "../tests/test-helper-state";

describe('action:createAlbumStarted', () => {
    it('should set loading status and clear error when dialog is a CreateDialog', () => {
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
                isLoading: false,
                error: "Some error",
            },
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
            dialog: {
                type: "DeleteDialog",
                deletableAlbums: [],
                isLoading: false,
            },
        };

        const action = createAlbumStarted();
        const newState = action.reducer(state, action);

        expect(newState.dialog).toEqual(state.dialog);
    });

    it('should ignore when dialog is closed', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: undefined,
        };

        const action = createAlbumStarted();
        const newState = action.reducer(state, action);

        expect(newState.dialog).toBeUndefined();
    });
});
