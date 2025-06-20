import {reduceAlbumDatesUpdated, albumDatesUpdated} from "./action-albumDatesUpdated";
import {loadedStateWithTwoAlbums, twoAlbums, someMedias, myselfUser} from "../tests/test-helper-state";
import {CatalogViewerState} from "../language";

describe("action:albumDatesUpdated", () => {
    it("updates albums and medias, closes dialog, and refreshes filters", () => {
        const updatedAlbum = {
            ...twoAlbums[0],
            start: new Date("2023-07-10T00:00:00Z"),
            end: new Date("2023-07-21T00:00:00Z"),
        };
        const updatedAlbums = [updatedAlbum, twoAlbums[1]];

        const stateWithEditDialog: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            editDatesDialog: {
                albumId: twoAlbums[0].albumId,
                albumName: twoAlbums[0].name,
                startDate: twoAlbums[0].start,
                endDate: twoAlbums[0].end,
                isLoading: true,
            },
        };

        const got = reduceAlbumDatesUpdated(
            stateWithEditDialog,
            albumDatesUpdated({
                albums: updatedAlbums,
                medias: someMedias[0].medias,
                updatedAlbumId: twoAlbums[0].albumId,
            })
        );

        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            allAlbums: updatedAlbums,
            albums: updatedAlbums,
            medias: someMedias,
            mediasLoadedFromAlbumId: twoAlbums[0].albumId,
            albumsLoaded: true,
            mediasLoaded: true,
            editDatesDialog: undefined,
        });
    });
});
