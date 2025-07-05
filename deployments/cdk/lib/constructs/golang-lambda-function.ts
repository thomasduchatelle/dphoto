import * as cdk from 'aws-cdk-lib';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import * as apigatewayv2_integrations from 'aws-cdk-lib/aws-apigatewayv2-integrations';
import * as logs from 'aws-cdk-lib/aws-logs';
import {Construct} from 'constructs';

export interface GoLangLambdaFunctionProps {
    environmentName: string;
    functionName: string;
    artifactPath: string;
    timeout?: cdk.Duration;
    memorySize?: number;
    environment?: Record<string, string>;
}

export class GoLangLambdaFunction extends Construct {
    public readonly function: lambda.Function;

    constructor(scope: Construct, id: string, props: GoLangLambdaFunctionProps) {
        super(scope, id);

        this.function = new lambda.Function(this, 'Function', {
            functionName: `dphoto-${props.environmentName}-${props.functionName}`,
            runtime: lambda.Runtime.PROVIDED_AL2,
            architecture: lambda.Architecture.ARM_64,
            handler: 'bootstrap',
            code: lambda.Code.fromAsset(props.artifactPath),
            timeout: props.timeout || cdk.Duration.seconds(30),
            memorySize: props.memorySize || 256,
            environment: props.environment || {},
            logRetention: logs.RetentionDays.ONE_WEEK
        });
    }

    public addToApiGateway(
        httpApi: apigatewayv2.HttpApi,
        path: string,
        method: apigatewayv2.HttpMethod
    ): void {
        const integration = new apigatewayv2_integrations.HttpLambdaIntegration(
            `${this.node.id}Integration`,
            this.function
        );

        new apigatewayv2.HttpRoute(this, `${this.node.id}Route`, {
            httpApi,
            routeKey: apigatewayv2.HttpRouteKey.with(path, method),
            integration
        });
    }
}
