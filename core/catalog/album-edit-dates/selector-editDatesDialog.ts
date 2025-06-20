export interface EditDatesDialogSelection {
    isOpen: boolean;
    isLoading: boolean;
    albumName: string;
    currentStartDate: Date;
    currentEndDate: Date;
}

export function editDatesDialogSelector(state: CatalogViewerState): EditDatesDialogSelection {
    if (!state.editDatesDialog) {
        return {
            isOpen: false,
            isLoading: false,
            albumName: "",
            currentStartDate: new Date(),
            currentEndDate: new Date()
        };
    }

    return {
        isOpen: true,
        isLoading: state.editDatesDialog.isLoading,
        albumName: state.editDatesDialog.albumName,
        currentStartDate: state.editDatesDialog.startDate,
        currentEndDate: state.editDatesDialog.endDate
    };
}
