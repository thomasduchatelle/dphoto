import {createAction} from "src/libs/daction";
import {CatalogViewerState, editNameDialogNoError, isEditNameDialog} from "../language";

export const customFolderNameChanged = createAction<CatalogViewerState, string>(
    "CustomFolderNameChanged",
    (current: CatalogViewerState, folderName: string) => {
        if (!isEditNameDialog(current.dialog)) {
            return current;
        }

        return {
            ...current,
            dialog: {
                ...current.dialog,
                customFolderName: folderName,
                error: !!folderName ? editNameDialogNoError : {folderNameError: "Folder name is mandatory"},
            },
        };
    }
);

