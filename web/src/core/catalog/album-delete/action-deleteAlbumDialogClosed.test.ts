import {deleteAlbumDialogClosed, reduceDeleteAlbumDialogClosed} from "./action-deleteAlbumDialogClosed";
import {CatalogViewerState} from "../language";
import {loadedStateWithTwoAlbums} from "../tests/test-helper-state";
import {deleteDialogSelector} from "./selector-deleteDialogSelector";

describe("action:deleteAlbumDialogClosed", () => {
    it("closes the dialog when it was open, no matter its state", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            deleteDialog: {deletableAlbums: [], isLoading: false, error: "some error"},
        };
        const got = reduceDeleteAlbumDialogClosed(state, deleteAlbumDialogClosed());
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
        expect(reduceDeleteAlbumDialogClosed(state, deleteAlbumDialogClosed())).toBe(state);
    });
});
