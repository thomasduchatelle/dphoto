import {displayedAlbumSelector} from "./selector-displayedAlbum";
import {albumsLoaded} from "../navigation/action-albumsLoaded";
import {loadedStateWithTwoAlbums, twoAlbums, displayedAlbumForLoadedState} from "../tests/test-helper-state";
import {Album} from "./catalog-state";

const jan25Album = twoAlbums[0];
const feb25Album = twoAlbums[1];

describe("selector:displayedAlbumSelector", () => {
    it("should return canDeleteAlbum as true when at least one album is owned by current user", () => {
        const state = loadedStateWithTwoAlbums;

        const action = albumsLoaded({albums: twoAlbums});
        const got = action.reducer(state, action);

        expect(displayedAlbumSelector(got)).toEqual(displayedAlbumForLoadedState);
    });

    it("should return canDeleteAlbum as false when no albums are owned by current user", () => {
        const albumOwnedByOther: Album = {
            ...jan25Album,
            ownedBy: {name: "Other Owner", users: [{name: "Other", email: "other@example.com"}]},
        };

        const state = {
            ...loadedStateWithTwoAlbums,
            mediasLoadedFromAlbumId: albumOwnedByOther.albumId,
        };

        const action = albumsLoaded({albums: [albumOwnedByOther]});
        const got = action.reducer(state, action);

        expect(displayedAlbumSelector(got)).toEqual({
            ...displayedAlbumForLoadedState,
            canDeleteAlbum: false,
            displayedAlbumIdIsOwned: false,
        });
    });

    it("should return canDeleteAlbum as true even when displayed album is not owned", () => {
        const state = {
            ...loadedStateWithTwoAlbums,
            mediasLoadedFromAlbumId: feb25Album.albumId,
        };

        const action = albumsLoaded({albums: twoAlbums});
        const got = action.reducer(state, action);

        expect(displayedAlbumSelector(got)).toEqual({
            ...displayedAlbumForLoadedState,
            displayedAlbumId: feb25Album.albumId,
            displayedAlbumIdIsOwned: false,
        });
    });

    it("should return canDeleteAlbum as false when there are no albums", () => {
        const state = {
            ...loadedStateWithTwoAlbums,
            mediasLoadedFromAlbumId: undefined,
        };

        const action = albumsLoaded({albums: []});
        const got = action.reducer(state, action);

        expect(displayedAlbumSelector(got)).toEqual({
            ...displayedAlbumForLoadedState,
            displayedAlbumId: undefined,
            displayedAlbumIdIsOwned: false,
            canDeleteAlbum: false,
        });
    });

    it("should return canDeleteAlbum based on all albums even when no album is displayed", () => {
        const state = {
            ...loadedStateWithTwoAlbums,
            mediasLoadedFromAlbumId: undefined,
            loadingMediasFor: undefined,
        };

        const action = albumsLoaded({albums: twoAlbums});
        const got = action.reducer(state, action);

        expect(displayedAlbumSelector(got)).toEqual({
            ...displayedAlbumForLoadedState,
            displayedAlbumId: undefined,
            displayedAlbumIdIsOwned: false,
        });
    });
});
