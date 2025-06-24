import {mediaLoadFailed} from "./action-mediaLoadFailed";
import {loadedStateWithTwoAlbums, myselfUser, selectionForLoadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";

import {initialCatalogState} from "../language";
import {catalogViewerPageSelector} from "./selector-catalog-viewer-page";

describe("action:mediaLoadFailed", () => {
    it("should set the errors and clears medias and media loading status when reducing MediaLoadFailed", () => {
        const testError = new Error("TEST loading error");
        const action = mediaLoadFailed({
            error: testError,
            albums: twoAlbums,
            displayedAlbumId: twoAlbums[0].albumId,
        });
        const got = action.reducer(
            initialCatalogState(myselfUser),
            action
        );
        expect(catalogViewerPageSelector(got)).toEqual({
            ...selectionForLoadedStateWithTwoAlbums,
            medias: [],
            error: testError,
        });
    });

    it("should set the errors and clears medias and media loading status when reducing MediaLoadFailed on a already loaded state", () => {
        const testError = new Error("TEST loading error");
        const action = mediaLoadFailed({
            error: testError,
            displayedAlbumId: twoAlbums[0].albumId,
        });
        const got = action.reducer(
            loadedStateWithTwoAlbums,
            action
        );
        expect(catalogViewerPageSelector(got)).toEqual({
            ...selectionForLoadedStateWithTwoAlbums,
            medias: [],
            error: testError,
            albumsLoaded: true,
            mediasLoaded: true,
        });
    });

    it("change the filter so the loaded album is visible", () => {
        const testError = new Error("TEST loading error");
        const directlyOwnedFilter = loadedStateWithTwoAlbums.albumFilterOptions[0];

        const action = mediaLoadFailed({
            error: testError,
            albums: twoAlbums,
            displayedAlbumId: twoAlbums[1].albumId,
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
            medias: [],
            error: testError,
        });
    });
});
