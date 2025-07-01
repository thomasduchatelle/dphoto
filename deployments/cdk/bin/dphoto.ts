#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import {environments} from '../lib/config/environments';
import {DPhotoInfrastructureStack} from '../lib/stacks/dphoto-infrastructure-stack';

const app = new cdk.App();

const envName = app.node.tryGetContext('environment') || 'next';
const config = environments[envName];

if (!config) {
    throw new Error(`Unknown environment: ${envName}. Available: ${Object.keys(environments).join(', ')}`);
}

console.log(`Initializing CDK for environment: ${envName}`);
console.log(`Configuration:`, config);

// Create infrastructure stack
new DPhotoInfrastructureStack(app, `dphoto-${envName}-infra`, {
    environmentName: envName,
    config: config,
    env: {
        account: process.env.CDK_DEFAULT_ACCOUNT,
        region: process.env.CDK_DEFAULT_REGION || 'eu-west-1'
    },
    description: `DPhoto infrastructure stack for ${envName} environment`
});
