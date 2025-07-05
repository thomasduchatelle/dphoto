import * as cdk from 'aws-cdk-lib';
import {DPhotoInfrastructureStack} from '../../lib/stacks/dphoto-infrastructure-stack';
import {environments} from '../../lib/config/environments';
import {DPhotoApplicationStack} from "../../lib/stacks/dphoto-application-stack";

describe('Deployment Integration Tests', () => {
    test('infrastructure stack can be synthesized for test environment', () => {
        const app = new cdk.App();

        // This should not throw
        const infraStack = new DPhotoInfrastructureStack(app, 'TestDeploymentStack', {
            environmentName: 'test',
            config: environments.test,
            env: {
                account: '123456789012',
                region: 'eu-west-1',
            },
        });
        const appStack = new DPhotoApplicationStack(app, 'TestApplicationStack',
            {
                config: environments.test,
                environmentName: "test",
                env: {
                    account: '123456789012',
                    region: 'eu-west-1',
                },
            })

        // Verify stack can be synthesized
        const assembly = app.synth();
        const infraStackArtifact = assembly.getStackByName(infraStack.stackName);

        expect(infraStackArtifact).toBeDefined();
        expect(infraStackArtifact.template).toBeDefined();

        const appStackArtifact = assembly.getStackByName(appStack.stackName);

        expect(appStackArtifact).toBeDefined();
        expect(appStackArtifact.template).toBeDefined();
    });
});
