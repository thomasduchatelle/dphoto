import {CatalogViewerState} from "../language";
import {createActionWithPayload} from "../common/action-factory";

export const editDatesDialogStartDateUpdated = createActionWithPayload(
    "EditDatesDialogStartDateUpdated",
    (current: CatalogViewerState, action) => {
        if (!current.editDatesDialog) {
            return current;
        }
        return {
            ...current,
            editDatesDialog: {
                ...current.editDatesDialog,
                startDate: action.payload!,
            },
        };
    }
);

export type EditDatesDialogStartDateUpdated = ReturnType<typeof editDatesDialogStartDateUpdated>;
