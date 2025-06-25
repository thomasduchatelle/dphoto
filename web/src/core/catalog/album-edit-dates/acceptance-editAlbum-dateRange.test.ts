import {CatalogViewerState, EditDatesDialog} from "../language";
import {loadedStateWithTwoAlbums} from "../tests/test-helper-state";
import {editDatesDialogOpened} from "./action-editDatesDialogOpened";
import {editDatesDialogSelector} from "./selector-editDatesDialogSelector";
import {startDateUpdated} from "../date-range/action-startDateUpdated";
import {endDateUpdated} from "../date-range/action-endDateUpdated";
import {atStartDayToggled} from "../date-range/action-atStartDayToggled";
import {atEndDayToggled} from "../date-range/action-atEndDayToggled";

describe("acceptance:editAlbum-dateRange", () => {
    let state: CatalogViewerState;

    beforeEach(() => {
        state = {
            ...loadedStateWithTwoAlbums,
            mediasLoadedFromAlbumId: loadedStateWithTwoAlbums.albums[0].albumId,
        };
        // Open edit dialog
        const openAction = editDatesDialogOpened();
        state = openAction.reducer(state, openAction);
    });

    it("should handle complete date range workflow with shared actions", () => {
        // Initial state should have album dates
        let selection = editDatesDialogSelector(state);
        expect(selection.isOpen).toBe(true);
        expect(selection.isSaveEnabled).toBe(true);
        expect(selection.dateRangeError).toBeUndefined();

        // Update start date
        const startDateAction = startDateUpdated(new Date("2023-06-01"));
        state = startDateAction.reducer(state, startDateAction);
        
        selection = editDatesDialogSelector(state);
        expect(selection.startDate).toEqual(new Date("2023-06-01"));
        expect(selection.dateRangeError).toBeUndefined();
        expect(selection.errorCode).toBeUndefined();

        // Update end date
        const endDateAction = endDateUpdated(new Date("2023-06-30"));
        state = endDateAction.reducer(state, endDateAction);
        
        selection = editDatesDialogSelector(state);
        expect(selection.endDate).toEqual(new Date("2023-06-30"));
        expect(selection.dateRangeError).toBeUndefined();
        expect(selection.errorCode).toBeUndefined();

        // Toggle start at day start
        const startAtDayStartAction = atStartDayToggled(false);
        state = startAtDayStartAction.reducer(state, startAtDayStartAction);
        
        selection = editDatesDialogSelector(state);
        expect(selection.startAtDayStart).toBe(false);
        expect(selection.errorCode).toBeUndefined();

        // Toggle end at day end
        const endAtDayEndAction = atEndDayToggled(false);
        state = endAtDayEndAction.reducer(state, endAtDayEndAction);
        
        selection = editDatesDialogSelector(state);
        expect(selection.endAtDayEnd).toBe(false);
        expect(selection.errorCode).toBeUndefined();
    });

    it("should detect invalid date range when end date is before start date", () => {
        // Set start date after end date
        const startDateAction = startDateUpdated(new Date("2023-06-30"));
        state = startDateAction.reducer(state, startDateAction);
        
        const endDateAction = endDateUpdated(new Date("2023-06-01"));
        state = endDateAction.reducer(state, endDateAction);
        
        const selection = editDatesDialogSelector(state);
        expect(selection.dateRangeError).toBe("The end date cannot be before the start date");
        expect(selection.isSaveEnabled).toBe(false);
    });

    it("should clear date range error when dates become valid", () => {
        // Create invalid range first
        const startDateAction = startDateUpdated(new Date("2023-06-30"));
        state = startDateAction.reducer(state, startDateAction);
        
        const endDateAction = endDateUpdated(new Date("2023-06-01"));
        state = endDateAction.reducer(state, endDateAction);
        
        let selection = editDatesDialogSelector(state);
        expect(selection.dateRangeError).toBe("The end date cannot be before the start date");

        // Fix the range
        const fixEndDateAction = endDateUpdated(new Date("2023-07-31"));
        state = fixEndDateAction.reducer(state, fixEndDateAction);
        
        selection = editDatesDialogSelector(state);
        expect(selection.dateRangeError).toBeUndefined();
        expect(selection.isSaveEnabled).toBe(true);
    });

    it("should clear errors when updating dates", () => {
        const dialog = state.dialog as EditDatesDialog

        // Set an error first
        state = {
            ...state,
            dialog: {
                ...dialog,
                error: "Some previous error",
            },
        };

        let selection = editDatesDialogSelector(state);
        expect(selection.errorCode).toBe("Some previous error");

        // Update start date should clear error
        const startDateAction = startDateUpdated(new Date("2023-06-01"));
        state = startDateAction.reducer(state, startDateAction);
        
        selection = editDatesDialogSelector(state);
        expect(selection.errorCode).toBeUndefined();
    });
});
