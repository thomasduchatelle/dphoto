#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import {environments} from '../lib/config/environments';
import {InfrastructureStack} from '../lib/stacks/infrastructure-stack';
import {ApplicationStack} from "../lib/stacks/application-stack";
import {CognitoCertificateStack} from "../lib/stacks/cognito-certificate-stack";
import {computeLetsEncryptHash} from "../lib/utils/letsencrypt-certificate-construct";

export default async function main(
    defaultEnvName: string = "next",
    account: string | undefined = process.env.CDK_DEFAULT_ACCOUNT,
    region: string | undefined = process.env.CDK_DEFAULT_REGION || 'eu-west-1',
) {
    await computeLetsEncryptHash()

    const app = new cdk.App();

    const envName = app.node.tryGetContext('environment') || defaultEnvName
    const config = environments[envName];

    if (!config) {
        throw new Error(`Unknown environment: ${envName}. Available: ${Object.keys(environments).join(', ')}`);
    }

    console.log(`Initializing CDK for environment: ${envName}`);
    console.log(`Configuration:`, config);

    // Infrastructure Stack has all the data that must be retained and secured (buckets with the medias, storages, queuses, ...)
    const infrastructureStack = new InfrastructureStack(app, `dphoto-${envName}-infra`, {
        environmentName: envName,
        config: config,
        env: {
            account: account,
            region: region
        },
        description: `DPhoto infrastructure stack for ${envName} environment`
    });

    // Stack required by cognito hosted UI, it installs a certificate in us-east-1
    const cognitoCertificateStack = new CognitoCertificateStack(app, `dphoto-${envName}-cognito-cert`, {
        environmentName: envName,
        config: config,
        env: {
            account: account,
            region: 'us-east-1'
        },
        description: `DPhoto Cognito certificate stack for ${envName} environment (us-east-1)`
    });

    // Infrastructure Stack has everything else, it can be destroyed and recreated at any time (gateway, workload deployments, UI, ...)
    const applicationStack = new ApplicationStack(app, `dphoto-${envName}-application`, {
        environmentName: envName,
        config,
        archiveStore: infrastructureStack.archiveStore,
        catalogStore: infrastructureStack.catalogStore,
        archivist: infrastructureStack.archivist,
        cognitoUserPool: infrastructureStack.cognitoUserPool,
        cognitoCertificate: cognitoCertificateStack.cognitoCertificate,
        env: {
            account: account,
            region: region
        }
    });

    applicationStack.addDependency(infrastructureStack);
    applicationStack.addDependency(cognitoCertificateStack);

    return app;
}

if (typeof jest === 'undefined' && require.main === module) {
    main();
}
