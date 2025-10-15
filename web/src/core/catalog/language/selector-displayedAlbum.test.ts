import {displayedAlbumSelector} from "./selector-displayedAlbum";
import {albumsLoaded} from "../navigation/action-albumsLoaded";
import {loadedStateWithTwoAlbums, twoAlbums, herselfOwner} from "../tests/test-helper-state";
import {Album, CatalogViewerState} from "./catalog-state";

const jan25Album = twoAlbums[0];
const feb25Album = twoAlbums[1];

describe("selector:displayedAlbumSelector", () => {
    it("should return canDeleteAlbum as true when at least one album is owned by current user", () => {
        const action = albumsLoaded({albums: twoAlbums});
        const state = action.reducer(loadedStateWithTwoAlbums, action);

        const selection = displayedAlbumSelector(state);

        expect(selection.canDeleteAlbum).toBe(true);
        expect(selection.displayedAlbumIdIsOwned).toBe(true);
        expect(selection.displayedAlbumId).toEqual(jan25Album.albumId);
    });

    it("should return canDeleteAlbum as false when no albums are owned by current user", () => {
        const albumOwnedByOther: Album = {
            ...jan25Album,
            ownedBy: {name: "Other Owner", users: [{name: "Other", email: "other@example.com"}]},
        };

        const action = albumsLoaded({albums: [albumOwnedByOther]});
        const state = action.reducer({
            ...loadedStateWithTwoAlbums,
            mediasLoadedFromAlbumId: albumOwnedByOther.albumId,
        }, action);

        const selection = displayedAlbumSelector(state);

        expect(selection.canDeleteAlbum).toBe(false);
        expect(selection.displayedAlbumIdIsOwned).toBe(false);
    });

    it("should return canDeleteAlbum as true even when displayed album is not owned", () => {
        const action = albumsLoaded({albums: twoAlbums});
        const state = action.reducer({
            ...loadedStateWithTwoAlbums,
            mediasLoadedFromAlbumId: feb25Album.albumId,
        }, action);

        const selection = displayedAlbumSelector(state);

        expect(selection.canDeleteAlbum).toBe(true);
        expect(selection.displayedAlbumIdIsOwned).toBe(false);
        expect(selection.displayedAlbumId).toEqual(feb25Album.albumId);
    });

    it("should return canDeleteAlbum as false when there are no albums", () => {
        const action = albumsLoaded({albums: []});
        const state = action.reducer({
            ...loadedStateWithTwoAlbums,
            mediasLoadedFromAlbumId: undefined,
        }, action);

        const selection = displayedAlbumSelector(state);

        expect(selection.canDeleteAlbum).toBe(false);
        expect(selection.displayedAlbumIdIsOwned).toBe(false);
        expect(selection.displayedAlbumId).toBeUndefined();
    });

    it("should return canDeleteAlbum based on all albums even when no album is displayed", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            mediasLoadedFromAlbumId: undefined,
            loadingMediasFor: undefined,
        };

        const selection = displayedAlbumSelector(state);

        expect(selection.canDeleteAlbum).toBe(true);
        expect(selection.displayedAlbumIdIsOwned).toBe(false);
        expect(selection.displayedAlbumId).toBeUndefined();
    });
});
