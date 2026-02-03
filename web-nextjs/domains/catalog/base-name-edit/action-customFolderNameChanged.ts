import {createAction} from "@/libs/daction";
import {CatalogViewerState, isNameEditBase} from "../language";

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
                nameError: {
                    ...current.dialog.nameError,
                    folderNameError: !!folderName ? undefined : "Folder name is mandatory",
                },
            },
        };
    }
);

