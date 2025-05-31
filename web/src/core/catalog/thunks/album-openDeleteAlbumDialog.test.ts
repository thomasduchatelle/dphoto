import {openDeleteAlbumDialogThunk} from "./album-openDeleteAlbumDialog";
import {catalogActions} from "../domain";

describe("openDeleteAlbumDialogThunk", () => {
    it("should dispatch openDeleteAlbumDialogAction", async () => {
        const dispatched: any[] = [];
        openDeleteAlbumDialogThunk(dispatched.push.bind(dispatched));
        expect(dispatched).toEqual([
            catalogActions.openDeleteAlbumDialogAction()
        ]);
    });
});
