import {ShareController, SharingAPI} from "./ShareController";
import {AlbumId, CatalogViewerAction, SharingType, UserDetails} from "../index";

class FakeSharingAPI implements SharingAPI {
    public revokeRequests: { albumId: AlbumId, email: string }[] = [];
    public revokeError: any = null;

    public grantRequests: { albumId: AlbumId, email: string, role: SharingType }[] = [];
    public grantError: any = null;

    public userDetails: { [email: string]: UserDetails } = {};
    public userDetailsError: any = null;

    revokeSharingAlbum(albumId: AlbumId, email: string): Promise<void> {
        this.revokeRequests.push({albumId, email});
        if (this.revokeError) {
            return Promise.reject(this.revokeError);
        }
        return Promise.resolve();
    }

    grantAccessToAlbum(albumId: AlbumId, email: string, role: SharingType): Promise<void> {
        this.grantRequests.push({albumId, email, role});
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

describe("ShareController.revokeAccess", () => {
    const albumId: AlbumId = {owner: "owner", folderName: "album"};
    const email = "user@example.com";

    it("should call the sharingAPI.revokeSharingAlbum with appropriate parameters and dispatch a RemoveSharingAction action", async () => {
        const fakeAPI = new FakeSharingAPI();
        const dispatched: CatalogViewerAction[] = [];
        const controller = new ShareController(action => dispatched.push(action), fakeAPI);

        await controller.revokeAccess(albumId, email);

        expect(fakeAPI.revokeRequests).toEqual([{albumId, email}]);
        expect(dispatched).toEqual([{type: "RemoveSharingAction", email}]);
    });

    it("should dispatch a SharingModalErrorAction if the call to sharingAPI.revokeSharingAlbum raise an error", async () => {
        const fakeAPI = new FakeSharingAPI();
        fakeAPI.revokeError = new Error("[TEST] fail!");
        const dispatched: CatalogViewerAction[] = [];
        const controller = new ShareController(action => dispatched.push(action), fakeAPI);

        await controller.revokeAccess(albumId, email);

        expect(dispatched).toEqual([
            {
                type: "SharingModalErrorAction",
                error: {
                    type: "general",
                    message: `Couldn't revoke access of user ${email}, try again later`
                }
            }
        ]);
    });
});

describe("ShareController.grantAccess", () => {
    const albumId: AlbumId = {owner: "owner", folderName: "album"};
    const email = "user@example.com";
    const role = SharingType.contributor;
    const userDetails: UserDetails = {email, name: "User Name", picture: "pic.jpg"};

    it("should call the sharingAPI.grantAccessToAlbum and dispatch a AddSharingAction with values from sharingAPI.loadUserDetails", async () => {
        const fakeAPI = new FakeSharingAPI();
        fakeAPI.userDetails[email] = userDetails;
        const dispatched: CatalogViewerAction[] = [];
        const controller = new ShareController(action => dispatched.push(action), fakeAPI);

        await controller.grantAccess(albumId, email, role);

        expect(fakeAPI.grantRequests).toEqual([{albumId, email, role}]);
        expect(dispatched).toEqual([
            {
                type: "AddSharingAction",
                sharing: {
                    user: userDetails,
                    role,
                }
            }
        ]);
    });

    it("should use the email as name if the call to sharingAPI.loadUserDetails failed and dispatch the AddSharingAction", async () => {
        const fakeAPI = new FakeSharingAPI();
        fakeAPI.userDetailsError = new Error("[TEST] fail user details");
        const dispatched: CatalogViewerAction[] = [];
        const controller = new ShareController(action => dispatched.push(action), fakeAPI);

        await controller.grantAccess(albumId, email, role);

        expect(dispatched).toEqual([
            {
                type: "AddSharingAction",
                sharing: {
                    user: {
                        email,
                        name: email,
                    },
                    role,
                }
            }
        ]);
    });

    it("should dispatch SharingModalErrorAction if the sharingAPI.grantAccessToAlbum failed", async () => {
        const fakeAPI = new FakeSharingAPI();
        fakeAPI.grantError = new Error("[TEST] fail grant");
        const dispatched: CatalogViewerAction[] = [];
        const controller = new ShareController(action => dispatched.push(action), fakeAPI);

        await controller.grantAccess(albumId, email, role);

        expect(fakeAPI.grantRequests).toEqual([{albumId, email, role}]);
        expect(dispatched).toEqual([
            {
                type: "SharingModalErrorAction",
                error: {
                    type: "adding",
                    message: "Failed to grant access, verify the email address or contact maintainers"
                }
            }
        ]);
    });
});
