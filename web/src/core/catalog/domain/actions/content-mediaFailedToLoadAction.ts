import {Album, CatalogViewerState} from "../catalog-state";
import {refreshFilters} from "./catalog-common-modifiers";

export interface MediaFailedToLoadAction {
    type: 'MediaFailedToLoadAction'
    albums?: Album[]
    selectedAlbum?: Album
    error: Error
}

export function mediaFailedToLoadAction(props: Omit<MediaFailedToLoadAction, "type">): MediaFailedToLoadAction {
    return {
        ...props,
        type: "MediaFailedToLoadAction",
    };
}

export function reduceMediaFailedToLoad(
    current: CatalogViewerState,
    action: MediaFailedToLoadAction,
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

export function mediaFailedToLoadReducerRegistration(handlers: any) {
    handlers["MediaFailedToLoadAction"] = reduceMediaFailedToLoad as (
        state: CatalogViewerState,
        action: MediaFailedToLoadAction
    ) => CatalogViewerState;
}
