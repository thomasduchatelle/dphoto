import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import {Construct} from 'constructs';
import {GoLangLambdaFunction} from './golang-lambda-function';
import {ApiGatewayConstruct} from './api-gateway';

export interface NotFoundRoutesProps {
    environmentName: string;
}

export class NotFoundRoutesConstruct extends Construct {
    private readonly props: NotFoundRoutesProps;

    constructor(scope: Construct, id: string, props: NotFoundRoutesProps) {
        super(scope, id);
        this.props = props;
    }

    public addToApiGateway(apiGateway: ApiGatewayConstruct): void {
        const notFoundLambda = new GoLangLambdaFunction(this, 'NotFoundLambda', {
            environmentName: this.props.environmentName,
            functionName: 'not-found',
            artifactPath: '../../bin/not-found.zip'
        });

        notFoundLambda.addToApiGateway(
            apiGateway.httpApi,
            '/api/{path+}',
            apigatewayv2.HttpMethod.ANY
        );
    }
}
