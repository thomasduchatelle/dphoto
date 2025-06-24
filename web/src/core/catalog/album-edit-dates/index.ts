import {openEditDatesDialogDeclaration} from "./thunk-openEditDatesDialog";
import {closeEditDatesDialogDeclaration} from "./thunk-closeEditDatesDialog";
import {updateAlbumDatesDeclaration} from "./thunk-updateAlbumDates";
import {updateEditDatesDialogStartDateDeclaration} from "./thunk-updateEditDatesDialogStartDate";
import {updateEditDatesDialogEndDateDeclaration} from "./thunk-updateEditDatesDialogEndDate";
import {updateEditDatesDialogStartAtDayStartDeclaration} from "./thunk-updateEditDatesDialogStartAtDayStart";
import {updateEditDatesDialogEndAtDayEndDeclaration} from "./thunk-updateEditDatesDialogEndAtDayEnd";

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
    updateEditDatesDialogStartDate: updateEditDatesDialogStartDateDeclaration,
    updateEditDatesDialogEndDate: updateEditDatesDialogEndDateDeclaration,
    updateEditDatesDialogStartAtDayStart: updateEditDatesDialogStartAtDayStartDeclaration,
    updateEditDatesDialogEndAtDayEnd: updateEditDatesDialogEndAtDayEndDeclaration,
};
