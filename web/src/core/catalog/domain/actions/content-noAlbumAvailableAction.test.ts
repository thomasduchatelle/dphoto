import {noAlbumAvailableAction, reduceNoAlbumAvailable} from "./content-noAlbumAvailableAction";
import {CurrentUserInsight} from "../catalog-state";
import {loadedStateWithTwoAlbums} from "../tests/test-helper-state";
import {initialCatalogState} from "../initial-catalog-state";

describe("reduceNoAlbumAvailable", () => {
    it("should return the state when no album is available", () => {
        const myselfUser: CurrentUserInsight = {picture: "my-face.jpg"};
        const got = reduceNoAlbumAvailable(
            loadedStateWithTwoAlbums,
            noAlbumAvailableAction()
        );
        expect(got).toEqual({
            ...initialCatalogState(myselfUser),
            "albumNotFound": true,
            "albumsLoaded": true,
            "mediasLoaded": true,
        });
    });
});
