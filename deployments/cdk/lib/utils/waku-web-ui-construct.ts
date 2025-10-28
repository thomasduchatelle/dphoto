import * as cdk from 'aws-cdk-lib';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import {Construct} from 'constructs';
import * as apigatewayv2 from "aws-cdk-lib/aws-apigatewayv2";
import * as apigatewayv2_integrations from "aws-cdk-lib/aws-apigatewayv2-integrations";
import {HttpLambdaIntegration} from "aws-cdk-lib/aws-apigatewayv2-integrations";
import * as logs from "aws-cdk-lib/aws-logs";
import * as cognito from 'aws-cdk-lib/aws-cognito';

export interface WakuWebUiConstructProps {
    environmentName: string;
    httpApi: apigatewayv2.HttpApi;
    userPool: cognito.IUserPool;
    userPoolClient: cognito.UserPoolClient;
    cognitoDomainName: string;
    googleLoginClientId: string;
}

export class WakuWebUiConstruct extends Construct {
    private readonly lambda: lambda.Function;
    private readonly integration: HttpLambdaIntegration;

    constructor(scope: Construct, id: string, {httpApi, environmentName, userPool, userPoolClient, cognitoDomainName, googleLoginClientId}: WakuWebUiConstructProps) {
        super(scope, id);

        const logGroup = new logs.LogGroup(this, 'LogGroup', {
            logGroupName: `/aws/lambda/dphoto-${environmentName}-web`,
            retention: logs.RetentionDays.ONE_WEEK,
            removalPolicy: cdk.RemovalPolicy.DESTROY
        });

        this.lambda = new lambda.Function(this, 'Lambda', {
            functionName: `dphoto-${environmentName}-web`,
            code: lambda.Code.fromAsset('../../web/dist/'),
            handler: 'serve-aws-lambda.handler',
            runtime: lambda.Runtime.NODEJS_20_X,
            memorySize: 256,
            timeout: cdk.Duration.seconds(10),
            logGroup: logGroup,
            environment: {
                NODE_ENV: 'production',
                COGNITO_USER_POOL_ID: userPool.userPoolId,
                COGNITO_CLIENT_ID: userPoolClient.userPoolClientId,
                COGNITO_CLIENT_SECRET: userPoolClient.userPoolClientSecret.unsafeUnwrap(),
                COGNITO_DOMAIN: `https://${cognitoDomainName}`,
                COGNITO_ISSUER: `https://cognito-idp.${cdk.Stack.of(this).region}.amazonaws.com/${userPool.userPoolId}`,
                GOOGLE_LOGIN_CLIENT_ID: googleLoginClientId,
            },
        });

        this.integration = new apigatewayv2_integrations.HttpLambdaIntegration(
            `WakuIntegration`,
            this.lambda,
        );

        new apigatewayv2.HttpRoute(this, 'RouteProxy', {
            httpApi,
            routeKey: apigatewayv2.HttpRouteKey.with('/{proxy+}', apigatewayv2.HttpMethod.ANY),
            integration: this.integration
        });

        new apigatewayv2.HttpRoute(this, 'Route', {
            httpApi,
            routeKey: apigatewayv2.HttpRouteKey.DEFAULT,
            integration: this.integration
        });
    }
}
