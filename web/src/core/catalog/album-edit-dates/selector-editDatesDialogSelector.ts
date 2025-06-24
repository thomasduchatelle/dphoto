import {CatalogViewerState, isEditDatesDialog} from "../language";

export interface EditDatesDialogSelection {
    isOpen: boolean;
    albumName: string;
    startDate: Date | null;
    endDate: Date | null;
    startAtDayStart: boolean;
    endAtDayEnd: boolean;
    isLoading: boolean;
    errorCode?: string;
    dateRangeError?: string;
    isSaveEnabled: boolean;
}

export const DEFAULT_EDIT_DATES_DIALOG_SELECTION: EditDatesDialogSelection = {
    isOpen: false,
    albumName: "",
    startDate: new Date(),
    endDate: new Date(),
    startAtDayStart: true,
    endAtDayEnd: true,
    isLoading: false,
    errorCode: undefined,
    dateRangeError: undefined,
    isSaveEnabled: true,
};

export function editDatesDialogSelector(state: CatalogViewerState): EditDatesDialogSelection {
    const dialog = state.dialog;
    if (!isEditDatesDialog(dialog)) {
        return DEFAULT_EDIT_DATES_DIALOG_SELECTION;
    }

    const areDatesValid = dialog.startDate !== null && dialog.endDate !== null;
    const isDateRangeValid = !(dialog.startDate !== null && dialog.endDate !== null) || dialog.startDate <= dialog.endDate;
    const dateRangeError = isDateRangeValid ? undefined : "The end date cannot be before the start date";

    return {
        ...DEFAULT_EDIT_DATES_DIALOG_SELECTION,
        isOpen: true,
        albumName: dialog.albumName,
        startDate: dialog.startDate,
        endDate: dialog.endDate,
        startAtDayStart: dialog.startAtDayStart,
        endAtDayEnd: dialog.endAtDayEnd,
        isLoading: dialog.isLoading,
        errorCode: dialog.error,
        dateRangeError,
        isSaveEnabled: areDatesValid && isDateRangeValid && !dialog.isLoading,
    };
}
