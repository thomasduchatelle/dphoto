import {noAlbumAvailable, reduceNoAlbumAvailable} from "./action-noAlbumAvailable";
import {CurrentUserInsight, initialCatalogState} from "../language";
import {loadedStateWithTwoAlbums} from "../tests/test-helper-state";

describe("action:noAlbumAvailable", () => {
    it("should return the state when no album is available", () => {
        const myselfUser: CurrentUserInsight = {picture: "my-face.jpg"};
        const got = reduceNoAlbumAvailable(
            loadedStateWithTwoAlbums,
            noAlbumAvailable()
        );
        expect(got).toEqual({
            ...initialCatalogState(myselfUser),
            "albumNotFound": true,
            "albumsLoaded": true,
            "mediasLoaded": true,
        });
    });
});
