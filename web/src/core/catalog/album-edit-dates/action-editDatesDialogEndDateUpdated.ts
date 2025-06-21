import {CatalogViewerState} from "../language";

export interface EditDatesDialogEndDateUpdated {
    type: "EditDatesDialogEndDateUpdated";
    endDate: Date;
}

export function editDatesDialogEndDateUpdated(endDate: Date): EditDatesDialogEndDateUpdated {
    return {
        type: "EditDatesDialogEndDateUpdated",
        endDate,
    };
}

export function reduceEditDatesDialogEndDateUpdated(
    current: CatalogViewerState,
    {endDate}: EditDatesDialogEndDateUpdated,
): CatalogViewerState {
    if (!current.editDatesDialog) {
        return current;
    }
    return {
        ...current,
        editDatesDialog: {
            ...current.editDatesDialog,
            endDate,
        },
    };
}

export function editDatesDialogEndDateUpdatedReducerRegistration(handlers: any) {
    handlers["EditDatesDialogEndDateUpdated"] = reduceEditDatesDialogEndDateUpdated as (
        state: CatalogViewerState,
        action: EditDatesDialogEndDateUpdated
    ) => CatalogViewerState;
}
