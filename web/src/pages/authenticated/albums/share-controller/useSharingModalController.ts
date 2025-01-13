import {useMemo, useReducer} from "react";
import {ShareController} from "./ShareController";
import {ShareState, sharingModalReducer} from "./sharingModalReducer";
import {useCatalogAPIAdapter} from "../../../../core/catalog-react";

export function useSharingModalController(): ShareState & ShareController {
    const [state, dispatch] = useReducer(sharingModalReducer, {
        open: false,
        sharedWith: [],
    })
    const catalogAPIAdapter = useCatalogAPIAdapter()

    const ctrl = useMemo(() => new ShareController(dispatch, catalogAPIAdapter), [dispatch, catalogAPIAdapter])

    return {
        ...ctrl,
        ...state
    }
}

