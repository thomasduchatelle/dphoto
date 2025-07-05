#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import {environments} from '../lib/config/environments';
import {DPhotoInfrastructureStack} from '../lib/stacks/dphoto-infrastructure-stack';
import {DPhotoApplicationStack} from "../lib/stacks/dphoto-application-stack";

export default function main(
    defaultEnvName: string = "next",
    account: string | undefined = process.env.CDK_DEFAULT_ACCOUNT,
    region: string | undefined = process.env.CDK_DEFAULT_REGION || 'eu-west-1',
) {
    const app = new cdk.App();

    const envName = app.node.tryGetContext('environment') || defaultEnvName
    const config = environments[envName];

    if (!config) {
        throw new Error(`Unknown environment: ${envName}. Available: ${Object.keys(environments).join(', ')}`);
    }

    console.log(`Initializing CDK for environment: ${envName}`);
    console.log(`Configuration:`, config);

// Create infrastructure stack
    const infrastructureStack = new DPhotoInfrastructureStack(app, `dphoto-${envName}-infra`, {
        environmentName: envName,
        config: config,
        env: {
            account: account,
            region: region
        },
        description: `DPhoto infrastructure stack for ${envName} environment`
    });

    const applicationStack = new DPhotoApplicationStack(app, `dphoto-${envName}-application`, {
        environmentName: envName,
        config,
        env: {
            account: account,
            region: region
        }
    });

    applicationStack.addDependency(infrastructureStack);

    return app;
}

if (typeof jest === 'undefined' && require.main === module) {
    main();
}