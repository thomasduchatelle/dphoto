import {onAlbumFilterChangeDeclaration} from "./thunk-onAlbumFilterChange";
import {onPageRefreshDeclaration} from "./thunk-onPageRefresh";
import {loadAlbumPageDeclaration} from "./thunk-loadAlbumPage";

export * from "./action-albumsAndMediasLoaded";
export * from "./action-albumsLoaded";
export * from "./action-albumsFiltered";
export * from "./action-mediasLoaded";
export * from "./action-mediaLoadFailed";
export * from "./action-noAlbumAvailable";
export type {FetchAlbumsAndMediasPort} from "./thunk-onPageRefresh";
export type {LoadAlbumPagePort} from "./thunk-loadAlbumPage";

/**
 * Thunks related to catalog navigation.
 *
 * Expected handler types:
 * - `onAlbumFilterChange`: `(criterion: AlbumFilterCriterion) => void`
 * - `onPageRefresh`: `(albumId?: AlbumId) => Promise<void>`
 * - `loadAlbumPage`: `(albumId: AlbumId) => Promise<void>`
 */
export const navigationThunks = {
    onAlbumFilterChange: onAlbumFilterChangeDeclaration,
    onPageRefresh: onPageRefreshDeclaration,
    loadAlbumPage: loadAlbumPageDeclaration,
};
export * from "./selector-catalog-viewer-page";
export * from "./selector-albumListActions";
