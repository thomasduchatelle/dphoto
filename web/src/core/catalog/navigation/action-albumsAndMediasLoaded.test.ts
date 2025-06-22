import {albumsAndMediasLoaded} from "./action-albumsAndMediasLoaded";
import {loadedStateWithTwoAlbums, myselfUser, someMedias, twoAlbums} from "../tests/test-helper-state";

import {Album, initialCatalogState} from "../language";
import {groupByDay} from "./group-by-day";

describe("action:albumsAndMediasLoaded", () => {

    it("should add the loaded albums and medias to the state, and reset all status when receiving AlbumsAndMediasLoaded", () => {
        const action = albumsAndMediasLoaded({
            albums: twoAlbums,
            medias: someMedias.flatMap(m => m.medias), // Pass raw medias
            selectedAlbum: twoAlbums[0],
        });
        const got = action.reducer({
            ...initialCatalogState(myselfUser),
            albumNotFound: true,
            albumsLoaded: false,
            mediasLoaded: false,
        }, action);

        expect(got).toEqual(loadedStateWithTwoAlbums);
    });

    it("should use 'All albums' filter even when it's the only selection available (only directly owned albums) when receiving AlbumsAndMediasLoaded", () => {
        const action = albumsAndMediasLoaded({
            albums: [twoAlbums[0]],
            medias: someMedias.flatMap(m => m.medias), // Pass raw medias
            selectedAlbum: twoAlbums[0],
        });
        const got = action.reducer(initialCatalogState(myselfUser), action);

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

    it("applies the current filter to only display albums matching", () => {
        const directlyOwnedFilter = loadedStateWithTwoAlbums.albumFilterOptions[0];
        const newDirectlyOwnedAlbum: Album = {
            albumId: {
                owner: "myself",
                folderName: "/mar-2025"
            },
            name: "March 2025",
            start: new Date(2025, 2, 1),
            end: new Date(2025, 3, 1),
            totalCount: 0,
            temperature: 0,
            relativeTemperature: 0,
            sharedWith: []
        };

        const action = albumsAndMediasLoaded({
            albums: [...twoAlbums, newDirectlyOwnedAlbum],
            medias: someMedias.flatMap(m => m.medias), // Pass raw medias
            selectedAlbum: twoAlbums [0],
        });
        const got = action.reducer(
            {
                ...loadedStateWithTwoAlbums,
                albumFilter: directlyOwnedFilter,
                albums: [loadedStateWithTwoAlbums.albums[0]],
            },
            action
        );

        // The filter should remain unchanged, and albums should contain both directly owned albums
        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            albumFilter: directlyOwnedFilter,
            albums: [loadedStateWithTwoAlbums.albums[0], newDirectlyOwnedAlbum],
            allAlbums: [...loadedStateWithTwoAlbums.allAlbums, newDirectlyOwnedAlbum],
        });
    });
});
