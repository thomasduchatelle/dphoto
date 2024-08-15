import {useMemo, useReducer} from "react";
import {ShareController} from "./ShareController";
import {ShareState, sharingModalReducer} from "./sharingModalReducer";
import {useApplication} from "../../../../core/application";
import {CatalogAPIAdapter} from "../../../../core/catalog/adapters/api";

export function useSharingModalController(): ShareState & ShareController {
    const [state, dispatch] = useReducer(sharingModalReducer, {
        open: false,
        sharedWith: [],
    })
    const app = useApplication()

    const ctrl = useMemo(() => new ShareController(dispatch, new CatalogAPIAdapter(app.axiosInstance, app)), [dispatch, app])

    return {
        ...ctrl,
        ...state
    }
}

