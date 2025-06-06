import {reduceRemoveSharing, removeSharingAction} from "./sharing-removeSharingAction";
import {herselfUser, loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";
import {sharingDialogSelector} from "./selector-sharingDialogSelector";
import {catalogActions} from "./index";

describe("removeSharingAction", () => {
    it("removes a sharing entry by email, and adds it in the suggestions, while keeping the modal open", () => {
        const initial = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: twoAlbums[0].sharedWith,
                suggestions: [],
            }
        };
        const state = reduceRemoveSharing(initial, removeSharingAction(herselfUser.email));
        expect(sharingDialogSelector(state)).toEqual({
            open: true,
            sharedWith: [],
            suggestions: [herselfUser],
        });
    });

    it("keeps consistent the shares in the visible list of albums", () => {
        const initial = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: twoAlbums[0].sharedWith,
                suggestions: [],
            }
        };
        const state = reduceRemoveSharing(initial, removeSharingAction(herselfUser.email));
        expect(state.albums).toEqual([
            {
                ...twoAlbums[0],
                sharedWith: [],
            },
            twoAlbums[1]
        ]);
    });

    it("should not change state when RemoveSharingAction is received and shareModal is undefined", () => {
        const result = reduceRemoveSharing(loadedStateWithTwoAlbums, removeSharingAction(herselfUser.email));
        expect(sharingDialogSelector(result)).toEqual({
            open: false,
            sharedWith: [],
            suggestions: [],
        });
    });

    it("should not change state when RemoveSharingAction is received with an email not in sharedWith", () => {
        const initial = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: twoAlbums[0].sharedWith,
                suggestions: [],
            }
        };
        const state = reduceRemoveSharing(initial, removeSharingAction("notfound@example.com"));
        expect(sharingDialogSelector(state)).toEqual({
            open: true,
            sharedWith: twoAlbums[0].sharedWith,
            suggestions: [],
        });
    });

    it("keeps the suggestions already present when revoking another email, and respect the suggestion order", () => {
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
                    sharedWith: [{user: aliceUserDetails}],
                },
                {
                    ...twoAlbums[1],
                    sharedWith: [{user: charlieUserDetails}],
                    ownedBy: undefined,
                }
            ],
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [{user: aliceUserDetails}],
                suggestions: [charlieUserDetails, bobUserDetails]
            }
        };

        const state = reduceRemoveSharing(initial, catalogActions.removeSharingAction(aliceUserDetails.email));

        expect(sharingDialogSelector(state)).toEqual({
            open: true,
            sharedWith: [],
            suggestions: [charlieUserDetails, aliceUserDetails, bobUserDetails],
            error: undefined,
        });
    });
});
