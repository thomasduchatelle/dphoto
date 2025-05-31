import {closeDeleteAlbumDialogThunk} from "./album-closeDeleteAlbumDialog";
import {catalogActions} from "../domain";

describe("closeDeleteAlbumDialogThunk", () => {
    it("should dispatch closeDeleteAlbumDialogAction", () => {
        const dispatched: any[] = [];
        closeDeleteAlbumDialogThunk(dispatched.push.bind(dispatched));
        expect(dispatched).toEqual([
            catalogActions.closeDeleteAlbumDialogAction()
        ]);
    });
});
