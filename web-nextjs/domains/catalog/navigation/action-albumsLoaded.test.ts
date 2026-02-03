import {albumsLoaded} from "./action-albumsLoaded";
import {
    albumListActionsPropsForLoadedState,
    createDialogPrefilledForMar25,
    loadedStateWithTwoAlbums,
    selectionForLoadedStateWithTwoAlbums,
    twoAlbums
} from "../tests/test-helper-state";
import {catalogViewerPageSelector} from "./selector-catalog-viewer-page";
import {albumListActionsSelector} from "./selector-albumListActions";

describe("action:albumsLoaded", () => {


    it("should update the list of albums and clear errors when AlbumsLoaded is received", () => {
        const action = albumsLoaded({albums: twoAlbums});
        const got = action.reducer({
            ...loadedStateWithTwoAlbums,
            allAlbums: [twoAlbums[0]],
            albums: [twoAlbums[0]],
            error: new Error("TEST previous error to clear"),
            albumsLoaded: false,
        }, action);

        expect(catalogViewerPageSelector(got)).toEqual(selectionForLoadedStateWithTwoAlbums);
        expect(albumListActionsSelector(got)).toEqual(albumListActionsPropsForLoadedState);
    });

    it("should update the available filters and re-apply the selected filter when receiving AlbumsLoaded", () => {
        const action = albumsLoaded({albums: twoAlbums, redirectTo: twoAlbums[0].albumId});
        const got = action.reducer({
            ...loadedStateWithTwoAlbums,
            albumFilterOptions: [loadedStateWithTwoAlbums.albumFilterOptions[0]],
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[0],
            allAlbums: [twoAlbums[1]],
            albums: [],
            mediasLoadedFromAlbumId: twoAlbums[0].albumId, // no effect
        }, action);

        expect(catalogViewerPageSelector(got)).toEqual({
            ...selectionForLoadedStateWithTwoAlbums,
            albums: [twoAlbums[0]],
        });
        expect(albumListActionsSelector(got)).toEqual({
            ...albumListActionsPropsForLoadedState,
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[0],
        });
    });

    it("should remove the album filter if the redirectTo in AlbumsLoaded wouldn't be displayed", () => {
        const action = albumsLoaded({albums: twoAlbums, redirectTo: twoAlbums[1].albumId});
        const got = action.reducer({
            ...loadedStateWithTwoAlbums,
            allAlbums: [twoAlbums[0]],
            albums: [],
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[0],
        }, action);

        expect(catalogViewerPageSelector(got)).toEqual({
            ...selectionForLoadedStateWithTwoAlbums,
            albums: twoAlbums,
        });
        expect(albumListActionsSelector(got)).toEqual({
            ...albumListActionsPropsForLoadedState,
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[1],
        });
    });

    it("should close all dialogs when albums are reloaded", () => {
        const action = albumsLoaded({albums: twoAlbums});
        const got = action.reducer({
            ...loadedStateWithTwoAlbums,
            dialog: createDialogPrefilledForMar25,
        }, action);

        expect(got.dialog).toBeUndefined();
    });
});
