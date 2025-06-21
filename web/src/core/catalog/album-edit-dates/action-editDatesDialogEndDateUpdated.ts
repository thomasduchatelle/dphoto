import {CatalogViewerState} from "../language";
import {createActionWithPayload} from "../common/action-factory";

export const editDatesDialogEndDateUpdated = createActionWithPayload<CatalogViewerState, Date>(
    "EditDatesDialogEndDateUpdated",
    (current: CatalogViewerState, action) => {
        if (!current.editDatesDialog) {
            return current;
        }
        return {
            ...current,
            editDatesDialog: {
                ...current.editDatesDialog,
                endDate: action.payload!,
            },
        };
    }
);

export type EditDatesDialogEndDateUpdated = ReturnType<typeof editDatesDialogEndDateUpdated>;
