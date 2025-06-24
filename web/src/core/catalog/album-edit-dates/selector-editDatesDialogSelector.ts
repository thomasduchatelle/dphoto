import {CatalogViewerState} from "../language";

export interface EditDatesDialogSelection {
    isOpen: boolean;
    albumName: string;
    startDate: Date;
    endDate: Date;
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
    if (!state.editDatesDialog) {
        return DEFAULT_EDIT_DATES_DIALOG_SELECTION;
    }

    const isDateRangeValid = state.editDatesDialog.startDate <= state.editDatesDialog.endDate;
    const dateRangeError = isDateRangeValid ? undefined : "The end date cannot be before the start date";

    return {
        ...DEFAULT_EDIT_DATES_DIALOG_SELECTION,
        isOpen: true,
        albumName: state.editDatesDialog.albumName,
        startDate: state.editDatesDialog.startDate,
        endDate: state.editDatesDialog.endDate,
        startAtDayStart: state.editDatesDialog.startAtDayStart,
        endAtDayEnd: state.editDatesDialog.endAtDayEnd,
        isLoading: state.editDatesDialog.isLoading,
        errorCode: state.editDatesDialog.error,
        dateRangeError,
        isSaveEnabled: isDateRangeValid && !state.editDatesDialog.isLoading,
    };
}
