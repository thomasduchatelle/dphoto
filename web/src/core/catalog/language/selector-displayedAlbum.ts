import {AlbumId, albumIsOwnedByCurrentUser, CatalogViewerState} from "./catalog-state";
import {albumIdEquals} from "./utils-albumIdEquals";

export interface DisplayedAlbumSelection {
    displayedAlbumId?: AlbumId;
    displayedAlbumIdIsOwned: boolean;
    canDeleteAlbum: boolean;
}

export function getDisplayedAlbumId(state: CatalogViewerState) {
    return state.loadingMediasFor || state.mediasLoadedFromAlbumId;
}

export function displayedAlbumSelector(state: CatalogViewerState): DisplayedAlbumSelection {
    const targetAlbumId: AlbumId | undefined = getDisplayedAlbumId(state);
    const canDeleteAlbum = state.allAlbums.some(album => albumIsOwnedByCurrentUser(album));

    if (!targetAlbumId) {
        return {displayedAlbumIdIsOwned: false, canDeleteAlbum};
    }

    const selectedAlbum = state.allAlbums.find(album =>
        albumIdEquals(album.albumId, targetAlbumId)
    );

    if (!selectedAlbum) {
        return {displayedAlbumIdIsOwned: false, canDeleteAlbum};
    }

    return {
        displayedAlbumId: selectedAlbum.albumId,
        displayedAlbumIdIsOwned: albumIsOwnedByCurrentUser(selectedAlbum),
        canDeleteAlbum,
    };
}
