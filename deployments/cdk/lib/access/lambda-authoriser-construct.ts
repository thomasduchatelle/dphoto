import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import {HttpLambdaAuthorizer, HttpLambdaResponseType} from 'aws-cdk-lib/aws-apigatewayv2-authorizers';
import {Construct} from 'constructs';
import {GoLangLambdaFunction} from '../utils/golang-lambda-function';
import {Duration} from 'aws-cdk-lib';
import {CatalogStoreConstruct} from '../catalog/catalog-store-construct';

export interface LambdaAuthoriserConstructProps {
    environmentName: string;
    catalogStore: CatalogStoreConstruct;
}

export class LambdaAuthoriserConstruct extends Construct {
    public readonly authorizer: HttpLambdaAuthorizer;
    private readonly authorizerLambda: GoLangLambdaFunction;

    constructor(scope: Construct, id: string, props: LambdaAuthoriserConstructProps) {
        super(scope, id);

        // Create the Lambda function for the authorizer
        this.authorizerLambda = new GoLangLambdaFunction(this, 'AuthorizerLambda', {
            environmentName: props.environmentName,
            functionName: 'lambda-authoriser',
            timeout: Duration.seconds(10),
            memorySize: 256,
        });

        // Grant read access to catalog store (for permission checks)
        props.catalogStore.grantReadAccess(this.authorizerLambda);

        // Create the HTTP Lambda Authorizer
        this.authorizer = new HttpLambdaAuthorizer('LambdaAuthorizer', this.authorizerLambda.function, {
            authorizerName: `dphoto-${props.environmentName}-lambda-authoriser`,
            identitySource: ['$request.header.Authorization', '$request.querystring.access_token'],
            responseTypes: [HttpLambdaResponseType.SIMPLE],
            resultsCacheTtl: Duration.seconds(300),
        });
    }
}
