import {createAction} from "src/libs/daction";
import {CatalogViewerState, DateRangeState} from "../language";
import {setDateToEndOfDay} from "./date-helper";

function isDateRangeDialog(dialog: any): dialog is DateRangeState & { error?: string } {
    return dialog && 'startDate' in dialog && 'endDate' in dialog && 'startAtDayStart' in dialog && 'endAtDayEnd' in dialog;
}

export const atEndDayToggled = createAction<CatalogViewerState, boolean>(
    "AtEndDayToggled",
    (current: CatalogViewerState, endAtDayEnd: boolean) => {
        const dialog = current.dialog;

        if (!isDateRangeDialog(dialog) || !dialog.endDate) {
            return current;
        }

        let updatedEndDate = dialog.endDate;
        if (!endAtDayEnd && dialog.endDate) {
            updatedEndDate = setDateToEndOfDay(dialog.endDate);
        }

        return {
            ...current,
            dialog: {
                ...dialog,
                endAtDayEnd,
                endDate: updatedEndDate,
                error: undefined,
            },
        };
    }
);

export type AtEndDayToggled = ReturnType<typeof atEndDayToggled>;
