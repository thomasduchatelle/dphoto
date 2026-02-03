import {deleteAlbumDialogOpened} from "./action-deleteAlbumDialogOpened";
import {createSimpleThunkDeclaration} from "@/libs/dthunks";

export const openDeleteAlbumDialogDeclaration = createSimpleThunkDeclaration(deleteAlbumDialogOpened);
