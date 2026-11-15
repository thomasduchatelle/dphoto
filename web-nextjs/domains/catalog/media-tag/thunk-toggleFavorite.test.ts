import {ToggleFavoritePort, toggleFavoriteThunk} from "./thunk-toggleFavorite";
import {mediaTagsUpdated} from "./action-mediaTagsUpdated";
import {AlbumId, FAVORITE_TAG, MediaId} from "../language";

class FakeToggleFavoritePort implements ToggleFavoritePort {
    public added: Array<{owner: string; mediaId: MediaId; tag: string}> = [];
    public removed: Array<{owner: string; mediaId: MediaId; tag: string}> = [];

    constructor(private readonly shouldFail: boolean = false) {
    }

    async addTag(owner: string, mediaId: MediaId, tag: string): Promise<void> {
        if (this.shouldFail) {
            return Promise.reject(new Error("TEST add failed"));
        }
        this.added.push({owner, mediaId, tag});
    }

    async removeTag(owner: string, mediaId: MediaId, tag: string): Promise<void> {
        if (this.shouldFail) {
            return Promise.reject(new Error("TEST remove failed"));
        }
        this.removed.push({owner, mediaId, tag});
    }
}

const albumId: AlbumId = {owner: "myself", folderName: "jan-25"};

describe("thunk:toggleFavoriteThunk", () => {
    it("optimistically adds favorite and calls addTag when the media is not yet favourited", async () => {
        const port = new FakeToggleFavoritePort();
        const dispatched: any[] = [];

        await toggleFavoriteThunk(dispatched.push.bind(dispatched), port, albumId, "media-1", []);

        expect(dispatched).toEqual([mediaTagsUpdated({mediaId: "media-1", tags: [FAVORITE_TAG]})]);
        expect(port.added).toEqual([{owner: "myself", mediaId: "media-1", tag: FAVORITE_TAG}]);
        expect(port.removed).toEqual([]);
    });

    it("optimistically removes favorite and calls removeTag when the media is already favourited", async () => {
        const port = new FakeToggleFavoritePort();
        const dispatched: any[] = [];

        await toggleFavoriteThunk(dispatched.push.bind(dispatched), port, albumId, "media-1", [FAVORITE_TAG, "people:luke"]);

        expect(dispatched).toEqual([mediaTagsUpdated({mediaId: "media-1", tags: ["people:luke"]})]);
        expect(port.removed).toEqual([{owner: "myself", mediaId: "media-1", tag: FAVORITE_TAG}]);
        expect(port.added).toEqual([]);
    });

    it("reverts to the original tags when the add call fails", async () => {
        const port = new FakeToggleFavoritePort(true);
        const dispatched: any[] = [];

        await toggleFavoriteThunk(dispatched.push.bind(dispatched), port, albumId, "media-1", ["people:luke"]);

        expect(dispatched).toEqual([
            mediaTagsUpdated({mediaId: "media-1", tags: ["people:luke", FAVORITE_TAG]}),
            mediaTagsUpdated({mediaId: "media-1", tags: ["people:luke"]}),
        ]);
    });

    it("reverts to the original tags when the remove call fails", async () => {
        const port = new FakeToggleFavoritePort(true);
        const dispatched: any[] = [];

        await toggleFavoriteThunk(dispatched.push.bind(dispatched), port, albumId, "media-1", [FAVORITE_TAG]);

        expect(dispatched).toEqual([
            mediaTagsUpdated({mediaId: "media-1", tags: []}),
            mediaTagsUpdated({mediaId: "media-1", tags: [FAVORITE_TAG]}),
        ]);
    });

    it("does nothing when no album is displayed", async () => {
        const port = new FakeToggleFavoritePort();
        const dispatched: any[] = [];

        await toggleFavoriteThunk(dispatched.push.bind(dispatched), port, undefined, "media-1", []);

        expect(dispatched).toEqual([]);
        expect(port.added).toEqual([]);
        expect(port.removed).toEqual([]);
    });
});
