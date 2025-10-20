# Cognito User Migration Guide

## Overview

This guide documents the process for migrating existing users from the DynamoDB-only user management system to AWS Cognito with Google SSO authentication. As outlined in `specs/2025-10_cognito-authentication-migration.md`, the migration involves recreating users in Cognito while maintaining their existing permissions in DynamoDB.

## Prerequisites

- AWS CLI configured with appropriate permissions
- Access to the Cognito User Pool (created via CDK deployment)
- List of existing users and their current permissions
- Cognito User Pool ID from CDK output

## Migration Strategy

The system has been designed to support **gradual migration** with **zero downtime**:

1. **Dual System Operation**: Users can exist in both DynamoDB (current system) and Cognito (future system) simultaneously
2. **Backward Compatibility**: The system continues to work without Cognito configured
3. **Optional Cognito**: Set `COGNITO_USER_POOL_ID` environment variable only when ready to enable Cognito integration

## Pre-Migration Checklist

Before starting the migration, ensure:

- [ ] Cognito User Pool is deployed via CDK
- [ ] Google OAuth credentials are configured in Cognito
- [ ] Users have been notified about the migration
- [ ] User email addresses are verified (Google accounts must match exactly)
- [ ] Backup of current DynamoDB scopes table exists

## Migration Steps

### Step 1: Identify Existing Users

Query DynamoDB to list all existing users and their permissions:

```bash
# List all users with MainOwnerScope
aws dynamodb scan \
    --table-name <DYNAMODB_TABLE_NAME> \
    --filter-expression "begins_with(SK, :scope_prefix)" \
    --expression-attribute-values '{":scope_prefix":{"S":"SCOPE#owner:main"}}' \
    --projection-expression "PK,SK,ResourceOwner" \
    --output table
```

```bash
# List all users with AlbumVisitorScope
aws dynamodb scan \
    --table-name <DYNAMODB_TABLE_NAME> \
    --filter-expression "begins_with(SK, :scope_prefix)" \
    --expression-attribute-values '{":scope_prefix":{"S":"SCOPE#album:visitor"}}' \
    --projection-expression "PK,SK,ResourceOwner,ResourceId" \
    --output table
```

Create a CSV file with the following format:

```csv
Email,UserType,Owner,AlbumId
tony@stark.com,owner,tony@stark.com,
pepper@stark.com,visitor,tony@stark.com,/2024/01-family-vacation
```

Where:
- **Email**: User's email address (must match their Google account)
- **UserType**: Either "owner" or "visitor"
- **Owner**: The owner tenant (for owners, usually same as email; for visitors, the album owner)
- **AlbumId**: For visitors only, the folder name of the shared album

### Step 2: Deploy Infrastructure with Cognito

Deploy the CDK stacks which includes the Cognito User Pool:

```bash
cd deployments/cdk
cdk deploy --context environment=next --all
```

Note the Cognito User Pool ID from the output:

```
Outputs:
  dphoto-next-infrastructure.UserPoolId = us-east-1_XXXXXXXXX
```

### Step 3: Test Cognito Configuration

Before migrating all users, test with a single test user:

```bash
# Create a test user in Cognito
aws cognito-idp admin-create-user \
    --user-pool-id <USER_POOL_ID> \
    --username test@example.com \
    --user-attributes Name=email,Value=test@example.com Name=email_verified,Value=true \
    --message-action SUPPRESS

# Add to owners group
aws cognito-idp admin-add-user-to-group \
    --user-pool-id <USER_POOL_ID> \
    --username test@example.com \
    --group-name owners
```

Verify the user can:
1. Authenticate via Google SSO
2. Access their albums
3. Upload photos (if they are an owner)

### Step 4: Migrate Owners

For each owner in your CSV file, create a Cognito user in the `owners` group.

**Using the CLI** (recommended for automated migration):

```bash
# Set environment variable
export COGNITO_USER_POOL_ID=<USER_POOL_ID>

# For each owner, run:
dphoto create-user <EMAIL> --owner <OWNER_ID>
```

The CLI will now automatically create the user in both DynamoDB and Cognito.

