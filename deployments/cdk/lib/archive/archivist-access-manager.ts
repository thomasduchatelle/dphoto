import {Workload} from "../utils/workload";

export interface ArchivistAccessManager {

    grantAccessToAsyncArchivist(lambda: Workload): void;
}