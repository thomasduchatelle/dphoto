import {mediasLoaded} from "./action-mediasLoaded";
import {loadedStateWithTwoAlbums, someMedias, twoAlbums} from "../tests/test-helper-state";

describe("action:mediasLoaded", () => {
    it("should only change the medias and loading status when reducing MediasLoaded, and clear errors", () => {
        const action = mediasLoaded({
            albumId: twoAlbums[1].albumId,
            medias: someMedias,
        });
        expect(action.reducer({
            ...loadedStateWithTwoAlbums,
            medias: [],
            mediasLoaded: false,
            loadingMediasFor: twoAlbums[1].albumId,
            albumNotFound: true,
            error: new Error("TEST previous error to clear"),
        }, action)).toEqual({
            ...loadedStateWithTwoAlbums,
            medias: someMedias,
            mediasLoadedFromAlbumId: twoAlbums[1].albumId,
        })
    });

    it("should ignore MediasLoaded if the medias are not for the expected album", () => {
        const action = mediasLoaded({
            albumId: twoAlbums[1].albumId,
            medias: someMedias,
        });
        expect(action.reducer({
            ...loadedStateWithTwoAlbums,
            medias: [],
            loadingMediasFor: twoAlbums[0].albumId,
        }, action)).toEqual({
            ...loadedStateWithTwoAlbums,
            medias: [],
            loadingMediasFor: twoAlbums[0].albumId,
        })
    });
});
