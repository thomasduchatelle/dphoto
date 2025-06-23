import {albumDatesUpdateStarted} from "./action-albumDatesUpdateStarted";
import {loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";
import {CatalogViewerState} from "../language";

describe("action:albumDatesUpdateStarted", () => {
    it("sets loading state to true when edit dates dialog is open", () => {
        const stateWithEditDialog: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            editDatesDialog: {
                albumId: twoAlbums[0].albumId,
                albumName: twoAlbums[0].name,
                startDate: twoAlbums[0].start,
                endDate: twoAlbums[0].end,
                isLoading: false,
            },
        };

        const action = albumDatesUpdateStarted();
        const got = action.reducer(stateWithEditDialog, action);

        expect(got).toEqual({
            ...stateWithEditDialog,
            editDatesDialog: {
                ...stateWithEditDialog.editDatesDialog,
                isLoading: true,
            },
        });
    });

    it("returns unchanged state when edit dates dialog is not open", () => {
        const action = albumDatesUpdateStarted();
        const got = action.reducer(loadedStateWithTwoAlbums, action);

        expect(got).toEqual(loadedStateWithTwoAlbums);
    });

    it("clears error when starting update", () => {
        const stateWithError: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            editDatesDialog: {
                albumId: twoAlbums[0].albumId,
                albumName: twoAlbums[0].name,
                startDate: twoAlbums[0].start,
                endDate: twoAlbums[0].end,
                isLoading: false,
                error: "Previous error",
            },
        };

        const action = albumDatesUpdateStarted();
        const got = action.reducer(stateWithError, action);

        expect(got).toEqual({
            ...stateWithError,
            editDatesDialog: {
                ...stateWithError.editDatesDialog,
                isLoading: true,
                error: undefined,
            },
        });
    });
});
