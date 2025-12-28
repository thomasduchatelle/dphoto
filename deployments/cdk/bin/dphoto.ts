#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import {environments} from '../lib/config/environments';
import {InfrastructureStack} from '../lib/stacks/infrastructure-stack';
import {ApplicationStack} from "../lib/stacks/application-stack";
import {CognitoCertificateStack} from "../lib/stacks/cognito-certificate-stack";
import {CognitoStack} from "../lib/stacks/cognito-stack";
import {computeLetsEncryptHash} from "../lib/utils/letsencrypt-certificate-construct";
import {CognitoCustomDomainStack} from "../lib/stacks/cognito-custom-domain-stack";

export default async function main(
    defaultEnvName: string = "next",
    account: string | undefined = process.env.CDK_DEFAULT_ACCOUNT,
    region: string | undefined = process.env.CDK_DEFAULT_REGION || 'eu-west-1',
) {
    await computeLetsEncryptHash()

    const app = new cdk.App();

    const environmentName = app.node.tryGetContext('environment') || defaultEnvName
    const config = environments[environmentName];

    if (!config) {
        throw new Error(`Unknown environment: ${environmentName}. Available: ${Object.keys(environments).join(', ')}`);
    }

    console.log(`Initializing CDK for environment: ${environmentName}`);
    console.log(`Configuration:`, config);

    // Infrastructure Stack has all the data that must be retained and secured (buckets with the medias, storages, queues, ...)
    const infrastructureStack = new InfrastructureStack(app, `dphoto-${environmentName}-infra`, {
        description: `DPhoto ${environmentName}: persistent infrastructure layer`,
        environmentName,
        config,
        env: {
            account: account,
            region: region
        },
    });

    // Cognito Stack has all authentication resources (user pool, client, custom domain with managed UI)
    const cognitoStack = new CognitoStack(app, `dphoto-${environmentName}-cognito`, {
        description: `DPhoto ${environmentName}: create Cognito User Pool and Client (semi-persistent)`,
        environmentName: environmentName,
        config: config,
        env: {
            account: account,
            region: region
        },
    });

    // Application Stack has everything else, it can be destroyed and recreated at any time (gateway, workload deployments, UI, ...)
    const applicationStack = new ApplicationStack(app, `dphoto-${environmentName}-application`, {
        description: `DPhoto ${environmentName}: core APIs and business logic management (stateless)`,
        environmentName,
        config,
        archiveAccessManager: infrastructureStack.archiveStore,
        catalogAccessManager: infrastructureStack.catalogStore,
        archivistAccessManager: infrastructureStack.archivist,
        oauth2ClientConfig: cognitoStack.getWebEnvironmentVariables(),
        env: {
            account: account,
            region: region
        },
    });
    applicationStack.addDependency(infrastructureStack);
    applicationStack.addDependency(cognitoStack);

    // CognitoCertificateStack creates certificates required by CognitoCustomDomainStack in us-east-1 region.
    const cognitoCertificateStack = new CognitoCertificateStack(app, `dphoto-${environmentName}-cognito-cert`, {
        description: `DPhoto ${environmentName}: companion stack of CognitoCustomDomainStack to create certificate in us-east-1`,
        environmentName: environmentName,
        config: config,
        env: {
            account: account,
            region: 'us-east-1'
        },
    });

    // CognitoCustomDomainStack configures the custom domain for Cognito User Pool. It is SLOW to delete and create (hence being created independently). A ARecord on the parent domain is required so it must be created after the API Gateway from ApplicationStack.
    const cognitoCustomDomainStack = new CognitoCustomDomainStack(app, `dphoto-${environmentName}-cognito-domain`, {
        description: `DPhoto ${environmentName}: configure Cognito custom domain (slow to create and delete)`,
        crossRegionReferences: true,
        userPool: cognitoStack.userPool,
        cognitoDomainName: config.cognitoDomainName,
        rootDomain: config.rootDomain,
        cognitoCertificate: cognitoCertificateStack.cognitoCertificate,
        env: {
            account: account,
            region: region,
        },
    })
    cognitoCustomDomainStack.addDependency(cognitoCertificateStack);
    cognitoCustomDomainStack.addDependency(cognitoStack);

    return app;
}

if (typeof jest === 'undefined' && require.main === module) {
    main();
}
