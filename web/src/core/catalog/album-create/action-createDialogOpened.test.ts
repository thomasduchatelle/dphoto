import {CatalogViewerState} from "../language";
import {createDialogOpened} from "./action-createDialogOpened";
import {loadedStateWithTwoAlbums} from "../tests/test-helper-state";

describe('action:createDialogOpened', () => {
    it('should open with the last week open from Saturday to Monday (9 days) both at start of the day', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: undefined,
        };

        // Mock Date to control the "current" date for consistent test results
        const mockDate = new Date("2024-03-13T10:00:00Z"); // A Wednesday
        jest.spyOn(global, 'Date').mockImplementation(() => mockDate);

        const action = createDialogOpened();
        const newState = action.reducer(state, action);

        const expectedStartDate = new Date("2024-03-02T00:00:00Z"); // Saturday of previous week
        const expectedEndDate = new Date("2024-03-11T00:00:00Z"); // Monday of current week

        expect(newState.dialog).toEqual({
            type: "CreateDialog",
            name: "",
            startDate: expectedStartDate,
            endDate: expectedEndDate,
            startAtDayStart: true,
            endAtDayEnd: true,
            forceFolderName: "",
            withCustomFolderName: false,
            isLoading: false,
        });

        // Restore original Date object
        jest.restoreAllMocks();
    });
});
