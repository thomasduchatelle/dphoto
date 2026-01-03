import {CatalogViewerState} from "../language";
import {createDialogOpened} from "./action-createDialogOpened";
import {loadedStateWithTwoAlbums} from "../tests/test-helper-state";
import {vi} from "vitest";

describe('action:createDialogOpened', () => {
    it('should open with the last week open from Saturday to Monday (9 days) both at start of the day', () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: undefined,
        };

        const mockDate = new Date("2024-03-13T10:00:00Z");
        const OriginalDate = Date;
        vi.spyOn(global, 'Date').mockImplementation(function (this: any, ...args: any[]) {
            if (args.length === 0) {
                return mockDate;
            }
            return new OriginalDate(...args);
        } as any);

        const action = createDialogOpened();
        const newState = action.reducer(state, action);

        const expectedStartDate = new Date("2024-03-02T00:00:00Z");
        const expectedEndDate = new Date("2024-03-11T00:00:00Z");

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

        vi.restoreAllMocks();
    });
});