**Using AWS CLI directly** (manual process):

```bash
# Create user
aws cognito-idp admin-create-user \
    --user-pool-id <USER_POOL_ID> \
    --username <EMAIL> \
    --user-attributes Name=email,Value=<EMAIL> Name=email_verified,Value=true \
    --message-action SUPPRESS

# Add to owners group
aws cognito-idp admin-add-user-to-group \
    --user-pool-id <USER_POOL_ID> \
    --username <EMAIL> \
    --group-name owners
```

**Important Notes:**
- User creation in DynamoDB (scopes) is still required and is already in place
- The `create-user` CLI command now handles both DynamoDB and Cognito
- Existing DynamoDB scopes are NOT automatically migrated - they remain as-is

### Step 5: Migrate Visitors

For each visitor in your CSV file, create a Cognito user in the `visitors` group.

**Using the CLI** (recommended):

```bash
# Set environment variable
export COGNITO_USER_POOL_ID=<USER_POOL_ID>

# For each visitor, run:
dphoto share-album --owner <OWNER> --album <ALBUM_ID> --email <EMAIL>
```

The CLI will now automatically create the visitor in both DynamoDB and Cognito.

**Using AWS CLI directly** (manual process):

```bash
# Create user
aws cognito-idp admin-create-user \
    --user-pool-id <USER_POOL_ID> \
    --username <EMAIL> \
    --user-attributes Name=email,Value=<EMAIL> Name=email_verified,Value=true \
    --message-action SUPPRESS

# Add to visitors group
aws cognito-idp admin-add-user-to-group \
    --user-pool-id <USER_POOL_ID> \
    --username <EMAIL> \
    --group-name visitors
```

### Step 6: Update Lambda Environment Variables

Update the Lambda functions to use Cognito:

**For share-album Lambda:**

```bash
aws lambda update-function-configuration \
    --function-name <FUNCTION_NAME> \
    --environment "Variables={COGNITO_USER_POOL_ID=<USER_POOL_ID>,...}"
```

The CDK deployment should handle this automatically, but verify that the environment variable is set.

### Step 7: Verification

After migration, verify each user:

1. **Authentication**: User can log in via Google SSO
2. **Authorization**: User has appropriate access to albums
3. **Group Membership**: User is in the correct Cognito group (admins, owners, or visitors)

**Verify user in Cognito:**

```bash
# Check user exists
aws cognito-idp admin-get-user \
    --user-pool-id <USER_POOL_ID> \
    --username <EMAIL>

# List user's groups
aws cognito-idp admin-list-groups-for-user \
    --user-pool-id <USER_POOL_ID> \
    --username <EMAIL>
```

**Verify user has DynamoDB scopes:**

```bash
aws dynamodb get-item \
    --table-name <DYNAMODB_TABLE_NAME> \
    --key '{"PK":{"S":"USER#<EMAIL>"},"SK":{"S":"SCOPE#owner:main#<OWNER>#"}}'
```

### Step 8: Monitor and Troubleshoot

After migration, monitor CloudWatch logs for:

- Authentication failures
- Authorization errors
- Missing Cognito users

Common issues and solutions:

**Issue**: User cannot authenticate
- **Solution**: Verify email address matches Google account exactly
- **Solution**: Check user exists in Cognito User Pool
- **Solution**: Verify user is in appropriate group

**Issue**: User gets "access denied" 
- **Solution**: Verify DynamoDB scopes are still present
- **Solution**: Check user's group membership in Cognito
- **Solution**: Verify Lambda has appropriate Cognito permissions

**Issue**: "User not found" error during share-album
- **Solution**: Ensure `COGNITO_USER_POOL_ID` is set in Lambda environment
- **Solution**: Verify Lambda has permissions to create users in Cognito

## Rollback Plan

If issues arise, the system can operate without Cognito:

1. **Remove environment variable**: Unset `COGNITO_USER_POOL_ID` from Lambda and CLI configuration
2. **DynamoDB continues to work**: All existing scopes in DynamoDB remain functional
3. **No data loss**: Cognito users can be deleted without affecting DynamoDB scopes

