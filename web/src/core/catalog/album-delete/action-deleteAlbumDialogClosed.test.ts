import {deleteAlbumDialogClosed} from "./action-deleteAlbumDialogClosed";
import {CatalogViewerState, DeleteDialog} from "../language";
import {loadedStateWithTwoAlbums} from "../tests/test-helper-state";
import {deleteDialogSelector} from "./selector-deleteDialogSelector";

describe("action:deleteAlbumDialogClosed", () => {
    it("closes the dialog when it was open, no matter its state", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {type: "DeleteDialog", deletableAlbums: [], isLoading: false, error: "some error"},
        };
        const action = deleteAlbumDialogClosed();
        const got = action.reducer(state, action);
        expect(got).toEqual({
            ...state,
            dialog: undefined,
        });

        expect(deleteDialogSelector(got)).toEqual({
            albums: [],
            initialSelectedAlbumId: undefined,
            isOpen: false,
            isLoading: false,
            error: undefined,
        })
    });

    it("ignores when the dialog is already closed", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: undefined,
        };
        const action = deleteAlbumDialogClosed();
        expect(action.reducer(state, action)).toBe(state);
    });
});
