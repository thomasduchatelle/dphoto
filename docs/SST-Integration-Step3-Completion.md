# SST Integration Step 3 - Completion Guide

This document describes what needs to be done to complete the SST integration after steps 1 and 2 are implemented.

## Current Status

Step 3 has been implemented with the following components:

1. **CDK Infrastructure** (`deployments/cdk/`):
   - `NextjsStack.ts` now exports CloudFormation outputs for:
     - CloudFront Distribution ID
     - Cognito Issuer URL
     - Cognito Client ID
     - Cognito Client Secret
   - The NextJS stack is conditionally created (skipped in tests until build artifacts exist)

2. **GitHub Workflow** (`.github/workflows/job-deploy.yml`):
   - Extracts CDK outputs after deployment
   - Generates `.env.<stage>` file with SST environment variables
   - Runs SST deployment with the proper configuration

## Dependencies on Other Steps

### Step 1: CloudFront Distribution in ApplicationStack

**Current State**: NextjsStack creates its own CloudFront distribution via `cdk-nextjs-standalone`.

**Required Changes**:

1. Once ApplicationStack creates a CloudFront distribution with `/api` routing:
   - Update `NextjsStack.ts` to accept the existing distribution as a parameter
   - Modify the `Nextjs` construct configuration to use the provided distribution instead of creating a new one
   - Update the CDK output to export the distribution ID from the ApplicationStack distribution

2. Example change in `NextjsStack.ts`:
   ```typescript
   export interface AppRouterStackProps extends StackProps {
       cognitoConfig: CognitoStackExports;
       distribution: cloudfront.IDistribution;  // ADD THIS
   }
   
   export class AppRouterStack extends Stack {
       constructor(scope: Construct, id: string, props: AppRouterStackProps) {
           super(scope, id, props);
   
           const nextjs = new Nextjs(this, 'nextjs', {
               nextjsPath: '../../web-nextjs',
               skipBuild: true,
               streaming: true,
               distribution: props.distribution,  // USE PROVIDED DISTRIBUTION
           });
           
           // Export the distribution ID from the provided distribution
           new CfnOutput(this, "CloudFrontDistributionId", {
               value: props.distribution.distributionId,
               description: 'CloudFront Distribution ID for SST deployment',
               exportName: `${this.stackName}-DistributionId`,
           });
           // ... rest of outputs
       }
   }
   ```

3. Update `bin/dphoto.ts` to pass the distribution from ApplicationStack to NextjsStack:
   ```typescript
   const appRouterStack = new AppRouterStack(app, `dphoto-${envName}-nextjs`, {
       cognitoConfig: cognitoStack.getWebEnvironmentVariables(),
       distribution: applicationStack.cloudFrontDistribution,  // ADD THIS
       env: {
           account: account,
           region: region,
       },
   });
   ```

### Step 2: NextJS Build Workflow and Artifacts

**Current State**: The workflow assumes build artifacts don't exist yet. Tests are skipped if `.open-next/server-functions/default` doesn't exist.

**Required Changes**:

1. Add a download step in `.github/workflows/job-deploy.yml` before CDK deployment:
   ```yaml
   - name: Download NextJS build
     uses: actions/download-artifact@v4
     with:
       name: dist-nextjs
       path: web-nextjs/.open-next
   ```

2. Update the artifact upload in the build workflow (from step 2) to create `dist-nextjs`:
   ```yaml
   - name: NextJS Build Artifact
     uses: actions/upload-artifact@v4
     with:
       name: dist-nextjs
       path: web-nextjs/.open-next
   ```

3. Remove the conditional stack creation in `bin/dphoto.ts`:
   - Once the build artifacts are always available, remove the `shouldCreateNextJS` conditional
   - The NextJS stack should always be created in production
   
   Change from:
   ```typescript
   if (shouldCreateNextJS) {
       const appRouterStack = new AppRouterStack(app, ...);
       appRouterStack.addDependency(cognitoStack);
   }
   ```
   
   To:
   ```typescript
   const appRouterStack = new AppRouterStack(app, ...);
   appRouterStack.addDependency(cognitoStack);
   ```

## Testing the Complete Integration

Once both steps 1 and 2 are completed:

1. **Build Phase**:
   - Trigger a build on a feature branch
   - Verify the `dist-nextjs` artifact is created with the `.open-next` directory structure

2. **Deployment Phase**:
   - Deploy to the `next` environment
   - Verify CDK deploys successfully with all stacks including NextJS
   - Check the GitHub Actions log for:
     - Successful CDK output extraction
     - Proper `.env.next` file generation (with secrets redacted)
     - Successful SST deployment
   - Verify the NextJS application is accessible via the CloudFront distribution

3. **Verify SST Integration**:
   - Check that the SST deployment attaches to the existing CloudFront distribution
   - Verify environment variables are properly set in the deployed Lambda functions
   - Test authentication flow with Cognito

## Rollback Plan

If issues occur during integration:

1. **CloudFront Issues** (Step 1):
   - The NextjsStack can continue to use its own distribution temporarily
   - Comment out the distribution parameter passing
   - Deploy and verify functionality before attempting integration again

2. **Build Issues** (Step 2):
   - Keep the conditional stack creation in place
   - The deployment will skip the NextJS stack if artifacts are missing
   - This allows CDK deployment to succeed even if the NextJS build fails

## Environment Variables Reference

The following environment variables are generated and used by SST:

- `SST_DISTRIBUTION_ID`: CloudFront distribution ID from CDK
- `SST_COGNITO_ISSUER`: Cognito issuer URL (format: `https://cognito-idp.{region}.amazonaws.com/{poolId}`)
- `SST_COGNITO_CLIENT_ID`: Cognito user pool client ID
- `SST_COGNITO_CLIENT_SECRET`: Cognito user pool client secret

These are:
1. Extracted from CloudFormation outputs after CDK deployment
2. Written to `web-nextjs/.env.<stage>` (gitignored)
3. Sourced and exported before running `sst deploy`
4. Used by `sst.config.ts` to configure the NextJS deployment

## Additional Considerations

### Security

- Ensure Cognito client secret is never logged in plaintext
- The `.env.*` files are already gitignored
- Consider using AWS Secrets Manager for production secrets in the future

### Cost Optimization

- SST deployments create Lambda functions for server-side rendering
- Monitor Lambda invocations and CloudFront cache hit rates
- Adjust cache policies if needed to stay within the <$5/month budget

### Monitoring

After deployment, monitor:
- Lambda function errors and durations
- CloudFront distribution metrics
- API Gateway request counts
- NextJS application performance

## Questions or Issues

If you encounter problems during integration:

1. Check CloudFormation stack events for deployment errors
2. Review GitHub Actions logs for workflow failures
3. Verify all environment variables are correctly set
4. Test the CloudFront distribution routing rules
5. Ensure the `.open-next` build directory structure is correct
