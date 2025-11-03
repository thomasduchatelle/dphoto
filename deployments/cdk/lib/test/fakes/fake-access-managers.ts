import {ArchiveAccessManager} from '../../archive/archive-access-manager';
import {ArchivistAccessManager} from '../../archive/archivist-access-manager';
import {CatalogAccessManager} from '../../catalog/catalog-access-manager';
import {Workload} from '../../utils/workload';

function getLambdaName(workload: Workload): string {
    return workload.function?.functionName || 'not-a-lambda';
}

export class FakeArchiveAccessManager implements ArchiveAccessManager {
    private readAccess: Set<string> = new Set();

    grantReadAccessToRawAndCacheMedias(workload: Workload): void {
        this.readAccess.add(getLambdaName(workload));
    }

    hasBeenGrantedForRawAndCacheMedias(...lambdaNames: string[]): string {
        const missing = lambdaNames.filter(name => !this.readAccess.has(name));
        if (missing.length > 0) {
            return `Read access to raw/cache medias NOT granted for: ${missing.join(', ')}`;
        }
        return '';
    }
}

export class FakeArchivistAccessManager implements ArchivistAccessManager {
    private asyncArchivistAccess: Set<string> = new Set();

    grantAccessToAsyncArchivist(lambda: Workload): void {
        this.asyncArchivistAccess.add(getLambdaName(lambda));
    }

    hasBeenGrantedForAsyncArchivist(...lambdaNames: string[]): string {
        const missing = lambdaNames.filter(name => !this.asyncArchivistAccess.has(name));
        if (missing.length > 0) {
            return `Async archivist access NOT granted for: ${missing.join(', ')}`;
        }
        return '';
    }
}

export class FakeCatalogAccessManager implements CatalogAccessManager {
    private readAccess: Set<string> = new Set();
    private readWriteAccess: Set<string> = new Set();

    grantCatalogReadAccess(authorizerLambda: Workload): void {
        this.readAccess.add(getLambdaName(authorizerLambda));
    }

    grantCatalogReadWriteAccess(grantee: Workload): void {
        this.readWriteAccess.add(getLambdaName(grantee));
    }

    hasBeenGrantedForCatalogRead(...lambdaNames: string[]): string {
        const missing = lambdaNames.filter(name => !this.readAccess.has(name));
        if (missing.length > 0) {
            return `Catalog read access NOT granted for: ${missing.join(', ')}`;
        }
        return '';
    }

    hasOnlyBeenGrantedCatalogReadWriteTo(...lambdaNames: string[]): string {
        const granted = Array.from(this.readWriteAccess);
        const missing = lambdaNames.filter(name => !this.readWriteAccess.has(name));
        const extra = granted.filter(name => !lambdaNames.includes(name));
        if (missing.length > 0) {
            return `Missing read-write grant for: ${missing.join(', ')}`;
        }
        if (extra.length > 0) {
            return `Unexpected read-write grant for: ${extra.join(', ')}`;
        }
        return '';
    }
}
