import {CatalogViewerState} from "./catalog-state";
import {DEFAULT_ALBUM_FILTER_ENTRY} from "./catalog-common-modifiers";

export interface NoAlbumAvailableAction {
    type: 'NoAlbumAvailableAction'
}

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

export function noAlbumAvailableAction(): NoAlbumAvailableAction {
    return { type: "NoAlbumAvailableAction" };
}

export function noAlbumAvailableReducerRegistration(handlers: any) {
    handlers["NoAlbumAvailableAction"] = reduceNoAlbumAvailable as (
        state: CatalogViewerState,
        action: NoAlbumAvailableAction
    ) => CatalogViewerState;
}
