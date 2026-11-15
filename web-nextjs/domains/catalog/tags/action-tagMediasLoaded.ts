import {createAction} from "@/libs/daction";
import {TagMedia, TagPageState} from "./language";
import {groupTagMediasByDay} from "./group-tag-medias-by-day";

interface TagMediasLoadedPayload {
    tag: string
    medias: TagMedia[]
}

export const tagMediasLoaded = createAction<TagPageState, TagMediasLoadedPayload>(
    "tagMediasLoaded",
    (current: TagPageState, {tag, medias}: TagMediasLoadedPayload): TagPageState => {
        if (current.tag !== tag) {
            return current;
        }

        return {
            ...current,
            medias: groupTagMediasByDay(medias),
            loaded: true,
            error: undefined,
        };
    }
);

export type TagMediasLoaded = ReturnType<typeof tagMediasLoaded>;
