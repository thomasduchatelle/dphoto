import {mediaTagRemoved} from "./action-mediaTagRemoved";
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
            tags: [FAVORITE_TAG, "people:luke"],
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

describe("action:mediaTagRemoved", () => {
    it("should remove the tag from the targeted media only", () => {
        const action = mediaTagRemoved({mediaId: "media-1", tag: FAVORITE_TAG});

        const got = action.reducer({...loadedStateWithTwoAlbums, medias: twoMediasOnOneDay}, action);

        expect(got.medias[0].medias.map(m => ({id: m.id, tags: m.tags}))).toEqual([
            {id: "media-1", tags: ["people:luke"]},
            {id: "media-2", tags: [FAVORITE_TAG]},
        ]);
    });

    it("should leave the media unchanged when the tag is not present", () => {
        const action = mediaTagRemoved({mediaId: "media-1", tag: "people:leia"});

        const got = action.reducer({...loadedStateWithTwoAlbums, medias: twoMediasOnOneDay}, action);

        expect(got.medias[0].medias.find(m => m.id === "media-1")?.tags).toEqual([FAVORITE_TAG, "people:luke"]);
    });
});
