import {albumsAndMediasLoadedAction, reduceAlbumsAndMediasLoaded} from "./catalog-action-AlbumsAndMediasLoadedAction";
import {loadedStateWithTwoAlbums, twoAlbums, someMedias} from "./tests/test-helper-state";

describe("reduceAlbumsAndMediasLoaded", () => {

    it("should add the loaded albums and medias to the state, and reset all status when receiving AlbumsAndMediasLoadedAction", () => {
        const action = albumsAndMediasLoadedAction(twoAlbums, someMedias, twoAlbums[0]);
        const got = reduceAlbumsAndMediasLoaded({
            ...loadedStateWithTwoAlbums,
            albumNotFound: true,
            albumsLoaded: false,
            mediasLoaded: false,
        }, action);

        expect(got).toEqual(loadedStateWithTwoAlbums);
    });

    it("should use 'All albums' filter even when it's the only selection available (only directly owned albums) when receiving AlbumsAndMediasLoadedAction", () => {
        const action = albumsAndMediasLoadedAction([twoAlbums[0]], someMedias, twoAlbums[0]);
        const got = reduceAlbumsAndMediasLoaded({
            ...loadedStateWithTwoAlbums,
            allAlbums: [twoAlbums[0]],
            albums: [twoAlbums[0]],
            albumFilterOptions: [{
                criterion: {
                    owners: []
                },
                avatars: [loadedStateWithTwoAlbums.currentUser.picture ?? ""],
                name: "All albums",
            }],
            albumFilter: {
                criterion: {
                    owners: []
                },
                avatars: [loadedStateWithTwoAlbums.currentUser.picture ?? ""],
                name: "All albums",
            },
            medias: someMedias,
            mediasLoadedFromAlbumId: twoAlbums[0].albumId,
        }, action);

        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            allAlbums: [twoAlbums[0]],
            albums: [twoAlbums[0]],
            albumFilterOptions: [{
                criterion: {
                    owners: []
                },
                avatars: [loadedStateWithTwoAlbums.currentUser.picture ?? ""],
                name: "All albums",
            }],
            albumFilter: {
                criterion: {
                    owners: []
                },
                avatars: [loadedStateWithTwoAlbums.currentUser.picture ?? ""],
                name: "All albums",
            },
            medias: someMedias,
            mediasLoadedFromAlbumId: twoAlbums[0].albumId,
        });
    });
});
