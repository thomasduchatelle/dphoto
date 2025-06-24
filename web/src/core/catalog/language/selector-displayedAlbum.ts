import {AlbumId, albumIsOwnedByCurrentUser, CatalogViewerState} from "./catalog-state";

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
        album.albumId.owner === targetAlbumId.owner &&
        album.albumId.folderName === targetAlbumId.folderName
    );

    if (!selectedAlbum) {
        return {displayedAlbumIdIsOwned: false};
    }

    return {
        displayedAlbumId: selectedAlbum.albumId,
        displayedAlbumIdIsOwned: albumIsOwnedByCurrentUser(selectedAlbum),
    };
}
