import {mediaLoadFailed} from "./action-mediaLoadFailed";
import {loadedStateWithTwoAlbums, myselfUser, twoAlbums} from "../tests/test-helper-state";

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
        expect(catalogViewerPageSelector(got, twoAlbums[0].albumId)).toEqual({
            ...catalogViewerPageSelector(loadedStateWithTwoAlbums, twoAlbums[0].albumId),
            albums: twoAlbums, // albums are loaded even if medias fail
            albumsLoaded: true, // albums are loaded even if medias fail
            medias: [],
            selectedAlbum: twoAlbums[0],
        });
        expect(got.error).toEqual(testError); // Assert error directly as it's not part of the selector's output
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
        expect(catalogViewerPageSelector(got, twoAlbums[0].albumId)).toEqual({
            ...catalogViewerPageSelector(loadedStateWithTwoAlbums, twoAlbums[0].albumId),
            medias: [],
            selectedAlbum: twoAlbums[0],
        });
        expect(got.error).toEqual(testError); // Assert error directly as it's not part of the selector's output
    });
});
