import {LoadTagMediasPort} from "../thunk-loadTagMedias";
import {TagMedia} from "../language";
import {convertTagMediaFromRest, RestTagMedia} from "./rest-tag-media";
import {FAVORITE_TAG} from "../language";

const FAVORITE_MEDIAS: RestTagMedia[] = [
    {
        id: "co-1",
        type: "IMAGE",
        filename: "clair-obscur-1.jpg",
        time: "2025-04-25T09:00:00Z",
        source: "",
        albumId: {owner: "sandfall", folderName: "/clair-obscur"},
        tags: [FAVORITE_TAG],
    },
    {
        id: "co-3",
        type: "IMAGE",
        filename: "clair-obscur-3.jpg",
        time: "2025-04-25T11:30:00Z",
        source: "Canon EOS R6",
        albumId: {owner: "sandfall", folderName: "/clair-obscur"},
        tags: [FAVORITE_TAG, "people:luke"],
    },
    {
        id: "ab-2",
        type: "IMAGE",
        filename: "astro-bot-02.jpg",
        time: "2024-09-06T14:12:00Z",
        source: "",
        albumId: {owner: "sony", folderName: "/astro-bot"},
        tags: [FAVORITE_TAG],
    },
    {
        id: "ds1-1",
        type: "IMAGE",
        filename: "death-stranding-1-01.jpg",
        time: "2023-11-08T18:40:00Z",
        source: "iPhone 15 Pro",
        albumId: {owner: "kojima", folderName: "/death-stranding-1"},
        tags: [FAVORITE_TAG, "place:school"],
    },
];

export class MockTagsAdapter implements LoadTagMediasPort {
    public async fetchMediasByTag(tag: string): Promise<TagMedia[]> {
        if (tag !== FAVORITE_TAG) {
            return [];
        }
        return FAVORITE_MEDIAS.map(convertTagMediaFromRest);
    }
}
