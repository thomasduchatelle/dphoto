#!/usr/bin/env bash
# Usage: extract-cdk-outputs-for-sst.sh <environment>
# Extracts CDK outputs and creates .env file for SST deployment
# The script retrieves CloudFormation outputs from the application stack and generates
# a .env file in web-nextjs directory with SST-required variables

set -euo pipefail

if [ "$#" -ne 1 ] || [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
  echo "Usage: $(basename "$0") <environment>" >&2
  echo "" >&2
  echo "Example: $(basename "$0") next" >&2
  echo "" >&2
  echo "This script extracts CDK outputs from the application stack and creates" >&2
  echo "a .env.<environment> file in the web-nextjs directory for SST deployment." >&2
  exit 1
fi

ENVIRONMENT="$1"
STACK_NAME="dphoto-${ENVIRONMENT}-application"
ENV_FILE="web-nextjs/.env.${ENVIRONMENT}"
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "Extracting CDK outputs from stack: ${STACK_NAME}"

# Function to get CloudFormation output value by key
get_output() {
  local export_name="$1"
  aws cloudformation describe-stacks \
    --stack-name "${STACK_NAME}" \
    --query "Stacks[0].Outputs[?ExportName=='${export_name}'].OutputValue" \
    --output text 2>/dev/null || echo ""
}

# Extract the required outputs
SST_CLOUD_FRONT_DOMAIN=$(get_output "dphoto-${ENVIRONMENT}-sst-cloudfront-domain")
COGNITO_ISSUER=$(get_output "dphoto-${ENVIRONMENT}-sst-cognito-issuer")
COGNITO_CLIENT_ID=$(get_output "dphoto-${ENVIRONMENT}-sst-cognito-client-id")
COGNITO_CLIENT_SECRET=$(get_output "dphoto-${ENVIRONMENT}-sst-cognito-client-secret")

# Validate that all required outputs were found
if [ -z "$SST_CLOUD_FRONT_DOMAIN" ] || [ -z "$COGNITO_ISSUER" ] || [ -z "$COGNITO_CLIENT_ID" ] || [ -z "$COGNITO_CLIENT_SECRET" ]; then
  echo "Error: Failed to retrieve one or more required CDK outputs from stack ${STACK_NAME}" >&2
  echo "  SST_CLOUD_FRONT_DOMAIN: ${SST_CLOUD_FRONT_DOMAIN:-<missing>}" >&2
  echo "  COGNITO_ISSUER: ${COGNITO_ISSUER:-<missing>}" >&2
  echo "  COGNITO_CLIENT_ID: ${COGNITO_CLIENT_ID:-<missing>}" >&2
  if [ -z "$COGNITO_CLIENT_SECRET" ]; then
    echo "  COGNITO_CLIENT_SECRET: <missing>" >&2
  else
    echo "  COGNITO_CLIENT_SECRET: <hidden>" >&2
  fi
  exit 1
fi

# Create the .env file
cat > "${REPO_ROOT}/${ENV_FILE}" <<EOF
SST_CLOUD_FRONT_DOMAIN=${SST_CLOUD_FRONT_DOMAIN}
SST_COGNITO_ISSUER=${COGNITO_ISSUER}
SST_COGNITO_CLIENT_ID=${COGNITO_CLIENT_ID}
SST_COGNITO_CLIENT_SECRET=${COGNITO_CLIENT_SECRET}
EOF

# Restrict permissions to owner only for security
chmod 600 "${REPO_ROOT}/${ENV_FILE}"

echo "Successfully created ${ENV_FILE}:"
cat "${REPO_ROOT}/${ENV_FILE}" | sed 's/^\(SST_COGNITO_CLIENT_SECRET=\).*/\1<hidden>/'
