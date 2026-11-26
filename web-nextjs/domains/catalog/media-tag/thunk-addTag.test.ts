import {AddTagPort, addTagThunk} from "./thunk-addTag";
import {mediaTagAdded} from "./action-mediaTagAdded";
import {mediaTagRemoved} from "./action-mediaTagRemoved";
import {AlbumId, FAVORITE_TAG, MediaId} from "../language";

class FakeAddTagPort implements AddTagPort {
    public added: Array<{owner: string; mediaId: MediaId; tag: string}> = [];

    constructor(private readonly shouldFail: boolean = false) {
    }

    async addTag(owner: string, mediaId: MediaId, tag: string): Promise<void> {
        if (this.shouldFail) {
            return Promise.reject(new Error("TEST add failed"));
        }
        this.added.push({owner, mediaId, tag});
    }
}

const albumId: AlbumId = {owner: "myself", folderName: "jan-25"};

describe("thunk:addTagThunk", () => {
    it("optimistically adds the tag and calls addTag", async () => {
        const port = new FakeAddTagPort();
        const dispatched: any[] = [];

        await addTagThunk(dispatched.push.bind(dispatched), port, albumId, "media-1", FAVORITE_TAG);

        expect(dispatched).toEqual([mediaTagAdded({mediaId: "media-1", tag: FAVORITE_TAG})]);
        expect(port.added).toEqual([{owner: "myself", mediaId: "media-1", tag: FAVORITE_TAG}]);
    });

    it("removes the optimistically added tag when the add call fails", async () => {
        const port = new FakeAddTagPort(true);
        const dispatched: any[] = [];

        await addTagThunk(dispatched.push.bind(dispatched), port, albumId, "media-1", FAVORITE_TAG);

        expect(dispatched).toEqual([
            mediaTagAdded({mediaId: "media-1", tag: FAVORITE_TAG}),
            mediaTagRemoved({mediaId: "media-1", tag: FAVORITE_TAG}),
        ]);
    });

    it("does nothing when no album is displayed", async () => {
        const port = new FakeAddTagPort();
        const dispatched: any[] = [];

        await addTagThunk(dispatched.push.bind(dispatched), port, undefined, "media-1", FAVORITE_TAG);

        expect(dispatched).toEqual([]);
        expect(port.added).toEqual([]);
    });
});
