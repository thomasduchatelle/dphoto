import {createAction} from "@/libs/daction";
import {CatalogViewerState, editNameDialogNoError, isCreateDialog, NameError} from "../language";

interface CreateAlbumFailedPayload {
    nameError?: NameError;
    error?: string;
}

export const createAlbumFailed = createAction<CatalogViewerState, CreateAlbumFailedPayload>(
    "CreateAlbumFailed",
    (current: CatalogViewerState, {nameError, error}: CreateAlbumFailedPayload) => {
        if (!isCreateDialog(current.dialog)) {
            return current;
        }

        return {
            ...current,
            dialog: {
                ...current.dialog,
                isLoading: false,
                nameError: nameError ?? editNameDialogNoError,
                error,
            },
        };
    }
);

export type CreateAlbumFailed = ReturnType<typeof createAlbumFailed>;
