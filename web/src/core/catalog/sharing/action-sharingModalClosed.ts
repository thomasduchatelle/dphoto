import {CatalogViewerState} from "../language";
import {createAction} from "src/libs/daction";

export const sharingModalClosed = createAction<CatalogViewerState>(
    "sharingModalClosed",
    ({shareModal, ...rest}: CatalogViewerState) => {
        return rest;
    }
);

export type SharingModalClosed = ReturnType<typeof sharingModalClosed>;
