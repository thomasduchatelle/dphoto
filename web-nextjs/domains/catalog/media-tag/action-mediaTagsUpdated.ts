import {CatalogViewerState, MediaId} from "../language";
import {createAction} from "@/libs/daction";

interface MediaTagsUpdatedPayload {
    mediaId: MediaId
    tags: string[]
}

export const mediaTagsUpdated = createAction<CatalogViewerState, MediaTagsUpdatedPayload>(
    "mediaTagsUpdated",
    (current: CatalogViewerState, {mediaId, tags}: MediaTagsUpdatedPayload): CatalogViewerState => ({
        ...current,
        medias: current.medias.map(day => ({
            ...day,
            medias: day.medias.map(media => media.id === mediaId ? {...media, tags} : media),
        })),
    })
);

export type MediaTagsUpdated = ReturnType<typeof mediaTagsUpdated>;
