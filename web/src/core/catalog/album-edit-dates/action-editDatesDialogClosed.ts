import {CatalogViewerState, isEditDatesDialog} from "../language";
import {createAction} from "src/libs/daction";

export const editDatesDialogClosed = createAction<CatalogViewerState>(
    "EditDatesDialogClosed",
    (current: CatalogViewerState) => {
        if (!isEditDatesDialog(current.dialog)) {
            return current;
        }

        return {
            ...current,
            dialog: undefined,
        };
    }
);

export type EditDatesDialogClosed = ReturnType<typeof editDatesDialogClosed>;
