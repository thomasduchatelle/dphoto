You need to add the authoriser on the get-media endpoint.

Unlike the other ones, the access token is passed by a query parameter. You need to:

1. create a new authoriser (using the same lambda) but expecting the identity from the query parameter (CDK script)
2. update the authroiser to get the access token from the query parameter
3. update the get-media to remove the authorising code and get the information it requires from the context passed down by the authoriser
4. update the test `deployments/cdk/lib/stacks/application-stack.test.ts`: the get media do need an authroiser now.

You might find valudable the following files:

* `deployments/cdk/lib/archive/archive-endpoints-construct.ts`
* `api/lambdas/authorizer/main.go`
* `deployments/cdk/lib/access/lambda-authoriser-construct.ts`
* `api/lambdas/get-media/main.go`

Your acceptance criteria:

1. the authrorisation logic must move from the get media handler to the authoriser
2. do not introduce security hole: the authoriser must fulfill the same requirement
3. the golang code must compile `make build-go`
4. the cdk project must compile and pass tests
