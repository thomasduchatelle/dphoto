import {createAction} from "src/libs/daction";
import {CatalogViewerState, isCreateDialog, isEditDatesDialog, DateRangeState} from "../language";

function isDateRangeDialog(dialog: any): dialog is DateRangeState & { error?: string } {
    return dialog && 'startDate' in dialog && 'endDate' in dialog && 'startAtDayStart' in dialog && 'endAtDayEnd' in dialog;
}

export const endDateUpdated = createAction<CatalogViewerState, Date | null>(
    "EndDateUpdated",
    (current: CatalogViewerState, endDate: Date | null) => {
        const dialog = current.dialog;
        
        if (!isDateRangeDialog(dialog)) {
            return current;
        }
        
        if (isCreateDialog(dialog)) {
            return {
                ...current,
                dialog: {
                    ...dialog,
                    endDate,
                    error: dialog.error === "AlbumStartAndEndDateMandatoryErr" ? undefined : dialog.error,
                },
            };
        }
        
        if (isEditDatesDialog(dialog)) {
            return {
                ...current,
                dialog: {
                    ...dialog,
                    endDate,
                    error: undefined,
                },
            };
        }
        
        return current;
    }
);

export type EndDateUpdated = ReturnType<typeof endDateUpdated>;
