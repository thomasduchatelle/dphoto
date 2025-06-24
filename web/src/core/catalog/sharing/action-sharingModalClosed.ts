import {CatalogViewerState, isShareDialog} from "../language";
import {createAction} from "src/libs/daction";

export const sharingModalClosed = createAction<CatalogViewerState>(
    "sharingModalClosed",
    (current: CatalogViewerState) => {
        if (!isShareDialog(current.dialog)) {
            return current;
        }
        return {
            ...current,
            dialog: undefined,
        };
    }
);

export type SharingModalClosed = ReturnType<typeof sharingModalClosed>;
