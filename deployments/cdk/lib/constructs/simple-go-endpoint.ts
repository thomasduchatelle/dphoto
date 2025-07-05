import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import {Construct} from 'constructs';
import {GoLangLambdaFunction} from './golang-lambda-function';
import {Duration} from "aws-cdk-lib";

export interface SimpleGoEndpointProps {
    environmentName: string;
    functionName: string;
    httpApi: apigatewayv2.HttpApi;
    path: string;
    method: apigatewayv2.HttpMethod;
    artifactPath?: string;
    memorySize?: number;
    timeout?: Duration;
}

export class SimpleGoEndpoint extends Construct {
    public readonly lambda: GoLangLambdaFunction;

    constructor(scope: Construct, id: string, props: SimpleGoEndpointProps) {
        super(scope, id);

        this.lambda = new GoLangLambdaFunction(this, 'Lambda', {
            environmentName: props.environmentName,
            functionName: `dphoto-${props.environmentName}-${props.functionName}`,
            artifactPath: props.artifactPath || `../../bin/${props.functionName}.zip`,
            memorySize: props.memorySize || 256,
            timeout: props.timeout || Duration.minutes(1),
        });

        this.lambda.addToApiGateway(props.httpApi, props.path, props.method);
    }
}
