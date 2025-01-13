import {useApplication} from "../application";
import {CatalogAPIAdapter} from "../catalog/adapters/api";

export const useCatalogAPIAdapter = () => {
    const app = useApplication()
    return new CatalogAPIAdapter(app.axiosInstance, app)
}