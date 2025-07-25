import {albumDatesUpdateStarted} from "./action-albumDatesUpdateStarted";
import {loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";
import {CatalogViewerState} from "../language";

describe("action:albumDatesUpdateStarted", () => {
    it("sets loading state to true when edit dates dialog is open", () => {
        const stateWithEditDialog: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                type: "EditDatesDialog",
                albumId: twoAlbums[0].albumId,
                albumName: twoAlbums[0].name,
                startDate: twoAlbums[0].start,
                endDate: twoAlbums[0].end,
                isLoading: false,
                startAtDayStart: true,
                endAtDayEnd: true,
            },
        };

        const action = albumDatesUpdateStarted();
        const got = action.reducer(stateWithEditDialog, action);

        expect(got).toEqual({
            ...stateWithEditDialog,
            dialog: {
                ...stateWithEditDialog.dialog,
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
            dialog: {
                type: "EditDatesDialog",
                albumId: twoAlbums[0].albumId,
                albumName: twoAlbums[0].name,
                startDate: twoAlbums[0].start,
                endDate: twoAlbums[0].end,
                isLoading: false,
                error: "Previous error",
                startAtDayStart: true,
                endAtDayEnd: true,
            },
        };

        const action = albumDatesUpdateStarted();
        const got = action.reducer(stateWithError, action);

        expect(got).toEqual({
            ...stateWithError,
            dialog: {
                ...stateWithError.dialog,
                isLoading: true,
                error: undefined,
            },
        });
    });
});
