import {Album, AlbumId, CatalogViewerState, Media, RedirectToAlbumIdPayload} from "../language";

import {refreshFilters} from "../common/utils";
import {createAction} from "@/libs/daction";

import {groupByDay} from "./group-by-day";

interface AlbumsAndMediasLoadedPayload extends RedirectToAlbumIdPayload {
    albums: Album[]
    medias: Media[]
    mediasFromAlbumId?: AlbumId
}

export const albumsAndMediasLoaded = createAction<CatalogViewerState, AlbumsAndMediasLoadedPayload>(
    'albumsAndMediasLoaded',
    (current: CatalogViewerState, {albums: allAlbums, medias, mediasFromAlbumId}: AlbumsAndMediasLoadedPayload): CatalogViewerState => {
        const {
            albumFilterOptions,
            albumFilter,
            albums: filteredAlbums
        } = refreshFilters(current.currentUser, current.albumFilter, allAlbums, mediasFromAlbumId);

        return {
            currentUser: current.currentUser,
            albumNotFound: false,
            allAlbums,
            albums: filteredAlbums,
            mediasLoadedFromAlbumId: mediasFromAlbumId,
            medias: groupByDay(medias),
            albumsLoaded: true,
            mediasLoaded: true,
            albumFilterOptions,
            albumFilter: albumFilter,
        }
    }
);

export type AlbumsAndMediasLoaded = ReturnType<typeof albumsAndMediasLoaded>;
