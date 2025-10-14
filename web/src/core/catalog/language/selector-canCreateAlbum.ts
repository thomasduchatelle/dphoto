import {albumIsOwnedByCurrentUser, CatalogViewerState} from "./catalog-state";

export interface CanCreateAlbumSelection {
    canCreateAlbum: boolean;
}

export function canCreateAlbumSelector(state: CatalogViewerState): CanCreateAlbumSelection {
    const hasOwnedAlbum = state.allAlbums.some(album => albumIsOwnedByCurrentUser(album));

    return {
        canCreateAlbum: hasOwnedAlbum,
    };
}
