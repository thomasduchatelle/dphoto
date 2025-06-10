import {mediasLoaded, reduceMediasLoaded} from "./action-mediasLoaded";
import {loadedStateWithTwoAlbums, someMedias, twoAlbums} from "../tests/test-helper-state";

describe("action:mediasLoaded", () => {
    it("should only change the medias and loading status when reducing MediasLoaded, and clear errors", () => {
        expect(reduceMediasLoaded({
            ...loadedStateWithTwoAlbums,
            medias: [],
            mediasLoaded: false,
            loadingMediasFor: twoAlbums[1].albumId,
            albumNotFound: true,
            error: new Error("TEST previous error to clear"),
        }, mediasLoaded({
            albumId: twoAlbums[1].albumId,
            medias: someMedias,
        }))).toEqual({
            ...loadedStateWithTwoAlbums,
            medias: someMedias,
            mediasLoadedFromAlbumId: twoAlbums[1].albumId,
        })
    });

    it("should ignore MediasLoaded if the medias are not for the expected album", () => {
        expect(reduceMediasLoaded({
            ...loadedStateWithTwoAlbums,
            medias: [],
            loadingMediasFor: twoAlbums[0].albumId,
        }, mediasLoaded({
            albumId: twoAlbums[1].albumId,
            medias: someMedias,
        }))).toEqual({
            ...loadedStateWithTwoAlbums,
            medias: [],
            loadingMediasFor: twoAlbums[0].albumId,
        })
    });
});
