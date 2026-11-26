import {RemoveTagPort, removeTagThunk} from "./thunk-removeTag";
import {mediaTagAdded} from "./action-mediaTagAdded";
import {mediaTagRemoved} from "./action-mediaTagRemoved";
import {AlbumId, FAVORITE_TAG, MediaId} from "../language";

class FakeRemoveTagPort implements RemoveTagPort {
    public removed: Array<{owner: string; mediaId: MediaId; tag: string}> = [];

    constructor(private readonly shouldFail: boolean = false) {
    }

    async removeTag(owner: string, mediaId: MediaId, tag: string): Promise<void> {
        if (this.shouldFail) {
            return Promise.reject(new Error("TEST remove failed"));
        }
        this.removed.push({owner, mediaId, tag});
    }
}

const albumId: AlbumId = {owner: "myself", folderName: "jan-25"};

describe("thunk:removeTagThunk", () => {
    it("optimistically removes the tag and calls removeTag", async () => {
        const port = new FakeRemoveTagPort();
        const dispatched: any[] = [];

        await removeTagThunk(dispatched.push.bind(dispatched), port, albumId, "media-1", FAVORITE_TAG);

        expect(dispatched).toEqual([mediaTagRemoved({mediaId: "media-1", tag: FAVORITE_TAG})]);
        expect(port.removed).toEqual([{owner: "myself", mediaId: "media-1", tag: FAVORITE_TAG}]);
    });

    it("adds the optimistically removed tag back when the remove call fails", async () => {
        const port = new FakeRemoveTagPort(true);
        const dispatched: any[] = [];

        await removeTagThunk(dispatched.push.bind(dispatched), port, albumId, "media-1", FAVORITE_TAG);

        expect(dispatched).toEqual([
            mediaTagRemoved({mediaId: "media-1", tag: FAVORITE_TAG}),
            mediaTagAdded({mediaId: "media-1", tag: FAVORITE_TAG}),
        ]);
    });

    it("does nothing when no album is displayed", async () => {
        const port = new FakeRemoveTagPort();
        const dispatched: any[] = [];

        await removeTagThunk(dispatched.push.bind(dispatched), port, undefined, "media-1", FAVORITE_TAG);

        expect(dispatched).toEqual([]);
        expect(port.removed).toEqual([]);
    });
});
