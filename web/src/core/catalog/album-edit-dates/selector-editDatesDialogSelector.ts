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
};

export function editDatesDialogSelector(state: CatalogViewerState): EditDatesDialogSelection {
    if (!state.editDatesDialog) {
        return DEFAULT_EDIT_DATES_DIALOG_SELECTION;
    }

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
    };
}
