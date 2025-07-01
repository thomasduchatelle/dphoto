import * as cdk from 'aws-cdk-lib';
import {DPhotoInfrastructureStack} from '../../lib/stacks/dphoto-infrastructure-stack';
import {environments} from '../../lib/config/environments';

describe('Deployment Integration Tests', () => {
    test('infrastructure stack can be synthesized for test environment', () => {
        const app = new cdk.App();
        
        // This should not throw
        const stack = new DPhotoInfrastructureStack(app, 'TestDeploymentStack', {
            environmentName: 'test',
            config: environments.test,
            env: {
                account: '123456789012', // Mock account
                region: 'eu-west-1'
            }
        });

        // Verify stack can be synthesized
        const assembly = app.synth();
        const stackArtifact = assembly.getStackByName(stack.stackName);
        
        expect(stackArtifact).toBeDefined();
        expect(stackArtifact.template).toBeDefined();
    });

    test('all environments have valid configurations', () => {
        Object.entries(environments).forEach(([envName, config]) => {
            expect(config).toBeDefined();
            expect(typeof config.production).toBe('boolean');
            expect(Array.isArray(config.cliAccessKeys)).toBe(true);
            expect(typeof config.keybaseUser).toBe('string');
        });
    });
});
