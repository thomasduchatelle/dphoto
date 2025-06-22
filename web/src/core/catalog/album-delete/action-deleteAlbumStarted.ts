import {CatalogViewerState} from "../language";
import {createAction} from "@light-state";

export const deleteAlbumStarted = createAction<CatalogViewerState>(
    "deleteAlbumStarted",
    (current: CatalogViewerState) => {
        if (!current.deleteDialog) {
            return current;
        }
        return {
            ...current,
            deleteDialog: {
                ...current.deleteDialog,
                isLoading: true,
                error: undefined,
            },
        };
    }
);

export type DeleteAlbumStarted = ReturnType<typeof deleteAlbumStarted>;
