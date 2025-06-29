import {createAction} from "src/libs/daction";
import {CatalogViewerState, editNameDialogNoError, isEditNameDialog, NameError} from "../language";

interface AlbumRenamingFailedPayload {
    code?: string;
    message: string;
}

export const albumRenamingFailed = createAction<CatalogViewerState, AlbumRenamingFailedPayload>(
    "AlbumRenamingFailed",
    (current: CatalogViewerState, {code, message}: AlbumRenamingFailedPayload) => {
        if (!isEditNameDialog(current.dialog)) {
            return current;
        }

        let nameError: NameError = editNameDialogNoError
        let technicalError: string | undefined = undefined;

        if (code === "AlbumFolderNameAlreadyTakenErr") {
            if (current.dialog.isCustomFolderNameEnabled) {
                nameError = {folderNameError: message};
            } else {
                nameError = {nameError: message};
            }
        } else if (code === "AlbumNameMandatoryErr") {
            nameError = {nameError: message};
        } else {
            technicalError = message
        }

        return {
            ...current,
            dialog: {
                ...current.dialog,
                isLoading: false,
                nameError,
                technicalError,
            },
        };
    }
);

export type AlbumRenamingFailed = ReturnType<typeof albumRenamingFailed>;
