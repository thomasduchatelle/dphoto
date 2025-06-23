import {CatalogViewerState} from "../language";
import {createAction} from "src/libs/daction";

export const deleteAlbumDialogClosed = createAction<CatalogViewerState>(
    "deleteAlbumDialogClosed",
    (state: CatalogViewerState) => {
        if (!state.deleteDialog) return state;
        return {
            ...state,
            deleteDialog: undefined,
        };
    }
);

export type DeleteAlbumDialogClosed = ReturnType<typeof deleteAlbumDialogClosed>;
