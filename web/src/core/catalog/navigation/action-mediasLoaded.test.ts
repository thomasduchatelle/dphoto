import {mediasLoaded} from "./action-mediasLoaded";
import {loadedStateWithTwoAlbums, selectionForLoadedStateWithTwoAlbums, someMedias, twoAlbums} from "../tests/test-helper-state";
import {catalogViewerPageSelector} from "./selector-catalog-viewer-page";

describe("action:mediasLoaded", () => {
    it("should change the medias and loading status when reducing MediasLoaded, and clear errors", () => {
        const action = mediasLoaded({
            albumId: twoAlbums[1].albumId,
            medias: someMedias.flatMap(m => m.medias),
        });
        const got = action.reducer({
            ...loadedStateWithTwoAlbums,
            medias: [],
            mediasLoaded: false,
            loadingMediasFor: twoAlbums[1].albumId,
            albumNotFound: true,
            error: new Error("TEST previous error to clear"),
        }, action);

        expect(catalogViewerPageSelector(got)).toEqual({
            ...selectionForLoadedStateWithTwoAlbums,
            displayedAlbum: twoAlbums[1],
            medias: someMedias,
        });
    });

    it("should ignore MediasLoaded if the medias are not for the expected album", () => {
        const action = mediasLoaded({
            albumId: twoAlbums[1].albumId,
            medias: someMedias.flatMap(m => m.medias),
        });
        const state = {
            ...loadedStateWithTwoAlbums,
            medias: [],
            loadingMediasFor: twoAlbums[0].albumId,
        };
        const got = action.reducer(state, action);

        expect(got).toBe(state);
    });
});
