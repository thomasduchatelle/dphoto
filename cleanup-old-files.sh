#!/bin/bash

# Remove old construct files
rm -f deployments/cdk/lib/constructs-storages/media-storage.ts
rm -f deployments/cdk/lib/constructs-storages/catalog-dynamodb.ts
rm -f deployments/cdk/lib/constructs-storages/archive-messaging.ts
rm -f deployments/cdk/lib/constructs-cli/cli-user.ts
rm -f deployments/cdk/lib/constructs-web/metadata-endpoints-construct.ts
rm -f deployments/cdk/lib/constructs-web/static-website-endpoint.ts
rm -f deployments/cdk/lib/constructs-users/authentication-endpoints-construct.ts
rm -f deployments/cdk/lib/constructs-users/users-endpoints-construct.ts
rm -f deployments/cdk/lib/constructs-catalog/catalog-endpoints-construct.ts
rm -f deployments/cdk/lib/constructs-archive/archive-endpoints-construct.ts

# Remove old directories if empty
rmdir deployments/cdk/lib/constructs-storages 2>/dev/null || true
rmdir deployments/cdk/lib/constructs-cli 2>/dev/null || true
rmdir deployments/cdk/lib/constructs-web 2>/dev/null || true
rmdir deployments/cdk/lib/constructs-users 2>/dev/null || true
rmdir deployments/cdk/lib/constructs-catalog 2>/dev/null || true

echo "Old files cleaned up successfully"
