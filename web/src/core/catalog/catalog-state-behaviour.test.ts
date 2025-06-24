import {herselfUser, loadedStateWithTwoAlbums, twoAlbums} from "./tests/test-helper-state";
import {Album, albumIdEquals, UserDetails} from "./language";
import {albumsFiltered, SELF_OWNED_ALBUM_FILTER_CRITERION} from "./navigation";
import {albumAccessGranted, albumAccessRevoked, sharingDialogSelector, sharingModalClosed, sharingModalErrorOccurred, sharingModalOpened} from "./sharing";
import {catalogReducer} from "./actions";
import {editDatesDialogOpened} from "./album-edit-dates/action-editDatesDialogOpened";
import {updateAlbumDatesDeclaration, UpdateAlbumDatesPort, updateAlbumDatesThunk} from "./album-edit-dates/thunk-updateAlbumDates";
import {editDatesDialogSelector} from "./album-edit-dates";

describe("State: behaviour", () => {
    it("keeps the album shares consistent when closing and reopening the dialog", () => {
        const initialState = loadedStateWithTwoAlbums;

        const album = twoAlbums[0];
        const openAction = sharingModalOpened(album.albumId);
        let state = catalogReducer(initialState, openAction);

        const newUser: UserDetails = {email: "bob@example.com", name: "Bob", picture: "bob-face.jpg"};
        const addAction = albumAccessGranted({user: newUser});
        state = catalogReducer(state, addAction);

        const removeAction = albumAccessRevoked(herselfUser.email);
        state = catalogReducer(state, removeAction);

        const closeAction = sharingModalClosed();
        state = catalogReducer(state, closeAction);

        state = catalogReducer(state, openAction);

        expect(sharingDialogSelector(state)).toEqual({
            open: true,
            sharedWith: [
                {user: newUser}
            ],
            suggestions: [],
        });

        const updatedAlbum = state.allAlbums.find(a => albumIdEquals(a.albumId, album.albumId)) as Album;
        expect(updatedAlbum.sharedWith).toEqual([
            {user: newUser}
        ]);
    });

    it("keeps the album shares consistent when re-applying the album filter", () => {
        const initialState = loadedStateWithTwoAlbums;

        const album = twoAlbums[0];
        const openAction = sharingModalOpened(album.albumId);
        let state = catalogReducer(initialState, openAction);

        const newUser: UserDetails = {email: "bob@example.com", name: "Bob", picture: "bob-face.jpg"};
        const addAction = albumAccessGranted({user: newUser});
        state = catalogReducer(state, addAction);

        const removeAction = albumAccessRevoked(herselfUser.email);
        state = catalogReducer(state, removeAction);

        const closeAction = sharingModalClosed();
        state = catalogReducer(state, closeAction);

        const filterAction = albumsFiltered({criterion: SELF_OWNED_ALBUM_FILTER_CRITERION});
        state = catalogReducer(state, filterAction);

        // Only the updated album should be present, and its shares should be correct
        const filteredAlbum = state.albums.find(a => albumIdEquals(a.albumId, album.albumId)) as Album;
        expect(filteredAlbum.sharedWith).toEqual([
            {user: newUser}
        ]);
    });

    it("reverts the changes made by addSharingAction, dispatched optimistically, when handling an error", () => {
        const userToGrant: UserDetails = {email: "bob@example.com", name: "Bob", picture: "bob-face.jpg"};

        let state = catalogReducer(loadedStateWithTwoAlbums, sharingModalOpened(twoAlbums[0].albumId));
        state = catalogReducer(state, albumAccessGranted({user: userToGrant}));
        state = catalogReducer(state, sharingModalErrorOccurred({type: "grant", message: "Failed to add user", email: userToGrant.email}));
        state = catalogReducer(state, sharingModalClosed());

        const filterAction = albumsFiltered({criterion: SELF_OWNED_ALBUM_FILTER_CRITERION});
        const ownedAlbums = catalogReducer(state, filterAction).albums;

        expect(state).toEqual(loadedStateWithTwoAlbums);
        expect(ownedAlbums).toEqual([loadedStateWithTwoAlbums.albums[0]]);
    });

    const testCases = [
        {
            name: "full day album (start and end at midnight)",
            album: {
                ...twoAlbums[0],
                start: new Date("2025-01-01T00:00:00.000Z"),
                end: new Date("2025-02-01T00:00:00.000Z"),
            },
            displayedEndRegex: /^2025-01-31.*/,
        },
        {
            name: "album with round start time and precise end time",
            album: {
                ...twoAlbums[0],
                start: new Date("2025-01-01T00:00:00.000Z"),
                end: new Date("2025-01-31T14:43:00.000Z"),
            },
            displayedEndRegex: /^2025-01-31T14:42.*/,
        },
        {
            name: "album with precise start time and round end time",
            album: {
                ...twoAlbums[0],
                start: new Date("2025-01-01T00:00:00.000Z"),
                end: new Date("2025-01-31T16:00:00.000Z"),
            },
            displayedEndRegex: /^2025-01-31T16:00.*/,
        },
    ];

    for (const testCase of testCases) {
        it(`preserves original album dates when edit dialog is opened and submitted without changes: ${testCase.name}`, async () => {
            class UpdateAlbumDatesPortFake implements UpdateAlbumDatesPort {
                public updatedAlbums: { albumId: any, startDate: Date, endDate: Date }[] = [];

                async updateAlbumDates(albumId: any, startDate: Date, endDate: Date): Promise<void> {
                    this.updatedAlbums.push({albumId, startDate, endDate});
                }

                async fetchAlbums(): Promise<Album[]> {
                    return twoAlbums;
                }

                async fetchMedias(): Promise<any[]> {
                    return [];
                }
            }

            const fakePort = new UpdateAlbumDatesPortFake();
            const dispatched: any[] = [];

            const stateWithAlbum = {
                ...loadedStateWithTwoAlbums,
                albums: [testCase.album],
                allAlbums: [testCase.album],
                mediasLoadedFromAlbumId: testCase.album.albumId,
            };

            let state = catalogReducer(stateWithAlbum, editDatesDialogOpened());

            await updateAlbumDatesThunk(
                dispatched.push.bind(dispatched),
                fakePort,
                updateAlbumDatesDeclaration.selector(state),
            );

            const {endDate} = editDatesDialogSelector(state);
            expect(endDate.toISOString()).toMatch(testCase.displayedEndRegex)

            expect(fakePort.updatedAlbums).toHaveLength(1);
            expect(fakePort.updatedAlbums[0]).toEqual({
                albumId: testCase.album.albumId,
                startDate: testCase.album.start,
                endDate: testCase.album.end,
            });

            fakePort.updatedAlbums = [];
        });
    }
});
