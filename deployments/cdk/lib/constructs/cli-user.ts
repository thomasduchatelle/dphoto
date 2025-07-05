import * as iam from 'aws-cdk-lib/aws-iam';
import {Construct} from 'constructs';

export interface CliUserProps {
    environmentName: string;
    cliAccessKeys: string[];
    storageRwPolicyArn: string;
    cacheRwPolicyArn: string;
    indexRwPolicyArn: string;
    archiveSnsPublishPolicyArn: string;
    archiveSqsSendPolicyArn: string;
    archiveRelocatePolicyArn: string;
}

export class CliUserConstruct extends Construct {
    public readonly user: iam.User;
    public readonly accessKeys: Record<string, iam.AccessKey>;

    constructor(scope: Construct, id: string, props: CliUserProps) {
        super(scope, id);

        const prefix = `dphoto-${props.environmentName}`;

        // IAM User
        this.user = new iam.User(this, 'CliUser', {
            userName: `${prefix}-cli`,
            path: '/dphoto/'
        });

        // Attach policies
        this.user.addManagedPolicy(iam.ManagedPolicy.fromManagedPolicyArn(this, 'StorageRwPolicy', props.storageRwPolicyArn));
        this.user.addManagedPolicy(iam.ManagedPolicy.fromManagedPolicyArn(this, 'CacheRwPolicy', props.cacheRwPolicyArn));
        this.user.addManagedPolicy(iam.ManagedPolicy.fromManagedPolicyArn(this, 'IndexRwPolicy', props.indexRwPolicyArn));
        this.user.addManagedPolicy(iam.ManagedPolicy.fromManagedPolicyArn(this, 'ArchiveSnsPublishPolicy', props.archiveSnsPublishPolicyArn));
        this.user.addManagedPolicy(iam.ManagedPolicy.fromManagedPolicyArn(this, 'ArchiveSqsSendPolicy', props.archiveSqsSendPolicyArn));
        this.user.addManagedPolicy(iam.ManagedPolicy.fromManagedPolicyArn(this, 'ArchiveRelocatePolicy', props.archiveRelocatePolicyArn));

        // Create access keys
        this.accessKeys = {};
        props.cliAccessKeys.forEach(keyDate => {
            this.accessKeys[keyDate] = new iam.AccessKey(this, `AccessKey${keyDate.replace('-', '')}`, {
                user: this.user
            });
        });
    }
}
