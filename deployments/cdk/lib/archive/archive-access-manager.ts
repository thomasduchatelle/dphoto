import {Workload} from "../utils/workload";

export interface ArchiveAccessManager {
    grantReadAccessToRawAndCacheMedias(workload: Workload): void;
}