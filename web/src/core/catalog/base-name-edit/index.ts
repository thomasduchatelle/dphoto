import {changeAlbumNameDeclaration} from "./thunk-changeAlbumName";
import {changeFolderNameEnabledDeclaration} from "./thunk-changeFolderNameEnabled";
import {changeFolderNameDeclaration} from "./thunk-changeFolderName";

export const baseNameEditThunks = {
    changeAlbumName: changeAlbumNameDeclaration,
    changeFolderNameEnabled: changeFolderNameEnabledDeclaration,
    changeFolderName: changeFolderNameDeclaration,
};