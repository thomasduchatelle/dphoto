import {createAction} from "src/libs/daction";
import {CatalogViewerState, editNameDialogNoError, isEditNameDialog} from "../language";

export const customFolderNameToggled = createAction<CatalogViewerState, boolean>(
    "CustomFolderNameToggled",
    (current: CatalogViewerState, isFolderNameEnabled: boolean) => {
        if (!isEditNameDialog(current.dialog)) {
            return current;
        }

        return {
            ...current,
            dialog: {
                ...current.dialog,
                isCustomFolderNameEnabled: isFolderNameEnabled,
                customFolderName: isFolderNameEnabled ? current.dialog.albumId.folderName : "",
                error: editNameDialogNoError,
            },
        };
    }
);
