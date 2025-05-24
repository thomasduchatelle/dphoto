import { CatalogViewerState, AlbumId } from "./catalog-state";

/**
 * StartLoadingMediasAction interface
 */
export interface StartLoadingMediasAction {
    type: 'StartLoadingMediasAction'
    albumId: AlbumId
}

/**
 * Reducer fragment for StartLoadingMediasAction.
 * Resets medias, error, and sets loadingMediasFor, albumNotFound, mediasLoaded.
 */
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

/**
 * Action creator for StartLoadingMediasAction.
 */
export function startLoadingMediasAction(albumId: AlbumId): StartLoadingMediasAction {
    return { type: "StartLoadingMediasAction", albumId };
}
