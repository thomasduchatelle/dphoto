import {noAlbumAvailable} from "./action-noAlbumAvailable";
import {CurrentUserInsight, initialCatalogState} from "../language";
import {albumListActionsPropsForLoadedState} from "../tests/test-helper-state";
import {catalogViewerPageSelector} from "./selector-catalog-viewer-page";
import {DEFAULT_ALBUM_FILTER_ENTRY} from "../common/utils";
import {albumListActionsSelector} from "./selector-albumListActions";

describe("action:noAlbumAvailable", () => {
    it("should return the state when no album is available", () => {
        const myselfUser: CurrentUserInsight = {picture: "my-face.jpg", isOwner: true};
        const action = noAlbumAvailable(undefined);
        const got = action.reducer(
            initialCatalogState(myselfUser),
            action
        );
        expect(catalogViewerPageSelector(got)).toEqual({
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
        expect(albumListActionsSelector(got)).toEqual({
            ...albumListActionsPropsForLoadedState,
            albumFilter: DEFAULT_ALBUM_FILTER_ENTRY,
            albumFilterOptions: [DEFAULT_ALBUM_FILTER_ENTRY],
            hasAlbumsToDelete: false,
            displayedAlbumIdIsOwned: false,
        });
    });
});
