import {CatalogViewerState} from "../language";
import {createAction} from "@light-state";

export const editDatesDialogStartDateUpdated = createAction<CatalogViewerState, Date>(
    "EditDatesDialogStartDateUpdated",
    (current: CatalogViewerState, startDate: Date) => {
        if (!current.editDatesDialog) {
            return current;
        }
        return {
            ...current,
            editDatesDialog: {
                ...current.editDatesDialog,
                startDate,
            },
        };
    }
);

export type EditDatesDialogStartDateUpdated = ReturnType<typeof editDatesDialogStartDateUpdated>;
