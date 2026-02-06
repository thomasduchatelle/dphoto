import {editNameDialogOpened} from "./action-editNameDialogOpened";
import {editNameDialogSelector} from "./selector-editNameDialogSelector";
import {editJanAlbumNameSelection, loadedStateWithTwoAlbums, myselfUser, twoAlbums} from "../tests/test-helper-state";
import {initialCatalogState} from "../language";

describe('action:editNameDialogOpened', () => {
    const albumId = twoAlbums[0].albumId;
    const albumName = twoAlbums[0].name;

    it('should open the edit name dialog with album details', () => {
        const state = {
            ...loadedStateWithTwoAlbums,
        };

        const action = editNameDialogOpened();
        const got = action.reducer(state, action);

        expect(editNameDialogSelector(got)).toEqual({
            ...editJanAlbumNameSelection,
        });
    });

    it('should not change state when no displayed album', () => {
        const state = initialCatalogState(myselfUser);

        const action = editNameDialogOpened();
        const got = action.reducer(state, action);

        expect(state).toBe(state);
    });

    it('should not change state when displayed album is not found in allAlbums', () => {
        const nonExistentAlbumId = {owner: "unknown", folderName: "unknown"};
        const state = {
            ...loadedStateWithTwoAlbums,
            mediasLoadedFromAlbumId: nonExistentAlbumId,
        };

        const action = editNameDialogOpened();
        const got = action.reducer(state, action);

        expect(state).toBe(state);
    });
});
