import {createAction} from "src/libs/daction";
import {CatalogViewerState, isCreateDialog} from "../language";

export const createDialogNameChanged = createAction<CatalogViewerState, string>(
    "CreateDialogNameChanged",
    (current: CatalogViewerState, name: string) => {
        if (!isCreateDialog(current.dialog)) {
            return current;
        }

        return {
            ...current,
            dialog: {
                ...current.dialog,
                name,
                error: current.dialog.error === "AlbumFolderNameAlreadyTakenErr" ? undefined : current.dialog.error,
            },
        };
    }
);

export type CreateDialogNameChanged = ReturnType<typeof createDialogNameChanged>;
