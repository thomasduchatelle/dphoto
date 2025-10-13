import * as cdk from 'aws-cdk-lib';
import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import * as iam from 'aws-cdk-lib/aws-iam';
import {Construct} from 'constructs';
import {GoLangLambdaFunction} from '../utils/golang-lambda-function';
import {CognitoUserPoolConstruct} from './cognito-user-pool-construct';

export interface CognitoAuthorizerConstructProps {
    environmentName: string;
    httpApi: apigatewayv2.HttpApi;
    cognitoUserPool: CognitoUserPoolConstruct;
}

export class CognitoAuthorizerConstruct extends Construct {
    public readonly authorizer: apigatewayv2.CfnAuthorizer;
    public readonly authorizerLambda: GoLangLambdaFunction;

    constructor(scope: Construct, id: string, props: CognitoAuthorizerConstructProps) {
        super(scope, id);

        // Create the authorizer Lambda function
        this.authorizerLambda = new GoLangLambdaFunction(this, 'AuthorizerLambda', {
            environmentName: props.environmentName,
            functionName: 'cognito-authorizer',
            environment: {
                COGNITO_REGION: cdk.Stack.of(this).region,
                COGNITO_USER_POOL_ID: props.cognitoUserPool.userPool.userPoolId,
            },
        });

        // Create the API Gateway authorizer using L1 construct
        this.authorizer = new apigatewayv2.CfnAuthorizer(this, 'Authorizer', {
            apiId: props.httpApi.apiId,
            authorizerType: 'REQUEST',
            identitySource: ['$request.header.Authorization', '$request.header.Cookie'],
            name: `dphoto-${props.environmentName}-cognito-authorizer`,
            authorizerUri: `arn:aws:apigateway:${cdk.Stack.of(this).region}:lambda:path/2015-03-31/functions/${this.authorizerLambda.function.functionArn}/invocations`,
            authorizerPayloadFormatVersion: '2.0',
            enableSimpleResponses: true,
            authorizerResultTtlInSeconds: 3600, // 1 hour cache
        });

        // Grant API Gateway permission to invoke the authorizer Lambda
        this.authorizerLambda.function.grantInvoke(
            new iam.ServicePrincipal('apigateway.amazonaws.com', {
                conditions: {
                    ArnLike: {
                        'aws:SourceArn': `arn:aws:execute-api:${cdk.Stack.of(this).region}:${cdk.Stack.of(this).account}:${props.httpApi.apiId}/*`,
                    },
                },
            })
        );
    }
}
