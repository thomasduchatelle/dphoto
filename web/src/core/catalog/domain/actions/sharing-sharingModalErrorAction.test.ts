import {reduceSharingModalError, sharingModalErrorAction} from "./sharing-sharingModalErrorAction";
import {herselfUser, loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";
import {ShareError, SharingType} from "../catalog-state";

describe("reduceSharingModalError", () => {
    it("should set the error field when receiving SharingModalErrorAction", () => {
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
        const error: ShareError = {type: "adding", message: "Failed to add user"};
        const action = sharingModalErrorAction(error);
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
                error,
            }
        };
        expect(reduceSharingModalError(initial, action)).toEqual(expected);
    });
});
