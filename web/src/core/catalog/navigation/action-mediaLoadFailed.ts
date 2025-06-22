import {Album, CatalogViewerState} from "../language";
import {refreshFilters} from "../common/utils";
import {createAction} from "src/light-state-lib";

interface MediaLoadFailedPayload {
    albums?: Album[]
    selectedAlbum?: Album
    error: Error
}

export const mediaLoadFailed = createAction<CatalogViewerState, MediaLoadFailedPayload>(
    'mediaLoadFailed',
    (current: CatalogViewerState, {albums, selectedAlbum, error}: MediaLoadFailedPayload): CatalogViewerState => {
        const allAlbums = albums ?? current.allAlbums;

        const {albumFilterOptions, albumFilter, albums: filteredAlbums} = refreshFilters(current.currentUser, current.albumFilter, allAlbums);

        return {
            currentUser: current.currentUser,
            allAlbums,
            albumFilterOptions,
            albumFilter,
            mediasLoadedFromAlbumId: selectedAlbum?.albumId,
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
