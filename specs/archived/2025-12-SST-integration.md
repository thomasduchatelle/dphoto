As a senior developer, you need to deploy the new NextJS application as part of the Github Action workflow. You need to implement one of these stage, assuming
the others are already implemented or will be implemented by a different agent.

1. Update `deployments/cdk` project to:
    * Create a CloudFront distribution as part of the ApplicationStack
    * Configure the distribution to serve the API Gateway apis until `/api` route
    * Set a policy on the distribution to NEVER cache the output of `/api`

2. Build the NextJS application as part of the workflows
    * create the file `.github/workflows/job-build-nextjs.yml` which runs `cd web-nextjs && npm run build` with NodeJS 24.
    * update the workflows `.github/workflows/workflow-feature-branch.yml` and `.github/workflows/workflow-main.yml` to execute the NextJS build
    * make sure the cache is optimised in the build (similar from `.github/workflows/job-build-waku.yml`, and use your NextJS experience to make it even better)
    * and the dist is exposed (it's the folder `web-nextjs/.open-next`)

3. Update Github action workflows in `.github/workflows/job-deploy.yml`
    * after running CDK, generate the file `web-nextjs/.env.next` (where `next` is the stage or environment) with the content:
      ```
      SST_DISTRIBUTION_ID=<the value from CDK>
      SST_COGNITO_ISSUER=<the value from CDK>
      SST_COGNITO_CLIENT_ID=<the value from CDK>
      SST_COGNITO_CLIENT_SECRET=<the value from CDK>
      ```
      You need to update the `deployments/cdk` project to get access to these values after a `cdk deploy ...`

    * make sure the .env files are in the .gitignore list, it must not be checked out
    * then (after `cdk deploy` have run and the env file has been generated), run the SST deployment with the command
      `cd web-nextjs && sst --stage <environment> deploy`
