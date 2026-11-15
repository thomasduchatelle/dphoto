import {MediaType} from "../../language";
import {mediaUrl} from "@/libs/requests/media-url";
import {TagMedia} from "../language";

export interface RestTagMedia {
    id: string
    type: string
    filename: string
    time: string
    source: string
    albumId: {
        owner: string
        folderName: string
    }
    tags: string[]
}

function convertToType(type: string): MediaType {
    switch ((type ?? "").toLowerCase()) {
        case "image":
            return MediaType.IMAGE;
        case "video":
            return MediaType.VIDEO;
        default:
            return MediaType.OTHER;
    }
}

export function convertTagMediaFromRest(media: RestTagMedia): TagMedia {
    const owner = media.albumId.owner;
    const folderName = media.albumId.folderName.replace(/^\//, "");
    const type = convertToType(media.type);
    const contentPath = `/api/v1/owners/${owner}/medias/${media.id}/${media.filename}`;

    return {
        id: media.id,
        source: media.source,
        type,
        time: new Date(media.time),
        albumId: {owner, folderName},
        tags: media.tags ?? [],
        uiRelativePath: `/albums/${owner}/${folderName}/${media.id}/${media.filename}`,
        contentPath,
        thumbnailUrl: type === MediaType.VIDEO
            ? "/video-placeholder.png"
            : mediaUrl(contentPath, 360),
    };
}
