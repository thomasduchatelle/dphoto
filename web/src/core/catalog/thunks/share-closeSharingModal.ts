import {catalogActions, CatalogViewerAction, CatalogViewerState} from "../domain";
import {ThunkDeclaration} from "../../thunk-engine";
import {CatalogFactoryArgs} from "./catalog-factory-args";

export function closeSharingModalThunk(dispatch: (action: CatalogViewerAction) => void): void {
    dispatch(catalogActions.closeSharingModalAction());
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
