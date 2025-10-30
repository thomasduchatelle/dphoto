#!/bin/bash

# Script to read resource identifiers from Terraform outputs
# Usage: ./read-identifiers-from-terraform.sh [workspace_name]
# Default workspace is 'dev'

set -e

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Get workspace name from argument or default to 'dev'
WORKSPACE=${1:-dev}

echo "Reading identifiers from Terraform workspace: $WORKSPACE"
echo "=================================================="

# Change to the infra-data directory relative to the script location
cd "$SCRIPT_DIR/../../infra-data"

# Select the terraform workspace
terraform workspace select "$WORKSPACE" || {
    echo "Error: Could not select workspace '$WORKSPACE'"
    echo "Available workspaces:"
    terraform workspace list
    exit 1
}

# Get terraform outputs in JSON format
OUTPUTS=$(terraform output -json)

# Extract and display the resource identifiers
echo ""
echo "Resource Identifiers:"
echo "===================="

# Storage RW Policy ARN
STORAGE_RW_ARN=$(echo "$OUTPUTS" | jq -r '.delegate_secret_access_key_decrypt_cmd.value // empty' 2>/dev/null || echo "")
if [ -n "$STORAGE_RW_ARN" ]; then
    # Get from SSM parameter instead since it's not in direct outputs
    STORAGE_RW_ARN=$(aws ssm get-parameter --name "/dphoto/$WORKSPACE/iam/policies/storageRWArn" --query 'Parameter.Value' --output text 2>/dev/null || echo "Not found")
fi
echo "MediaStorage/StorageRwPolicy ARN: $STORAGE_RW_ARN"

# Storage RO Policy ARN
STORAGE_RO_ARN=$(aws ssm get-parameter --name "/dphoto/$WORKSPACE/iam/policies/storageROArn" --query 'Parameter.Value' --output text 2>/dev/null || echo "Not found")
echo "MediaStorage/StorageRoPolicy ARN: $STORAGE_RO_ARN"

# Cache RW Policy ARN
CACHE_RW_ARN=$(aws ssm get-parameter --name "/dphoto/$WORKSPACE/iam/policies/cacheRWArn" --query 'Parameter.Value' --output text 2>/dev/null || echo "Not found")
echo "MediaStorage/CacheRwPolicy ARN: $CACHE_RW_ARN"

# Index RW Policy ARN
INDEX_RW_ARN=$(aws ssm get-parameter --name "/dphoto/$WORKSPACE/iam/policies/indexRWArn" --query 'Parameter.Value' --output text 2>/dev/null || echo "Not found")
echo "CatalogDb/IndexRwPolicy ARN: $INDEX_RW_ARN"

# Archive Topic ARN
ARCHIVE_TOPIC_ARN=$(echo "$OUTPUTS" | jq -r '.sns_archive_arn.value // empty')
echo "ArchiveMessaging/ArchiveTopic ARN: $ARCHIVE_TOPIC_ARN"

# Archive Queue URL
ARCHIVE_QUEUE_URL=$(echo "$OUTPUTS" | jq -r '.sqs_archive_url.value // empty')
echo "ArchiveMessaging/ArchiveQueue URL: $ARCHIVE_QUEUE_URL"

# Archive Queue ARN (for policy ID)
ARCHIVE_QUEUE_ARN=$(echo "$OUTPUTS" | jq -r '.sqs_async_archive_jobs_arn.value // empty')
echo "ArchiveMessaging/ArchiveQueue ARN: $ARCHIVE_QUEUE_ARN"

# Archive Relocate Queue URL
RELOCATE_QUEUE_URL=$(aws ssm get-parameter --name "/dphoto/$WORKSPACE/sqs/archive_relocate/url" --query 'Parameter.Value' --output text 2>/dev/null || echo "Not found")
echo "ArchiveMessaging/ArchiveRelocateQueue URL: $RELOCATE_QUEUE_URL"

# Archive SNS Publish Policy ARN
SNS_PUBLISH_ARN=$(aws ssm get-parameter --name "/dphoto/$WORKSPACE/iam/policies/archive_sns_publish/arn" --query 'Parameter.Value' --output text 2>/dev/null || echo "Not found")
echo "ArchiveMessaging/ArchiveSnsPublishPolicy ARN: $SNS_PUBLISH_ARN"

# Archive SQS Send Policy ARN
SQS_SEND_ARN=$(aws ssm get-parameter --name "/dphoto/$WORKSPACE/iam/policies/archive_sqs_send/arn" --query 'Parameter.Value' --output text 2>/dev/null || echo "Not found")
echo "ArchiveMessaging/ArchiveSqsSendPolicy ARN: $SQS_SEND_ARN"

# Archive Relocate Policy ARN
RELOCATE_POLICY_ARN=$(aws ssm get-parameter --name "/dphoto/$WORKSPACE/iam/policies/archive_relocate_send/arn" --query 'Parameter.Value' --output text 2>/dev/null || echo "Not found")
echo "ArchiveMessaging/ArchiveRelocatePolicy ARN: $RELOCATE_POLICY_ARN"

# CLI User Name
CLI_USER_NAME="dphoto-$WORKSPACE-cli"
echo "CliUser/CliUser Name: $CLI_USER_NAME"

echo ""
echo "Summary for CDK import:"
echo "======================"
echo "dphoto-$WORKSPACE-infra/MediaStorage/StorageRwPolicy/Resource: $STORAGE_RW_ARN"
echo "dphoto-$WORKSPACE-infra/MediaStorage/StorageRoPolicy/Resource: $STORAGE_RO_ARN"
echo "dphoto-$WORKSPACE-infra/MediaStorage/CacheRwPolicy/Resource: $CACHE_RW_ARN"
echo "dphoto-$WORKSPACE-infra/CatalogDb/IndexRwPolicy/Resource: $INDEX_RW_ARN"
echo "dphoto-$WORKSPACE-infra/ArchiveMessaging/ArchiveTopic/Resource: $ARCHIVE_TOPIC_ARN"
echo "dphoto-$WORKSPACE-infra/ArchiveMessaging/ArchiveQueue/Resource: $ARCHIVE_QUEUE_URL"
echo "dphoto-$WORKSPACE-infra/ArchiveMessaging/ArchiveQueue/Policy/Resource: $ARCHIVE_QUEUE_ARN"
echo "dphoto-$WORKSPACE-infra/ArchiveMessaging/ArchiveRelocateQueue/Resource: $RELOCATE_QUEUE_URL"
echo "dphoto-$WORKSPACE-infra/ArchiveMessaging/ArchiveSnsPublishPolicy/Resource: $SNS_PUBLISH_ARN"
echo "dphoto-$WORKSPACE-infra/ArchiveMessaging/ArchiveSqsSendPolicy/Resource: $SQS_SEND_ARN"
echo "dphoto-$WORKSPACE-infra/ArchiveMessaging/ArchiveRelocatePolicy/Resource: $RELOCATE_POLICY_ARN"
echo "dphoto-$WORKSPACE-infra/CliUser/CliUser/Resource: $CLI_USER_NAME"
