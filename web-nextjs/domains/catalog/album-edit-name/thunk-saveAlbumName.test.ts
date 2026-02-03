import {SaveAlbumNamePort, saveAlbumNameThunk} from "./thunk-saveAlbumName";
import {albumRenamingStarted} from "./action-albumRenamingStarted";
import {albumRenamed} from "./action-albumRenamed";
import {albumRenamingFailed} from "./action-albumRenamingFailed";
import {AlbumId, CatalogError, CatalogViewerState} from "../language";
import {Action} from "@/libs/daction";

class FakeSaveAlbumNamePort implements SaveAlbumNamePort {
    public renameRequests: { albumId: AlbumId; newName: string; newFolderName?: string }[] = [];
    public shouldFail = false;
    public errorCode = "";
    public errorMessage = "";

    async renameAlbum(albumId: AlbumId, newName: string, newFolderName?: string): Promise<AlbumId> {
        this.renameRequests.push({albumId, newName, newFolderName});

        if (this.shouldFail) {
            if (this.errorCode) {
                throw new CatalogError(this.errorCode, this.errorMessage);
            } else {
                throw new Error(this.errorMessage);
            }
        }


        return {...albumId, folderName: `/${newFolderName ?? newName.toLowerCase().replace(/\s+/g, '-')}`};
    }
}

describe("thunk:saveAlbumNameThunk", () => {
    const albumId = {owner: "myself", folderName: "old-folder"};
    const newName = "New Album Name";
    const newAlbumId = {owner: "myself", folderName: "/new-album-name"};

    it("should dispatch albumRenamingStarted, call port with display name only, and dispatch albumRenamed on success", async () => {

        const fakePort = new FakeSaveAlbumNamePort();

        const dispatchedActions: Action<CatalogViewerState, any>[] = [];

        await saveAlbumNameThunk(
            dispatchedActions.push.bind(dispatchedActions),
            fakePort,
            {
                albumId,
                albumName: newName,
                customFolderName: "ignored",
                isCustomFolderNameEnabled: false
            }
        );

        expect(fakePort.renameRequests).toEqual([{
            albumId,
            newName,
        }]);

        expect(dispatchedActions).toEqual([
            albumRenamingStarted(),
            albumRenamed({previousAlbumId: albumId, newAlbumId, newName, redirectTo: newAlbumId}),
        ]);
    });

    it("should call port with both display name and folder name when folder name is enabled", async () => {
        const newFolderName = "new-folder-name";
        const newAlbumId = {owner: "myself", folderName: `/${newFolderName}`};

        const fakePort = new FakeSaveAlbumNamePort();

        const dispatchedActions: Action<CatalogViewerState, any>[] = [];

        await saveAlbumNameThunk(
            dispatchedActions.push.bind(dispatchedActions),
            fakePort,
            {
                albumId,
                albumName: newName,
                customFolderName: newFolderName,
                isCustomFolderNameEnabled: true
            }
        );

        expect(fakePort.renameRequests).toEqual([{
            albumId,
            newName,
            newFolderName,
        }]);

        expect(dispatchedActions).toEqual([
            albumRenamingStarted(),
            albumRenamed({previousAlbumId: albumId, newAlbumId, newName, redirectTo: newAlbumId}),
        ]);
    });

    it("should dispatch albumRenamingFailed with error code when API returns AlbumFolderNameAlreadyTakenErr", async () => {
        const albumId = {owner: "myself", folderName: "old-folder"};

        const fakePort = new FakeSaveAlbumNamePort();
        fakePort.shouldFail = true;
        fakePort.errorCode = "AlbumFolderNameAlreadyTakenErr";
        fakePort.errorMessage = "Folder name already taken";

        const dispatchedActions: Action<CatalogViewerState, any>[] = [];

        await saveAlbumNameThunk(
            dispatchedActions.push.bind(dispatchedActions),
            fakePort,
            {
                albumId,
                albumName: "New Name",
                customFolderName: "existing-folder",
                isCustomFolderNameEnabled: true
            }
        );

        expect(dispatchedActions).toEqual([
            albumRenamingStarted(),
            albumRenamingFailed(new CatalogError("AlbumFolderNameAlreadyTakenErr", "Folder name already taken"))
        ]);
    });

    it("should dispatch albumRenamingFailed with error code when API returns AlbumNameMandatoryErr", async () => {
        const albumId = {owner: "myself", folderName: "old-folder"};

        const fakePort = new FakeSaveAlbumNamePort();
        fakePort.shouldFail = true;
        fakePort.errorCode = "AlbumNameMandatoryErr";
        fakePort.errorMessage = "Album name is mandatory";

        const dispatchedActions: Action<CatalogViewerState, any>[] = [];

        await saveAlbumNameThunk(
            dispatchedActions.push.bind(dispatchedActions),
            fakePort,
            {
                albumId,
                albumName: "",
                customFolderName: "",
                isCustomFolderNameEnabled: false
            }
        );

        expect(dispatchedActions).toEqual([
            albumRenamingStarted(),
            albumRenamingFailed(new CatalogError("AlbumNameMandatoryErr", "Album name is mandatory"))
        ]);
    });

    it("should dispatch albumRenamingFailed with error message when API returns unexpected CatalogError", async () => {
        const albumId = {owner: "myself", folderName: "old-folder"};

        const fakePort = new FakeSaveAlbumNamePort();
        fakePort.shouldFail = true;
        fakePort.errorCode = "UnexpectedError";
        fakePort.errorMessage = "Something unexpected happened";

        const dispatchedActions: Action<CatalogViewerState, any>[] = [];

        await saveAlbumNameThunk(
            dispatchedActions.push.bind(dispatchedActions),
            fakePort,
            {
                albumId,
                albumName: "New Name",
                customFolderName: "",
                isCustomFolderNameEnabled: false
            }
        );

        expect(dispatchedActions).toEqual([
            albumRenamingStarted(),
            albumRenamingFailed(new CatalogError("UnexpectedError", "Something unexpected happened")),
        ]);
    });

    it("should not redirect when album ID remains the same", async () => {
        const albumId = {owner: "myself", folderName: "same-folder"};
        const newName = "New Album Name";
        const sameAlbumId = {owner: "myself", folderName: "same-folder"};

        const fakePort = new FakeSaveAlbumNamePort();
        fakePort.renameAlbum = async () => sameAlbumId;

        const dispatchedActions: Action<CatalogViewerState, any>[] = [];

        await saveAlbumNameThunk(
            dispatchedActions.push.bind(dispatchedActions),
            fakePort,
            {
                albumId,
                albumName: newName,
                customFolderName: "",
                isCustomFolderNameEnabled: false
            }
        );

        expect(dispatchedActions).toEqual([
            albumRenamingStarted(),
            albumRenamed({previousAlbumId: albumId, newAlbumId: sameAlbumId, newName, redirectTo: undefined}),
        ]);
    });

    it("should dispatch albumRenamingFailed with generic message when API returns non-CatalogError", async () => {
        const albumId = {owner: "myself", folderName: "old-folder"};

        const fakePort = new FakeSaveAlbumNamePort();
        fakePort.shouldFail = true;
        fakePort.errorMessage = "Network error";

        const dispatchedActions: Action<CatalogViewerState, any>[] = [];

        await saveAlbumNameThunk(
            dispatchedActions.push.bind(dispatchedActions),
            fakePort,
            {
                albumId,
                albumName: "New Name",
                customFolderName: "",
                isCustomFolderNameEnabled: false
            }
        );

        expect(dispatchedActions).toEqual([
            albumRenamingStarted(),
            albumRenamingFailed({message: "Network error"}),
        ]);
    });
});
