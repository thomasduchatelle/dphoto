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
  local output_key="$1"
  aws cloudformation describe-stacks \
    --stack-name "${STACK_NAME}" \
    --query "Stacks[0].Outputs[?OutputKey=='${output_key}'].OutputValue" \
    --output text 2>/dev/null || echo ""
}

# Extract the required outputs
DISTRIBUTION_ID=$(get_output "SSTDistributionId")
COGNITO_ISSUER=$(get_output "SSTCognitoIssuer")
COGNITO_CLIENT_ID=$(get_output "SSTCognitoClientId")
COGNITO_CLIENT_SECRET=$(get_output "SSTCognitoClientSecret")

# Validate that all required outputs were found
if [ -z "$DISTRIBUTION_ID" ] || [ -z "$COGNITO_ISSUER" ] || [ -z "$COGNITO_CLIENT_ID" ] || [ -z "$COGNITO_CLIENT_SECRET" ]; then
  echo "Error: Failed to retrieve one or more required CDK outputs from stack ${STACK_NAME}" >&2
  echo "  DISTRIBUTION_ID: ${DISTRIBUTION_ID:-<missing>}" >&2
  echo "  COGNITO_ISSUER: ${COGNITO_ISSUER:-<missing>}" >&2
  echo "  COGNITO_CLIENT_ID: ${COGNITO_CLIENT_ID:-<missing>}" >&2
  if [ -z "$COGNITO_CLIENT_SECRET" ]; then
    echo "  COGNITO_CLIENT_SECRET: <missing>" >&2
  else
    echo "  COGNITO_CLIENT_SECRET: <found>" >&2
  fi
  exit 1
fi

# Create the .env file
cat > "${REPO_ROOT}/${ENV_FILE}" <<EOF
SST_DISTRIBUTION_ID=${DISTRIBUTION_ID}
SST_COGNITO_ISSUER=${COGNITO_ISSUER}
SST_COGNITO_CLIENT_ID=${COGNITO_CLIENT_ID}
SST_COGNITO_CLIENT_SECRET=${COGNITO_CLIENT_SECRET}
EOF

# Restrict permissions to owner only for security
chmod 600 "${REPO_ROOT}/${ENV_FILE}"

echo "Successfully created ${ENV_FILE}"
echo "  SST_DISTRIBUTION_ID=${DISTRIBUTION_ID}"
echo "  SST_COGNITO_ISSUER=${COGNITO_ISSUER}"
echo "  SST_COGNITO_CLIENT_ID=${COGNITO_CLIENT_ID}"
echo "  SST_COGNITO_CLIENT_SECRET=<hidden>"
