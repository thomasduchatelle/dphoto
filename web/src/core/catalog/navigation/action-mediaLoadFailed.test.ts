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
            selectedAlbum: twoAlbums[0],
        });
        const got = action.reducer(
            initialCatalogState(myselfUser),
            action
        );
        expect(catalogViewerPageSelector(got)).toEqual({
            ...selectionForLoadedStateWithTwoAlbums,
            albums: twoAlbums, // albums are loaded even if medias fail
            albumsLoaded: true, // albums are loaded even if medias fail
            medias: [],
            displayedAlbum: twoAlbums[0],
            mediasLoaded: true,
            mediasLoadedFromAlbumId: undefined,
            error: testError,
        });
    });

    it("should set the errors and clears medias and media loading status when reducing MediaLoadFailed that hasn't the albums", () => {
        const testError = new Error("TEST loading error");
        const action = mediaLoadFailed({
            error: testError,
            selectedAlbum: twoAlbums[0],
        });
        const got = action.reducer(
            {
                ...loadedStateWithTwoAlbums,
                allAlbums: twoAlbums,
            },
            action
        );
        expect(catalogViewerPageSelector(got)).toEqual({
            ...selectionForLoadedStateWithTwoAlbums,
            medias: [],
            mediasLoadedFromAlbumId: undefined,
            error: testError,
            albumsLoaded: true,
            mediasLoaded: true,
        });
    });
});
