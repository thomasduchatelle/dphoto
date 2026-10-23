import {createAction} from "@/libs/daction";
import {CatalogViewerState, editNameDialogNoError, isCreateDialog, NameError} from "../language";

interface CreateAlbumFailedPayload {
    code?: string;
    message: string;
}

export const createAlbumFailed = createAction<CatalogViewerState, CreateAlbumFailedPayload>(
    "CreateAlbumFailed",
    (current: CatalogViewerState, {code, message}: CreateAlbumFailedPayload) => {
        if (!isCreateDialog(current.dialog)) {
            return current;
        }

        let nameError: NameError = editNameDialogNoError;
        let error: string | undefined = undefined;

        if (code === "AlbumFolderNameAlreadyTakenErr") {
            if (current.dialog.isCustomFolderNameEnabled) {
                nameError = {folderNameError: message};
            } else {
                nameError = {nameError: message};
            }
        } else {
            error = message;
        }

        return {
            ...current,
            dialog: {
                ...current.dialog,
                isLoading: false,
                nameError,
                error,
            },
        };
    }
);

export type CreateAlbumFailed = ReturnType<typeof createAlbumFailed>;
