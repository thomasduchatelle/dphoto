import {CatalogViewerState} from "../language";
import {DEFAULT_ALBUM_FILTER_ENTRY} from "../common/utils";
import {createAction} from "src/libs/daction";

export const noAlbumAvailable = createAction<CatalogViewerState>(
    'noAlbumAvailable',
    (current: CatalogViewerState): CatalogViewerState => {
        return {
            currentUser: current.currentUser,
            albumNotFound: true,
            allAlbums: [],
            albums: [],
            medias: [],
            albumsLoaded: true,
            mediasLoaded: true,
            albumFilterOptions: [DEFAULT_ALBUM_FILTER_ENTRY],
            albumFilter: DEFAULT_ALBUM_FILTER_ENTRY,
        };
    }
);

export type NoAlbumAvailable = ReturnType<typeof noAlbumAvailable>;
