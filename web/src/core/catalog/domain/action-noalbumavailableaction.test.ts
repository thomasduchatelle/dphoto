import { reduceNoAlbumAvailable, noAlbumAvailableAction } from "./action-noalbumavailableaction";
import { initialCatalogState } from "./catalog-reducer";
import { CurrentUserInsight } from "./catalog-state";
import { DEFAULT_ALBUM_FILTER_ENTRY } from "./catalog-common-modifiers";

describe("reduceNoAlbumAvailable", () => {
    it("should return the state when no album is available", () => {
        const myselfUser: CurrentUserInsight = { picture: "my-face.jpg" };
        const got = reduceNoAlbumAvailable(
            initialCatalogState(myselfUser),
            noAlbumAvailableAction()
        );
        expect(got).toEqual({
            currentUser: myselfUser,
            albumNotFound: true,
            allAlbums: [],
            albums: [],
            medias: [],
            albumsLoaded: true,
            mediasLoaded: true,
            albumFilterOptions: [DEFAULT_ALBUM_FILTER_ENTRY],
            albumFilter: DEFAULT_ALBUM_FILTER_ENTRY,
        });
    });
});
