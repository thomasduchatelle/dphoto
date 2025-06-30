#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import {environments} from '../lib/config/environments';
import {CdkStack} from "../lib/cdk-stack";

const app = new cdk.App();

const envName = app.node.tryGetContext('environment') || 'dev';
const config = environments[envName];

if (!config) {
    throw new Error(`Unknown environment: ${envName}. Available: ${Object.keys(environments).join(', ')}`);
}

// Placeholder stacks - will be implemented in Milestone 2
console.log(`Initializing CDK for environment: ${envName}`);
console.log(`Configuration:`, config);

new CdkStack(app, 'CdkStack', {
    env: {
        account: process.env.CDK_DEFAULT_ACCOUNT,
        region: process.env.CDK_DEFAULT_REGION || 'eu-west-1'
    }
});
