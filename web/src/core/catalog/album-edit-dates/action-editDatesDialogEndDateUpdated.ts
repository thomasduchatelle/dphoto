import {CatalogViewerState, isEditDatesDialog} from "../language";
import {createAction} from "src/libs/daction";

export const editDatesDialogEndDateUpdated = createAction<CatalogViewerState, Date | null>(
    "EditDatesDialogEndDateUpdated",
    (current: CatalogViewerState, endDate: Date | null) => {
        const dialog = current.dialog;
        if (!isEditDatesDialog(dialog)) {
            return current;
        }
        return {
            ...current,
            dialog: {
                ...dialog,
                endDate,
                error: undefined,
            },
        };
    }
);

export type EditDatesDialogEndDateUpdated = ReturnType<typeof editDatesDialogEndDateUpdated>;
