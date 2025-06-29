import {createAction} from "src/libs/daction";
import {CatalogViewerState, editNameDialogNoError, isNameEditBase} from "../language";

export const customFolderNameToggled = createAction<CatalogViewerState, boolean>(
    "CustomFolderNameToggled",
    (current: CatalogViewerState, isFolderNameEnabled: boolean) => {
        if (!isNameEditBase(current.dialog)) {
            return current;
        }

        return {
            ...current,
            dialog: {
                ...current.dialog,
                isCustomFolderNameEnabled: isFolderNameEnabled,
                customFolderName: isFolderNameEnabled ? current.dialog.albumId.folderName : "",
                nameError: editNameDialogNoError,
            },
        };
    }
);
