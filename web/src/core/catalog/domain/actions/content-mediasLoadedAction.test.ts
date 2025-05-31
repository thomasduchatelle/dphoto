import {mediasLoadedAction, reduceMediasLoaded} from "./content-mediasLoadedAction";
import {loadedStateWithTwoAlbums, someMedias, twoAlbums} from "../tests/test-helper-state";

describe("reduceMediasLoaded", () => {
    it("should only change the medias and loading status when reducing MediasLoadedAction, and clear errors", () => {
        expect(reduceMediasLoaded({
            ...loadedStateWithTwoAlbums,
            medias: [],
            mediasLoaded: false,
            loadingMediasFor: twoAlbums[1].albumId,
            albumNotFound: true,
            error: new Error("TEST previous error to clear"),
        }, mediasLoadedAction({
            albumId: twoAlbums[1].albumId,
            medias: someMedias,
        }))).toEqual({
            ...loadedStateWithTwoAlbums,
            medias: someMedias,
            mediasLoadedFromAlbumId: twoAlbums[1].albumId,
        })
    });

    it("should ignore MediasLoadedAction if the medias are not for the expected album", () => {
        expect(reduceMediasLoaded({
            ...loadedStateWithTwoAlbums,
            medias: [],
            loadingMediasFor: twoAlbums[0].albumId,
        }, mediasLoadedAction({
            albumId: twoAlbums[1].albumId,
            medias: someMedias,
        }))).toEqual({
            ...loadedStateWithTwoAlbums,
            medias: [],
            loadingMediasFor: twoAlbums[0].albumId,
        })
    });
});
