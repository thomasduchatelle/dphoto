import {createAction} from "src/libs/daction";
import {CatalogViewerState, editNameDialogNoError, isNameEditBase} from "../language";

export const customFolderNameChanged = createAction<CatalogViewerState, string>(
    "CustomFolderNameChanged",
    (current: CatalogViewerState, folderName: string) => {
        if (!isNameEditBase(current.dialog)) {
            return current;
        }

        return {
            ...current,
            dialog: {
                ...current.dialog,
                customFolderName: folderName,
                nameError: !!folderName ? editNameDialogNoError : {folderNameError: "Folder name is mandatory"},
            },
        };
    }
);

