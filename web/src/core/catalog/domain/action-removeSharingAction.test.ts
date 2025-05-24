import {removeSharingAction, reduceRemoveSharing} from "./action-removeSharingAction";
import {herselfUser, loadedStateWithTwoAlbums, twoAlbums} from "./tests/test-helper-state";
import {SharingType, UserDetails} from "./catalog-state";

describe("reduceRemoveSharing", () => {
    it("should remove a sharing entry by email and keep the modal open when receiving RemoveSharingAction", () => {
        const bobEmail = "bob@example.com";
        const initial = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [
                    {
                        user: herselfUser,
                        role: SharingType.visitor,
                    },
                    {
                        user: {email: bobEmail, name: "Bob", picture: "bob-face.jpg"},
                        role: SharingType.contributor,
                    }
                ],
            }
        };
        const action = removeSharingAction(bobEmail);
        const expected = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [
                    {
                        user: herselfUser,
                        role: SharingType.visitor,
                    }
                ],
            }
        };
        expect(reduceRemoveSharing(initial, action)).toEqual(expected);
    });

    it("should not change state when RemoveSharingAction is received and shareModal is undefined", () => {
        const action = removeSharingAction(herselfUser.email);
        const initial = {
            ...loadedStateWithTwoAlbums,
            shareModal: undefined,
        };
        expect(reduceRemoveSharing(initial, action)).toEqual(initial);
    });

    it("should not change state when RemoveSharingAction is received with an email not in sharedWith", () => {
        const action = removeSharingAction("notfound@example.com");
        const initial = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [
                    {
                        user: herselfUser,
                        role: SharingType.visitor,
                    }
                ],
            }
        };
        expect(reduceRemoveSharing(initial, action)).toEqual(initial);
    });
});
