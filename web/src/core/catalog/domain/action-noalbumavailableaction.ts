import { CatalogViewerState } from "./catalog-state";
import { DEFAULT_ALBUM_FILTER_ENTRY } from "./catalog-common-modifiers";

/**
 * NoAlbumAvailableAction interface
 */
export interface NoAlbumAvailableAction {
    type: 'NoAlbumAvailableAction'
}

/**
 * Reducer fragment for NoAlbumAvailableAction.
 * Returns the state when no album is available.
 */
export function reduceNoAlbumAvailable(
    current: CatalogViewerState,
    action: NoAlbumAvailableAction
): CatalogViewerState {
    return {
        currentUser: current.currentUser,
        albumNotFound: true,
        allAlbums: [],
        albums: [],
        medias: [],
        albumsLoaded: true,
        mediasLoaded: true,
        albumFilterOptions: [DEFAULT_ALBUM_FILTER_ENTRY],
        albumFilter: DEFAULT_ALBUM_FILTER_ENTRY,
    };
}

/**
 * Action creator for NoAlbumAvailableAction.
 */
export function noAlbumAvailableAction(): NoAlbumAvailableAction {
    return { type: "NoAlbumAvailableAction" };
}
