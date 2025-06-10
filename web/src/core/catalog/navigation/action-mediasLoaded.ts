import {AlbumId, albumIdEquals, CatalogViewerState, MediaWithinADay} from "../language";

export interface MediasLoaded {
    type: "mediasLoaded"
    albumId: AlbumId
    medias: MediaWithinADay[]
}

export function mediasLoaded(props: Omit<MediasLoaded, "type">): MediasLoaded {
    return {
        ...props,
        type: "mediasLoaded",
    };
}

export function reduceMediasLoaded(
    current: CatalogViewerState,
    action: MediasLoaded,
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
    handlers["mediasLoaded"] = reduceMediasLoaded as (state: CatalogViewerState, action: MediasLoaded) => CatalogViewerState;
}
