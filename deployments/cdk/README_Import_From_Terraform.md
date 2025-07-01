Terraform Migration
=======================================

1. Create the stack (empty)
   ```shell
   cdk deploy "$STACK_NAME" \
        --context environment="$ENVIRONMENT" \
        --context importOnly=true \
        --require-approval never \
        --verbose
   ```

2. Get the identifiers from terraform state
   ```
   ./scripts/read-identifiers-from-terraform.sh live
   ```

3. Run the CDK import ; confirm names that have been found, give the missing ones from the script above
   ```shell
   cdk import --context environment="live"
   ```

4. Deploy CDK ; or merge to the main branch
   ```shell
   cdk deploy --context environment="live"
   ```