import {Album, CatalogViewerState} from "../language";
import {ThunkDeclaration} from "../../thunk-engine";
import {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {SharingModalOpened, sharingModalOpened} from "./action-sharingModalOpened";

export function openSharingModalThunk(dispatch: (action: SharingModalOpened) => void, album: Album): void {
    dispatch(sharingModalOpened(album.albumId));
}

export const openSharingModalDeclaration: ThunkDeclaration<
    CatalogViewerState,
    {},
    (album: Album) => void,
    CatalogFactoryArgs
> = {
    factory: ({dispatch}) => openSharingModalThunk.bind(null, dispatch),
    selector: (_state: CatalogViewerState) => ({}),
};
