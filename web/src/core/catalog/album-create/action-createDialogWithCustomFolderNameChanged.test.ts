import {CatalogViewerState} from "../language";
import {createDialogWithCustomFolderNameChanged} from "./action-createDialogWithCustomFolderNameChanged";
import {loadedStateWithTwoAlbums} from "../tests/test-helper-state";

describe('action:createDialogWithCustomFolderNameChanged', () => {
    it('should clear forceFolderName when set to true and dialog is a CreateDialog', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                type: "CreateDialog",
                name: "Test Album",
                startDate: new Date(),
                endDate: new Date(),
                startAtDayStart: true,
                endAtDayEnd: true,
                forceFolderName: "some-folder",
                withCustomFolderName: false,
                isLoading: false,
            },
        };

        const action = createDialogWithCustomFolderNameChanged(true);
        const newState = action.reducer(state, action);

        expect(newState.dialog).toEqual({
            ...state.dialog,
            withCustomFolderName: true,
            forceFolderName: "",
        });
    });

    it('should clear forceFolderName when set to false and dialog is a CreateDialog', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                type: "CreateDialog",
                name: "Test Album",
                startDate: new Date(),
                endDate: new Date(),
                startAtDayStart: true,
                endAtDayEnd: true,
                forceFolderName: "some-folder",
                withCustomFolderName: true,
                isLoading: false,
            },
        };

        const action = createDialogWithCustomFolderNameChanged(false);
        const newState = action.reducer(state, action);

        expect(newState.dialog).toEqual({
            ...state.dialog,
            withCustomFolderName: false,
            forceFolderName: "",
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

        const action = createDialogWithCustomFolderNameChanged(true);
        const newState = action.reducer(state, action);

        expect(newState.dialog).toEqual(state.dialog);
    });

    it('should ignore when dialog is closed', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: undefined,
        };

        const action = createDialogWithCustomFolderNameChanged(true);
        const newState = action.reducer(state, action);

        expect(newState.dialog).toBeUndefined();
    });
});
