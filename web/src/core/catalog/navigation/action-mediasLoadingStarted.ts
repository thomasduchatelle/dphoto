import {AlbumId, CatalogViewerState} from "../language";
import {createAction} from "src/light-state-lib";

export const mediasLoadingStarted = createAction<CatalogViewerState, AlbumId>(
    'mediasLoadingStarted',
    (current: CatalogViewerState, albumId: AlbumId): CatalogViewerState => {
        return {
            ...current,
            medias: [],
            error: undefined,
            loadingMediasFor: albumId,
            albumNotFound: false,
            mediasLoaded: false,
        };
    }
);

export type MediasLoadingStarted = ReturnType<typeof mediasLoadingStarted>;
