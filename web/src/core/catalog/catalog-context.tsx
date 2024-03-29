import {CatalogState} from "./catalog-model";
import {createContext, Dispatch, ReactNode, useReducer} from "react";
import {CatalogAction, catalogReducer, initialCatalogState} from "./catalog-reducer";

export interface CatalogStateWithDispatch {
    catalog: CatalogState
    dispatch: Dispatch<CatalogAction>
}

export const CatalogContext = createContext<CatalogStateWithDispatch>({
    catalog: initialCatalogState, dispatch: () => {
    }
})

export const CatalogContextComponent = ({children}: {
    children?: ReactNode
}) => {
    const [catalog, dispatch] = useReducer(catalogReducer, initialCatalogState)

    return (
        <CatalogContext.Provider value={{catalog, dispatch}}>
            {children}
        </CatalogContext.Provider>
    )
}
