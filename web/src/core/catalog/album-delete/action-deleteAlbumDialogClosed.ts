import {CatalogViewerState, isDeleteDialog} from "../language";
import {createAction} from "src/libs/daction";

export const deleteAlbumDialogClosed = createAction<CatalogViewerState>(
    "deleteAlbumDialogClosed",
    (state: CatalogViewerState) => {
        if (!isDeleteDialog(state.dialog)) return state;
        return {
            ...state,
            dialog: undefined,
        };
    }
);

export type DeleteAlbumDialogClosed = ReturnType<typeof deleteAlbumDialogClosed>;
