import {createAction} from "src/libs/daction";
import {AlbumId, albumIdEquals, albumMatchCriterion, CatalogViewerState} from "../language";

interface AlbumRenamedPayload {
    previousAlbumId: AlbumId;
    newAlbumId: AlbumId;
    newName: string;
}

export const albumRenamed = createAction<CatalogViewerState, AlbumRenamedPayload>(
    "AlbumRenamed",
    (current: CatalogViewerState, {previousAlbumId, newAlbumId, newName}: AlbumRenamedPayload) => {
        const allAlbums = current.allAlbums.map(album => {
            if (albumIdEquals(album.albumId, previousAlbumId)) {
                return {
                    ...album,
                    albumId: newAlbumId,
                    name: newName,
                };
            }
            return album;
        });

        return {
            ...current,
            allAlbums,
            albums: allAlbums.filter(albumMatchCriterion(current.albumFilter.criterion)),
            mediasLoadedFromAlbumId: albumIdEquals(current.mediasLoadedFromAlbumId, previousAlbumId) ? newAlbumId : current.mediasLoadedFromAlbumId,
            loadingMediasFor: albumIdEquals(current.loadingMediasFor, previousAlbumId) ? newAlbumId : current.loadingMediasFor,
            dialog: undefined,
        };
    }
);

export type AlbumRenamed = ReturnType<typeof albumRenamed>;
