import { mediaFailedToLoadAction, reduceMediaFailedToLoad } from "./catalog-action-MediaFailedToLoadAction";
import { loadedStateWithTwoAlbums, twoAlbums, herselfUser } from "./tests/test-helper-state";

describe("reduceMediaFailedToLoad", () => {
    it("should set the errors and clears medias and media loading status when reducing MediaFailedToLoadAction", () => {
        const testError = new Error("TEST loading error");
        const got = reduceMediaFailedToLoad(
            loadedStateWithTwoAlbums,
            {
                error: testError,
                albums: twoAlbums,
                selectedAlbum: twoAlbums[0],
            }
        );
        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            medias: [],
            mediasLoadedFromAlbumId: undefined,
            error: testError,
        });
    });

    it("should set the errors and clears medias and media loading status when reducing MediaFailedToLoadAction that hasn't the albums", () => {
        const testError = new Error("TEST loading error");
        const got = reduceMediaFailedToLoad(
            {
                ...loadedStateWithTwoAlbums,
                allAlbums: twoAlbums,
            },
            {
                error: testError,
                selectedAlbum: twoAlbums[0],
            }
        );
        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            medias: [],
            mediasLoadedFromAlbumId: undefined,
            error: testError,
        });
    });
});
