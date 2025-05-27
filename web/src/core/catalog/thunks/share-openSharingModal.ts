import {Album, catalogActions, CatalogViewerAction, CatalogViewerState} from "../domain";
import {ThunkDeclaration} from "../../thunk-engine";
import {CatalogFactoryArgs} from "./catalog-factory-args";

export function openSharingModalThunk(dispatch: (action: CatalogViewerAction) => void, album: Album): void {
    dispatch(catalogActions.openSharingModalAction({albumId: album.albumId}));
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
