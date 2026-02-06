import {AlbumId, CatalogViewerState} from "../language";
import {withOpenShareDialog} from "./sharing";
import {createAction} from "@/libs/daction";

export const sharingModalOpened = createAction<CatalogViewerState, AlbumId>(
    "sharingModalOpened",
    (current: CatalogViewerState, albumId: AlbumId) => {
        return withOpenShareDialog(current, albumId);
    }
);

export type SharingModalOpened = ReturnType<typeof sharingModalOpened>;
