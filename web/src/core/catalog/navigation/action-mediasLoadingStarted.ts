import {AlbumId, CatalogViewerState} from "../language";

export interface MediasLoadingStarted {
    type: 'mediasLoadingStarted'
    albumId: AlbumId
}

export function mediasLoadingStarted(albumId: AlbumId): MediasLoadingStarted {
    return {type: "mediasLoadingStarted", albumId};
}

export function reduceMediasLoadingStarted(
    current: CatalogViewerState,
    action: MediasLoadingStarted
): CatalogViewerState {
    return {
        ...current,
        medias: [],
        error: undefined,
        loadingMediasFor: action.albumId,
        albumNotFound: false,
        mediasLoaded: false,
    };
}

export function mediasLoadingStartedReducerRegistration(handlers: any) {
    handlers["mediasLoadingStarted"] = reduceMediasLoadingStarted as (
        state: CatalogViewerState,
        action: MediasLoadingStarted
    ) => CatalogViewerState;
}
