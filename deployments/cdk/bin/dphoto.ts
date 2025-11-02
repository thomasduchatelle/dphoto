#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import {environments} from '../lib/config/environments';
import {InfrastructureStack} from '../lib/stacks/infrastructure-stack';
import {ApplicationStack} from "../lib/stacks/application-stack";
import {CognitoCertificateStack} from "../lib/stacks/cognito-certificate-stack";
import {CognitoStack} from "../lib/stacks/cognito-stack";
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

    // Create certificate in us-east-1 which is required for Cognito custom domain.
    const cognitoCertificateStack = new CognitoCertificateStack(app, `dphoto-${envName}-cognito-cert`, {
        environmentName: envName,
        config: config,
        env: {
            account: account,
            region: 'us-east-1'
        },
        description: `Create certificate in us-east-1 which is required for Cognito custom domain.`
    });

    // Cognito Stack has all authentication resources (user pool, client, custom domain with managed UI)
    const cognitoStack = new CognitoStack(app, `dphoto-${envName}-cognito`, {
        environmentName: envName,
        config: config,
        cognitoCertificate: cognitoCertificateStack.cognitoCertificate,
        env: {
            account: account,
            region: region
        },
        description: `DPhoto Cognito stack for ${envName} environment`
    });
    cognitoStack.addDependency(cognitoCertificateStack);

    // Infrastructure Stack has everything else, it can be destroyed and recreated at any time (gateway, workload deployments, UI, ...)
    const applicationStack = new ApplicationStack(app, `dphoto-${envName}-application`, {
        environmentName: envName,
        config,
        archiveAccessManager: infrastructureStack.archiveStore,
        catalogAccessManager: infrastructureStack.catalogStore,
        archivistAccessManager: infrastructureStack.archivist,
        oauth2ClientConfig: cognitoStack.getWebEnvironmentVariables(),
        env: {
            account: account,
            region: region
        }
    });

    applicationStack.addDependency(infrastructureStack);
    applicationStack.addDependency(cognitoStack);

    return app;
}

if (typeof jest === 'undefined' && require.main === module) {
    main();
}
