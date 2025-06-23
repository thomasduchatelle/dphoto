import {CatalogViewerState} from "../language";
import {createAction} from "src/libs/daction";

export const editDatesDialogEndDateUpdated = createAction<CatalogViewerState, Date | null>(
    "EditDatesDialogEndDateUpdated",
    (current: CatalogViewerState, endDate: Date | null) => {
        if (!current.editDatesDialog || !endDate) {
            return current;
        }
        return {
            ...current,
            editDatesDialog: {
                ...current.editDatesDialog,
                endDate,
                error: undefined,
            },
        };
    }
);

export type EditDatesDialogEndDateUpdated = ReturnType<typeof editDatesDialogEndDateUpdated>;
