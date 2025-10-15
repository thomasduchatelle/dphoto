import {albumListActionsSelector} from "./selector-albumListActions";
import {albumListActionsPropsForLoadedState, loadedStateWithTwoAlbums} from "../tests/test-helper-state";
import {albumsLoaded} from "./action-albumsLoaded";
import {CatalogViewerState} from "../language";

describe("selector:albumListActionsSelector", () => {
    it("should return deleteButtonEnabled as true when at least one album is owned by the current user", () => {
        const got = albumListActionsSelector(loadedStateWithTwoAlbums);

        expect(got).toEqual(albumListActionsPropsForLoadedState);
    });

    it("should return deleteButtonEnabled as false when no albums are owned by the current user", () => {
        const stateWithOnlySharedAlbums: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            allAlbums: [loadedStateWithTwoAlbums.allAlbums[1]],
        };

        const got = albumListActionsSelector(stateWithOnlySharedAlbums);

        expect(got.deleteButtonEnabled).toBe(false);
    });

    it("should return deleteButtonEnabled as false when there are no albums", () => {
        const action = albumsLoaded({albums: [], filterOptions: []});
        const stateWithNoAlbums = action.reducer(loadedStateWithTwoAlbums, action);

        const got = albumListActionsSelector(stateWithNoAlbums);

        expect(got.deleteButtonEnabled).toBe(false);
    });

    it("should return deleteButtonEnabled as true when albums are loaded and at least one is owned", () => {
        const action = albumsLoaded({
            albums: loadedStateWithTwoAlbums.allAlbums,
            filterOptions: loadedStateWithTwoAlbums.albumFilterOptions,
        });
        const state = action.reducer(loadedStateWithTwoAlbums, action);

        const got = albumListActionsSelector(state);

        expect(got.deleteButtonEnabled).toBe(true);
    });
});
