import {Construct} from 'constructs';
import * as apigatewayv2 from "aws-cdk-lib/aws-apigatewayv2";
import {MappingValue, ParameterMapping} from "aws-cdk-lib/aws-apigatewayv2";
import * as apigatewayv2_integrations from "aws-cdk-lib/aws-apigatewayv2-integrations";

export interface WakuWebUiConstructProps {
    nextjsDomainName: string
    httpApi: apigatewayv2.HttpApi;
}

export class NextJsRoutingConstruct extends Construct {
    constructor(scope: Construct, id: string, {
        httpApi,
        nextjsDomainName,
    }: WakuWebUiConstructProps) {
        super(scope, id);

        const nextJsHttpIntegration = new apigatewayv2_integrations.HttpUrlIntegration(
            `NextJsWithBasePathIntegration`,
            `https://${nextjsDomainName}`,
            {
                parameterMapping: new ParameterMapping().overwritePath(MappingValue.requestPath()),
            },
        );
        new apigatewayv2.HttpRoute(this, 'NextJsEagerRoute', {
            httpApi: httpApi,
            routeKey: apigatewayv2.HttpRouteKey.with('/nextjs/{proxy+}', apigatewayv2.HttpMethod.ANY),
            integration: nextJsHttpIntegration
        });
        new apigatewayv2.HttpRoute(this, 'NextJsBaseRoute', {
            httpApi: httpApi,
            routeKey: apigatewayv2.HttpRouteKey.with('/nextjs', apigatewayv2.HttpMethod.ANY),
            integration: nextJsHttpIntegration
        });
    }
}
