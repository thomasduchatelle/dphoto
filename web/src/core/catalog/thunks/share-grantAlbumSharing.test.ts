import {GrantAlbumSharingAPI, grantAlbumSharingThunk} from "./share-grantAlbumSharing";
import {AlbumId, catalogActions, UserDetails} from "../domain";

class FakeSharingAPI implements GrantAlbumSharingAPI {
    public grantRequests: { albumId: AlbumId, email: string }[] = [];
    public grantError: any = null;

    public userDetails: { [email: string]: UserDetails } = {};
    public userDetailsError: any = null;

    grantAccessToAlbum(albumId: AlbumId, email: string): Promise<void> {
        this.grantRequests.push({albumId, email});
        if (this.grantError) {
            return Promise.reject(this.grantError);
        }
        return Promise.resolve();
    }

    loadUserDetails(email: string): Promise<UserDetails> {
        if (this.userDetailsError) {
            return Promise.reject(this.userDetailsError);
        }
        if (this.userDetails[email]) {
            return Promise.resolve(this.userDetails[email]);
        }
        return Promise.reject(new Error("User not found"));
    }
}

describe("grantAccessThunk", () => {
    const albumId: AlbumId = {owner: "owner", folderName: "album"};
    const email = "user@example.com";
    const userDetails: UserDetails = {email, name: "User Name", picture: "pic.jpg"};

    it("should call the sharingAPI.grantAccessToAlbum and dispatch a AddSharingAction with values from sharingAPI.loadUserDetails", async () => {
        const fakeAPI = new FakeSharingAPI();
        fakeAPI.userDetails[email] = userDetails;
        const dispatched: any[] = [];
        await grantAlbumSharingThunk(dispatched.push.bind(dispatched), fakeAPI, albumId, email);

        expect(fakeAPI.grantRequests).toEqual([{albumId, email}]);
        expect(dispatched).toEqual([
            catalogActions.addSharingAction({
                sharing: {
                    user: userDetails,
                }
            })
        ]);
    });

    it("should use the email as name if the call to sharingAPI.loadUserDetails failed and dispatch the AddSharingAction", async () => {
        const fakeAPI = new FakeSharingAPI();
        fakeAPI.userDetailsError = new Error("[TEST] fail user details");
        const dispatched: any[] = [];
        await grantAlbumSharingThunk(dispatched.push.bind(dispatched), fakeAPI, albumId, email);

        expect(dispatched).toEqual([
            catalogActions.addSharingAction({
                sharing: {
                    user: {
                        email,
                        name: email,
                    },
                }
            })
        ]);
    });

    it("should dispatch SharingModalErrorAction if the sharingAPI.grantAccessToAlbum failed", async () => {
        const fakeAPI = new FakeSharingAPI();
        fakeAPI.grantError = new Error("[TEST] fail grant");
        const dispatched: any[] = [];
        await grantAlbumSharingThunk(dispatched.push.bind(dispatched), fakeAPI, albumId, email);

        expect(fakeAPI.grantRequests).toEqual([{albumId, email}]);
        expect(dispatched).toEqual([
            catalogActions.sharingModalErrorAction({
                error: {
                    type: "adding",
                    message: "Failed to grant access, verify the email address or contact maintainers"
                }
            })
        ]);
    });
});
