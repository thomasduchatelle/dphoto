import {CatalogViewerState, MediaId} from "../language";
import {createAction} from "@/libs/daction";

interface MediaTagRemovedPayload {
    mediaId: MediaId
    tag: string
}

export const mediaTagRemoved = createAction<CatalogViewerState, MediaTagRemovedPayload>(
    "mediaTagRemoved",
    (current: CatalogViewerState, {mediaId, tag}: MediaTagRemovedPayload): CatalogViewerState => ({
        ...current,
        medias: current.medias.map(day => ({
            ...day,
            medias: day.medias.map(media =>
                media.id === mediaId
                    ? {...media, tags: media.tags.filter(t => t !== tag)}
                    : media
            ),
        })),
    })
);

export type MediaTagRemoved = ReturnType<typeof mediaTagRemoved>;
