import {createAction} from "@/libs/daction";
import {CatalogViewerState, isCreateDialog} from "../language";

export const createAlbumFailed = createAction<CatalogViewerState, string>(
    "CreateAlbumFailed",
    (current: CatalogViewerState, error: string) => {
        if (!isCreateDialog(current.dialog)) {
            return current;
        }

        return {
            ...current,
            dialog: {
                ...current.dialog,
                isLoading: false,
                error,
            },
        };
    }
);

export type CreateAlbumFailed = ReturnType<typeof createAlbumFailed>;
