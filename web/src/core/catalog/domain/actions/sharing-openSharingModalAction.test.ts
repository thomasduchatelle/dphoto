import {openSharingModalAction, reduceOpenSharingModal} from "./sharing-openSharingModalAction";
import {herselfUser, loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";
import {SharingDialogFrag, sharingDialogSelector} from "./selector-sharingDialogSelector";
import {UserDetails} from "../catalog-state";

describe("action:openSharingModal", () => {

    it("should close the sharing modal by clearing the shareModel property if album not found", () => {
        const action = openSharingModalAction({
            albumId: {owner: "notfound", folderName: "notfound"},
        });

        const expected: SharingDialogFrag = {
            open: false,
            sharedWith: [],
            suggestions: [],
        };
        expect(sharingDialogSelector(reduceOpenSharingModal(loadedStateWithTwoAlbums, action))).toEqual(expected);
    });

    it("should suggest all known users not already shared with the album, sorted by popularity then name", () => {
        // Add a third album with two users, one of which is herselfUser (already has access to album 0)
        const userA: UserDetails = {email: "a@a.com", name: "Alice"};
        const userB: UserDetails = {email: "b@b.com", name: "Bob"};
        const userC: UserDetails = {email: "c@c.com", name: "Charlie"};
        const albums = [
            {...twoAlbums[0]}, // sharedWith: [herselfUser]
            {...twoAlbums[1], sharedWith: [{user: userA}, {user: userB}], ownedBy: undefined},
            {
                albumId: {owner: twoAlbums[0].albumId.owner, folderName: "mar-25"},
                name: "March 2025",
                start: new Date(2025, 2, 1),
                end: new Date(2025, 2, 31),
                totalCount: 0,
                temperature: 0,
                relativeTemperature: 0,
                sharedWith: [{user: userA}, {user: userB}, {user: userC}],
            }
        ];
        // userA: 2 albums, userB: 2 albums, userC: 1 album, herselfUser: 1 album
        const state = {...loadedStateWithTwoAlbums, allAlbums: albums, albums};
        const action = openSharingModalAction(albums[0].albumId);

        const expected: SharingDialogFrag = {
            open: true,
            sharedWith: [
                {
                    user: herselfUser,
                }
            ],
            suggestions: [
                userA, // 2 albums
                userB, // 2 albums
                userC, // 1 album
            ],
        };
        expect(sharingDialogSelector(reduceOpenSharingModal(state, action))).toEqual(expected);
    });

    it("should open the dialog with an empty suggestion list if all known users already have access", () => {
        // All albums share with the same user
        const userA: UserDetails = {email: "a@a.com", name: "Alice"};
        const albums = [
            {...twoAlbums[0], sharedWith: [{user: userA}]},
            {...twoAlbums[1], sharedWith: [{user: userA}]},
        ];
        const state = {...loadedStateWithTwoAlbums, allAlbums: albums, albums};
        const action = openSharingModalAction(albums[0].albumId);

        const expected: SharingDialogFrag = {
            open: true,
            sharedWith: [
                {user: userA}
            ],
            suggestions: [],
        };
        expect(sharingDialogSelector(reduceOpenSharingModal(state, action))).toEqual(expected);
    });

    it("should return empty suggestions if there are no known users", () => {
        // No album has any sharedWith users
        const albums = [
            {...twoAlbums[0], sharedWith: []},
            {...twoAlbums[1], sharedWith: []},
        ];
        const state = {...loadedStateWithTwoAlbums, allAlbums: albums, albums};
        const action = openSharingModalAction(albums[0].albumId);

        const expected: SharingDialogFrag = {
            open: true,
            sharedWith: [],
            suggestions: [],
        };
        expect(sharingDialogSelector(reduceOpenSharingModal(state, action))).toEqual(expected);
    });
});
