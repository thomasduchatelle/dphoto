import {AlbumId, albumIsOwnedByCurrentUser, CatalogViewerState} from "./catalog-state";

export interface DisplayedAlbumSelection {
    albumId?: AlbumId;
    isOwned: boolean;
}

export function displayedAlbumSelector(state: CatalogViewerState): DisplayedAlbumSelection {
    const targetAlbumId: AlbumId | undefined = state.loadingMediasFor || state.mediasLoadedFromAlbumId;

    if (!targetAlbumId) {
        return {isOwned: false};
    }

    const selectedAlbum = state.albums.find(album =>
        album.albumId.owner === targetAlbumId.owner &&
        album.albumId.folderName === targetAlbumId.folderName
    );

    if (!selectedAlbum) {
        return {isOwned: false};
    }

    return {
        albumId: selectedAlbum.albumId,
        isOwned: albumIsOwnedByCurrentUser(selectedAlbum),
    };
}
