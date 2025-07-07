import {Template} from 'aws-cdk-lib/assertions';
import * as cxapi from "aws-cdk-lib/cx-api";
import main from './dphoto';

jest.mock('aws-cdk-lib/aws-lambda', () => {
    const actual = jest.requireActual('aws-cdk-lib/aws-lambda');

    return {
        ...actual,
        Code: {
            ...actual.Code,
            fromAsset: jest.fn().mockImplementation(() => {
                const validAsset = "bin/";
                return actual.Code.fromAsset(validAsset);
            }),
        }
    };
});

jest.mock('aws-cdk-lib/aws-s3-deployment', () => {
    const actual = jest.requireActual('aws-cdk-lib/aws-s3-deployment');

    return {
        ...actual,
        Source: {
            ...actual.Source,
            asset: jest.fn().mockImplementation(() => {
                const validAsset = "bin/";
                return actual.Source.asset(validAsset);
            }),
        }
    };
});

describe('CDK Integration Tests', () => {
    test.each(['next', 'live'])('deployment has required resources for env %s', (envName) => {
        // Call your main function to create the app with stacks
        const app = main(envName, "01234567890123456");

        const assembly = app.synth();
        expect(assembly).toBeDefined();
        expect(assembly.stacks.length).toBeGreaterThan(0);

        const matcher = createAssemblyMatcher(assembly);

        // Test that required resources exist somewhere in the deployment
        matcher.hasResource('AWS::S3::Bucket', {
            BucketName: `dphoto-${envName}-storage`
        });

        matcher.hasResource('AWS::ApiGatewayV2::ApiMapping', {}); // properties are references, not actual values.
    });
});

const createAssemblyMatcher = (assembly: cxapi.CloudAssembly) => ({
    hasResource: (resourceType: string, properties?: any) => {
        const inspection = assembly.stacks.map(stack => {
            const template = Template.fromJSON(stack.template);
            try {
                template.hasResourceProperties(resourceType, properties || {});
                return {name: stack.stackName, error: undefined};
            } catch (e) {
                return {name: stack.stackName, error: e};
            }
        })

        const resourceHasBeenFound = inspection.some(i => !i.error);
        if (resourceHasBeenFound) {
            return
        }

        let message = `Resource ${resourceType} ${JSON.stringify(properties)} not found in any stack (${assembly.stacks.length}):`;
        inspection.forEach(i => {
            if (i.error) {
                message += `\n- Stack: ${i.name}, Error: ${i.error}`;
            }
        });

        assembly.stacks.forEach(stack => {
            const template = Template.fromJSON(stack.template);
            console.log(`${stack.stackName}:\n${JSON.stringify(template.toJSON(), null, 2)}`);
        })

        throw new Error(message);
    }
});