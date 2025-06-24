import {Album, AlbumId, CatalogViewerState} from "../language";
import {refreshFilters} from "../common/utils";
import {createAction} from "src/libs/daction";

interface MediaLoadFailedPayload {
    albums?: Album[]
    displayedAlbumId?: AlbumId
    error: Error
}

export const mediaLoadFailed = createAction<CatalogViewerState, MediaLoadFailedPayload>(
    'mediaLoadFailed',
    (current: CatalogViewerState, {albums, displayedAlbumId, error}: MediaLoadFailedPayload): CatalogViewerState => {
        const allAlbums = albums ?? current.allAlbums;

        const {albumFilterOptions, albumFilter, albums: filteredAlbums} = refreshFilters(current.currentUser, current.albumFilter, allAlbums, displayedAlbumId);

        return {
            currentUser: current.currentUser,
            allAlbums,
            albumFilterOptions,
            albumFilter,
            mediasLoadedFromAlbumId: displayedAlbumId,
            albums: filteredAlbums,
            albumNotFound: false,
            medias: [],
            error: error,
            albumsLoaded: true,
            mediasLoaded: true,
        };
    }
);

export type MediaLoadFailed = ReturnType<typeof mediaLoadFailed>;
