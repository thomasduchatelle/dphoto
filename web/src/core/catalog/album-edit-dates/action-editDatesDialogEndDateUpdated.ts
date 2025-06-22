import {CatalogViewerState} from "../language";
import {createAction} from "@light-state";

export const editDatesDialogEndDateUpdated = createAction<CatalogViewerState, Date>(
    "EditDatesDialogEndDateUpdated",
    (current: CatalogViewerState, endDate: Date) => {
        if (!current.editDatesDialog) {
            return current;
        }
        return {
            ...current,
            editDatesDialog: {
                ...current.editDatesDialog,
                endDate,
            },
        };
    }
);

export type EditDatesDialogEndDateUpdated = ReturnType<typeof editDatesDialogEndDateUpdated>;
