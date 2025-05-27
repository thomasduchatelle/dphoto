import {revokeAlbumSharingAPI, revokeAlbumSharingThunk} from "./share-revokeAlbumSharing";
import {AlbumId, catalogActions} from "../domain";

class FakeSharingAPI implements revokeAlbumSharingAPI {
    public revokeRequests: { albumId: AlbumId, email: string }[] = [];
    public revokeError: any = null;

    revokeSharingAlbum(albumId: AlbumId, email: string): Promise<void> {
        this.revokeRequests.push({albumId, email});
        if (this.revokeError) {
            return Promise.reject(this.revokeError);
        }
        return Promise.resolve();
    }
}

describe("revokeAccessThunk", () => {
    const albumId: AlbumId = {owner: "owner", folderName: "album"};
    const email = "user@example.com";

    it("should call the sharingAPI.revokeSharingAlbum with appropriate parameters and dispatch a RemoveSharingAction action", async () => {
        const fakeAPI = new FakeSharingAPI();
        const dispatched: any[] = [];
        await revokeAlbumSharingThunk(dispatched.push.bind(dispatched), fakeAPI, albumId, email);

        expect(fakeAPI.revokeRequests).toEqual([{albumId, email}]);
        expect(dispatched).toEqual([catalogActions.removeSharingAction(email)]);
    });

    it("should dispatch a SharingModalErrorAction if the call to sharingAPI.revokeSharingAlbum raise an error", async () => {
        const fakeAPI = new FakeSharingAPI();
        fakeAPI.revokeError = new Error("[TEST] fail!");
        const dispatched: any[] = [];
        await revokeAlbumSharingThunk(dispatched.push.bind(dispatched), fakeAPI, albumId, email);

        expect(dispatched).toEqual([
            catalogActions.sharingModalErrorAction({
                error: {type: "general", message: `Couldn't revoke access of user ${email}, try again later`}
            })
        ]);
    });
});
