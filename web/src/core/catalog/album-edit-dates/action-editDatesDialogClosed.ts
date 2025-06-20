import {CatalogViewerState} from "../language";

export interface EditDatesDialogClosed {
    type: "EditDatesDialogClosed";
}

export function editDatesDialogClosed(): EditDatesDialogClosed {
    return {
        type: "EditDatesDialogClosed",
    };
}

export function reduceEditDatesDialogClosed(
    current: CatalogViewerState,
    _: EditDatesDialogClosed,
): CatalogViewerState {
    if (!current.editDatesDialog) {
        return current;
    }

    return {
        ...current,
        editDatesDialog: undefined,
    };
}

export function editDatesDialogClosedReducerRegistration(handlers: any) {
    handlers["EditDatesDialogClosed"] = reduceEditDatesDialogClosed as (
        state: CatalogViewerState,
        action: EditDatesDialogClosed
    ) => CatalogViewerState;
}
