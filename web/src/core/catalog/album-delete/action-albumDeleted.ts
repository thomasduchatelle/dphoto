import {Album, AlbumId, albumIdEquals, CatalogViewerState, RedirectToAlbumIdAction} from "../language";
import {albumFilterAreCriterionEqual, ALL_ALBUMS_FILTER_CRITERION, DEFAULT_ALBUM_FILTER_ENTRY, refreshFilters} from "../navigation";

export interface AlbumDeleted extends RedirectToAlbumIdAction {
    type: "albumDeleted";
    albums: Album[];
    redirectTo?: AlbumId;
}

export function albumDeleted(props: Omit<AlbumDeleted, "type">): AlbumDeleted {
    return {
        ...props,
        type: "albumDeleted",
    };
}

export function reduceAlbumDeleted(
    current: CatalogViewerState,
    action: AlbumDeleted,
): CatalogViewerState {
    let {albumFilterOptions, albumFilter, albums} = refreshFilters(current.currentUser, current.albumFilter, action.albums);

    if (
        action.redirectTo &&
        !albums.some(album => albumIdEquals(album.albumId, action.redirectTo))
    ) {
        albumFilter =
            albumFilterOptions.find(option =>
                albumFilterAreCriterionEqual(option.criterion, ALL_ALBUMS_FILTER_CRITERION)
            ) ?? DEFAULT_ALBUM_FILTER_ENTRY;
        albums = action.albums
    }

    return {
        ...current,
        albumFilterOptions,
        albumFilter,
        allAlbums: action.albums,
        albums: albums,
        error: undefined,
        albumsLoaded: true,
        deleteDialog: undefined,
    };
}

export function albumDeletedReducerRegistration(handlers: any) {
    handlers["albumDeleted"] = reduceAlbumDeleted as (
        state: CatalogViewerState,
        action: AlbumDeleted
    ) => CatalogViewerState;
}
