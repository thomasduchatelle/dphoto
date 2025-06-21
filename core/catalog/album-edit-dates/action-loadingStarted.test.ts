import {reduceLoadingStarted, loadingStarted} from "./action-loadingStarted";
import {CatalogViewerState, EditDatesDialogState} from "../language/CatalogViewerState";
import {initialCatalogState} from "../tests/test-helper-state";
import {myselfUser} from "../tests/test-helper-state";

describe("action:loadingStarted", () => {
    const editDatesDialog: EditDatesDialogState = {
        albumId: {owner: "myself", folderName: "album1"},
        albumName: "Test Album",
        startDate: new Date("2023-01-01"),
        endDate: new Date("2023-01-31"),
        isLoading: false
    };

    it("sets isLoading to true when editDatesDialog exists", () => {
        const initialState: CatalogViewerState = {
            ...initialCatalogState(myselfUser),
            editDatesDialog
        };

        const result = reduceLoadingStarted(initialState, loadingStarted());

        expect(result).toEqual({
            ...initialState,
            editDatesDialog: {
                ...editDatesDialog,
                isLoading: true
            }
        });
    });

    it("returns unchanged state when editDatesDialog does not exist", () => {
        const initialState = initialCatalogState(myselfUser);

        const result = reduceLoadingStarted(initialState, loadingStarted());

        expect(result).toEqual(initialState);
    });
});
