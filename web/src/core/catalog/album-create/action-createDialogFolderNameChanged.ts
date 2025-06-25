import {createAction} from "src/libs/daction";
import {CatalogViewerState, isCreateDialog} from "../language";

export const createDialogFolderNameChanged = createAction<CatalogViewerState, string>(
    "CreateDialogFolderNameChanged",
    (current: CatalogViewerState, forceFolderName: string) => {
        if (!isCreateDialog(current.dialog)) {
            return current;
        }

        return {
            ...current,
            dialog: {
                ...current.dialog,
                forceFolderName,
                error: current.dialog.error === "AlbumFolderNameAlreadyTakenErr" ? undefined : current.dialog.error,
            },
        };
    }
);

export type CreateDialogFolderNameChanged = ReturnType<typeof createDialogFolderNameChanged>;
