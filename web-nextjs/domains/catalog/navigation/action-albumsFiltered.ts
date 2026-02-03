import {AlbumFilterCriterion, albumMatchCriterion, CatalogViewerState, RedirectToAlbumIdPayload} from "../language";
import {albumFilterAreCriterionEqual, ALL_ALBUMS_FILTER_CRITERION, DEFAULT_ALBUM_FILTER_ENTRY} from "../common/utils";
import {createAction} from "@/libs/daction";

interface AlbumsFilteredPayload extends RedirectToAlbumIdPayload {
    criterion: AlbumFilterCriterion
}

export const albumsFiltered = createAction<CatalogViewerState, AlbumsFilteredPayload>(
    'albumsFiltered',
    (current: CatalogViewerState, {criterion}: AlbumsFilteredPayload): CatalogViewerState => {
        let albumFilter = current.albumFilterOptions.find(option => albumFilterAreCriterionEqual(option.criterion, criterion));
        if (!albumFilter) {
            albumFilter = current.albumFilterOptions.find(option => albumFilterAreCriterionEqual(option.criterion, ALL_ALBUMS_FILTER_CRITERION));
        }
        if (!albumFilter) {
            albumFilter = DEFAULT_ALBUM_FILTER_ENTRY;
        }

        const filteredAlbums = current.allAlbums.filter(albumMatchCriterion(criterion))

        return {
            ...current,
            albums: filteredAlbums,
            albumFilter,
        }
    }
);

export type AlbumsFiltered = ReturnType<typeof albumsFiltered>;
