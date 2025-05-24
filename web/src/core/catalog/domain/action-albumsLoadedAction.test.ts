import {albumsLoadedAction, reduceAlbumsLoaded} from "./action-albumsLoadedAction";
import {loadedStateWithTwoAlbums, twoAlbums} from "./tests/test-helper-state";

describe("reduceAlbumsLoaded", () => {


    it("should update the list of albums and clear errors when AlbumsLoadedAction is received", () => {
        const action = albumsLoadedAction(twoAlbums);
        const got = reduceAlbumsLoaded({
            ...loadedStateWithTwoAlbums,
            allAlbums: [twoAlbums[0]],
            albums: [twoAlbums[0]],
            error: new Error("TEST previous error to clear"),
            albumsLoaded: false,
        }, action);

        expect(got).toEqual(loadedStateWithTwoAlbums);
    });

    it("should update the available filters and re-apply the selected filter when receiving AlbumsLoadedAction", () => {
        const action = albumsLoadedAction({albums: twoAlbums, redirectTo: twoAlbums[0].albumId});
        const got = reduceAlbumsLoaded({
            ...loadedStateWithTwoAlbums,
            albumFilterOptions: [loadedStateWithTwoAlbums.albumFilterOptions[0]],
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[0],
            allAlbums: [twoAlbums[1]],
            albums: [],
            mediasLoadedFromAlbumId: twoAlbums[0].albumId, // no effect
        }, action);

        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[0],
            allAlbums: twoAlbums,
            albums: [twoAlbums[0]],
        });
    });

    it("should remove the album filter if the redirectTo in AlbumsLoadedAction wouldn't be displayed", () => {
        const action = albumsLoadedAction({albums: twoAlbums, redirectTo: twoAlbums[1].albumId});
        const got = reduceAlbumsLoaded({
            ...loadedStateWithTwoAlbums,
            allAlbums: [twoAlbums[0]],
            albums: [],
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[0],
        }, action);

        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[1],
            allAlbums: twoAlbums,
            albums: twoAlbums,
        });
    });
});
