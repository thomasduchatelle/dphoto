import {CatalogViewerState} from "../language";
import {DEFAULT_ALBUM_FILTER_ENTRY} from "../common/utils";

export interface NoAlbumAvailable {
    type: 'noAlbumAvailable'
}

export function reduceNoAlbumAvailable(
    current: CatalogViewerState,
    action: NoAlbumAvailable
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

export function noAlbumAvailable(): NoAlbumAvailable {
    return {type: "noAlbumAvailable"};
}

export function noAlbumAvailableReducerRegistration(handlers: any) {
    handlers["noAlbumAvailable"] = reduceNoAlbumAvailable as (
        state: CatalogViewerState,
        action: NoAlbumAvailable
    ) => CatalogViewerState;
}
