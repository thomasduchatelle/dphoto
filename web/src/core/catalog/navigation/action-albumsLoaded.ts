import {Album, albumIdEquals, albumMatchCriterion, CatalogViewerState, RedirectToAlbumIdAction} from "../language";
import {albumFilterAreCriterionEqual, DEFAULT_ALBUM_FILTER_ENTRY, refreshFilters} from "../common/utils";

export interface AlbumsLoaded extends RedirectToAlbumIdAction {
    type: 'albumsLoaded'
    albums: Album[]
}

export function albumsLoaded(props: Omit<AlbumsLoaded, "type">): AlbumsLoaded {
    return {
        ...props,
        type: 'albumsLoaded',
    };
}

export function reduceAlbumsLoaded(
    current: CatalogViewerState,
    action: AlbumsLoaded,
): CatalogViewerState {
    const {albumFilterOptions, albumFilter, albums} = refreshFilters(current.currentUser, current.albumFilter, action.albums);

    let staging: CatalogViewerState = {
        ...current,
        albumFilterOptions,
        albumFilter,
        allAlbums: action.albums,
        albums: albums,
        error: undefined,
        albumsLoaded: true,
    }

    if (!staging.albums.find(album => albumIdEquals(album.albumId, action.redirectTo))) {
        const albumFilter = current.albumFilterOptions.find(option => albumFilterAreCriterionEqual(option.criterion, DEFAULT_ALBUM_FILTER_ENTRY.criterion)) ?? DEFAULT_ALBUM_FILTER_ENTRY
        staging = {
            ...staging,
            albumFilter,
            albums: action.albums.filter(albumMatchCriterion(albumFilter.criterion)),
        }
    }

    return staging
}

export function albumsLoadedReducerRegistration(handlers: any) {
    handlers["albumsLoaded"] = reduceAlbumsLoaded as (
        state: CatalogViewerState,
        action: AlbumsLoaded
    ) => CatalogViewerState;
}
