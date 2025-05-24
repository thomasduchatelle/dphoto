import {reduceStartLoadingMedias, startLoadingMediasAction} from "./action-startLoadingMediasAction";
import {loadedStateWithTwoAlbums, twoAlbums} from "./tests/test-helper-state";

describe("reduceStartLoadingMedias", () => {
    it("should only change the medias and loading status when reducing StartLoadingMediasAction, and clear errors", () => {
        const got = reduceStartLoadingMedias(
            {
                ...loadedStateWithTwoAlbums,
                albumNotFound: true,
                error: new Error("TEST previous error to clear"),
            },
            startLoadingMediasAction(twoAlbums[1].albumId)
        );
        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            medias: [],
            loadingMediasFor: twoAlbums[1].albumId,
            mediasLoaded: false,
        });
    });
});
