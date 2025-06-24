import {CatalogViewerState, isEditDatesDialog} from "../language";
import {createAction} from "src/libs/daction";

export const editDatesDialogStartDateUpdated = createAction<CatalogViewerState, Date | null>(
    "EditDatesDialogStartDateUpdated",
    (current: CatalogViewerState, startDate: Date | null) => {
        const dialog = current.dialog;
        if (!isEditDatesDialog(dialog)) {
            return current;
        }
        return {
            ...current,
            dialog: {
                ...dialog,
                startDate,
                error: undefined,
            },
        };
    }
);

export type EditDatesDialogStartDateUpdated = ReturnType<typeof editDatesDialogStartDateUpdated>;
