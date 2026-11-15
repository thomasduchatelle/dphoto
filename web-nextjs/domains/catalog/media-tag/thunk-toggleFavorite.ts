import {ThunkDeclaration} from "@/libs/dthunks";
import {AlbumId, CatalogViewerState, FAVORITE_TAG, MediaId} from "../language";
import {CatalogDispatch} from "../common/catalog-dispatch";
import {getDisplayedAlbumId} from "../language/selector-displayedAlbum";
import {mediaTagsUpdated} from "./action-mediaTagsUpdated";

export interface ToggleFavoritePort {
    addTag(owner: string, mediaId: MediaId, tag: string): Promise<void>;

    removeTag(owner: string, mediaId: MediaId, tag: string): Promise<void>;
}

export type ToggleFavoriteThunk = (mediaId: MediaId, currentTags: string[]) => Promise<void>;

export async function toggleFavoriteThunk(
    dispatch: (action: ReturnType<typeof mediaTagsUpdated>) => void,
    port: ToggleFavoritePort,
    displayedAlbumId: AlbumId | undefined,
    mediaId: MediaId,
    currentTags: string[]
): Promise<void> {
    if (!displayedAlbumId) {
        return;
    }

    const wasFavorite = currentTags.includes(FAVORITE_TAG);
    const optimisticTags = wasFavorite
        ? currentTags.filter(tag => tag !== FAVORITE_TAG)
        : [...currentTags, FAVORITE_TAG];

    dispatch(mediaTagsUpdated({mediaId, tags: optimisticTags}));

    try {
        if (wasFavorite) {
            await port.removeTag(displayedAlbumId.owner, mediaId, FAVORITE_TAG);
        } else {
            await port.addTag(displayedAlbumId.owner, mediaId, FAVORITE_TAG);
        }
    } catch {
        dispatch(mediaTagsUpdated({mediaId, tags: currentTags}));
    }
}

export const toggleFavoriteDeclaration: ThunkDeclaration<
    CatalogViewerState,
    { displayedAlbumId: AlbumId | undefined },
    ToggleFavoriteThunk,
    CatalogDispatch & { adapter: ToggleFavoritePort }
> = {
    selector: (state: CatalogViewerState) => ({
        displayedAlbumId: getDisplayedAlbumId(state),
    }),
    factory: ({dispatch, adapter, partialState: {displayedAlbumId}}) => {
        return toggleFavoriteThunk.bind(null, dispatch, adapter, displayedAlbumId);
    },
};
