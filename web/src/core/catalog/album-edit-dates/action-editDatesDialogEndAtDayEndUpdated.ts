import {CatalogViewerState, isEditDatesDialog} from "../language";
import {createAction} from "src/libs/daction";

export const editDatesDialogEndAtDayEndUpdated = createAction<CatalogViewerState, boolean>(
    "EditDatesDialogEndAtDayEndUpdated",
    (current: CatalogViewerState, endAtDayEnd: boolean) => {
        if (!isEditDatesDialog(current.dialog)) {
            return current;
        }

        let updatedEndDate = current.dialog.endDate;
        if (!endAtDayEnd && current.dialog.endDate) {
            updatedEndDate = new Date(current.dialog.endDate);
            updatedEndDate.setUTCHours(23, 59, 0, 0);
        }
        
        return {
            ...current,
            dialog: {
                ...current.dialog,
                endAtDayEnd,
                endDate: updatedEndDate,
                error: undefined,
            },
        };
    }
);

export type EditDatesDialogEndAtDayEndUpdated = ReturnType<typeof editDatesDialogEndAtDayEndUpdated>;
