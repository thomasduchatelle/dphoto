import {mediaLoadFailed} from "./action-mediaLoadFailed";
import {loadedStateWithTwoAlbums, myselfUser, twoAlbums} from "../tests/test-helper-state";

import {initialCatalogState} from "../language";

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
        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            medias: [],
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
        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            medias: [],
            mediasLoadedFromAlbumId: undefined,
            error: testError,
        });
    });
});
