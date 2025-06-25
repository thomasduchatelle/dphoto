import {createAction} from "src/libs/daction";
import {CatalogViewerState, isCreateDialog} from "../language";

export const createAlbumStarted = createAction<CatalogViewerState>(
    "CreateAlbumStarted",
    (current: CatalogViewerState) => {
        if (!isCreateDialog(current.dialog)) {
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

export type CreateAlbumStarted = ReturnType<typeof createAlbumStarted>;
