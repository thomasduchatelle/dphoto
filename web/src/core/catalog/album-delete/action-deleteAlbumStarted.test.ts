import {CatalogViewerState, DeleteDialogState} from "../language";
import {deleteAlbumStarted, reduceDeleteAlbumStarted} from "./action-deleteAlbumStarted";
import {loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";

describe("action:deleteAlbumStarted", () => {
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

        const got = reduceDeleteAlbumStarted(state, deleteAlbumStarted());

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

        const got = reduceDeleteAlbumStarted(state, deleteAlbumStarted());

        expect(got).toBe(state);
    });
});
