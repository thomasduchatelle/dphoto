import {CatalogViewerState} from "../language";
import {createActionWithoutPayload} from "../common/action-factory";

export const editDatesDialogClosed = createActionWithoutPayload(
    "EditDatesDialogClosed",
    (current: CatalogViewerState) => {
        if (!current.editDatesDialog) {
            return current;
        }

        return {
            ...current,
            editDatesDialog: undefined,
        };
    }
);

export type EditDatesDialogClosed = ReturnType<typeof editDatesDialogClosed>;
