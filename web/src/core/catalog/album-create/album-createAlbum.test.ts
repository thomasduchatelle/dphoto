import {createAlbum, CreateAlbumPort, CreateAlbumPayload} from "./album-createAlbum";
import {Album, AlbumId, albumIdEquals} from "../language";
import {albumsLoaded} from "../navigation";

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

    async createAlbum({name, start, end, forcedFolderName}: CreateAlbumPayload): Promise<AlbumId> {
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

describe("createAlbum action", () => {
    it("should call store the new Album and dispatch a albumsLoadedAction action with an existing album with the newly created one", async () => {
        const expectedAlbumId: AlbumId = {owner: "myself", folderName: "/album:new-album"};
        const request: CreateAlbumPayload = {
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

        // Simulate the thunk execution
        const albumId = await createAlbum.reducer(
            {
                dispatch: dispatched.push.bind(dispatched),
                createAlbumPort: fakePort
            },
            request
        );

        expect(albumId).toEqual(expectedAlbumId);

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
            albumsLoaded({
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

    it("supports action comparison for testing", () => {
        const payload1: CreateAlbumPayload = {
            name: "Album A",
            start: new Date("2025-01-01"),
            end: new Date("2025-01-02"),
            forcedFolderName: "/album:album-a",
        };
        const payload2: CreateAlbumPayload = {
            name: "Album A",
            start: new Date("2025-01-01"),
            end: new Date("2025-01-02"),
            forcedFolderName: "/album:album-a",
        };
        const payload3: CreateAlbumPayload = {
            name: "Album B",
            start: new Date("2025-01-01"),
            end: new Date("2025-01-02"),
            forcedFolderName: "/album:album-b",
        };

        const action1 = createAlbum(payload1);
        const action2 = createAlbum(payload2);
        const action3 = createAlbum(payload3);

        expect(action1).toEqual(action2);
        expect(action1).not.toEqual(action3);
        expect([action1]).toContainEqual(action2);
    });
});
