import {deleteAlbumDialogOpened} from "./action-deleteAlbumDialogOpened";
import {createSimpleThunkDeclaration} from "src/libs/dthunks";

export const openDeleteAlbumDialogDeclaration = createSimpleThunkDeclaration(deleteAlbumDialogOpened);
