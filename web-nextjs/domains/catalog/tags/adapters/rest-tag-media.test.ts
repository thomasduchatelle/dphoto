import {describe, expect, it} from "vitest";
import {MediaType} from "../../language";
import {convertTagMediaFromRest, RestTagMedia} from "./rest-tag-media";

const restMedia: RestTagMedia = {
    id: "co-1",
    type: "IMAGE",
    filename: "clair-obscur-1.jpg",
    time: "2025-04-25T09:00:00Z",
    source: "Canon EOS R6",
    albumId: {owner: "sandfall", folderName: "/clair-obscur"},
    tags: ["favorite", "people:luke"],
};

describe("convertTagMediaFromRest", () => {
    it("maps the pinned list-medias-by-tag DTO to a TagMedia linking back to its album", () => {
        expect(convertTagMediaFromRest(restMedia)).toEqual({
            id: "co-1",
            type: MediaType.IMAGE,
            time: new Date("2025-04-25T09:00:00Z"),
            source: "Canon EOS R6",
            albumId: {owner: "sandfall", folderName: "clair-obscur"},
            tags: ["favorite", "people:luke"],
            uiRelativePath: "/albums/sandfall/clair-obscur/co-1/clair-obscur-1.jpg",
            contentPath: "/api/v1/owners/sandfall/medias/co-1/clair-obscur-1.jpg",
            thumbnailUrl: "/api/v1/owners/sandfall/medias/co-1/clair-obscur-1.jpg?w=360",
        });
    });

    it("uses the video placeholder thumbnail for videos", () => {
        const video = convertTagMediaFromRest({...restMedia, type: "VIDEO"});

        expect(video.type).toEqual(MediaType.VIDEO);
        expect(video.thumbnailUrl).toEqual("/video-placeholder.png");
    });
});
