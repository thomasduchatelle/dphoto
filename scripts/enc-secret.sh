#!/usr/bin/env bash
set -euo pipefail

# enc-secret.sh
# Usage: enc-secret.sh <env-name> <secret> [--key-id KMS_KEY_ID] [--raw]
# Encrypts <secret> using AWS KMS and prints the base64 ciphertext to stdout.
# With --raw the script prints only the base64 ciphertext (no extra comments),
# which is suitable for scripting and storing in environment variables.
# The environment name is retained for naming/consistency but isn't used by KMS.
# By default this script will use the production KMS alias 'alias/dphoto-production-archive'.
# You can override the key by setting the KMS_KEY_ID environment variable or by passing --key-id.
# You must have the AWS CLI configured and a KMS key id available via
# the KMS_KEY_ID environment variable or by passing --key-id.

usage() {
  cat <<EOF
Usage: $(basename "$0") <env-name> <secret> [--key-id KMS_KEY_ID] [--raw]

Examples:
  # use default production KMS alias (no KMS_KEY_ID required)
  ./scripts/enc-secret.sh next "my-secret"

  # use a specific KMS key id / alias and print only the base64 ciphertext
  ./scripts/enc-secret.sh next "my-secret" --key-id arn:aws:kms:us-east-1:123456789012:key/abcd1234 --raw
  ./scripts/enc-secret.sh next "my-secret" --key-id alias/dphoto-production-archive --raw

Notes:
 - By default the script uses the production KMS alias: alias/dphoto-production-archive
 - The script prints a base64 ciphertext to stdout. Do NOT commit secrets to the git repository.
 - Requires AWS CLI v2 configured with permissions to call kms:Encrypt for the provided key.
EOF
  exit 1
}

if [ "$#" -lt 2 ]; then
  usage
fi

ENV_NAME=$1
SECRET=$2
shift 2

# parse optional arguments
DEFAULT_KMS_ALIAS='alias/dphoto-production-archive'
KMS_KEY_ID="${KMS_KEY_ID:-}" # allow env override
RAW_OUTPUT=0
while [ "$#" -gt 0 ]; do
  case "$1" in
    --key-id)
      shift
      if [ -z "${1:-}" ]; then
        echo "--key-id requires a value" >&2
        usage
      fi
      KMS_KEY_ID="$1"
      ;;
    --raw)
      RAW_OUTPUT=1
      ;;
    *)
      echo "Unknown option: $1" >&2
      usage
      ;;
  esac
  shift
done

if ! command -v aws >/dev/null 2>&1; then
  echo "aws CLI not found. Install and configure AWS CLI (https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html)" >&2
  exit 2
fi

# If no KMS key provided, default to the production alias baked into the deployment
if [ -z "$KMS_KEY_ID" ]; then
  KMS_KEY_ID="$DEFAULT_KMS_ALIAS"
  echo "No KMS key id provided; defaulting to production KMS alias: $KMS_KEY_ID" >&2
fi

# Use fileb://- to provide plaintext from stdin to aws kms encrypt
# The output is the CiphertextBlob in base64.
# We use printf to avoid adding a trailing newline to the secret.
CIPHERTEXT_B64=$(printf '%s' "$SECRET" | aws kms encrypt --key-id "$KMS_KEY_ID" --plaintext fileb://- --query CiphertextBlob --output text) || {
  echo "Failed to encrypt secret with KMS (key: $KMS_KEY_ID)." >&2
  exit 3
}

if [ "$RAW_OUTPUT" -eq 1 ]; then
  # Print only the base64 ciphertext for automation
  printf '%s' "$CIPHERTEXT_B64"
  exit 0
fi

# Print a simple label with the environment and the ciphertext
cat <<EOF
# Encrypted secret for environment: $ENV_NAME
# Ciphertext (base64):
$CIPHERTEXT_B64
EOF

# Helpful instruction for storing the value (example for CDK context or env var)
cat <<'USAGE'

How to use the result:
 - Save the base64 ciphertext in a secure place (e.g., your secrets manager), or
   add it to your deployment configuration where you expect `googleClientSecretEncrypted`.
 - Example: export GOOGLE_CLIENT_SECRET_ENCRYPTED="<the base64 string>"
 - The CDK or deployment must know how to decrypt this ciphertext before passing a plaintext
   to the runtime where required. This script only performs the KMS encryption step.

Security note: Do NOT commit secrets (encrypted or plaintext) to source control unless you
understand and accept the risks. Prefer storing encrypted values in a secrets manager.
USAGE
