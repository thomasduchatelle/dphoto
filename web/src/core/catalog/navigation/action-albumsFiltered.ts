import {AlbumFilterCriterion, albumMatchCriterion, CatalogViewerState, RedirectToAlbumIdAction} from "../language";
import {albumFilterAreCriterionEqual, ALL_ALBUMS_FILTER_CRITERION, DEFAULT_ALBUM_FILTER_ENTRY} from "../common/utils";

export interface AlbumsFiltered extends RedirectToAlbumIdAction {
    type: 'albumsFiltered'
    criterion: AlbumFilterCriterion
}

export function albumsFiltered(props: Omit<AlbumsFiltered, "type">): AlbumsFiltered {
    return {...props, type: "albumsFiltered"};
}

export function reduceAlbumsFiltered(
    current: CatalogViewerState,
    action: AlbumsFiltered
): CatalogViewerState {
    const filteredAlbums = current.allAlbums.filter(albumMatchCriterion(action.criterion))

    const allAlbumFilter = current.albumFilterOptions.find(option => albumFilterAreCriterionEqual(option.criterion, ALL_ALBUMS_FILTER_CRITERION)) ?? DEFAULT_ALBUM_FILTER_ENTRY

    return {
        ...current,
        albums: filteredAlbums,
        albumFilter: current.albumFilterOptions.find(option => albumFilterAreCriterionEqual(option.criterion, action.criterion)) ?? allAlbumFilter,
    }
}

export function albumsFilteredReducerRegistration(handlers: any) {
    handlers["albumsFiltered"] = reduceAlbumsFiltered as (
        state: CatalogViewerState,
        action: AlbumsFiltered
    ) => CatalogViewerState;
}
