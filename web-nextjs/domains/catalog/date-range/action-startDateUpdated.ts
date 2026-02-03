import {createAction} from "@/libs/daction";
import {CatalogViewerState, DateRangeState, isCreateDialog, isEditDatesDialog} from "../language";

function isDateRangeDialog(dialog: any): dialog is DateRangeState & { error?: string } {
    return dialog && 'startDate' in dialog && 'endDate' in dialog && 'startAtDayStart' in dialog && 'endAtDayEnd' in dialog;
}

export const startDateUpdated = createAction<CatalogViewerState, Date | null>(
    "StartDateUpdated",
    (current: CatalogViewerState, startDate: Date | null) => {
        const dialog = current.dialog;

        if (!isDateRangeDialog(dialog)) {
            return current;
        }

        if (isCreateDialog(dialog)) {
            return {
                ...current,
                dialog: {
                    ...dialog,
                    startDate,
                    error: dialog.error === "AlbumStartAndEndDateMandatoryErr" ? undefined : dialog.error,
                },
            };
        }

        if (isEditDatesDialog(dialog)) {
            return {
                ...current,
                dialog: {
                    ...dialog,
                    startDate,
                    error: undefined,
                },
            };
        }

        return current;
    }
);

export type StartDateUpdated = ReturnType<typeof startDateUpdated>;
