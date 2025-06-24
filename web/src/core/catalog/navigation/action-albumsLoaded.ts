import {Album, CatalogViewerState, RedirectToAlbumIdPayload} from "../language";
import {refreshFilters} from "../common/utils";
import {createAction} from "src/libs/daction";

interface AlbumsLoadedPayload extends RedirectToAlbumIdPayload {
    albums: Album[]
}

export const albumsLoaded = createAction<CatalogViewerState, AlbumsLoadedPayload>(
    'albumsLoaded',
    (current: CatalogViewerState, {albums, redirectTo}: AlbumsLoadedPayload): CatalogViewerState => {
        const {albumFilterOptions, albumFilter, albums: filteredAlbums} = refreshFilters(current.currentUser, current.albumFilter, albums, redirectTo);

        return {
            ...current,
            albumFilterOptions,
            albumFilter,
            allAlbums: albums,
            albums: filteredAlbums,
            error: undefined,
            albumsLoaded: true,
        }
    }
);

export type AlbumsLoaded = ReturnType<typeof albumsLoaded>;
