import {createAction} from "src/libs/daction";
import {CatalogViewerState, editNameDialogNoError, isEditNameDialog} from "../language";

export const folderNameEnabledChanged = createAction<CatalogViewerState, boolean>(
    "FolderNameEnabledChanged",
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

export type FolderNameEnabledChanged = ReturnType<typeof folderNameEnabledChanged>;
