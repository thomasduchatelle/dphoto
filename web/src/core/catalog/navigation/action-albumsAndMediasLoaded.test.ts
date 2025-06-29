import {albumsAndMediasLoaded} from "./action-albumsAndMediasLoaded";
import {loadedStateWithTwoAlbums, myselfUser, selectionForLoadedStateWithTwoAlbums, someMediasByDays, twoAlbums} from "../tests/test-helper-state";

import {Album, initialCatalogState} from "../language";
import {catalogViewerPageSelector} from "./selector-catalog-viewer-page";
import {groupByDay} from "./group-by-day";

describe("action:albumsAndMediasLoaded", () => {

    it("should add the loaded albums and medias to the state, and reset all status when receiving AlbumsAndMediasLoaded", () => {
        const action = albumsAndMediasLoaded({
            albums: twoAlbums,
            medias: someMediasByDays.flatMap(m => m.medias), // Pass raw medias
            mediasFromAlbumId: twoAlbums[0].albumId,
        });
        const got = action.reducer({
            ...initialCatalogState(myselfUser),
            albumNotFound: true,
            albumsLoaded: false,
            mediasLoaded: false,
        }, action);

        expect(catalogViewerPageSelector(got)).toEqual(selectionForLoadedStateWithTwoAlbums);
    });

    it("should use 'All albums' filter even when it's the only selection available (only directly owned albums) when receiving AlbumsAndMediasLoaded", () => {
        const action = albumsAndMediasLoaded({
            albums: [twoAlbums[0]],
            medias: someMediasByDays.flatMap(m => m.medias), // Pass raw medias
            mediasFromAlbumId: twoAlbums[0].albumId,
        });
        const got = action.reducer(initialCatalogState(myselfUser), action);

        const allAlbumFilter = {
            criterion: {
                owners: []
            },
            avatars: [myselfUser.picture ?? ""],
            name: "All albums",
        };
        expect(catalogViewerPageSelector(got)).toEqual({
            ...selectionForLoadedStateWithTwoAlbums,
            albums: [twoAlbums[0]],
            albumFilter: allAlbumFilter,
            albumFilterOptions: [allAlbumFilter],
            displayedAlbum: twoAlbums[0],
            medias: someMediasByDays,
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
            medias: someMediasByDays.flatMap(m => m.medias), // Pass raw medias
            mediasFromAlbumId: twoAlbums [0].albumId,
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
        expect(catalogViewerPageSelector(got)).toEqual({
            ...selectionForLoadedStateWithTwoAlbums,
            albumFilter: directlyOwnedFilter,
            albums: [loadedStateWithTwoAlbums.albums[0], newDirectlyOwnedAlbum],
            displayedAlbum: twoAlbums[0],
            medias: groupByDay(someMediasByDays.flatMap(m => m.medias)),
        });
    });

    it("change the filter so the loaded album is visible", () => {
        const directlyOwnedFilter = loadedStateWithTwoAlbums.albumFilterOptions[0];

        const action = albumsAndMediasLoaded({
            albums: twoAlbums,
            medias: someMediasByDays.flatMap(m => m.medias),
            mediasFromAlbumId: twoAlbums[1].albumId,
        });
        const got = action.reducer(
            {
                ...loadedStateWithTwoAlbums,
                albumFilter: directlyOwnedFilter,
                albums: [loadedStateWithTwoAlbums.albums[0]],
            },
            action
        );

        expect(catalogViewerPageSelector(got)).toEqual({
            ...selectionForLoadedStateWithTwoAlbums,
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[1],
            albums: twoAlbums,
            displayedAlbum: twoAlbums[1],
            medias: someMediasByDays,
        });
    });
});
