import {editDatesDialogSelector} from "./selector-editDatesDialog";
import {CatalogViewerState, EditDatesDialogState} from "../language/CatalogViewerState";
import {initialCatalogState} from "../tests/test-helper-state";
import {myselfUser} from "../tests/test-helper-state";

describe("selector:editDatesDialog", () => {
    const editDatesDialog: EditDatesDialogState = {
        albumId: {owner: "myself", folderName: "album1"},
        albumName: "Test Album",
        startDate: new Date("2023-01-01"),
        endDate: new Date("2023-01-31"),
        isLoading: true
    };

    it("returns dialog data when dialog is open", () => {
        const state: CatalogViewerState = {
            ...initialCatalogState(myselfUser),
            editDatesDialog
        };

        const result = editDatesDialogSelector(state);

        expect(result).toEqual({
            isOpen: true,
            isLoading: true,
            albumName: "Test Album",
            currentStartDate: new Date("2023-01-01"),
            currentEndDate: new Date("2023-01-31")
        });
    });

    it("returns closed dialog data when dialog is not open", () => {
        const state = initialCatalogState(myselfUser);

        const result = editDatesDialogSelector(state);

        expect(result).toEqual({
            isOpen: false,
            isLoading: false,
            albumName: "",
            currentStartDate: expect.any(Date),
            currentEndDate: expect.any(Date)
        });
    });
});
