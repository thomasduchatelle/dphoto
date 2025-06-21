import {reduceAlbumDatesUpdated, albumDatesUpdated} from "./action-albumDatesUpdated";
import {CatalogViewerState, EditDatesDialogState} from "../language/CatalogViewerState";
import {initialCatalogState, twoAlbums, twoMedias} from "../tests/test-helper-state";
import {myselfUser} from "../tests/test-helper-state";

describe("action:albumDatesUpdated", () => {
    const editDatesDialog: EditDatesDialogState = {
        albumId: {owner: "myself", folderName: "album1"},
        albumName: "Test Album",
        startDate: new Date("2023-01-01"),
        endDate: new Date("2023-01-31"),
        isLoading: true
    };

    it("updates albums, medias and closes dialog", () => {
        const initialState: CatalogViewerState = {
            ...initialCatalogState(myselfUser),
            editDatesDialog
        };

        const result = reduceAlbumDatesUpdated(
            initialState,
            albumDatesUpdated({
                albums: twoAlbums,
                medias: twoMedias,
                updatedAlbumId: {owner: "myself", folderName: "album1"}
            })
        );

        expect(result).toEqual({
            ...initialState,
            albums: twoAlbums,
            allAlbums: twoAlbums,
            medias: twoMedias,
            mediasLoaded: true,
            mediasLoadedFromAlbumId: {owner: "myself", folderName: "album1"},
            editDatesDialog: undefined
        });
    });
});
