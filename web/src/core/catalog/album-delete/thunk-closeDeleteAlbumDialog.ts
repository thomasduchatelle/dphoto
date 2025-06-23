import {deleteAlbumDialogClosed} from "./action-deleteAlbumDialogClosed";
import {createSimpleThunkDeclaration} from "src/libs/dthunks";

export const closeDeleteAlbumDialogDeclaration = createSimpleThunkDeclaration(deleteAlbumDialogClosed);
