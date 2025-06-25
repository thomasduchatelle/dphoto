import {openEditDatesDialogDeclaration} from "./thunk-openEditDatesDialog";
import {closeEditDatesDialogDeclaration} from "./thunk-closeEditDatesDialog";
import {updateAlbumDatesDeclaration} from "./thunk-updateAlbumDates";
import {dateRangeThunks} from "../date-range";

export * from "./selector-editDatesDialogSelector";
export type {UpdateAlbumDatesPort} from "./thunk-updateAlbumDates";

/**
 * Thunks related to album date editing.
 *
 * Expected handler types:
 * - `openEditDatesDialog`: `() => void`
 * - `closeEditDatesDialog`: `() => void`
 * - `updateAlbumDates`: `() => Promise<void>`
 * - `updateEditDatesDialogStartDate`: `(startDate: Date | null) => void`
 * - `updateEditDatesDialogEndDate`: `(endDate: Date | null) => void`
 * - `updateEditDatesDialogStartAtDayStart`: `(startAtDayStart: boolean) => void`
 * - `updateEditDatesDialogEndAtDayEnd`: `(endAtDayEnd: boolean) => void`
 */
export const albumEditDatesThunks = {
    openEditDatesDialog: openEditDatesDialogDeclaration,
    closeEditDatesDialog: closeEditDatesDialogDeclaration,
    updateAlbumDates: updateAlbumDatesDeclaration,
    updateEditDatesDialogStartDate: dateRangeThunks.updateDateRangeStartDate,
    updateEditDatesDialogEndDate: dateRangeThunks.updateDateRangeEndDate,
    updateEditDatesDialogStartAtDayStart: dateRangeThunks.updateDateRangeStartAtDayStart,
    updateEditDatesDialogEndAtDayEnd: dateRangeThunks.updateDateRangeEndAtDayEnd,
};
