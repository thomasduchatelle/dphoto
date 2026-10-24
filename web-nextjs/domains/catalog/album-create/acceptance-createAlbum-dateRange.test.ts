import {Album, AlbumId, CatalogError, CatalogViewerAction, CatalogViewerState} from "../language";
import {loadedStateWithTwoAlbums} from "../tests/test-helper-state";
import {createDialogOpened} from "./action-createDialogOpened";
import {createDialogSelector} from "./selector-createDialogSelector";
import {startDateUpdated} from "../date-range/action-startDateUpdated";
import {endDateUpdated} from "../date-range/action-endDateUpdated";
import {atStartDayToggled} from "../date-range/action-atStartDayToggled";
import {atEndDayToggled} from "../date-range/action-atEndDayToggled";
import {albumNameChanged} from "../base-name-edit/action-albumNameChanged";
import {CreateAlbumPort, CreateAlbumRequest, submitCreateAlbumThunk} from "./thunk-submitCreateAlbum";

describe("acceptance:createAlbum-dateRange", () => {
    let state: CatalogViewerState;

    beforeEach(() => {
        state = loadedStateWithTwoAlbums;
        // Open create dialog
        const openAction = createDialogOpened();
        state = openAction.reducer(state, openAction);
    });

    it("should handle complete date range workflow with shared actions", () => {
        // Initial state should have default dates
        let selection = createDialogSelector(state);
        expect(selection.open).toBe(true);
        expect(selection.canSubmit).toBe(false); // No album name yet
        expect(selection.dateRangeError).toBeUndefined();

        // Update start date
        const startDateAction = startDateUpdated(new Date("2023-06-01"));
        state = startDateAction.reducer(state, startDateAction);

        selection = createDialogSelector(state);
        expect(selection.start).toEqual(new Date("2023-06-01"));
        expect(selection.dateRangeError).toBeUndefined();

        // Update end date
        const endDateAction = endDateUpdated(new Date("2023-06-30"));
        state = endDateAction.reducer(state, endDateAction);

        selection = createDialogSelector(state);
        expect(selection.end).toEqual(new Date("2023-06-30"));
        expect(selection.dateRangeError).toBeUndefined();

        // Toggle start at day start
        const startAtDayStartAction = atStartDayToggled(false);
        state = startAtDayStartAction.reducer(state, startAtDayStartAction);

        selection = createDialogSelector(state);
        expect(selection.startsAtStartOfTheDay).toBe(false);

        // Toggle end at day end
        const endAtDayEndAction = atEndDayToggled(false);
        state = endAtDayEndAction.reducer(state, endAtDayEndAction);

        selection = createDialogSelector(state);
        expect(selection.endsAtEndOfTheDay).toBe(false);
    });

    it("should detect invalid date range when end date is before start date", () => {
        // Set start date after end date
        const startDateAction = startDateUpdated(new Date("2023-06-30"));
        state = startDateAction.reducer(state, startDateAction);

        const endDateAction = endDateUpdated(new Date("2023-06-01"));
        state = endDateAction.reducer(state, endDateAction);

        const selection = createDialogSelector(state);
        expect(selection.dateRangeError).toBe("The end date cannot be before the start date");
        expect(selection.canSubmit).toBe(false);
    });

    it("should clear date range error when dates become valid", () => {
        // Create invalid range first
        const startDateAction = startDateUpdated(new Date("2023-06-30"));
        state = startDateAction.reducer(state, startDateAction);

        const endDateAction = endDateUpdated(new Date("2023-06-01"));
        state = endDateAction.reducer(state, endDateAction);

        let selection = createDialogSelector(state);
        expect(selection.dateRangeError).toBe("The end date cannot be before the start date");

        // Fix the range
        const fixEndDateAction = endDateUpdated(new Date("2023-07-31"));
        state = fixEndDateAction.reducer(state, fixEndDateAction);

        selection = createDialogSelector(state);
        expect(selection.dateRangeError).toBeUndefined();
    });

    it("should surface a folder name clash as an error under the name field", async () => {
        const nameAction = albumNameChanged("Clair Obscur");
        state = nameAction.reducer(state, nameAction);

        const port: CreateAlbumPort = {
            createAlbum: (_request: CreateAlbumRequest): Promise<AlbumId> => {
                return Promise.reject(new CatalogError("AlbumFolderNameAlreadyTakenErr", "An album with this folder name already exists."));
            },
            fetchAlbums: (): Promise<Album[]> => Promise.reject(new Error("should not be called")),
        };

        await expect(
            submitCreateAlbumThunk((action: CatalogViewerAction) => {
                state = action.reducer(state, action);
            }, port, {
                albumName: "Clair Obscur",
                startDate: new Date("2023-06-01"),
                endDate: new Date("2023-06-30"),
                customFolderName: "",
                startAtDayStart: true,
                endAtDayEnd: true,
                isCustomFolderNameEnabled: false,
            })
        ).rejects.toThrow();

        const selection = createDialogSelector(state);
        expect(selection.nameError).toBe("An album with this folder name already exists.");
        expect(selection.error).toBeUndefined();
        expect(selection.canSubmit).toBe(false);
        expect(selection.isLoading).toBe(false);
    });
});
