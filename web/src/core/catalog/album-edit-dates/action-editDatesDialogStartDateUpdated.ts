import {CatalogViewerState} from "../language";

export interface EditDatesDialogStartDateUpdated {
    type: "EditDatesDialogStartDateUpdated";
    startDate: Date;
}

export function editDatesDialogStartDateUpdated(startDate: Date): EditDatesDialogStartDateUpdated {
    return {
        type: "EditDatesDialogStartDateUpdated",
        startDate,
    };
}

export function reduceEditDatesDialogStartDateUpdated(
    current: CatalogViewerState,
    {startDate}: EditDatesDialogStartDateUpdated,
): CatalogViewerState {
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

export function editDatesDialogStartDateUpdatedReducerRegistration(handlers: any) {
    handlers["EditDatesDialogStartDateUpdated"] = reduceEditDatesDialogStartDateUpdated as (
        state: CatalogViewerState,
        action: EditDatesDialogStartDateUpdated
    ) => CatalogViewerState;
}
