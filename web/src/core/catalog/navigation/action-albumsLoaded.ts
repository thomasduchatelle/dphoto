import {Album, albumIdEquals, albumMatchCriterion, CatalogViewerState, RedirectToAlbumIdAction} from "../language";
import {albumFilterAreCriterionEqual, DEFAULT_ALBUM_FILTER_ENTRY, refreshFilters} from "../common/utils";
import {createAction} from "src/libs/daction";

interface AlbumsLoadedPayload extends RedirectToAlbumIdAction {
    albums: Album[]
}

export const albumsLoaded = createAction<CatalogViewerState, AlbumsLoadedPayload>(
    'albumsLoaded',
    (current: CatalogViewerState, {albums, redirectTo}: AlbumsLoadedPayload): CatalogViewerState => {
        const {albumFilterOptions, albumFilter, albums: filteredAlbums} = refreshFilters(current.currentUser, current.albumFilter, albums);

        let staging: CatalogViewerState = {
            ...current,
            albumFilterOptions,
            albumFilter,
            allAlbums: albums,
            albums: filteredAlbums,
            error: undefined,
            albumsLoaded: true,
        }

        if (redirectTo && !staging.albums.find(album => albumIdEquals(album.albumId, redirectTo))) {
            const albumFilter = current.albumFilterOptions.find(option => albumFilterAreCriterionEqual(option.criterion, DEFAULT_ALBUM_FILTER_ENTRY.criterion)) ?? DEFAULT_ALBUM_FILTER_ENTRY
            staging = {
                ...staging,
                albumFilter,
                albums: albums.filter(albumMatchCriterion(albumFilter.criterion)),
            }
        }

        return staging
    }
);

export type AlbumsLoaded = ReturnType<typeof albumsLoaded>;
