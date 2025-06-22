import {Album, CatalogViewerState, MediaWithinADay, RedirectToAlbumIdAction} from "../language";

import {refreshFilters} from "../common/utils";
import {createAction} from "src/light-state-lib";

interface AlbumsAndMediasLoadedPayload extends RedirectToAlbumIdAction {
    albums: Album[]
    medias: MediaWithinADay[]
    selectedAlbum?: Album
}

export const albumsAndMediasLoaded = createAction<CatalogViewerState, AlbumsAndMediasLoadedPayload>(
    'albumsAndMediasLoaded',
    (current: CatalogViewerState, {albums, medias, selectedAlbum, redirectTo}: AlbumsAndMediasLoadedPayload): CatalogViewerState => {
        const {albumFilterOptions, albumFilter, albums: filteredAlbums} = refreshFilters(current.currentUser, current.albumFilter, albums);

        return {
            currentUser: current.currentUser,
            albumNotFound: false,
            allAlbums: albums,
            albums: filteredAlbums,
            mediasLoadedFromAlbumId: selectedAlbum?.albumId,
            medias: medias,
            albumsLoaded: true,
            mediasLoaded: true,
            albumFilterOptions,
            albumFilter: albumFilter,
        }
    }
);

export type AlbumsAndMediasLoaded = ReturnType<typeof albumsAndMediasLoaded>;
