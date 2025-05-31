import {AlbumId, CatalogViewerState, MediaWithinADay} from "../catalog-state";
import {albumIdEquals} from "../utils-albumIdEquals";

export interface MediasLoadedAction {
    type: "MediasLoadedAction"
    albumId: AlbumId
    medias: MediaWithinADay[]
}

export function mediasLoadedAction(props: Omit<MediasLoadedAction, "type">): MediasLoadedAction {
    return {
        ...props,
        type: "MediasLoadedAction",
    };
}

export function reduceMediasLoaded(
    current: CatalogViewerState,
    action: MediasLoadedAction,
): CatalogViewerState {
    if (current.loadingMediasFor && !albumIdEquals(current.loadingMediasFor, action.albumId)) {
        // concurrency management - ignore if not the last album requested
        return current
    }

    return {
        ...current,
        loadingMediasFor: undefined,
        mediasLoadedFromAlbumId: action.albumId,
        medias: action.medias,
        error: undefined,
        mediasLoaded: true,
        albumNotFound: false,
    }
}

export function mediasLoadedReducerRegistration(handlers: any) {
    handlers["MediasLoadedAction"] = reduceMediasLoaded as (state: CatalogViewerState, action: MediasLoadedAction) => CatalogViewerState;
}
