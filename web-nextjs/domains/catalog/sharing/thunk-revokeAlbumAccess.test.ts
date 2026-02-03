import {RevokeAlbumAccessAPI, revokeAlbumAccessThunk} from "./thunk-revokeAlbumAccess";
import {AlbumId} from "../language";
import {albumAccessRevoked} from "./action-albumAccessRevoked";
import {sharingModalErrorOccurred} from "./action-sharingModalErrorOccurred";

class FakeSharingAPI implements RevokeAlbumAccessAPI {
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
        await revokeAlbumAccessThunk(dispatched.push.bind(dispatched), fakeAPI, albumId, email);

        expect(fakeAPI.revokeRequests).toEqual([{albumId, email}]);
        expect(dispatched).toEqual([albumAccessRevoked(email)]);
    });

    it("should dispatch a SharingModalErrorAction if the call to sharingAPI.revokeSharingAlbum raise an error", async () => {
        const fakeAPI = new FakeSharingAPI();
        fakeAPI.revokeError = new Error("[TEST] fail!");
        const dispatched: any[] = [];
        await revokeAlbumAccessThunk(dispatched.push.bind(dispatched), fakeAPI, albumId, email);

        expect(dispatched).toEqual([
            albumAccessRevoked(email),
            sharingModalErrorOccurred({
                type: "revoke", email, message: `Couldn't revoke access of user ${email}, try again later`
            })
        ]);
    });
});
