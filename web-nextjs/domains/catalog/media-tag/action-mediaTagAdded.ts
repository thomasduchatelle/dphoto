import {CatalogViewerState, MediaId} from "../language";
import {createAction} from "@/libs/daction";

interface MediaTagAddedPayload {
    mediaId: MediaId
    tag: string
}

export const mediaTagAdded = createAction<CatalogViewerState, MediaTagAddedPayload>(
    "mediaTagAdded",
    (current: CatalogViewerState, {mediaId, tag}: MediaTagAddedPayload): CatalogViewerState => ({
        ...current,
        medias: current.medias.map(day => ({
            ...day,
            medias: day.medias.map(media =>
                media.id === mediaId && !media.tags.includes(tag)
                    ? {...media, tags: [...media.tags, tag]}
                    : media
            ),
        })),
    })
);

export type MediaTagAdded = ReturnType<typeof mediaTagAdded>;
