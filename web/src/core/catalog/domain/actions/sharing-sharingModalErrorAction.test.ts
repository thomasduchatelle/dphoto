import {reduceSharingModalError, sharingModalErrorAction} from "./sharing-sharingModalErrorAction";
import {herselfUser, loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";
import {ShareError} from "../catalog-state";
import {SharingDialogFrag, sharingDialogSelector} from "./selector-sharingDialogSelector";

describe("reduceSharingModalError", () => {
    it("should set the error field when receiving SharingModalErrorAction", () => {
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
        const error: ShareError = {type: "adding", message: "Failed to add user"};
        const action = sharingModalErrorAction(error);
        const expected: SharingDialogFrag = {
            open: true,
            sharedWith: [
                {
                    user: herselfUser,
                }
            ],
            error,
        };
        expect(sharingDialogSelector(reduceSharingModalError(initial, action))).toEqual(expected);
    });
});
