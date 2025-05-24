import {Album, albumMatchCriterion, CatalogViewerState} from "./catalog-state";
import {albumFilterAreCriterionEqual, DEFAULT_ALBUM_FILTER_ENTRY, generateAlbumFilterOptions} from "./catalog-common-modifiers";

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

    const albumFilterOptions = generateAlbumFilterOptions(current.currentUser, allAlbums);
    const albumFilter = albumFilterOptions.find(option => albumFilterAreCriterionEqual(option.criterion, current.albumFilter.criterion)) ?? DEFAULT_ALBUM_FILTER_ENTRY
    const albums = allAlbums.filter(albumMatchCriterion(albumFilter.criterion))

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
