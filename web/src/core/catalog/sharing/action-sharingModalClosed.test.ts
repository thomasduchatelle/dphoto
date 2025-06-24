import {sharingModalClosed} from "./action-sharingModalClosed";
import {herselfUser, loadedStateWithTwoAlbums, twoAlbums} from "../tests/test-helper-state";
import {sharingDialogSelector} from "./selector-sharingDialogSelector";
import {CatalogViewerState} from "../language";

describe("action:sharingModalClosed", () => {
    it("should close the sharing modal by clearing the shareModel property", () => {
        const initial: CatalogViewerState = {
            ...loadedStateWithTwoAlbums,
            dialog: {
                type: "ShareDialog",
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
