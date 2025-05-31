import {addSharingAction, reduceAddSharing} from "./sharing-addSharingAction";
import {CatalogViewerState, SharingType, UserDetails} from "../catalog-state";
import {herselfUser, loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";

describe("reduceAddSharing", () => {
    it("should add a new sharing entry and keep the modal open when receiving AddSharingAction", () => {
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
        const newUser: UserDetails = {email: "bob@example.com", name: "Bob", picture: "bob-face.jpg"};
        const action = addSharingAction({
            user: newUser,
            role: SharingType.contributor,
        });
        const expected = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [
                    {
                        user: newUser,
                        role: SharingType.contributor,
                    },
                    {
                        user: herselfUser,
                        role: SharingType.visitor,
                    }
                ]
            }
        };
        expect(reduceAddSharing(initial, action)).toEqual(expected);
    });

    it("should replace an existing sharing entry for the same user when receiving AddSharingAction", () => {
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
        // Add the same user with a different role: user is overridden and not added
        const action = addSharingAction({
            user: herselfUser,
            role: SharingType.contributor,
        });
        const expected: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [
                    {
                        user: herselfUser,
                        role: SharingType.contributor,
                    }
                ],
            }
        };
        expect(reduceAddSharing(initial, action)).toEqual(expected);
    });

    it("should not change state when AddSharingAction is received and shareModal is closed", () => {
        const action = addSharingAction(twoAlbums[0].sharedWith[0]);
        expect(reduceAddSharing(loadedStateWithTwoAlbums, action)).toEqual(loadedStateWithTwoAlbums);
    });
});
