import {sharingModalClosed} from "./action-sharingModalClosed";
import {herselfUser, loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";
import {sharingDialogSelector} from "./selector-sharingDialogSelector";

describe("action:sharingModalClosed", () => {
    it("should close the sharing modal by clearing the shareModel property", () => {
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
        const action = sharingModalClosed();
        expect(sharingDialogSelector(action.reducer(initial, action))).toEqual({
            open: false,
            sharedWith: [],
            suggestions: [],
        });
    });
});
