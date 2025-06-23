import {CatalogViewerState} from "../language";
import {createAction} from "src/libs/daction";

export const editDatesDialogStartDateUpdated = createAction<CatalogViewerState, Date | null>(
    "EditDatesDialogStartDateUpdated",
    (current: CatalogViewerState, startDate: Date | null) => {
        if (!current.editDatesDialog || !startDate) {
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
