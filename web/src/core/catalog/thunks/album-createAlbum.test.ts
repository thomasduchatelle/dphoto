import {CreateAlbumPort, CreateAlbumRequest, createAlbumThunk} from "./album-createAlbum";
import {Album, AlbumId, albumIdEquals, catalogActions} from "../domain";

class CreateAlbumPortFake implements CreateAlbumPort {
    defaultAlbumValues = {
        totalCount: 42,
        temperature: 8,
        relativeTemperature: 0.8,
        ownedBy: {name: "myself", users: []},
        sharedWith: [],
    }

    constructor(
        public albums: Album[] = []) {
    }

    async createAlbum({name, start, end, forcedFolderName}: CreateAlbumRequest): Promise<AlbumId> {
        const album = {
            ...this.defaultAlbumValues,
            name,
            start,
            end,
            albumId: {owner: "myself", folderName: forcedFolderName ?? `/album:${name.toLowerCase()}`},
        };
        this.albums.push(album);
        return album.albumId;
    }

    async fetchAlbums(): Promise<Album[]> {
        return this.albums;
    }
}

describe("createAlbumThunk", () => {
    it("should call store the new Album and dispatch a albumsLoadedAction action with an existing album with the newly created one", async () => {
        const expectedAlbumId: AlbumId = {owner: "myself", folderName: "/album:new-album"};
        const request: CreateAlbumRequest = {
            name: "Album 1",
            start: new Date("2025-01-01"),
            end: new Date("2025-01-02"),
            forcedFolderName: "/album:new-album",
        };

        const fakePort = new CreateAlbumPortFake([{
            albumId: {owner: "owner1", folderName: "/album:existing-album"},
            name: "Q1 2025",
            start: new Date("2025-01-01"),
            end: new Date("2025-04-01"),
            totalCount: 10,
            temperature: 0.5,
            relativeTemperature: 0.5,
            sharedWith: [],
            ownedBy: {name: "Owner", users: []},
        }]);

        const dispatched: any[] = [];

        await createAlbumThunk(dispatched.push.bind(dispatched), fakePort, request);

        expect(
            fakePort.albums
                .filter(a => albumIdEquals(a.albumId, expectedAlbumId))
                .map(({name, start, end}) => ({name, start, end}))
        ).toEqual([{
            name: "Album 1",
            start: new Date("2025-01-01"),
            end: new Date("2025-01-02"),
        }]);

        expect(dispatched).toEqual([
            catalogActions.albumsLoadedAction({
                albums: [
                    {
                        albumId: {owner: "owner1", folderName: "/album:existing-album"},
                        name: "Q1 2025",
                        start: new Date("2025-01-01"),
                        end: new Date("2025-04-01"),
                        totalCount: 10,
                        temperature: 0.5,
                        relativeTemperature: 0.5,
                        sharedWith: [],
                        ownedBy: {name: "Owner", users: []},
                    },
                    {
                        ...fakePort.defaultAlbumValues,
                        albumId: expectedAlbumId,
                        name: "Album 1",
                        start: new Date("2025-01-01"),
                        end: new Date("2025-01-02"),
                    }],
                redirectTo: expectedAlbumId
            })
        ]);
    });
});
