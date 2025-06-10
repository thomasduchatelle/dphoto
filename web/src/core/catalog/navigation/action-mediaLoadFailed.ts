import {Album, CatalogViewerState} from "../language";
import {refreshFilters} from "../common/utils";

export interface MediaLoadFailed {
    type: 'mediaLoadFailed'
    albums?: Album[]
    selectedAlbum?: Album
    error: Error
}

export function mediaLoadFailed(props: Omit<MediaLoadFailed, "type">): MediaLoadFailed {
    return {
        ...props,
        type: "mediaLoadFailed",
    };
}

export function reduceMediaLoadFailed(
    current: CatalogViewerState,
    action: MediaLoadFailed,
): CatalogViewerState {
    const allAlbums = action.albums ?? current.allAlbums;

    const {albumFilterOptions, albumFilter, albums} = refreshFilters(current.currentUser, current.albumFilter, allAlbums);

    return {
        currentUser: current.currentUser,
        allAlbums,
        albumFilterOptions,
        albumFilter,
        albums,
        albumNotFound: false,
        medias: [],
        error: action.error,
        albumsLoaded: true,
        mediasLoaded: true,
    };
}

export function mediaLoadFailedReducerRegistration(handlers: any) {
    handlers["mediaLoadFailed"] = reduceMediaLoadFailed as (
        state: CatalogViewerState,
        action: MediaLoadFailed
    ) => CatalogViewerState;
}
