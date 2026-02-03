import {AlbumId, albumIsOwnedByCurrentUser, CatalogViewerState} from "./catalog-state";
import {albumIdEquals} from "./utils-albumIdEquals";

export interface DisplayedAlbumSelection {
    displayedAlbumId?: AlbumId;
    displayedAlbumIdIsOwned: boolean;
}

export function getDisplayedAlbumId(state: CatalogViewerState) {
    return state.loadingMediasFor || state.mediasLoadedFromAlbumId;
}

export function displayedAlbumSelector(state: CatalogViewerState): DisplayedAlbumSelection {
    const targetAlbumId: AlbumId | undefined = getDisplayedAlbumId(state);

    if (!targetAlbumId) {
        return {displayedAlbumIdIsOwned: false};
    }

    const selectedAlbum = state.allAlbums.find(album =>
        albumIdEquals(album.albumId, targetAlbumId)
    );

    if (!selectedAlbum) {
        return {displayedAlbumIdIsOwned: false};
    }

    return {
        displayedAlbumId: selectedAlbum.albumId,
        displayedAlbumIdIsOwned: albumIsOwnedByCurrentUser(selectedAlbum),
    };
}
