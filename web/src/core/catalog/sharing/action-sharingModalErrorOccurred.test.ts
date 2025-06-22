import {sharingModalErrorOccurred} from "./action-sharingModalErrorOccurred";
import {herselfUser, loadedStateWithTwoAlbums, march2025, twoAlbums} from "../tests/test-helper-state";
import {ShareError} from "../language";
import {SharingDialogFrag, sharingDialogSelector} from "./selector-sharingDialogSelector";

describe("action:sharingModalErrorOccurred:grant", () => {
    it("ignores the error if the modal is closed", () => {
        const expected: SharingDialogFrag = {
            open: false,
            sharedWith: [],
            suggestions: [],
        };
        const action = sharingModalErrorOccurred({
            type: "grant",
            message: "Failed to add user",
            email: "foo@example.com"
        });
        expect(sharingDialogSelector(action.reducer(
            loadedStateWithTwoAlbums,
            action
        ))).toEqual(expected);
    });

    it("should set the error field when receiving SharingModalErrorOccurred even if the user not found in sharedWith", () => {
        const initial = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [
                    {
                        user: herselfUser,
                    }
                ],
                suggestions: [],
            }
        };
        const error: ShareError = {type: "grant", message: "Failed to add user", email: "foo@example.com"};
        const expected: SharingDialogFrag = {
            open: true,
            sharedWith: [
                {
                    user: herselfUser,
                }
            ],
            suggestions: [],
            error,
        };
        const action = sharingModalErrorOccurred(error);
        expect(sharingDialogSelector(action.reducer(initial, action))).toEqual(expected);
    });

    it("removes the user from sharedWith in the dialog and visible albums, and adds it to suggestions", () => {
        const herselfUserDetails = twoAlbums[0].sharedWith[0].user;
        const initial = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: twoAlbums[0].sharedWith,
                suggestions: []
            }
        };
        const error: ShareError = {type: "grant", message: "Failed to add user", email: herselfUserDetails.email};
        const action = sharingModalErrorOccurred(error);

        const state = action.reducer(initial, action)
        expect(sharingDialogSelector(state)).toEqual({
            open: true,
            sharedWith: [], // restored to before grant
            suggestions: [herselfUserDetails],
            error,
        } as SharingDialogFrag);

        expect(state.albums).toEqual([
            {
                ...twoAlbums[0],
                sharedWith: [], // restored to before grant
            },
            {
                ...twoAlbums[1],
                sharedWith: twoAlbums[1].sharedWith // no change
            }
        ]);
    });

    it("removes the user from sharedWith and adds it to suggestions while respecting the order (popularity desc then name asc)", () => {
        const aliceUserDetails = {
            name: "Alice",
            email: "alice@example.com",
            picture: "pic"
        };
        const bobUserDetails = {
            name: "Bob",
            email: "bob@example.com"
        };
        const charlieUserDetails = {
            name: "Charlie",
            email: "charlie@example.com"
        };

        // alice: 1 album, bob: 2 albums, charlie: 2 albums
        const albums = [
            {
                ...twoAlbums[0],
                sharedWith: [{user: aliceUserDetails}]
            },
            {
                ...twoAlbums[1],
                sharedWith: [{user: aliceUserDetails}, {user: bobUserDetails}, {user: charlieUserDetails}]
            },
            {
                ...march2025,
                sharedWith: [{user: bobUserDetails}]
            }
        ];
        const initial = {
            ...loadedStateWithTwoAlbums,
            allAlbums: albums,
            albums,
            shareModal: {
                sharedAlbumId: albums[0].albumId,
                sharedWith: [{user: aliceUserDetails}],
                suggestions: []
            }
        };
        const error: ShareError = {type: "grant", message: "Failed to add user", email: "alice@example.com"};
        const action = sharingModalErrorOccurred(error);

        const state = action.reducer(initial, action);
        expect(sharingDialogSelector(state)).toEqual({
            open: true,
            sharedWith: [],
            suggestions: [bobUserDetails, aliceUserDetails, charlieUserDetails],
            error,
        });
    });
});

describe("action:sharingModalErrorOccurred:revoke", () => {
    it("ignores the error if the modal is closed", () => {
        const action = sharingModalErrorOccurred({
            type: "revoke",
            message: "Failed to revoke user",
            email: "foo@example.com"
        });
        const state = action.reducer(
            loadedStateWithTwoAlbums,
            action
        );
        expect(sharingDialogSelector(state)).toEqual({
            open: false,
            sharedWith: [],
            suggestions: [],
        });
    });

    it("should set the error field and add the user in sharedWith even if the user hasn't been found in suggestions", () => {
        const initial = {
            ...loadedStateWithTwoAlbums,
            allAlbums: [
                {
                    ...twoAlbums[0],
                    sharedWith: []
                },
                {
                    ...twoAlbums[1],
                    sharedWith: [{user: herselfUser}]
                }
            ],
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [],
                suggestions: [herselfUser],
            }
        };
        const error: ShareError = {type: "revoke", message: "Failed to revoke user", email: "foo@example.com"};
        const action = sharingModalErrorOccurred(error);
        const state = action.reducer(initial, action);
        expect(sharingDialogSelector(state)).toEqual({
            open: true,
            sharedWith: [
                {
                    user: {
                        name: "foo@example.com",
                        email: "foo@example.com"
                    }
                }
            ],
            suggestions: [herselfUser],
            error,
        });
    });

    it("should set the error field, restore the user in sharedWith (both in share dialog and visible albums), and remove it from suggestion", () => {
        const aliceUserDetails = {
            name: "Alice",
            email: "alice@example.com",
            picture: "pic"
        };
        const bobUserDetails = {
            name: "Bob",
            email: "bob@example.com"
        };
        const charlieUserDetails = {
            name: "Charlie",
            email: "charlie@example.com"
        };

        const initial = {
            ...loadedStateWithTwoAlbums,
            allAlbums: [
                {
                    ...twoAlbums[0],
                    sharedWith: [{user: bobUserDetails}]
                },
                {
                    ...twoAlbums[1],
                    sharedWith: [{user: charlieUserDetails}]
                }
            ],
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [{user: bobUserDetails}],
                suggestions: [aliceUserDetails, charlieUserDetails]
            }
        };
        const error: ShareError = {type: "revoke", message: "Failed to revoke user", email: "alice@example.com"};
        const action = sharingModalErrorOccurred(error);
        const state = action.reducer(initial, action);

        expect(sharingDialogSelector(state)).toEqual({
            open: true,
            sharedWith: [
                {user: aliceUserDetails},
                {user: bobUserDetails}
            ],
            suggestions: [charlieUserDetails],
            error,
        });
        expect(state.albums).toEqual([
            {
                ...twoAlbums[0],
                sharedWith: [
                    {user: aliceUserDetails},
                    {user: bobUserDetails}
                ]
            },
            {
                ...twoAlbums[1],
                sharedWith: [{user: charlieUserDetails}]
            }
        ]);
    });
});
