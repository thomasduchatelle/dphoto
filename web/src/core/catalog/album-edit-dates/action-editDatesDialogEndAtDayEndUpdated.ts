import {CatalogViewerState} from "../language";
import {createAction} from "src/libs/daction";

export const editDatesDialogEndAtDayEndUpdated = createAction<CatalogViewerState, boolean>(
    "EditDatesDialogEndAtDayEndUpdated",
    (current: CatalogViewerState, endAtDayEnd: boolean) => {
        if (!current.editDatesDialog) {
            return current;
        }

        let updatedEndDate = current.editDatesDialog.endDate;
        if (!endAtDayEnd) {
            updatedEndDate = new Date(current.editDatesDialog.endDate);
            updatedEndDate.setUTCHours(23, 59, 0, 0);
        }
        
        return {
            ...current,
            editDatesDialog: {
                ...current.editDatesDialog,
                endAtDayEnd,
                endDate: updatedEndDate,
                error: undefined,
            },
        };
    }
);

export type EditDatesDialogEndAtDayEndUpdated = ReturnType<typeof editDatesDialogEndAtDayEndUpdated>;
