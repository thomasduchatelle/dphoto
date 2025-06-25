import {CatalogViewerState, isEditDatesDialog} from "../language";
import {validateDateRange} from "../date-range/date-helper";

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

    const validation = validateDateRange({
        startDate: dialog.startDate,
        endDate: dialog.endDate,
        startAtDayStart: dialog.startAtDayStart,
        endAtDayEnd: dialog.endAtDayEnd,
    });

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
        dateRangeError: validation.dateRangeError,
        isSaveEnabled: validation.areDatesValid && validation.isDateRangeValid && !dialog.isLoading,
    };
}
