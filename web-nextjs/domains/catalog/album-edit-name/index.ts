import {openEditNameDialogDeclaration} from "./thunk-openEditNameDialog";
import {closeEditNameDialogDeclaration} from "./thunk-closeEditNameDialog";
import {saveAlbumNameDeclaration} from "./thunk-saveAlbumName";

export * from "./selector-editNameDialogSelector";

export const albumEditNameThunks = {
    openEditNameDialog: openEditNameDialogDeclaration,
    closeEditNameDialog: closeEditNameDialogDeclaration,
    saveAlbumName: saveAlbumNameDeclaration,
};
