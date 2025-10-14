The code of the function `checkPermissions` in `api/lambdas/authorizer/main.go` is awful. You need to refactor it to make it easier to extends with new routes.

1. Use a library to match the path with a pattern (or route id) and extract the path parameters
2. (alternative of 1, if no library is available) use regex
3. move the logic to match the route in a different file with a function that looks like:
   ```golang
   type Route struct {
        RouteId string // could be a regex, or some pattern the library is using (like `/ablums/{name}`)
        Method string // prefer using the HTTP method than a string if available
   }
   
   type MatchedRoute struct {
        Route Route
        PathParams map[string]string // PathParams are the extracted values from the path {name: 'return-f-the-jedi'}
   }
   
   func MatchRoute(routes []Route, path string) (*MatchedRoute, err) {}
   ```
4. make sure the function is well tested - all cases must be covered. Use the test naming like "it should match a route which doesn't have any path
   parameters", ...
5. on the `main`, you should be left with an easy-to-read configuration of all the endpoints supported. You can find the list of the endpoints to support in the
   CDK project: `deployments/cdk/lib/archive/archive-endpoints-construct.ts`, `deployments/cdk/lib/catalog/catalog-endpoints-construct.ts`

Before handing over the code, make sure the new tests are passing and the code is compiling.
