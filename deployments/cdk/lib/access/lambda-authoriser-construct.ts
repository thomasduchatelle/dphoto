import {HttpLambdaAuthorizer, HttpLambdaResponseType} from 'aws-cdk-lib/aws-apigatewayv2-authorizers';
import {Construct} from 'constructs';
import {GoLangLambdaFunction} from '../utils/golang-lambda-function';
import {Duration, Stack} from 'aws-cdk-lib';
import {CatalogStoreConstruct} from '../catalog/catalog-store-construct';
import {CognitoUserPoolConstruct} from './cognito-user-pool-construct';

export interface LambdaAuthoriserConstructProps {
    environmentName: string;
    catalogStore: CatalogStoreConstruct;
    cognitoUserPool: CognitoUserPoolConstruct;
}

export class LambdaAuthoriserConstruct extends Construct {
    public readonly authorizer: HttpLambdaAuthorizer;
    public readonly queryParamAuthorizer: HttpLambdaAuthorizer;
    private readonly authorizerLambda: GoLangLambdaFunction;

    constructor(scope: Construct, id: string, props: LambdaAuthoriserConstructProps) {
        super(scope, id);

        // Extract region from stack
        const region = Stack.of(this).region;

        // Construct Cognito JWKS URL
        const cognitoJwksUrl = `https://cognito-idp.${region}.amazonaws.com/${props.cognitoUserPool.userPool.userPoolId}/.well-known/openid-configuration`;

        // Create the Lambda function for the authorizer
        this.authorizerLambda = new GoLangLambdaFunction(this, 'AuthorizerLambda', {
            environmentName: props.environmentName,
            functionName: 'authorizer',
            timeout: Duration.seconds(10),
            memorySize: 256,
            environment: {
                COGNITO_JWKS_URL: cognitoJwksUrl,
            },
        });

        // Grant read access to catalog store (for permission checks)
        props.catalogStore.grantReadAccess(this.authorizerLambda);

        // Create the HTTP Lambda Authorizer - "identitySource" must all be present in the request or authorizer will not be called.
        this.authorizer = new HttpLambdaAuthorizer('LambdaAuthorizer', this.authorizerLambda.function, {
            authorizerName: `dphoto-${props.environmentName}-authorizer`,
            // identitySource: ['$request.header.Authorization', '$request.querystring.access_token'],
            // identitySource: ['$request.header.Cookie'],
            responseTypes: [HttpLambdaResponseType.SIMPLE],
            resultsCacheTtl: Duration.seconds(1800),
        });

        // Create a second authorizer for endpoints that use query parameter authentication
        // This is needed for get-media which passes the token as access_token query parameter
        this.queryParamAuthorizer = new HttpLambdaAuthorizer('LambdaAuthorizerQueryParam', this.authorizerLambda.function, {
            authorizerName: `dphoto-${props.environmentName}-authorizer-queryparam`,
            identitySource: ['$request.querystring.access_token'],
            responseTypes: [HttpLambdaResponseType.SIMPLE],
            resultsCacheTtl: Duration.seconds(1800),
        });
    }
}
