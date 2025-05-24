import {albumsAndMediasLoadedAction, reduceAlbumsAndMediasLoaded} from "./action-albumsAndMediasLoadedAction";
import {loadedStateWithTwoAlbums, myselfUser, someMedias, twoAlbums} from "./tests/test-helper-state";

import {initialCatalogState} from "./initial-catalog-state";

describe("reduceAlbumsAndMediasLoaded", () => {

    it("should add the loaded albums and medias to the state, and reset all status when receiving AlbumsAndMediasLoadedAction", () => {
        const action = albumsAndMediasLoadedAction({
            albums: twoAlbums,
            medias: someMedias,
            selectedAlbum: twoAlbums[0],
        });
        const got = reduceAlbumsAndMediasLoaded({
            ...initialCatalogState(myselfUser),
            albumNotFound: true,
            albumsLoaded: false,
            mediasLoaded: false,
        }, action);

        expect(got).toEqual(loadedStateWithTwoAlbums);
    });

    it("should use 'All albums' filter even when it's the only selection available (only directly owned albums) when receiving AlbumsAndMediasLoadedAction", () => {
        const action = albumsAndMediasLoadedAction({
            albums: [twoAlbums[0]],
            medias: someMedias,
            selectedAlbum: twoAlbums[0],
        });
        const got = reduceAlbumsAndMediasLoaded(initialCatalogState(myselfUser), action);

        const allAlbumFilter = {
            criterion: {
                owners: []
            },
            avatars: [myselfUser.picture ?? ""],
            name: "All albums",
        };
        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            albums: [twoAlbums[0]],
            allAlbums: [twoAlbums[0]],
            albumFilter: allAlbumFilter,
            albumFilterOptions: [allAlbumFilter],
        });
    });
});
