import {deleteAlbumDialogClosed} from "./action-deleteAlbumDialogClosed";
import {createSimpleThunkDeclaration} from "@/libs/dthunks";

export const closeDeleteAlbumDialogDeclaration = createSimpleThunkDeclaration(deleteAlbumDialogClosed);
