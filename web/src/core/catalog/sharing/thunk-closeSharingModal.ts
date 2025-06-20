import {CatalogViewerState} from "../language";
import {ThunkDeclaration} from "../../thunk-engine";
import {CatalogFactoryArgs} from "../common/catalog-factory-args";
import {SharingModalClosed, sharingModalClosed} from "./action-sharingModalClosed";

export function closeSharingModalThunk(dispatch: (action: SharingModalClosed) => void): void {
    dispatch(sharingModalClosed());
}

export const closeSharingModalDeclaration: ThunkDeclaration<
    CatalogViewerState,
    {},
    () => void,
    CatalogFactoryArgs
> = {
    factory: ({dispatch}) => closeSharingModalThunk.bind(null, dispatch),
    selector: (_state: CatalogViewerState) => ({}),
};
