import {CatalogViewerState, isCreateDialog} from "../language";
import {validateDateRange} from "../date-range/date-helper";
import {baseEditNameSelector, BaseEditNameSelection} from "../base-name-edit/selector-baseEditNameSelector";

export interface CreateDialogSelection extends BaseEditNameSelection {
    open: boolean;
    start: Date | null;
    end: Date | null;
    startsAtStartOfTheDay: boolean;
    endsAtEndOfTheDay: boolean;
    isLoading: boolean;
    error?: string;
    canSubmit: boolean;
    dateRangeError?: string;
}

export function createDialogSelector(state: CatalogViewerState): CreateDialogSelection {
    const dialog = state.dialog;
    
    if (!isCreateDialog(dialog)) {
        return {
            open: false,
            albumName: "",
            originalName: "",
            customFolderName: "",
            isCustomFolderNameEnabled: false,
            start: null,
            end: null,
            startsAtStartOfTheDay: true,
            endsAtEndOfTheDay: true,
            isLoading: false,
            canSubmit: false,
            dateRangeError: undefined,
        };
    }

    const baseSelection = baseEditNameSelector(state, dialog);
    
    const validation = validateDateRange({
        startDate: dialog.startDate,
        endDate: dialog.endDate,
        startAtDayStart: dialog.startAtDayStart,
        endAtDayEnd: dialog.endAtDayEnd,
    });

    return {
        ...baseSelection,
        open: true,
        start: dialog.startDate,
        end: dialog.endDate,
        startsAtStartOfTheDay: dialog.startAtDayStart,
        endsAtEndOfTheDay: dialog.endAtDayEnd,
        isLoading: dialog.isLoading,
        error: dialog.error,
        dateRangeError: validation.dateRangeError,
        canSubmit: baseSelection.albumName.length > 0 && validation.areDatesValid && validation.isDateRangeValid && !dialog.isLoading,
    };
}
