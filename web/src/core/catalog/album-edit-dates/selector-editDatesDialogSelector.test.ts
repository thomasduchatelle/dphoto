import {editDatesDialogSelector} from "./selector-editDatesDialogSelector";
import {loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";
import {CatalogViewerState} from "../language";

describe("selector:editDatesDialogSelector", () => {
    it("returns closed dialog state when no edit dialog is open", () => {
        const got = editDatesDialogSelector(loadedStateWithTwoAlbums);

        expect(got).toEqual({
            isOpen: false,
            albumName: "",
            startDate: expect.any(Date),
            endDate: expect.any(Date),
            startAtDayStart: true,
            endAtDayEnd: true,
            isLoading: false,
        });
    });

    it("returns dialog state with loading indicator when dialog is open", () => {
        const stateWithEditDialog: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            editDatesDialog: {
                albumId: twoAlbums[0].albumId,
                albumName: twoAlbums[0].name,
                startDate: new Date("2023-07-10T00:00:00Z"),
                endDate: new Date("2023-07-21T00:00:00Z"),
                isLoading: true,
            },
        };

        const got = editDatesDialogSelector(stateWithEditDialog);

        expect(got).toEqual({
            isOpen: true,
            albumName: twoAlbums[0].name,
            startDate: new Date("2023-07-10T00:00:00Z"),
            endDate: new Date("2023-07-20T00:00:00Z"),
            startAtDayStart: true,
            endAtDayEnd: true,
            isLoading: true,
        });
    });

    it("converts exclusive end date back to display format", () => {
        const stateWithEditDialog: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            editDatesDialog: {
                albumId: twoAlbums[0].albumId,
                albumName: twoAlbums[0].name,
                startDate: new Date("2023-07-10T00:00:00Z"),
                endDate: new Date("2023-07-21T00:00:00Z"),
                isLoading: false,
            },
        };

        const got = editDatesDialogSelector(stateWithEditDialog);

        expect(got.endDate).toEqual(new Date("2023-07-20T00:00:00Z"));
    });
});
