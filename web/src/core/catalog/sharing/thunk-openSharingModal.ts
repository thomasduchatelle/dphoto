import {Album, CatalogViewerState} from "../language";
import {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {SharingModalOpened, sharingModalOpened} from "./action-sharingModalOpened";
import {createSimpleThunkDeclaration, ThunkDeclaration} from "src/libs/dthunks";

export function openSharingModalThunk(dispatch: (action: SharingModalOpened) => void, album: Album): void {
    dispatch(sharingModalOpened(album.albumId));
}

export const openSharingModalDeclaration = createSimpleThunkDeclaration(sharingModalOpened);
