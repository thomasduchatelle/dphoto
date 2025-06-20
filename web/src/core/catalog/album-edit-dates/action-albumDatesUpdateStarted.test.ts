import {reduceAlbumDatesUpdateStarted, albumDatesUpdateStarted} from "./action-albumDatesUpdateStarted";
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

        const got = reduceAlbumDatesUpdateStarted(stateWithEditDialog, albumDatesUpdateStarted());

        expect(got).toEqual({
            ...stateWithEditDialog,
            editDatesDialog: {
                ...stateWithEditDialog.editDatesDialog,
                isLoading: true,
            },
        });
    });

    it("returns unchanged state when edit dates dialog is not open", () => {
        const got = reduceAlbumDatesUpdateStarted(loadedStateWithTwoAlbums, albumDatesUpdateStarted());

        expect(got).toEqual(loadedStateWithTwoAlbums);
    });
});
