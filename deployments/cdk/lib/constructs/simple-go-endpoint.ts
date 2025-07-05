import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import {Construct} from 'constructs';
import {GoLangLambdaFunction} from './golang-lambda-function';

export interface SimpleGoEndpointProps {
    environmentName: string;
    functionName: string;
    path: string;
    method: apigatewayv2.HttpMethod;
    artifactPath?: string;
    memorySize?: number;
    timeout?: number;
}

export class SimpleGoEndpoint extends Construct {
    public readonly lambda: GoLangLambdaFunction;

    constructor(scope: Construct, id: string, props: SimpleGoEndpointProps) {
        super(scope, id);

        this.lambda = new GoLangLambdaFunction(this, 'Lambda', {
            environmentName: props.environmentName,
            functionName: props.functionName,
            artifactPath: props.artifactPath || `../../bin/${props.functionName}.zip`,
            memorySize: props.memorySize || 256,
            timeout: props.timeout
        });
    }

    public addToApiGateway(httpApi: apigatewayv2.HttpApi, path?: string, method?: apigatewayv2.HttpMethod): void {
        this.lambda.addToApiGateway(
            httpApi,
            path || this.getPath(),
            method || this.getMethod()
        );
    }

    protected getPath(): string {
        throw new Error('Path must be provided either in constructor or addToApiGateway method');
    }

    protected getMethod(): apigatewayv2.HttpMethod {
        throw new Error('Method must be provided either in constructor or addToApiGateway method');
    }
}
