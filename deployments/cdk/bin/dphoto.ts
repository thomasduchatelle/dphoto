#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import {environments} from '../lib/config/environments';
import {InfrastructureStack} from '../lib/stacks/infrastructure-stack';
import {ApplicationStack} from "../lib/stacks/application-stack";
import {CognitoCertificateStack} from "../lib/stacks/cognito-certificate-stack";
import {CloudFrontCertificateStack} from "../lib/stacks/cloudfront-certificate-stack";
import {CognitoStack} from "../lib/stacks/cognito-stack";
import {computeLetsEncryptHash} from "../lib/utils/letsencrypt-certificate-construct";
import {CognitoCustomDomainStack} from "../lib/access/CognitoCustomDomainStack";
import {AppRouterStack} from "../lib/stacks/NextjsStack";

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

    // Cognito Stack has all authentication resources (user pool, client, custom domain with managed UI)
    const cognitoStack = new CognitoStack(app, `dphoto-${envName}-cognito`, {
        environmentName: envName,
        config: config,
        env: {
            account: account,
            region: region
        },
        description: `DPhoto Cognito stack for ${envName} environment`
    });

    // Create a custom domain for Cognito User Pool. It requires the app to be provisioned (ARecord on the parent domain), and the certificate to be provisioned in US-EAST-1.
    const cognitoCertificateStack = new CognitoCertificateStack(app, `dphoto-${envName}-cognito-cert`, {
        environmentName: envName,
        config: config,
        env: {
            account: account,
            region: 'us-east-1'
        },
        description: `Create certificate in us-east-1 which is required for Cognito custom domain.`
    });

    // Create certificate in us-east-1 which is required for CloudFront custom domain.
    const cloudFrontCertificateStack = new CloudFrontCertificateStack(app, `dphoto-${envName}-cloudfront-cert`, {
        environmentName: envName,
        config: config,
        env: {
            account: account,
            region: 'us-east-1'
        },
        description: `Create certificate in us-east-1 which is required for CloudFront custom domain.`
    });

    // Infrastructure Stack has everything else, it can be destroyed and recreated at any time (gateway, workload deployments, UI, ...)
    const applicationStack = new ApplicationStack(app, `dphoto-${envName}-application`, {
        environmentName: envName,
        config,
        archiveAccessManager: infrastructureStack.archiveStore,
        catalogAccessManager: infrastructureStack.catalogStore,
        archivistAccessManager: infrastructureStack.archivist,
        oauth2ClientConfig: cognitoStack.getWebEnvironmentVariables(),
        cloudFrontCertificate: cloudFrontCertificateStack.cloudFrontCertificate,
        env: {
            account: account,
            region: region
        }
    });

    applicationStack.addDependency(infrastructureStack);
    applicationStack.addDependency(cognitoStack);
    applicationStack.addDependency(cloudFrontCertificateStack);

    const cognitoCustomDomainStack = new CognitoCustomDomainStack(app, `dphoto-${envName}-cognito-domain`, {
        userPool: cognitoStack.userPool,
        cognitoDomainName: config.cognitoDomainName,
        rootDomain: config.rootDomain,
        cognitoCertificate: cognitoCertificateStack.cognitoCertificate,
        env: {
            account: account,
            region: region,
        },
        crossRegionReferences: true,
    })
    cognitoCustomDomainStack.addDependency(cognitoCertificateStack);
    cognitoCustomDomainStack.addDependency(applicationStack);

    new AppRouterStack(app, `dphoto-${envName}-nextjs`, {
        env: {
            account: account,
            region: region,
        },
    })

    return app;
}

if (typeof jest === 'undefined' && require.main === module) {
    main();
}
