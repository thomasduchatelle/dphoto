import {CatalogViewerState} from "../language";
import {DEFAULT_ALBUM_FILTER_ENTRY} from "../common/utils";
import {createAction} from "@/libs/daction";

export const noAlbumAvailable = createAction<CatalogViewerState, Error | undefined>(
    'noAlbumAvailable',
    (current: CatalogViewerState, error): CatalogViewerState => {
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
            error,
        };
    }
);

export type NoAlbumAvailable = ReturnType<typeof noAlbumAvailable>;
