import { reduceStartLoadingMedias, startLoadingMediasAction } from "./action-startloadingmediasaction";
import { loadedStateWithTwoAlbums } from "./tests/test-helper-state";

describe("reduceStartLoadingMedias", () => {
    it("should only change the medias and loading status when reducing StartLoadingMediasAction, and clear errors", () => {
        const twoAlbums = loadedStateWithTwoAlbums.allAlbums;
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
