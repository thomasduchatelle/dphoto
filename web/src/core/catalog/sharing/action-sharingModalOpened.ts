import {AlbumId, CatalogViewerState} from "../language";
import {withOpenShareModal} from "./sharing";
import {createAction} from "src/libs/daction";

export const sharingModalOpened = createAction<CatalogViewerState, AlbumId>(
    "sharingModalOpened",
    (current: CatalogViewerState, albumId: AlbumId) => {
        return withOpenShareModal(current, albumId);
    }
);

export type SharingModalOpened = ReturnType<typeof sharingModalOpened>;
