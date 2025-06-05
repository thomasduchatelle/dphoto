import {addSharingAction, reduceAddSharing} from "./sharing-addSharingAction";
import {UserDetails} from "../catalog-state";
import {herselfUser, loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";
import {SharingDialogFrag, sharingDialogSelector} from "./selector-sharingDialogSelector";

describe("reduceAddSharing", () => {
    it("should add a new sharing entry and keep the modal open when receiving AddSharingAction", () => {
        const initial = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [
                    {
                        user: herselfUser,
                    }
                ],
            }
        };
        const newUser: UserDetails = {email: "bob@example.com", name: "Bob", picture: "bob-face.jpg"};
        const action = addSharingAction({
            user: newUser,
        });

        const expected: SharingDialogFrag = {
            open: true,
            sharedWith: [
                {user: newUser},
                {user: herselfUser}
            ],
        };
        expect(sharingDialogSelector(reduceAddSharing(initial, action))).toEqual(expected);
    });

    it("should update the album's sharedWith list in the state albums (allAlbums and [visible]albums) when a sharing is granted", () => {
        const initial = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [
                    {
                        user: herselfUser,
                    }
                ],
            }
        };
        const newUser: UserDetails = {email: "bob@example.com", name: "Bob", picture: "bob-face.jpg"};
        const action = addSharingAction({
            user: newUser,
        });
        const expected: SharingDialogFrag = {
            open: true,
            sharedWith: [
                {user: newUser},
                {user: herselfUser}
            ],
        };
        expect(sharingDialogSelector(reduceAddSharing(initial, action))).toEqual(expected);
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
        };
        expect(sharingDialogSelector(reduceAddSharing(initial, action))).toEqual(expected);
    });

    it("should not change state when AddSharingAction is received and shareModal is closed", () => {
        const action = addSharingAction(twoAlbums[0].sharedWith[0]);
        const result = reduceAddSharing(loadedStateWithTwoAlbums, action);
        expect(sharingDialogSelector(result)).toEqual({
            open: false,
            sharedWith: [],
        });
    });
});
