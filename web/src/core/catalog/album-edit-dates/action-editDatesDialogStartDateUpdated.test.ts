import {editDatesDialogStartDateUpdated} from "./action-editDatesDialogStartDateUpdated";
import {CatalogViewerState} from "../language";
import {loadedStateWithTwoAlbums} from "../tests/test-helper-state";

describe("action:editDatesDialogStartDateUpdated", () => {
    it("updates the start date when dialog is open", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            editDatesDialog: {
                albumId: {owner: "myself", folderName: "summer-trip"},
                albumName: "Summer Trip",
                startDate: new Date("2023-07-01T00:00:00"),
                endDate: new Date("2023-08-01T00:00:00"),
                isLoading: false,
            },
        };

        const newStartDate = new Date("2023-07-15T00:00:00");
        const action = editDatesDialogStartDateUpdated(newStartDate);
        const got = action.reducer(state, action);

        expect(got.editDatesDialog?.startDate).toEqual(newStartDate);
    });

    it("does nothing when dialog is closed", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            editDatesDialog: undefined,
        };

        const newStartDate = new Date("2023-07-15T00:00:00");
        const action = editDatesDialogStartDateUpdated(newStartDate);
        const got = action.reducer(state, action);

        expect(got).toBe(state);
    });

    it("supports action comparison for testing", () => {
        const date = new Date("2023-07-15T00:00:00");
        const action1 = editDatesDialogStartDateUpdated(date);
        const action2 = editDatesDialogStartDateUpdated(date);
        
        expect(action1).toEqual(action2);
        expect([action1]).toContainEqual(action2);
    });

    it("demonstrates the new simplified API with payload", () => {
        const date = new Date("2023-07-15T00:00:00");
        const action = editDatesDialogStartDateUpdated(date);
        
        expect(action.type).toBe("EditDatesDialogStartDateUpdated");
        expect(action.payload).toBe(date);
        expect(typeof action.reducer).toBe("function");
    });

    it("distinguishes between different payloads", () => {
        const date1 = new Date("2023-07-15T00:00:00");
        const date2 = new Date("2023-07-16T00:00:00");
        
        const action1 = editDatesDialogStartDateUpdated(date1);
        const action2 = editDatesDialogStartDateUpdated(date2);
        
        expect(action1).not.toEqual(action2);
        expect([action1]).not.toContainEqual(action2);
    });
});
