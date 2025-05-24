import {AlbumId, CatalogViewerState} from "./catalog-state";

export interface StartLoadingMediasAction {
    type: 'StartLoadingMediasAction'
    albumId: AlbumId
}

export function startLoadingMediasAction(albumId: AlbumId): StartLoadingMediasAction {
    return {type: "StartLoadingMediasAction", albumId};
}

export function reduceStartLoadingMedias(
    current: CatalogViewerState,
    action: StartLoadingMediasAction
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
