import * as cdk from 'aws-cdk-lib';
import {Duration} from 'aws-cdk-lib';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as logs from 'aws-cdk-lib/aws-logs';
import {Construct} from 'constructs';

export interface LambdaPermissions {
    cacheRw?: boolean;
    dynamodbRw?: boolean;
    storageRw?: boolean;
    storageRo?: boolean;
}

export interface GoLangLambdaFunctionProps {
    environmentName: string;
    functionName: string;
    artifactPath?: string;
    timeout?: cdk.Duration;
    memorySize?: number;
    environment?: Record<string, string>;
}

export class GoLangLambdaFunction extends Construct {
    public readonly function: lambda.Function;
    public readonly role: iam.Role;

    constructor(scope: Construct, id: string, props: GoLangLambdaFunctionProps) {
        super(scope, id);

        this.role = new iam.Role(this, 'Role', {
            roleName: `dphoto-${props.environmentName}-${props.functionName}-role`,
            path: `/dphoto/${props.environmentName}/`,
            assumedBy: new iam.ServicePrincipal('lambda.amazonaws.com'),
            managedPolicies: [
                iam.ManagedPolicy.fromAwsManagedPolicyName('service-role/AWSLambdaBasicExecutionRole')
            ]
        });

        const logGroup = new logs.LogGroup(this, 'LogGroup', {
            logGroupName: `/dphoto/${props.environmentName}/lambda/${props.functionName}`,
            retention: logs.RetentionDays.ONE_WEEK,
            removalPolicy: cdk.RemovalPolicy.DESTROY
        });

        this.function = new lambda.Function(this, 'Function', {
            functionName: `dphoto-${props.environmentName}-${props.functionName}`,
            runtime: lambda.Runtime.PROVIDED_AL2,
            architecture: lambda.Architecture.ARM_64,
            handler: 'bootstrap',
            code: lambda.Code.fromAsset(props.artifactPath || `../../bin/${props.functionName}.zip`),
            timeout: props.timeout || Duration.minutes(1),
            memorySize: props.memorySize || 256,
            environment: props.environment || {},
            logGroup: logGroup,
            role: this.role
        });
    }
}
