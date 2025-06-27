import {openEditNameDialogDeclaration} from "./thunk-openEditNameDialog";
import {closeEditNameDialogDeclaration} from "./thunk-closeEditNameDialog";
import {changeAlbumNameDeclaration} from "./thunk-changeAlbumName";
import {changeFolderNameEnabledDeclaration} from "./thunk-changeFolderNameEnabled";
import {changeFolderNameDeclaration} from "./thunk-changeFolderName";
import {saveAlbumNameDeclaration} from "./thunk-saveAlbumName";

export * from "./selector-editNameDialogSelector";

export const albumEditNameThunks = {
    openEditNameDialog: openEditNameDialogDeclaration,
    closeEditNameDialog: closeEditNameDialogDeclaration,
    changeAlbumName: changeAlbumNameDeclaration,
    changeFolderNameEnabled: changeFolderNameEnabledDeclaration,
    changeFolderName: changeFolderNameDeclaration,
    saveAlbumName: saveAlbumNameDeclaration,
};
