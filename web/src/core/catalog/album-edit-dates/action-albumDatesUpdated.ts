import {Album, CatalogViewerState, MediaWithinADay} from "../language";
import {refreshFilters} from "../common/utils";
import {createAction} from "@light-state";

interface AlbumDatesUpdatedPayload {
    albums: Album[];
    medias: MediaWithinADay[];
}

export const albumDatesUpdated = createAction<CatalogViewerState, AlbumDatesUpdatedPayload>(
    "AlbumDatesUpdated",
    (current: CatalogViewerState, {albums, medias}: AlbumDatesUpdatedPayload) => {
        const {albumFilterOptions, albumFilter, albums: filteredAlbums} = refreshFilters(
            current.currentUser,
            current.albumFilter,
            albums
        );

        return {
            ...current,
            allAlbums: albums,
            albumFilterOptions,
            albumFilter,
            albums: filteredAlbums,
            medias: medias,
            mediasLoadedFromAlbumId: current.editDatesDialog?.albumId, // Use albumId from dialog state
            albumsLoaded: true,
            mediasLoaded: true,
            editDatesDialog: undefined,
        };
    }
);

export type AlbumDatesUpdated = ReturnType<typeof albumDatesUpdated>;
