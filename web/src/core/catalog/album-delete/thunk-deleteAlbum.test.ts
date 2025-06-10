import {DeleteAlbumPort, deleteAlbumThunk} from "./thunk-deleteAlbum";
import {Album, AlbumId} from "../language";
import {deleteAlbumStarted} from "./action-deleteAlbumStarted";
import {albumDeleteFailed} from "./action-albumDeleteFailed";
import {albumDeleted} from "./action-albumDeleted";


function makeAlbum(id: string): Album {
    return {
        albumId: {owner: "owner", folderName: id},
        name: id,
        start: new Date(),
        end: new Date(),
        totalCount: 0,
        temperature: 0,
        relativeTemperature: 0,
        sharedWith: []
    };
}

class FakeDeleteAlbumPort implements DeleteAlbumPort {
    constructor(
        public albums: Album[] = [],
        public deleteShouldFail: boolean = false,
        public fetchShouldFail: boolean = false,
    ) {
    }

    async deleteAlbum(albumId: AlbumId): Promise<void> {
        if (this.deleteShouldFail) {
            return Promise.reject({code: "OrphanedMediasError"})
        }
        this.albums = this.albums.filter(a => !(a.albumId.owner === albumId.owner && a.albumId.folderName === albumId.folderName));
    }

    async fetchAlbums(): Promise<Album[]> {
        if (this.fetchShouldFail) {
            return Promise.reject("TEST fetch failed");
        }
        return this.albums;
    }
}

describe("deleteAlbumThunk", () => {
    const albumA = makeAlbum("A");
    const albumB = makeAlbum("B");

    it("dispatches startDeleteAlbum and albumDeleted without redirection when deletion succeed on non-selected album", async () => {
        const port = new FakeDeleteAlbumPort([albumA, albumB]);
        const dispatched: any[] = [];
        await deleteAlbumThunk(
            dispatched.push.bind(dispatched),
            port,
            albumB.albumId,
            albumA.albumId
        );
        expect(dispatched).toEqual([
            deleteAlbumStarted(),
            albumDeleted({albums: [albumB], redirectTo: undefined})
        ]);

        expect(port.albums.map(album => album.albumId)).toEqual([albumB.albumId]);
    });

    it("dispatches startDeleteAlbum and albumDeleted with redirection to first album when deletion succeed on selected album", async () => {
        const port = new FakeDeleteAlbumPort([albumA, albumB]);
        const dispatched: any[] = [];
        await deleteAlbumThunk(
            dispatched.push.bind(dispatched),
            port,
            albumA.albumId,
            albumA.albumId
        );

        expect(dispatched).toEqual([
            deleteAlbumStarted(),
            albumDeleted({albums: [albumB], redirectTo: albumB.albumId}),
        ]);
        expect(port.albums.map(album => album.albumId)).toEqual([albumB.albumId]);
    });

    it("dispatches startDeleteAlbum and albumFailedToDelete on delete error", async () => {
        const port = new FakeDeleteAlbumPort([albumA, albumB], true);
        const dispatched: any[] = [];
        await deleteAlbumThunk(
            dispatched.push.bind(dispatched),
            port,
            albumA.albumId,
            albumA.albumId
        );
        expect(dispatched).toEqual([
            deleteAlbumStarted(),
            albumDeleteFailed("A technical error prevented the album to be deleted, please report it to a developer."),
        ]);
    });

    it("dispatches startDeleteAlbum and albumFailedToDelete on fetch error", async () => {
        const port = new FakeDeleteAlbumPort([albumA, albumB], false, true);
        const dispatched: any[] = [];
        await deleteAlbumThunk(
            dispatched.push.bind(dispatched),
            port,
            albumA.albumId,
            albumA.albumId
        );
        expect(dispatched).toEqual([
            deleteAlbumStarted(),
            albumDeleteFailed("Failed to fetch albums after deletion: TEST fetch failed"),
        ]);
    });
});
