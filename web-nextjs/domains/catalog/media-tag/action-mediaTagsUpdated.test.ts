import {mediaTagsUpdated} from "./action-mediaTagsUpdated";
import {FAVORITE_TAG, MediaType, MediaWithinADay} from "../language";
import {loadedStateWithTwoAlbums} from "../tests/test-helper-state";

const twoMediasOnOneDay: MediaWithinADay[] = [{
    day: new Date(2025, 0, 5),
    medias: [
        {
            id: "media-1",
            type: MediaType.IMAGE,
            time: new Date("2025-01-05T12:42:00Z"),
            uiRelativePath: "media-1/image.jpg",
            contentPath: "/media-1.jpg",
            thumbnailUrl: "/media-1.jpg?w=360",
            source: "",
            tags: [],
        },
        {
            id: "media-2",
            type: MediaType.IMAGE,
            time: new Date("2025-01-05T13:00:00Z"),
            uiRelativePath: "media-2/image.jpg",
            contentPath: "/media-2.jpg",
            thumbnailUrl: "/media-2.jpg?w=360",
            source: "",
            tags: [FAVORITE_TAG],
        },
    ],
}];

describe("action:mediaTagsUpdated", () => {
    it("should replace the tags of the targeted media only", () => {
        const action = mediaTagsUpdated({mediaId: "media-1", tags: [FAVORITE_TAG]});

        const got = action.reducer({...loadedStateWithTwoAlbums, medias: twoMediasOnOneDay}, action);

        expect(got.medias[0].medias.map(m => ({id: m.id, tags: m.tags}))).toEqual([
            {id: "media-1", tags: [FAVORITE_TAG]},
            {id: "media-2", tags: [FAVORITE_TAG]},
        ]);
    });

    it("should remove a tag when the new list no longer contains it", () => {
        const action = mediaTagsUpdated({mediaId: "media-2", tags: []});

        const got = action.reducer({...loadedStateWithTwoAlbums, medias: twoMediasOnOneDay}, action);

        expect(got.medias[0].medias.find(m => m.id === "media-2")?.tags).toEqual([]);
    });
});
