import {CatalogViewerState} from "../language";
import {createDialogNameChanged} from "./action-createDialogNameChanged";
import {loadedStateWithTwoAlbums} from "../tests/test-helper-state";

describe('action:createDialogNameChanged', () => {
    it('should set the name and clear the error when dialog is a CreateDialog', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                type: "CreateDialog",
                name: "Old Name",
                startDate: new Date(),
                endDate: new Date(),
                startAtDayStart: true,
                endAtDayEnd: true,
                forceFolderName: "",
                withCustomFolderName: false,
                isLoading: false,
                error: "AlbumFolderNameAlreadyTakenErr",
            },
        };
        const newName = "New Album Name";

        const action = createDialogNameChanged(newName);
        const newState = action.reducer(state, action);

        expect(newState.dialog).toEqual({
            ...state.dialog,
            name: newName,
            error: undefined,
        });
    });

    it('should set the name and clear any error when dialog is a CreateDialog', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                type: "CreateDialog",
                name: "Old Name",
                startDate: new Date(),
                endDate: new Date(),
                startAtDayStart: true,
                endAtDayEnd: true,
                forceFolderName: "",
                withCustomFolderName: false,
                isLoading: false,
                error: "SomeOtherError",
            },
        };
        const newName = "New Album Name";

        const action = createDialogNameChanged(newName);
        const newState = action.reducer(state, action);

        expect(newState.dialog).toEqual({
            ...state.dialog,
            name: newName,
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
        const newName = "New Album Name";

        const action = createDialogNameChanged(newName);
        const newState = action.reducer(state, action);

        expect(newState.dialog).toEqual(state.dialog);
    });

    it('should ignore when dialog is closed', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: undefined,
        };
        const newName = "New Album Name";

        const action = createDialogNameChanged(newName);
        const newState = action.reducer(state, action);

        expect(newState.dialog).toBeUndefined();
    });
});
