import {createSimpleThunkDeclaration} from "src/libs/dthunks";
import {albumNameChanged} from "./action-albumNameChanged";

export const changeAlbumNameDeclaration = createSimpleThunkDeclaration(albumNameChanged);
