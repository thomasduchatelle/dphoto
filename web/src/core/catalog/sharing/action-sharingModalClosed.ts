import {CatalogViewerState} from "../language";
import {createAction} from "src/light-state-lib";

export const sharingModalClosed = createAction<CatalogViewerState>(
    "sharingModalClosed",
    ({shareModal, ...rest}: CatalogViewerState) => {
        return rest;
    }
);

export type SharingModalClosed = ReturnType<typeof sharingModalClosed>;
