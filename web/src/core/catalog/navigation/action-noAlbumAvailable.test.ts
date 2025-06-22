import {noAlbumAvailable} from "./action-noAlbumAvailable";
import {CurrentUserInsight, initialCatalogState} from "../language";
import {loadedStateWithTwoAlbums} from "../tests/test-helper-state";
import {catalogViewerPageSelector} from "./selector-catalog-viewer-page";
import {DEFAULT_ALBUM_FILTER_ENTRY} from "../common/utils";

describe("action:noAlbumAvailable", () => {
    it("should return the state when no album is available", () => {
        const myselfUser: CurrentUserInsight = {picture: "my-face.jpg"};
        const action = noAlbumAvailable();
        const got = action.reducer(
            loadedStateWithTwoAlbums,
            action
        );
        expect(catalogViewerPageSelector(got)).toEqual({
            albumFilter: DEFAULT_ALBUM_FILTER_ENTRY,
            albumFilterOptions: [DEFAULT_ALBUM_FILTER_ENTRY],
            albumsLoaded: true,
            albums: [],
            displayedAlbum: undefined,
            medias: [],
            mediasLoaded: true,
            mediasLoadedFromAlbumId: undefined,
            loadingMediasFor: undefined,
            albumNotFound: true,
            error: undefined,
        });
    });
});
