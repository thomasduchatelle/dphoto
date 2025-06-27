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

        let error = {};

        if (code === "AlbumFolderNameAlreadyTakenErr") {
            if (current.dialog.isCustomFolderNameEnabled) {
                error = {folderNameError: message};
            } else {
                error = {nameError: message};
            }
        } else if (code === "AlbumNameMandatoryErr") {
            error = {nameError: message};
        } else {
            error = {technicalError: message};
        }

        return {
            ...current,
            dialog: {
                ...current.dialog,
                isLoading: false,
                error,
            },
        };
    }
);

export type AlbumRenamingFailed = ReturnType<typeof albumRenamingFailed>;
