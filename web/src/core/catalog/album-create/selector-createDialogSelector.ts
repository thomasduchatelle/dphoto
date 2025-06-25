import {CatalogViewerState, isCreateDialog} from "../language";
import {validateDateRange} from "../date-range/date-helper";

export interface CreateDialogSelection {
    open: boolean;
    name: string;
    start: Date | null;
    end: Date | null;
    forceFolderName: string;
    startsAtStartOfTheDay: boolean;
    endsAtEndOfTheDay: boolean;
    withCustomFolderName: boolean;
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
            name: "",
            start: null,
            end: null,
            forceFolderName: "",
            startsAtStartOfTheDay: true,
            endsAtEndOfTheDay: true,
            withCustomFolderName: false,
            isLoading: false,
            canSubmit: false,
            dateRangeError: undefined,
        };
    }

    const validation = validateDateRange({
        startDate: dialog.startDate,
        endDate: dialog.endDate,
        startAtDayStart: dialog.startAtDayStart,
        endAtDayEnd: dialog.endAtDayEnd,
    });

    return {
        open: true,
        name: dialog.name,
        start: dialog.startDate,
        end: dialog.endDate,
        forceFolderName: dialog.forceFolderName,
        startsAtStartOfTheDay: dialog.startAtDayStart,
        endsAtEndOfTheDay: dialog.endAtDayEnd,
        withCustomFolderName: dialog.withCustomFolderName,
        isLoading: dialog.isLoading,
        error: dialog.error,
        dateRangeError: validation.dateRangeError,
        canSubmit: dialog.name.length > 0 && validation.areDatesValid && validation.isDateRangeValid && !dialog.isLoading,
    };
}
