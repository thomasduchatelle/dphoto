import {mediasLoadingStarted} from "./action-mediasLoadingStarted";
import {loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";

describe("action:mediasLoadingStarted", () => {
    it("should only change the medias and loading status when reducing MediasLoadingStarted, and clear errors", () => {
        const action = mediasLoadingStarted(twoAlbums[1].albumId);
        const got = action.reducer(
            {
                ...loadedStateWithTwoAlbums,
                albumNotFound: true,
                error: new Error("TEST previous error to clear"),
            },
            action
        );
        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            medias: [],
            loadingMediasFor: twoAlbums[1].albumId,
            mediasLoaded: false,
        });
    });
});
