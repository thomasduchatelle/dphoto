import {closeDeleteAlbumDialogAction, reduceCloseDeleteAlbumDialog} from "./delete-closeDeleteAlbumDialog";
import {CatalogViewerState} from "../catalog-state";
import {loadedStateWithTwoAlbums} from "../tests/test-helper-state";
import {deleteDialogSelector} from "./selector-deleteDialogSelector";

describe("reduceCloseDeleteAlbumDialog", () => {
    it("closes the dialog when it was open, no matter its state", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            deleteDialog: {deletableAlbums: [], isLoading: false, error: "some error"},
        };
        const got = reduceCloseDeleteAlbumDialog(state, closeDeleteAlbumDialogAction());
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
        expect(reduceCloseDeleteAlbumDialog(state, closeDeleteAlbumDialogAction())).toBe(state);
    });
});
