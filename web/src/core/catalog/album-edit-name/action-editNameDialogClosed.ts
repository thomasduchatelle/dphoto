import {createAction} from "src/libs/daction";
import {CatalogViewerState, isEditNameDialog} from "../language";

export const editNameDialogClosed = createAction<CatalogViewerState>(
    "EditNameDialogClosed",
    (current: CatalogViewerState) => {
        if (!isEditNameDialog(current.dialog)) {
            return current;
        }

        return {
            ...current,
            dialog: undefined,
        };
    }
);

export type EditNameDialogClosed = ReturnType<typeof editNameDialogClosed>;
