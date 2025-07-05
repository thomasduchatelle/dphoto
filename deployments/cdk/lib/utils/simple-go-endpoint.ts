import * as apigatewayv2 from 'aws-cdk-lib/aws-apigatewayv2';
import {Construct} from 'constructs';
import {GoLangLambdaFunction, GoLangLambdaFunctionProps} from './golang-lambda-function';
import * as apigatewayv2_integrations from "aws-cdk-lib/aws-apigatewayv2-integrations";

export interface RouteConfig {
    path: string;
    method: apigatewayv2.HttpMethod;
}

export interface SimpleGoEndpointProps extends GoLangLambdaFunctionProps {
    httpApi: apigatewayv2.HttpApi;
    routes: RouteConfig[];
}

export class SimpleGoEndpoint extends Construct {
    public readonly lambda: GoLangLambdaFunction;
    private readonly integration: apigatewayv2_integrations.HttpLambdaIntegration;

    constructor(scope: Construct, id: string, {httpApi, routes, ...props}: SimpleGoEndpointProps) {
        super(scope, id);

        this.lambda = new GoLangLambdaFunction(this, 'Lambda', {
            ...props,
        });

        this.integration = new apigatewayv2_integrations.HttpLambdaIntegration(
            `${this.node.id}Integration`,
            this.lambda.function,
        );

        routes.forEach((route, index) => {
            this.addRoute(httpApi, route);
        });
    }

    public addRoute(
        httpApi: apigatewayv2.HttpApi,
        route: RouteConfig
    ): void {
        const routeId = `Route${route.method}${route.path.replace(/[^a-zA-Z0-9]/g, '')}`;

        new apigatewayv2.HttpRoute(this, routeId, {
            httpApi,
            routeKey: apigatewayv2.HttpRouteKey.with(route.path, route.method),
            integration: this.integration
        });
    }
}

export function createSingleRouteEndpoint(
    scope: Construct,
    id: string,
    props: Omit<SimpleGoEndpointProps, 'routes'> & {
        path: string;
        method: apigatewayv2.HttpMethod;
    }
): SimpleGoEndpoint {
    const {path, method, ...endpointProps} = props;
    return new SimpleGoEndpoint(scope, id, {
        ...endpointProps,
        routes: [{path, method}]
    });
}