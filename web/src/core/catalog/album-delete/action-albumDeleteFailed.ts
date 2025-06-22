import {CatalogViewerState} from "../language";
import {createAction} from "src/libs/daction";

export const albumDeleteFailed = createAction<CatalogViewerState, string>(
    "albumDeleteFailed",
    (current: CatalogViewerState, error: string) => {
        if (!error || error.trim() === "") {
            throw new Error("AlbumDeleteFailed requires a non-empty error message");
        }

        if (!current.deleteDialog) {
            return current;
        }

        return {
            ...current,
            deleteDialog: {
                ...current.deleteDialog,
                error: error,
                isLoading: false,
            },
        };
    }
);

export type AlbumDeleteFailed = ReturnType<typeof albumDeleteFailed>;