To rollback:

```bash
# Remove environment variable from Lambda
aws lambda update-function-configuration \
    --function-name <FUNCTION_NAME> \
    --environment "Variables={COGNITO_USER_POOL_ID=}"

# In CLI configuration, remove or comment out
# cognito.user.pool.id=<USER_POOL_ID>
```

## Post-Migration

After successful migration and verification:

1. **Update documentation**: Inform users about the new Google SSO login
2. **Monitor metrics**: Watch authentication and authorization success rates
3. **Plan phase 2**: Consider migrating to Cognito-based authorization (Amazon Verified Permissions)

## User Count Estimation

For reference, the design document mentions only **6 existing users**. This small number makes manual migration straightforward and low-risk.

## Support and Troubleshooting

### Logs to Check

- **Lambda logs**: CloudWatch Logs for share-album and authentication functions
- **Application logs**: Search for "Cognito" to see user creation events

### Useful Commands

```bash
# List all users in Cognito
aws cognito-idp list-users \
    --user-pool-id <USER_POOL_ID> \
    --output table

# Count users by group
aws cognito-idp list-users-in-group \
    --user-pool-id <USER_POOL_ID> \
    --group-name owners \
    | jq '.Users | length'

# Disable a user (emergency)
aws cognito-idp admin-disable-user \
    --user-pool-id <USER_POOL_ID> \
    --username <EMAIL>
```

## Automated Migration Script

For convenience, here's a sample bash script to automate the migration:

```bash
#!/bin/bash

# Configuration
USER_POOL_ID="<YOUR_USER_POOL_ID>"
CSV_FILE="users.csv"

# Read CSV and create users
tail -n +2 "$CSV_FILE" | while IFS=, read -r email userType owner albumId; do
    echo "Processing $email ($userType)..."
    
    # Create user in Cognito
    aws cognito-idp admin-create-user \
        --user-pool-id "$USER_POOL_ID" \
        --username "$email" \
        --user-attributes Name=email,Value="$email" Name=email_verified,Value=true \
        --message-action SUPPRESS \
        2>/dev/null || echo "User $email might already exist"
    
    # Add to appropriate group
    if [ "$userType" = "owner" ]; then
        GROUP="owners"
    else
        GROUP="visitors"
    fi
    
    aws cognito-idp admin-add-user-to-group \
        --user-pool-id "$USER_POOL_ID" \
        --username "$email" \
        --group-name "$GROUP" \
        2>/dev/null || echo "User $email might already be in group $GROUP"
    
    echo "✓ $email migrated as $userType"
done

echo "Migration complete!"
```

Save this as `migrate-users.sh`, make it executable with `chmod +x migrate-users.sh`, and run it with your users.csv file.

## Timeline and Communication

**Recommended Timeline:**

1. **Week 1**: Deploy Cognito infrastructure, test with development environment
2. **Week 2**: Test with 1-2 power users, gather feedback
3. **Week 3**: Migrate remaining users, monitor closely
4. **Week 4**: Full production deployment, remove old authentication system

**User Communication Template:**

> Subject: DPhoto Authentication Upgrade - Action Required
> 
> We're upgrading the authentication system to use Google Single Sign-On (SSO) for improved security and convenience.
> 
> **What's changing:**
> - You'll now log in using your Google account
> - Your email address must match the one we have on file
> - All your photos and albums will remain accessible
> 
> **Action required:**
> - Verify your email address: [YOUR_EMAIL]
> - Ensure you have access to the Google account with this email
> - After [DATE], you'll need to log in using "Sign in with Google"
> 
> **Timeline:**
> - [DATE]: Migration begins
> - [DATE]: Old login method disabled
> 
> If you have any questions or concerns, please contact support.

## Conclusion

This migration guide provides a safe, gradual approach to moving users from DynamoDB-only authentication to Cognito with Google SSO. The backward-compatible design ensures zero downtime and easy rollback if needed.

For questions or issues during migration, refer to the design document at `specs/2025-10_cognito-authentication-migration.md` or contact the development team.
