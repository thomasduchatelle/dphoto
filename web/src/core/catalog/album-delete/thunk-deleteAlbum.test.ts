import {DeleteAlbumPort, deleteAlbumThunk} from "./thunk-deleteAlbum";
import {Album, AlbumId, albumIdEquals, Media, MediaType} from "../language";
import {deleteAlbumStarted} from "./action-deleteAlbumStarted";
import {albumDeleteFailed} from "./action-albumDeleteFailed";
import {albumsAndMediasLoaded, noAlbumAvailable} from "../navigation";


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

function newMedia(mediaId: string, dateTime: string): Media {
    return {
        id: mediaId,
        type: MediaType.IMAGE,
        time: new Date(dateTime),
        uiRelativePath: `${mediaId}/image-${mediaId}.jpg`,
        contentPath: `/content/$\{id}/image-${mediaId}.jpg`,
        source: 'Samsung Galaxy S24'
    };
}

const fetchError = new Error("TEST fetch failed");

class FakeDeleteAlbumPort implements DeleteAlbumPort {
    constructor(
        public albums: Album[] = [],
        public medias: Map<AlbumId, Media[]> = new Map(),
        public deleteShouldFail: boolean = false,
        public fetchShouldFail: boolean = false,
    ) {
    }

    async deleteAlbum(albumId: AlbumId): Promise<void> {
        if (this.deleteShouldFail) {
            return Promise.reject({code: "OrphanedMediasError"})
        }
        this.albums = this.albums.filter(a => !albumIdEquals(a.albumId, albumId));
    }

    async fetchAlbums(): Promise<Album[]> {
        if (this.fetchShouldFail) {
            return Promise.reject(fetchError);
        }
        return this.albums;
    }

    async fetchMedias(albumId: AlbumId): Promise<Media[]> {
        if (this.fetchShouldFail) {
            return Promise.reject(fetchError);
        }
        return this.medias.get(albumId) || [];
    }
}

describe("deleteAlbumThunk", () => {
    const albumA = makeAlbum("A");
    const albumB = makeAlbum("B");
    const mediasB = [newMedia('b1', "2024-01-01T00:00:00Z")];

    it("dispatches startDeleteAlbum and albumsAndMediasLoaded without redirection when deletion succeed on non-selected album", async () => {
        const port = new FakeDeleteAlbumPort([albumA, albumB], new Map([[albumB.albumId, mediasB]]));
        const dispatched: any[] = [];
        await deleteAlbumThunk(
            dispatched.push.bind(dispatched),
            port,
            albumB.albumId,
            albumA.albumId
        );
        expect(dispatched).toEqual([
            deleteAlbumStarted(),
            albumsAndMediasLoaded({albums: [albumB], medias: mediasB, mediasFromAlbumId: albumB.albumId, redirectTo: albumB.albumId})
        ]);

        expect(port.albums.map(album => album.albumId)).toEqual([albumB.albumId]);
    });

    it("dispatches startDeleteAlbum and albumsAndMediasLoaded with redirection to first album when deletion succeed on selected album", async () => {
        const port = new FakeDeleteAlbumPort([albumA, albumB], new Map([[albumB.albumId, mediasB]]));
        const dispatched: any[] = [];
        await deleteAlbumThunk(
            dispatched.push.bind(dispatched),
            port,
            albumA.albumId,
            albumA.albumId
        );

        expect(dispatched).toEqual([
            deleteAlbumStarted(),
            albumsAndMediasLoaded({albums: [albumB], medias: mediasB, mediasFromAlbumId: albumB.albumId, redirectTo: albumB.albumId}),
        ]);
        expect(port.albums.map(album => album.albumId)).toEqual([albumB.albumId]);
    });

    it("dispatches startDeleteAlbum and albumsAndMediasLoaded with empty medias if no album to redirect to", async () => {
        const port = new FakeDeleteAlbumPort([albumA], new Map());
        const dispatched: any[] = [];
        await deleteAlbumThunk(
            dispatched.push.bind(dispatched),
            port,
            albumA.albumId,
            albumA.albumId
        );

        expect(dispatched).toEqual([
            deleteAlbumStarted(),
            noAlbumAvailable(undefined),
        ]);
        expect(port.albums.map(album => album.albumId)).toEqual([]);
    });

    it("dispatches startDeleteAlbum and albumFailedToDelete on delete error", async () => {
        const port = new FakeDeleteAlbumPort([albumA, albumB], new Map(), true);
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

    it("dispatches startDeleteAlbum and albumFailedToDelete on fetch albums error", async () => {
        const port = new FakeDeleteAlbumPort([albumA, albumB], new Map(), false, true);
        const dispatched: any[] = [];
        await deleteAlbumThunk(
            dispatched.push.bind(dispatched),
            port,
            albumA.albumId,
            albumA.albumId
        );
        expect(dispatched).toEqual([
            deleteAlbumStarted(),
            noAlbumAvailable(fetchError),
        ]);
    });

    it("dispatches startDeleteAlbum and albumFailedToDelete on fetch medias error", async () => {
        const port = new FakeDeleteAlbumPort([albumA, albumB], new Map([[albumB.albumId, mediasB]]), false, true);
        const dispatched: any[] = [];
        await deleteAlbumThunk(
            dispatched.push.bind(dispatched),
            port,
            albumA.albumId,
            albumA.albumId
        );
        expect(dispatched).toEqual([
            deleteAlbumStarted(),
            noAlbumAvailable(fetchError),
        ]);
    });
});
