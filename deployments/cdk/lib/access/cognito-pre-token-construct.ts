import {Construct} from 'constructs';
import {Duration} from 'aws-cdk-lib';
import * as cognito from 'aws-cdk-lib/aws-cognito';
import {GoLangLambdaFunction} from '../utils/golang-lambda-function';
import {CatalogAccessManager} from '../catalog/catalog-access-manager';

export interface CognitoPreTokenConstructProps {
    environmentName: string;
    userPool: cognito.UserPool;
    catalogStore: CatalogAccessManager;
}

/**
 * CognitoPreTokenConstruct attaches the Pre Token Generation (V2) trigger to the Cognito User Pool. The Lambda
 * resolves the signing-in user's DPhoto scopes and owner from the catalog table and stamps them as custom claims
 * into the access token, so the Lambda Authorizer can trust the token without any UserInfo call nor database lookup.
 */
export class CognitoPreTokenConstruct extends Construct {
    private readonly preTokenLambda: GoLangLambdaFunction;

    constructor(scope: Construct, id: string, props: CognitoPreTokenConstructProps) {
        super(scope, id);

        this.preTokenLambda = new GoLangLambdaFunction(this, 'PreTokenLambda', {
            environmentName: props.environmentName,
            functionName: 'cognito-pre-token',
            timeout: Duration.seconds(10),
            memorySize: 256,
        });

        props.catalogStore.grantCatalogReadAccess(this.preTokenLambda);

        props.userPool.addTrigger(
            cognito.UserPoolOperation.PRE_TOKEN_GENERATION_CONFIG,
            this.preTokenLambda.function,
            cognito.LambdaVersion.V2_0,
        );
    }
}
