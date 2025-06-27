import {albumRenamed} from "./action-albumRenamed";
import {editNameDialogSelector} from "./selector-editNameDialogSelector";
import {catalogViewerPageSelector, SELF_OWNED_ALBUM_FILTER_CRITERION} from "../navigation";
import {albumsFiltered} from "../navigation/action-albumsFiltered";
import {CatalogViewerState} from "../language";
import {editJanAlbumNameDialog, loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";

describe('action:albumRenamed', () => {
    const originalAlbumId = twoAlbums[0].albumId;
    const newAlbumId = {owner: "myself", folderName: "new-jan-25"};
    const newName = "Updated January 2025";

    it('should update the album in the state and close the dialog when only the name is changed', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: editJanAlbumNameDialog,
        };

        const action = albumRenamed({previousAlbumId: originalAlbumId, newAlbumId: originalAlbumId, newName});
        const got = action.reducer(state, action);

        expect(editNameDialogSelector(got).isOpen).toEqual(false);

        const pageSelection = catalogViewerPageSelector(got);
        expect(pageSelection.albums).toEqual([
            {...twoAlbums[0], albumId: originalAlbumId, name: newName},
            twoAlbums[1]
        ]);
    });

    it('should update the album id on the albums state, and the mediaLoaded if using the previous id', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            mediasLoadedFromAlbumId: originalAlbumId,
        };

        const action = albumRenamed({previousAlbumId: originalAlbumId, newAlbumId, newName});
        const got = action.reducer(state, action);

        const pageSelection = catalogViewerPageSelector(got);
        expect(pageSelection.albums).toEqual([
            {...twoAlbums[0], albumId: newAlbumId, name: newName},
            twoAlbums[1]
        ]);
        
        expect(got.mediasLoadedFromAlbumId).toEqual(newAlbumId);
    });

    it('should not update mediasLoadedFromAlbumId when it does not match previous album', () => {
        const differentAlbumId = {owner: "other", folderName: "other-album"};
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            mediasLoadedFromAlbumId: differentAlbumId,
        };

        const action = albumRenamed({previousAlbumId: originalAlbumId, newAlbumId, newName});
        const got = action.reducer(state, action);

        expect(got.mediasLoadedFromAlbumId).toEqual(differentAlbumId);
    });

    it('should update name and id of the album even after being re-filtered', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
        };

        const renameAction = albumRenamed({previousAlbumId: originalAlbumId, newAlbumId, newName});
        const afterRename = renameAction.reducer(state, renameAction);

        const filterAction = albumsFiltered({criterion: SELF_OWNED_ALBUM_FILTER_CRITERION});
        const afterFilter = filterAction.reducer(afterRename, filterAction);

        const pageSelection = catalogViewerPageSelector(afterFilter);
        expect(pageSelection.albums).toEqual([
            {...twoAlbums[0], albumId: newAlbumId, name: newName},
        ]);
    });
});
