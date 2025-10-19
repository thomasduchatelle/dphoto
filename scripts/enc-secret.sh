#!/usr/bin/env bash
# Usage: enc-secret.sh <secret>
# Encrypts secret using AWS KMS (alias/dphoto-live-archive)
if [ "$#" -ne 1 ] || [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
  echo "Usage: $(basename "$0") <secret>" >&2
  exit 1
fi
aws kms encrypt --key-id alias/dphoto-live-archive --plaintext "$(echo -n "$1" | base64)" --query CiphertextBlob --output text
