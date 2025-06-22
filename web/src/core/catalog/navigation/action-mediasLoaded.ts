import {AlbumId, albumIdEquals, CatalogViewerState, MediaWithinADay} from "../language";
import {createAction} from "src/light-state-lib";

interface MediasLoadedPayload {
    albumId: AlbumId
    medias: MediaWithinADay[]
}

export const mediasLoaded = createAction<CatalogViewerState, MediasLoadedPayload>(
    "mediasLoaded",
    (current: CatalogViewerState, {albumId, medias}: MediasLoadedPayload): CatalogViewerState => {
        if (current.loadingMediasFor && !albumIdEquals(current.loadingMediasFor, albumId)) {
            // concurrency management - ignore if not the last album requested
            return current
        }

        return {
            ...current,
            loadingMediasFor: undefined,
            mediasLoadedFromAlbumId: albumId,
            medias: medias,
            error: undefined,
            mediasLoaded: true,
            albumNotFound: false,
        }
    }
);

export type MediasLoaded = ReturnType<typeof mediasLoaded>;
