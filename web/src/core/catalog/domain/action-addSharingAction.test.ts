import {addSharingAction, reduceAddSharing} from "./action-addSharingAction";
import {SharingType, UserDetails} from "./catalog-state";
import {herselfUser, loadedStateWithTwoAlbums, twoAlbums} from "./tests/test-helper-state";

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
            sharing: {
                user: newUser,
                role: SharingType.contributor,
            }
        });
        const expected = {
            ...loadedStateWithTwoAlbums,
            shareModal: {
                sharedAlbumId: twoAlbums[0].albumId,
                sharedWith: [
                    {
                        user: herselfUser,
                        role: SharingType.visitor,
                    },
                    {
                        user: newUser,
                        role: SharingType.contributor,
                    }
                ].sort((a, b) => {
                    const nameA = a.user.name?.trim() || "";
                    const nameB = b.user.name?.trim() || "";
                    if (nameA && nameB) {
                        const cmp = nameA.localeCompare(nameB);
                        if (cmp !== 0) return cmp;
                        return a.user.email.localeCompare(b.user.email);
                    }
                    if (!nameA && !nameB) {
                        return a.user.email.localeCompare(b.user.email);
                    }
                    if (!nameA) return 1;
                    if (!nameB) return -1;
                    return 0;
                }),
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
            sharing: {
                user: herselfUser,
                role: SharingType.contributor,
            }
        });
        const expected = {
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
        const initial = {
            ...loadedStateWithTwoAlbums,
            shareModal: undefined,
        };
        const action = addSharingAction({
            sharing: {
                user: herselfUser,
                role: SharingType.visitor,
            }
        });
        expect(reduceAddSharing(initial, action)).toEqual(initial);
    });
});
