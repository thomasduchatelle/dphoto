export * from "./action-editDatesDialogOpened";
export * from "./action-editDatesDialogClosed";
export * from "./action-albumDatesUpdateStarted";
export * from "./action-albumDatesUpdated";
export * from "./action-editDatesDialogStartDateUpdated";
export * from "./action-editDatesDialogEndDateUpdated";
export * from "./selector-editDatesDialogSelector";
export * from "./thunk-openEditDatesDialog";
export * from "./thunk-closeEditDatesDialog";
export * from "./thunk-updateAlbumDates";
export * from "./thunk-updateEditDatesDialogStartDate";
export * from "./thunk-updateEditDatesDialogEndDate";

/**
 * Thunks related to album date editing.
 *
 * Expected handler types:
 * - `openEditDatesDialog`: `() => void`
 * - `closeEditDatesDialog`: `() => void`
 * - `updateAlbumDates`: `() => Promise<void>`
 * - `updateEditDatesDialogStartDate`: `(startDate: Date | null) => void`
 * - `updateEditDatesDialogEndDate`: `(endDate: Date | null) => void`
 */
export const albumEditDatesThunks = {
    openEditDatesDialog: openEditDatesDialogDeclaration,
    closeEditDatesDialog: closeEditDatesDialogDeclaration,
    updateAlbumDates: updateAlbumDatesDeclaration,
    updateEditDatesDialogStartDate: updateEditDatesDialogStartDateDeclaration,
    updateEditDatesDialogEndDate: updateEditDatesDialogEndDateDeclaration,
};
