import {albumIsOwnedByCurrentUser, CatalogViewerState} from "./catalog-state";

export interface CanCreateAlbumSelection {
    canCreateAlbum: boolean;
}

export function canCreateAlbumSelector(state: CatalogViewerState): CanCreateAlbumSelection {
    if (state.allAlbums.length === 0) {
        return {canCreateAlbum: true};
    }
    const hasOwnedAlbum = state.allAlbums.some(album => albumIsOwnedByCurrentUser(album));
    return {
        canCreateAlbum: hasOwnedAlbum,
    };
}
