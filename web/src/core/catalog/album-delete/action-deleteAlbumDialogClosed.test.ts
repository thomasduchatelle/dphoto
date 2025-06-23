import {deleteAlbumDialogClosed} from "./action-deleteAlbumDialogClosed";
import {CatalogViewerState} from "../language";
import {loadedStateWithTwoAlbums} from "../tests/test-helper-state";
import {deleteDialogSelector} from "./selector-deleteDialogSelector";

describe("action:deleteAlbumDialogClosed", () => {
    it("closes the dialog when it was open, no matter its state", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            deleteDialog: {deletableAlbums: [], isLoading: false, error: "some error"},
        };
        const action = deleteAlbumDialogClosed();
        const got = action.reducer(state, action);
        expect(got).toEqual({
            ...state,
            deleteDialog: undefined,
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
            deleteDialog: undefined,
        };
        const action = deleteAlbumDialogClosed();
        expect(action.reducer(state, action)).toBe(state);
    });

    it("supports action comparison for testing", () => {
        const action1 = deleteAlbumDialogClosed();
        const action2 = deleteAlbumDialogClosed();

        expect(action1).toEqual(action2);
        expect([action1]).toContainEqual(action2);
    });
});
