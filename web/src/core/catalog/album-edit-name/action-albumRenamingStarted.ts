import {createAction} from "src/libs/daction";
import {CatalogViewerState, editNameDialogNoError, isEditNameDialog} from "../language";

export const albumRenamingStarted = createAction<CatalogViewerState>(
    "AlbumRenamingStarted",
    (current: CatalogViewerState) => {
        if (!isEditNameDialog(current.dialog)) {
            return current;
        }

        return {
            ...current,
            dialog: {
                ...current.dialog,
                isLoading: true,
                error: editNameDialogNoError,
            },
        };
    }
);

export type AlbumRenamingStarted = ReturnType<typeof albumRenamingStarted>;
