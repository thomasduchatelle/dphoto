import {CatalogViewerState} from "../language";
import {createAction} from "../common/action-factory";

export const editDatesDialogClosed = createAction<CatalogViewerState>(
    "EditDatesDialogClosed",
    (current: CatalogViewerState) => {
        if (!current.editDatesDialog) {
            return current;
        }

        return {
            ...current,
            editDatesDialog: undefined,
        };
    }
);

export type EditDatesDialogClosed = ReturnType<typeof editDatesDialogClosed>;
