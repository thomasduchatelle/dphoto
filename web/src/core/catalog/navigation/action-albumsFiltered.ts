import {AlbumFilterCriterion, albumMatchCriterion, CatalogViewerState, RedirectToAlbumIdAction} from "../language";
import {albumFilterAreCriterionEqual, ALL_ALBUMS_FILTER_CRITERION, DEFAULT_ALBUM_FILTER_ENTRY} from "../common/utils";
import {createAction} from "src/light-state-lib";

interface AlbumsFilteredPayload extends RedirectToAlbumIdAction {
    criterion: AlbumFilterCriterion
}

export const albumsFiltered = createAction<CatalogViewerState, AlbumsFilteredPayload>(
    'albumsFiltered',
    (current: CatalogViewerState, {criterion, redirectTo}: AlbumsFilteredPayload): CatalogViewerState => {
        const filteredAlbums = current.allAlbums.filter(albumMatchCriterion(criterion))

        const allAlbumFilter = current.albumFilterOptions.find(option => albumFilterAreCriterionEqual(option.criterion, ALL_ALBUMS_FILTER_CRITERION)) ?? DEFAULT_ALBUM_FILTER_ENTRY

        return {
            ...current,
            albums: filteredAlbums,
            albumFilter: current.albumFilterOptions.find(option => albumFilterAreCriterionEqual(option.criterion, criterion)) ?? allAlbumFilter,
            redirectTo: redirectTo,
        }
    }
);

export type AlbumsFiltered = ReturnType<typeof albumsFiltered>;
