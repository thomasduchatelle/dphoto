import {albumsFiltered} from "./action-albumsFiltered";
import {
    albumListActionsPropsForLoadedState,
    herselfOwner,
    loadedStateWithTwoAlbums,
    selectionForLoadedStateWithTwoAlbums,
    twoAlbums
} from "../tests/test-helper-state";
import {catalogViewerPageSelector} from "./selector-catalog-viewer-page";
import {ALL_ALBUMS_FILTER_CRITERION, SELF_OWNED_ALBUM_FILTER_CRITERION} from "../common/utils";
import {albumListActionsSelector} from "./selector-albumListActions";

describe("action:albumsFiltered", () => {
    it("should show only directly owned album after the AlbumsFiltered", () => {
        const action = albumsFiltered({criterion: SELF_OWNED_ALBUM_FILTER_CRITERION});
        const got = action.reducer(
            loadedStateWithTwoAlbums,
            action
        );
        expect(catalogViewerPageSelector(got)).toEqual({
            ...selectionForLoadedStateWithTwoAlbums,
            albums: [twoAlbums[0]],
        });
        expect(albumListActionsSelector(got)).toEqual({
            ...albumListActionsPropsForLoadedState,
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[0],
        });
    });

    it("should show all albums when the filter moves back to 'All albums'", () => {
        const action = albumsFiltered({criterion: ALL_ALBUMS_FILTER_CRITERION});
        const got = action.reducer(
            {
                ...loadedStateWithTwoAlbums,
                albums: [],
            },
            action
        );
        expect(catalogViewerPageSelector(got)).toEqual({
            ...selectionForLoadedStateWithTwoAlbums,
            albums: twoAlbums,
        });
        expect(albumListActionsSelector(got)).toEqual({
            ...albumListActionsPropsForLoadedState,
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[1],
        });
    });

    it("should filter albums to match selected owner", () => {
        const action = albumsFiltered({criterion: {owners: [herselfOwner]}});
        const got = action.reducer(
            loadedStateWithTwoAlbums,
            action
        );
        expect(catalogViewerPageSelector(got)).toEqual({
            ...selectionForLoadedStateWithTwoAlbums,
            albums: [twoAlbums[1]],
            // displayedAlbum: twoAlbums[1], // Displayed album remains the same
        });
        expect(albumListActionsSelector(got)).toEqual({
            ...albumListActionsPropsForLoadedState,
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[2],
        });
    });
});
