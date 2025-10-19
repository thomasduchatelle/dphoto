#!/usr/bin/env bash
# Usage: enc-secret.sh <secret>
# Encrypts secret using AWS KMS (alias/dphoto-production-archive)
if [ "$#" -ne 1 ] || [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
  echo "Usage: $(basename "$0") <secret>" >&2
  exit 1
fi
KEY_ID=$(aws kms describe-key --key-id alias/dphoto-production-archive --query KeyMetadata.KeyId --output text)
aws kms encrypt --key-id "$KEY_ID" --plaintext "$(echo -n "$1" | base64)" --query CiphertextBlob --output text
