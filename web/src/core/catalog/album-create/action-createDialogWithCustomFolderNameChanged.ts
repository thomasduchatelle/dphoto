import {createAction} from "src/libs/daction";
import {CatalogViewerState, isCreateDialog} from "../language";

export const createDialogWithCustomFolderNameChanged = createAction<CatalogViewerState, boolean>(
    "CreateDialogWithCustomFolderNameChanged",
    (current: CatalogViewerState, withCustomFolderName: boolean) => {
        if (!isCreateDialog(current.dialog)) {
            return current;
        }

        return {
            ...current,
            dialog: {
                ...current.dialog,
                withCustomFolderName,
                forceFolderName: "", // Clear forceFolderName when custom folder name status changes
            },
        };
    }
);

export type CreateDialogWithCustomFolderNameChanged = ReturnType<typeof createDialogWithCustomFolderNameChanged>;
