import {herselfUser, loadedStateWithTwoAlbums, twoAlbums} from "./tests/test-helper-state";
import {openSharingModalAction} from "./actions/sharing-openSharingModalAction";
import {Album, UserDetails} from "./catalog-state";
import {catalogActions, catalogReducer, SELF_OWNED_ALBUM_FILTER_CRITERION, sharingDialogSelector} from "./actions";
import {albumIdEquals} from "./utils-albumIdEquals";

describe("State: behaviour", () => {
    it("keeps the album shares consistent when closing and reopening the dialog", () => {
        const initialState = loadedStateWithTwoAlbums;

        const album = twoAlbums[0];
        const openAction = catalogActions.openSharingModalAction(album.albumId);
        let state = catalogReducer(initialState, openAction);

        const newUser: UserDetails = {email: "bob@example.com", name: "Bob", picture: "bob-face.jpg"};
        const addAction = catalogActions.addSharingAction({user: newUser});
        state = catalogReducer(state, addAction);

        const removeAction = catalogActions.removeSharingAction(herselfUser.email);
        state = catalogReducer(state, removeAction);

        const closeAction = catalogActions.closeSharingModalAction();
        state = catalogReducer(state, closeAction);

        state = catalogReducer(state, openAction);

        expect(sharingDialogSelector(state)).toEqual({
            open: true,
            sharedWith: [
                {user: newUser}
            ],
            suggestions: [],
        });

        const updatedAlbum = state.allAlbums.find(a => albumIdEquals(a.albumId, album.albumId)) as Album;
        expect(updatedAlbum.sharedWith).toEqual([
            {user: newUser}
        ]);
    });

    it("keeps the album shares consistent when re-applying the album filter", () => {
        const initialState = loadedStateWithTwoAlbums;

        const album = twoAlbums[0];
        const openAction = openSharingModalAction(album.albumId);
        let state = catalogReducer(initialState, openAction);

        const newUser: UserDetails = {email: "bob@example.com", name: "Bob", picture: "bob-face.jpg"};
        const addAction = catalogActions.addSharingAction({user: newUser});
        state = catalogReducer(state, addAction);

        const removeAction = catalogActions.removeSharingAction(herselfUser.email);
        state = catalogReducer(state, removeAction);

        const closeAction = catalogActions.closeSharingModalAction();
        state = catalogReducer(state, closeAction);

        const filterAction = catalogActions.albumsFilteredAction({criterion: SELF_OWNED_ALBUM_FILTER_CRITERION});
        state = catalogReducer(state, filterAction);

        // Only the updated album should be present, and its shares should be correct
        const filteredAlbum = state.albums.find(a => albumIdEquals(a.albumId, album.albumId)) as Album;
        expect(filteredAlbum.sharedWith).toEqual([
            {user: newUser}
        ]);
    });

    it("reverts the changes made by addSharingAction, dispatched optimistically, when handling an error", () => {
        const userToGrant: UserDetails = {email: "bob@example.com", name: "Bob", picture: "bob-face.jpg"};

        let state = catalogReducer(loadedStateWithTwoAlbums, catalogActions.openSharingModalAction(twoAlbums[0].albumId));
        state = catalogReducer(state, catalogActions.addSharingAction({user: userToGrant}));
        state = catalogReducer(state, catalogActions.sharingModalErrorAction({type: "grant", message: "Failed to add user", email: userToGrant.email}));
        state = catalogReducer(state, catalogActions.closeSharingModalAction());

        const filterAction = catalogActions.albumsFilteredAction({criterion: SELF_OWNED_ALBUM_FILTER_CRITERION});
        const ownedAlbums = catalogReducer(state, filterAction).albums;

        expect(state).toEqual(loadedStateWithTwoAlbums);
        expect(ownedAlbums).toEqual([loadedStateWithTwoAlbums.albums[0]]);
    });
});
