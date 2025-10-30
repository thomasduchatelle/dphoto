import * as iam from 'aws-cdk-lib/aws-iam';
import {Construct} from 'constructs';
import {pinLogicalId} from '../utils/override-logical-ids';
import {ArchiveStoreConstruct} from "../archive/archive-store-construct";
import {CatalogStoreConstruct} from "../catalog/catalog-store-construct";
import {ArchivistConstruct} from "../archive/archivist-construct";

export interface AccessStoreConstructProps {
    environmentName: string;
    cliAccessKeys: string[];
    archiveStore: ArchiveStoreConstruct
    catalogStore: CatalogStoreConstruct;
    archivist: ArchivistConstruct;
}

export class CliUserAccessConstruct extends Construct {
    public readonly user: iam.User;
    public readonly accessKeys: Record<string, iam.AccessKey>;

    constructor(scope: Construct, id: string, props: AccessStoreConstructProps) {
        super(scope, id);

        const prefix = `dphoto-${props.environmentName}`;

        // IAM User
        this.user = new iam.User(this, 'CliUser', {
            userName: `${prefix}-cli`,
            path: '/dphoto/'
        });
        pinLogicalId(this.user, "CliUserA7F35037");

        // Attach policies
        props.catalogStore.grantCatalogReadWriteAccess({role: this.user});
        props.archiveStore.grantWriteAccessToRawAndCachedMedias({role: this.user});
        props.archivist.grantAccessToAsyncArchivist({role: this.user})

        // Create access keys
        this.accessKeys = {};
        props.cliAccessKeys.forEach(keyDate => {
            this.accessKeys[keyDate] = new iam.AccessKey(this, `AccessKey${keyDate.replace('-', '')}`, {
                user: this.user
            });
            if (keyDate === "2025-07" || keyDate === '2024-04') {
                pinLogicalId(this.accessKeys[keyDate], "CliUserAccessKey202507148E1156");
            }
        });
    }
}
