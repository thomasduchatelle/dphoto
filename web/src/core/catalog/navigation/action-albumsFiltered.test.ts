import {albumsFiltered} from "./action-albumsFiltered";
import {herselfOwner, loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";

describe("action:albumsFiltered", () => {
    it("should show only directly owned album after the AlbumsFiltered", () => {
        const action = albumsFiltered({criterion: {selfOwned: true, owners: []}});
        const got = action.reducer(
            loadedStateWithTwoAlbums,
            action
        );
        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[0],
            albums: [twoAlbums[0]],
        });
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
        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[1],
            albums: twoAlbums,
        });
    });

    it("should filter albums to those with a certain owner when the filter with that owner is selected", () => {
        const action = albumsFiltered({criterion: {owners: [herselfOwner]}});
        const got = action.reducer(
            {
                ...loadedStateWithTwoAlbums,
                albums: [],
            },
            action
        );
        expect(got).toEqual({
            ...loadedStateWithTwoAlbums,
            albumFilter: loadedStateWithTwoAlbums.albumFilterOptions[2],
            albums: [twoAlbums[1]],
        });
    });
});
