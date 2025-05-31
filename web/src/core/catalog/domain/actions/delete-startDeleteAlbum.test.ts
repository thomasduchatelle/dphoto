import {CatalogViewerState, DeleteDialogState} from "../catalog-state";
import {reduceStartDeleteAlbum, startDeleteAlbumAction} from "./delete-startDeleteAlbum";
import {loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";

describe("reduceStartDeleteAlbum", () => {
    it("sets the deleteDialog.loading to true and clear the error when deleteDialog is defined", () => {
        const deleteDialog: DeleteDialogState = {
            deletableAlbums: twoAlbums,
            isLoading: false,
            error: "Some error",
        };
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            deleteDialog,
        };

        const got = reduceStartDeleteAlbum(state, startDeleteAlbumAction());

        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            deleteDialog: {
                ...deleteDialog,
                isLoading: true,
                error: undefined,
            },
        });
    });

    it("ignores when the dialog is closed", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            deleteDialog: undefined,
        };

        const got = reduceStartDeleteAlbum(state, startDeleteAlbumAction());

        expect(got).toBe(state);
    });
});
