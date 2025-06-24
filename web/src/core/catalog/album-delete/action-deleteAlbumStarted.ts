import {CatalogViewerState, isDeleteDialog} from "../language";
import {createAction} from "src/libs/daction";

export const deleteAlbumStarted = createAction<CatalogViewerState>(
    "deleteAlbumStarted",
    (current: CatalogViewerState) => {
        if (!isDeleteDialog(current.dialog)) {
            return current;
        }
        return {
            ...current,
            dialog: {
                ...current.dialog,
                isLoading: true,
                error: undefined,
            },
        };
    }
);

export type DeleteAlbumStarted = ReturnType<typeof deleteAlbumStarted>;
