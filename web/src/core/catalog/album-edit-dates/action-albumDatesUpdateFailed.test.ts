import {albumDatesUpdateFailed} from "./action-albumDatesUpdateFailed";
import {CatalogViewerState} from "../language";
import {loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";
import {editDatesDialogSelector} from "./selector-editDatesDialogSelector";

describe("action:albumDatesUpdateFailed", () => {
    it("sets error and stops loading when dialog is open", () => {
        const stateWithEditDialog: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            editDatesDialog: {
                albumId: twoAlbums[0].albumId,
                albumName: twoAlbums[0].name,
                startDate: twoAlbums[0].start,
                endDate: twoAlbums[0].end,
                isLoading: true,
            },
        };

        const errorMessage = "Network error occurred";
        const action = albumDatesUpdateFailed({error: errorMessage});
        const got = action.reducer(stateWithEditDialog, action);

        expect(editDatesDialogSelector(got)).toEqual({
            isOpen: true,
            albumName: twoAlbums[0].name,
            startDate: twoAlbums[0].start,
            endDate: twoAlbums[0].end,
            startAtDayStart: true,
            endAtDayEnd: true,
            isLoading: false,
            errorCode: errorMessage,
        });
    });

    it("does nothing when dialog is closed", () => {
        const state: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            editDatesDialog: undefined,
        };

        const action = albumDatesUpdateFailed({error: "Some error"});
        const got = action.reducer(state, action);

        expect(got).toBe(state);
    });
});
