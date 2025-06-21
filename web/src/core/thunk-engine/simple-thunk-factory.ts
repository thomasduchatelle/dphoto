import {ThunkDeclaration} from "../thunk-engine";
import {CatalogFactoryArgs} from "../catalog/common/catalog-factory-args";
import {CatalogViewerState} from "../catalog/language";

type ActionCreator<TAction, TArgs extends any[]> = (...args: TArgs) => TAction;

export function createSimpleThunkDeclaration<TAction, TArgs extends any[]>(
    actionCreator: ActionCreator<TAction, TArgs>
): ThunkDeclaration<CatalogViewerState, {}, (...args: TArgs) => void, CatalogFactoryArgs> {
    return {
        selector: () => ({}),
        factory: ({dispatch}) => {
            return (...args: TArgs) => {
                dispatch(actionCreator(...args));
            };
        },
    };
}
