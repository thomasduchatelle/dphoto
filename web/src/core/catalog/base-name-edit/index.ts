import {changeAlbumNameDeclaration} from "./thunk-changeAlbumName";
import {changeFolderNameEnabledDeclaration} from "./thunk-changeFolderNameEnabled";
import {changeFolderNameDeclaration} from "./thunk-changeFolderName";

/**
 * Catalog's base-name-edit feature exposes the handlers:
 *
 * - `changeAlbumName`: `(albumName: string) => void`
 * - `changeFolderNameEnabled`: `(isFolderNameEnabled: boolean) => void`
 * - `changeFolderName`: `(folderName: string) => void`
 */
export const baseNameEditThunks = {
    changeAlbumName: changeAlbumNameDeclaration,
    changeFolderNameEnabled: changeFolderNameEnabledDeclaration,
    changeFolderName: changeFolderNameDeclaration,
};
