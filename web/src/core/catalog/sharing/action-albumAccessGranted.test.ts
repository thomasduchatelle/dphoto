import {albumAccessGranted, reduceAlbumAccessGranted} from "./action-albumAccessGranted";
import {UserDetails} from "../language";
import {herselfUser, loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";
import {SharingDialogFrag, sharingDialogSelector} from "./selector-sharingDialogSelector";


describe("action:albumAccessGranted", () => {
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

        const state = reduceAlbumAccessGranted(initial, albumAccessGranted({
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

    it("should replace an existing sharing entry for the same user when receiving AlbumAccessGranted", () => {
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
        const action = albumAccessGranted({
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
        expect(sharingDialogSelector(reduceAlbumAccessGranted(initial, action))).toEqual(expected);
    });

    it("should not change state when AlbumAccessGranted is received and shareModal is closed", () => {
        const action = albumAccessGranted(twoAlbums[0].sharedWith[0]);
        const result = reduceAlbumAccessGranted(loadedStateWithTwoAlbums, action);
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

        const state = reduceAlbumAccessGranted(initial, albumAccessGranted({
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
