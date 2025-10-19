#!/usr/bin/env bash
# Usage: enc-secret.sh <environment> <secret>
# Stores secret in SSM Parameter Store under dphoto/cdk-input/googleClientSecret/<environment>
if [ "$#" -ne 2 ] || [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
  echo "Usage: $(basename "$0") <environment> <secret>" >&2
  exit 1
fi
aws ssm put-parameter --name "dphoto/cdk-input/googleClientSecret/$1" --value "$2" --type SecureString --overwrite
