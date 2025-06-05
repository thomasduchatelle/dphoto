import {reduceRemoveSharing, removeSharingAction} from "./sharing-removeSharingAction";
import {herselfUser, loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";
import {SharingDialogFrag, sharingDialogSelector} from "./selector-sharingDialogSelector";

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
                    },
                    {
                        user: {email: bobEmail, name: "Bob", picture: "bob-face.jpg"},
                    }
                ],
            }
        };
        const action = removeSharingAction(bobEmail);
        const expected: SharingDialogFrag = {
            open: true,
            sharedWith: [
                {
                    user: herselfUser,
                }
            ]
        };
        expect(sharingDialogSelector(reduceRemoveSharing(initial, action))).toEqual(expected);
    });

    it("should update the album's sharedWith list in the state when a sharing is revoked", () => {
        const bobEmail = "bob@example.com";
        const initial = {
            ...loadedStateWithTwoAlbums,
            albums: [
                {
                    ...twoAlbums[0],
                    sharedWith: [
                        {
                            user: herselfUser,
                        },
                        {
                            user: {email: bobEmail, name: "Bob", picture: "bob-face.jpg"},
                        }
                    ]
                },
                twoAlbums[1]
            ],
            allAlbums: [
                {
                    ...twoAlbums[0],
                    sharedWith: [
                        {
                            user: herselfUser,
                        },
                        {
                            user: {email: bobEmail, name: "Bob", picture: "bob-face.jpg"},
                        }
                    ]
                },
                twoAlbums[1]
            ],
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [
                    {
                        user: herselfUser,
                    },
                    {
                        user: {email: bobEmail, name: "Bob", picture: "bob-face.jpg"},
                    }
                ],
            }
        };
        const action = removeSharingAction(bobEmail);
        const expected: SharingDialogFrag = {
            open: true,
            sharedWith: [
                {
                    user: herselfUser,
                }
            ]
        };
        expect(sharingDialogSelector(reduceRemoveSharing(initial, action))).toEqual(expected);
    });

    it("should not change state when RemoveSharingAction is received and shareModal is undefined", () => {
        const action = removeSharingAction(herselfUser.email);
        const result = reduceRemoveSharing(loadedStateWithTwoAlbums, action);
        expect(sharingDialogSelector(result)).toEqual({
            open: false,
            sharedWith: [],
        });
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
                    }
                ],
            }
        };
        expect(sharingDialogSelector(reduceRemoveSharing(initial, action))).toEqual({
            open: true,
            sharedWith: [
                {
                    user: herselfUser,
                }
            ]
        });
    });
});
