import {CatalogViewerState} from "../language";
import {createAction} from "src/libs/daction";

export const editDatesDialogStartAtDayStartUpdated = createAction<CatalogViewerState, boolean>(
    "EditDatesDialogStartAtDayStartUpdated",
    (current: CatalogViewerState, startAtDayStart: boolean) => {
        if (!current.editDatesDialog) {
            return current;
        }
        
        let updatedStartDate = current.editDatesDialog.startDate;
        if (startAtDayStart) {
            updatedStartDate = new Date(current.editDatesDialog.startDate);
            updatedStartDate.setHours(0, 0, 0, 0);
        }
        
        return {
            ...current,
            editDatesDialog: {
                ...current.editDatesDialog,
                startAtDayStart,
                startDate: updatedStartDate,
                error: undefined,
            },
        };
    }
);

export type EditDatesDialogStartAtDayStartUpdated = ReturnType<typeof editDatesDialogStartAtDayStartUpdated>;
