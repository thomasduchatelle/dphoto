import {mediasLoaded} from "./action-mediasLoaded";
import {loadedStateWithTwoAlbums, someMedias, twoAlbums} from "../tests/test-helper-state";
import {groupByDay} from "./group-by-day";
import {catalogViewerPageSelector} from "./selector-catalog-viewer-page";

describe("action:mediasLoaded", () => {
    it("should only change the medias and loading status when reducing MediasLoaded, and clear errors", () => {
        const action = mediasLoaded({
            albumId: twoAlbums[1].albumId,
            medias: someMedias.flatMap(m => m.medias), // Pass raw medias
        });
        const got = action.reducer({
            ...loadedStateWithTwoAlbums,
            medias: [],
            mediasLoaded: false,
            loadingMediasFor: twoAlbums[1].albumId,
            albumNotFound: true,
            error: new Error("TEST previous error to clear"),
        }, action);

        expect(catalogViewerPageSelector(got, twoAlbums[1].albumId)).toEqual({
            ...catalogViewerPageSelector(loadedStateWithTwoAlbums, twoAlbums[1].albumId),
            medias: groupByDay(someMedias.flatMap(m => m.medias)),
            selectedAlbum: twoAlbums[1],
        });
        expect(got.mediasLoadedFromAlbumId).toEqual(twoAlbums[1].albumId);
        expect(got.error).toBeUndefined();
    });

    it("should ignore MediasLoaded if the medias are not for the expected album", () => {
        const action = mediasLoaded({
            albumId: twoAlbums[1].albumId,
            medias: someMedias.flatMap(m => m.medias), // Pass raw medias
        });
        const got = action.reducer({
            ...loadedStateWithTwoAlbums,
            medias: [],
            loadingMediasFor: twoAlbums[0].albumId,
        }, action);

        expect(catalogViewerPageSelector(got, twoAlbums[0].albumId)).toEqual({
            ...catalogViewerPageSelector(loadedStateWithTwoAlbums, twoAlbums[0].albumId),
            medias: [],
            selectedAlbum: twoAlbums[0],
        });
        expect(got.loadingMediasFor).toEqual(twoAlbums[0].albumId);
    });
});
