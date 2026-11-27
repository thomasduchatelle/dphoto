import {describe, expect, it} from "vitest";
import {MediaType} from "../language";
import {TagMedia, initialTagPageState, FAVORITE_TAG} from "./language";
import {loadTagMedias, LoadTagMediasPort} from "./thunk-loadTagMedias";
import {tagPageReducer, TagPageAction} from "./reducer";
import {tagPageSelector} from "./selector-tagPage";

const favouriteAcrossTwoDays: TagMedia[] = [
    {
        id: "m2",
        type: MediaType.IMAGE,
        time: new Date("2025-04-25T11:00:00Z"),
        albumId: {owner: "sandfall", folderName: "clair-obscur"},
        tags: [FAVORITE_TAG],
        uiRelativePath: "m2.jpg",
        contentPath: "/m2.jpg",
        thumbnailUrl: "/m2.jpg?w=360",
        source: "",
    },
    {
        id: "m1",
        type: MediaType.IMAGE,
        time: new Date("2025-04-24T09:00:00Z"),
        albumId: {owner: "sony", folderName: "astro-bot"},
        tags: [FAVORITE_TAG],
        uiRelativePath: "m1.jpg",
        contentPath: "/m1.jpg",
        thumbnailUrl: "/m1.jpg?w=360",
        source: "",
    },
];

class FakeTagsAdapter implements LoadTagMediasPort {
    constructor(private readonly mediasByTag: Map<string, TagMedia[]>) {}

    fetchMediasByTag(tag: string): Promise<TagMedia[]> {
        return Promise.resolve(this.mediasByTag.get(tag) ?? []);
    }
}

class FailingTagsAdapter implements LoadTagMediasPort {
    constructor(private readonly error: Error) {}

    fetchMediasByTag(): Promise<TagMedia[]> {
        return Promise.reject(this.error);
    }
}

function runThunkThenSelect(adapter: LoadTagMediasPort, tag: string) {
    return (async () => {
        const dispatched: TagPageAction[] = [];
        await loadTagMedias(a => dispatched.push(a), adapter, tag);
        const state = dispatched.reduce(tagPageReducer, initialTagPageState(tag));
        return tagPageSelector(state);
    })();
}

describe("tag page", () => {
    it("groups the medias carrying the tag by day, newest day first", async () => {
        const adapter = new FakeTagsAdapter(new Map([[FAVORITE_TAG, favouriteAcrossTwoDays]]));

        const selection = await runThunkThenSelect(adapter, FAVORITE_TAG);

        expect(selection.tag).toEqual(FAVORITE_TAG);
        expect(selection.loaded).toBe(true);
        expect(selection.isEmpty).toBe(false);
        expect(selection.medias).toEqual([
            {day: new Date(2025, 3, 25), medias: [favouriteAcrossTwoDays[0]]},
            {day: new Date(2025, 3, 24), medias: [favouriteAcrossTwoDays[1]]},
        ]);
    });

    it("reports an empty page when no media carries the tag", async () => {
        const adapter = new FakeTagsAdapter(new Map());

        const selection = await runThunkThenSelect(adapter, FAVORITE_TAG);

        expect(selection.loaded).toBe(true);
        expect(selection.isEmpty).toBe(true);
        expect(selection.medias).toEqual([]);
    });

    it("surfaces an error when the medias cannot be loaded", async () => {
        const cause = new Error("boom");
        const adapter = new FailingTagsAdapter(cause);

        const selection = await runThunkThenSelect(adapter, FAVORITE_TAG);

        expect(selection.loaded).toBe(true);
        expect(selection.isEmpty).toBe(false);
        expect(selection.error).toEqual(new Error(`failed to load medias of tag '${FAVORITE_TAG}'`, {cause}));
    });
});
