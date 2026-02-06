import {createAction} from "@/libs/daction";
import {CatalogViewerState, DateRangeState} from "../language";
import {setDateToStartOfDay} from "./date-helper";

function isDateRangeDialog(dialog: any): dialog is DateRangeState & { error?: string } {
    return dialog && 'startDate' in dialog && 'endDate' in dialog && 'startAtDayStart' in dialog && 'endAtDayEnd' in dialog;
}

export const atStartDayToggled = createAction<CatalogViewerState, boolean>(
    "AtStartDayToggled",
    (current: CatalogViewerState, startAtDayStart: boolean) => {
        const dialog = current.dialog;

        if (!isDateRangeDialog(dialog) || !dialog.startDate) {
            return current;
        }

        return {
            ...current,
            dialog: {
                ...dialog,
                startAtDayStart,
                startDate: startAtDayStart ? dialog.startDate : setDateToStartOfDay(dialog.startDate),
                error: undefined,
            },
        };
    }
);

export type AtStartDayToggled = ReturnType<typeof atStartDayToggled>;
