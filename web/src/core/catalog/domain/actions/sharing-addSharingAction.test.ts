import {addSharingAction, reduceAddSharing} from "./sharing-addSharingAction";
import {UserDetails} from "../catalog-state";
import {herselfUser, loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";
import {SharingDialogFrag, sharingDialogSelector} from "./selector-sharingDialogSelector";


describe("reduceAddSharing", () => {
    it("should add a new sharing entry and keep the list of suggestions if user was not suggested", () => {
        const initial = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [
                    {
                        user: herselfUser,
                    }
                ],
                suggestions: [
                    {email: "alice@example.com", name: "Alice", picture: "alice-face.jpg"},
                ],
            }
        };
        const newUser: UserDetails = {email: "bob@example.com", name: "Bob", picture: "bob-face.jpg"};

        const state = reduceAddSharing(initial, addSharingAction({
            user: newUser,
        }));
        expect(sharingDialogSelector(state)).toEqual({
            open: true,
            sharedWith: [
                {user: newUser},
                {user: herselfUser}
            ],
            suggestions: [
                {email: "alice@example.com", name: "Alice", picture: "alice-face.jpg"},
            ],
        });
    });

    it("should replace an existing sharing entry for the same user when receiving AddSharingAction", () => {
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
        // Add the same user again: user is overridden and not added
        const action = addSharingAction({
            user: herselfUser,
        });
        const expected: SharingDialogFrag = {
            open: true,
            sharedWith: [
                {
                    user: herselfUser,
                }
            ],
            suggestions: [],
        };
        expect(sharingDialogSelector(reduceAddSharing(initial, action))).toEqual(expected);
    });

    it("should not change state when AddSharingAction is received and shareModal is closed", () => {
        const action = addSharingAction(twoAlbums[0].sharedWith[0]);
        const result = reduceAddSharing(loadedStateWithTwoAlbums, action);
        expect(sharingDialogSelector(result)).toEqual({
            open: false,
            sharedWith: [],
            suggestions: [],
        });
    });

    it("should remove the newly granted email from the suggestions list if it was present", () => {
        const newUser: UserDetails = {email: "bob@example.com", name: "Bob", picture: "bob-face.jpg"};
        const initial = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: twoAlbums[0].sharedWith,
                suggestions: [
                    newUser,
                    {email: "alice@example.com", name: "Alice", picture: "alice-face.jpg"},
                ],
            }
        };

        const state = reduceAddSharing(initial, addSharingAction({
            user: newUser,
        }));
        expect(sharingDialogSelector(state)).toEqual({
            open: true,
            sharedWith: [
                {user: newUser},
                {user: herselfUser},
            ],
            suggestions: [
                {email: "alice@example.com", name: "Alice", picture: "alice-face.jpg"},
            ],
        });
    });
});
