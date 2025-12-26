import {CatalogViewerState} from "../language";
import {createDialogOpened} from "./action-createDialogOpened";
import {loadedStateWithTwoAlbums} from "../tests/test-helper-state";
import {beforeEach, vi} from "vitest";

describe('action:createDialogOpened', () => {
    beforeEach(() => {
        // Reset time mocks before each test
        vi.useRealTimers();
    });

    it('should open with the last week open from Saturday to Monday (9 days) both at start of the day', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: undefined,
        };

        // Mock Date to control the "current" date for consistent test results
        const mockDate = new Date("2024-03-13T10:00:00Z"); // A Wednesday
        vi.setSystemTime(mockDate);

        const action = createDialogOpened();
        const newState = action.reducer(state, action);

        const expectedStartDate = new Date("2024-03-02T00:00:00Z"); // Saturday of previous week
        const expectedEndDate = new Date("2024-03-11T00:00:00Z"); // Monday of current week

        expect(newState.dialog).toEqual({
            type: "CreateDialog",
            albumId: {owner: "", folderName: ""},
            albumName: "",
            customFolderName: "",
            isCustomFolderNameEnabled: false,
            nameError: {},
            startDate: expectedStartDate,
            endDate: expectedEndDate,
            startAtDayStart: true,
            endAtDayEnd: true,
            isLoading: false,
        });

        // Restore original Date object
        vi.useRealTimers();
    });
});
