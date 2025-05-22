import {useApplication} from "../application";
import {CatalogAPIAdapter} from "../catalog/adapters/api";
import {useMemo} from "react";

export const useCatalogAPIAdapter = () => {
    const app = useApplication()
    return useMemo(() => new CatalogAPIAdapter(app.axiosInstance, app), [app])
}