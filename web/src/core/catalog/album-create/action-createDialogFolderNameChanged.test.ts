import {CatalogViewerState} from "../language";
import {createDialogFolderNameChanged} from "./action-createDialogFolderNameChanged";
import {loadedStateWithTwoAlbums} from "../tests/test-helper-state";

describe('action:createDialogFolderNameChanged', () => {
    it('should set the forced folder name and clear the error when dialog is a CreateDialog', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                type: "CreateDialog",
                name: "Test Album",
                startDate: new Date(),
                endDate: new Date(),
                startAtDayStart: true,
                endAtDayEnd: true,
                forceFolderName: "old-folder",
                withCustomFolderName: true,
                isLoading: false,
                error: "AlbumFolderNameAlreadyTakenErr",
            },
        };
        const newFolderName = "new-folder";

        const action = createDialogFolderNameChanged(newFolderName);
        const newState = action.reducer(state, action);

        expect(newState.dialog).toEqual({
            ...state.dialog,
            forceFolderName: newFolderName,
            error: undefined,
        });
    });

    it('should set the forced folder name and clear any error when dialog is a CreateDialog', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                type: "CreateDialog",
                name: "Test Album",
                startDate: new Date(),
                endDate: new Date(),
                startAtDayStart: true,
                endAtDayEnd: true,
                forceFolderName: "old-folder",
                withCustomFolderName: true,
                isLoading: false,
                error: "SomeOtherError",
            },
        };
        const newFolderName = "new-folder";

        const action = createDialogFolderNameChanged(newFolderName);
        const newState = action.reducer(state, action);

        expect(newState.dialog).toEqual({
            ...state.dialog,
            forceFolderName: newFolderName,
            error: "SomeOtherError", // Only "AlbumFolderNameAlreadyTakenErr" is cleared
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
        const newFolderName = "new-folder";

        const action = createDialogFolderNameChanged(newFolderName);
        const newState = action.reducer(state, action);

        expect(newState.dialog).toEqual(state.dialog);
    });

    it('should ignore when dialog is closed', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: undefined,
        };
        const newFolderName = "new-folder";

        const action = createDialogFolderNameChanged(newFolderName);
        const newState = action.reducer(state, action);

        expect(newState.dialog).toBeUndefined();
    });
});
