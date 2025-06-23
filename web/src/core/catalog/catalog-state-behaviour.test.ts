import {herselfUser, loadedStateWithTwoAlbums, twoAlbums} from "./tests/test-helper-state";
import {Album, albumIdEquals, UserDetails} from "./language";
import {albumsFiltered, SELF_OWNED_ALBUM_FILTER_CRITERION} from "./navigation";
import {albumAccessGranted, albumAccessRevoked, sharingDialogSelector, sharingModalClosed, sharingModalErrorOccurred, sharingModalOpened} from "./sharing";
import {catalogReducer} from "./actions";

describe("State: behaviour", () => {
    it("keeps the album shares consistent when closing and reopening the dialog", () => {
        const initialState = loadedStateWithTwoAlbums;

        const album = twoAlbums[0];
        const openAction = sharingModalOpened(album.albumId);
        let state = catalogReducer(initialState, openAction);

        const newUser: UserDetails = {email: "bob@example.com", name: "Bob", picture: "bob-face.jpg"};
        const addAction = albumAccessGranted({user: newUser});
        state = catalogReducer(state, addAction);

        const removeAction = albumAccessRevoked(herselfUser.email);
        state = catalogReducer(state, removeAction);

        const closeAction = sharingModalClosed();
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
        const openAction = sharingModalOpened(album.albumId);
        let state = catalogReducer(initialState, openAction);

        const newUser: UserDetails = {email: "bob@example.com", name: "Bob", picture: "bob-face.jpg"};
        const addAction = albumAccessGranted({user: newUser});
        state = catalogReducer(state, addAction);

        const removeAction = albumAccessRevoked(herselfUser.email);
        state = catalogReducer(state, removeAction);

        const closeAction = sharingModalClosed();
        state = catalogReducer(state, closeAction);

        const filterAction = albumsFiltered({criterion: SELF_OWNED_ALBUM_FILTER_CRITERION});
        state = catalogReducer(state, filterAction);

        // Only the updated album should be present, and its shares should be correct
        const filteredAlbum = state.albums.find(a => albumIdEquals(a.albumId, album.albumId)) as Album;
        expect(filteredAlbum.sharedWith).toEqual([
            {user: newUser}
        ]);
    });

    it("reverts the changes made by addSharingAction, dispatched optimistically, when handling an error", () => {
        const userToGrant: UserDetails = {email: "bob@example.com", name: "Bob", picture: "bob-face.jpg"};

        let state = catalogReducer(loadedStateWithTwoAlbums, sharingModalOpened(twoAlbums[0].albumId));
        state = catalogReducer(state, albumAccessGranted({user: userToGrant}));
        state = catalogReducer(state, sharingModalErrorOccurred({type: "grant", message: "Failed to add user", email: userToGrant.email}));
        state = catalogReducer(state, sharingModalClosed());

        const filterAction = albumsFiltered({criterion: SELF_OWNED_ALBUM_FILTER_CRITERION});
        const ownedAlbums = catalogReducer(state, filterAction).albums;

        expect(state).toEqual(loadedStateWithTwoAlbums);
        expect(ownedAlbums).toEqual([loadedStateWithTwoAlbums.albums[0]]);
    });
});
