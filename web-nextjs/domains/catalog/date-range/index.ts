import {updateDateRangeStartDateDeclaration} from "./thunk-updateDateRangeStartDate";
import {updateDateRangeEndDateDeclaration} from "./thunk-updateDateRangeEndDate";
import {updateDateRangeStartAtDayStartDeclaration} from "./thunk-updateDateRangeStartAtDayStart";
import {updateDateRangeEndAtDayEndDeclaration} from "./thunk-updateDateRangeEndAtDayEnd";

/**
 * Shared thunks for date range functionality.
 *
 * Expected handler types:
 * - `updateDateRangeStartDate`: `(startDate: Date | null) => void`
 * - `updateDateRangeEndDate`: `(endDate: Date | null) => void`
 * - `updateDateRangeStartAtDayStart`: `(startAtDayStart: boolean) => void`
 * - `updateDateRangeEndAtDayEnd`: `(endAtDayEnd: boolean) => void`
 */
export const dateRangeThunks = {
    updateDateRangeStartDate: updateDateRangeStartDateDeclaration,
    updateDateRangeEndDate: updateDateRangeEndDateDeclaration,
    updateDateRangeStartAtDayStart: updateDateRangeStartAtDayStartDeclaration,
    updateDateRangeEndAtDayEnd: updateDateRangeEndAtDayEndDeclaration,
};
