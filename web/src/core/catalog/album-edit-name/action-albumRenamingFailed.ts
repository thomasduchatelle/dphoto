import {createAction} from "src/libs/daction";
import {CatalogViewerState, isEditNameDialog} from "../language";

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

        let nameError = {};

        if (code === "AlbumFolderNameAlreadyTakenErr") {
            if (current.dialog.isCustomFolderNameEnabled) {
                nameError = {folderNameError: message};
            } else {
                nameError = {nameError: message};
            }
        } else if (code === "AlbumNameMandatoryErr") {
            nameError = {nameError: message};
        } else {
            nameError = {technicalError: message};
        }

        return {
            ...current,
            dialog: {
                ...current.dialog,
                isLoading: false,
                nameError,
            },
        };
    }
);

export type AlbumRenamingFailed = ReturnType<typeof albumRenamingFailed>;
