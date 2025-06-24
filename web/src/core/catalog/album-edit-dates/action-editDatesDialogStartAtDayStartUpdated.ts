import {CatalogViewerState, isEditDatesDialog} from "../language";
import {createAction} from "src/libs/daction";

export const editDatesDialogStartAtDayStartUpdated = createAction<CatalogViewerState, boolean>(
    "EditDatesDialogStartAtDayStartUpdated",
    (current: CatalogViewerState, startAtDayStart: boolean) => {
        const dialog = current.dialog;
        if (!isEditDatesDialog(dialog)) {
            return current;
        }
        
        let updatedStartDate = dialog.startDate;
        if (startAtDayStart && dialog.startDate) {
            updatedStartDate = new Date(dialog.startDate);
            updatedStartDate.setHours(0, 0, 0, 0);
        }
        
        return {
            ...current,
            dialog: {
                ...dialog,
                startAtDayStart,
                startDate: updatedStartDate,
                error: undefined,
            },
        };
    }
);

export type EditDatesDialogStartAtDayStartUpdated = ReturnType<typeof editDatesDialogStartAtDayStartUpdated>;
