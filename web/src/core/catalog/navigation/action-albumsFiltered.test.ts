import {albumsFiltered} from "./action-albumsFiltered";
import {loadedStateWithTwoAlbums, selectionForLoadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";

describe("action:albumsFiltered", () => {
    it("should show only directly owned album after the AlbumsFiltered", () => {
        const action = albumsFiltered({criterion: {selfOwned: true, owners: []}});
        const got = action.reducer(
            loadedStateWithTwoAlbums,
            action
        );
        expect(got.albums).toEqual([twoAlbums[0]]);
        expect(got.albumFilter).toEqual(loadedStateWithTwoAlbums.albumFilterOptions[0]);
    });

    it("should show all albums when the filter moves back to 'All albums'", () => {
        const action = albumsFiltered({criterion: {owners: []}});
        const got = action.reducer(
            {
                ...loadedStateWithTwoAlbums,
                albums: [],
            },
            action
        );
        expect(got.albums).toEqual(twoAlbums);
        expect(got.albumFilter).toEqual(loadedStateWithTwoAlbums.albumFilterOptions[1]);
    });

    it("should filter albums to those with a certain owner when the filter with that owner is selected", () => {
        const action = albumsFiltered({criterion: {owners: [{name: "Herself", users: []}]}});
        const got = action.reducer(
            {
                ...loadedStateWithTwoAlbums,
                albums: [],
            },
            action
        );
        expect(got.albums).toEqual([twoAlbums[1]]);
        expect(got.albumFilter).toEqual(loadedStateWithTwoAlbums.albumFilterOptions[2]);
    });
});
