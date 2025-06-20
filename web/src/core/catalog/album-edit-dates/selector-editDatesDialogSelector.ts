import {CatalogViewerState} from "../language";

export interface EditDatesDialogSelection {
    isOpen: boolean;
    albumName: string;
    startDate: Date;
    endDate: Date;
    startAtDayStart: boolean;
    endAtDayEnd: boolean;
    isLoading: boolean;
}

export function editDatesDialogSelector(state: CatalogViewerState): EditDatesDialogSelection {
    if (!state.editDatesDialog) {
        return {
            isOpen: false,
            albumName: "",
            startDate: new Date(),
            endDate: new Date(),
            startAtDayStart: true,
            endAtDayEnd: true,
            isLoading: false,
        };
    }

    const displayEndDate = new Date(state.editDatesDialog.endDate);
    if (displayEndDate.getHours() === 0 && displayEndDate.getMinutes() === 0 && displayEndDate.getSeconds() === 0) {
        displayEndDate.setDate(displayEndDate.getDate() - 1);
    }

    return {
        isOpen: true,
        albumName: state.editDatesDialog.albumName,
        startDate: state.editDatesDialog.startDate,
        endDate: displayEndDate,
        startAtDayStart: true,
        endAtDayEnd: true,
        isLoading: state.editDatesDialog.isLoading,
    };
}
