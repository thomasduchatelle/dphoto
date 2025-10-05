import {onAlbumFilterChangeDeclaration} from "./thunk-onAlbumFilterChange";
import {onPageRefreshDeclaration} from "./thunk-onPageRefresh";

export * from "./action-albumsAndMediasLoaded";
export * from "./action-albumsLoaded";
export * from "./action-albumsFiltered";
export * from "./action-mediasLoaded";
export * from "./action-mediaLoadFailed";
export * from "./action-noAlbumAvailable";
export type {FetchAlbumsAndMediasPort} from "./thunk-onPageRefresh";

/**
 * Thunks related to catalog navigation.
 *
 * Expected handler types:
 * - `onAlbumFilterChange`: `(criterion: AlbumFilterCriterion) => void`
 * - `onPageRefresh`: `(albumId?: AlbumId) => Promise<void>`
 */
export const navigationThunks = {
    onAlbumFilterChange: onAlbumFilterChangeDeclaration,
    onPageRefresh: onPageRefreshDeclaration,
};
export * from "./selector-catalog-viewer-page";
