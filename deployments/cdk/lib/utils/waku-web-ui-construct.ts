import * as cdk from 'aws-cdk-lib';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import {NodejsFunction} from 'aws-cdk-lib/aws-lambda-nodejs';
import {Construct} from 'constructs';
import * as apigatewayv2 from "aws-cdk-lib/aws-apigatewayv2";
import * as apigatewayv2_integrations from "aws-cdk-lib/aws-apigatewayv2-integrations";
import {HttpLambdaIntegration} from "aws-cdk-lib/aws-apigatewayv2-integrations";

export interface WakuWebUiConstructProps {
    environmentName: string;
    httpApi: apigatewayv2.HttpApi;
}

export class WakuWebUiConstruct extends Construct {
    private readonly lambda: NodejsFunction;
    private readonly integration: HttpLambdaIntegration;

    constructor(scope: Construct, id: string, {httpApi}: WakuWebUiConstructProps) {
        super(scope, id);

        this.lambda = new NodejsFunction(this, 'Lambda', {
            entry: '../../web-waku/dist/serve-aws-lambda.js',
            handler: 'handler',
            runtime: lambda.Runtime.NODEJS_20_X,
            memorySize: 512,
            timeout: cdk.Duration.seconds(10),
            environment: {
                NODE_ENV: 'production',
            },
        });


        this.integration = new apigatewayv2_integrations.HttpLambdaIntegration(
            `${this.node.id}Integration`,
            this.lambda,
        );

        new apigatewayv2.HttpRoute(this, 'Route', {
            httpApi,
            routeKey: apigatewayv2.HttpRouteKey.with('/waku', apigatewayv2.HttpMethod.ANY),
            integration: this.integration
        });
    }
}
