You need to prepare a plan to update the file `web/src/middleware/authentication.tsx` to get this code production ready.

You'll use the `TODO AGENT` comment in the file to generate your plan, and the following cases must be covered.

1. happy path: being redirected, then come cack with authorisation code, and being redirected to the originally requested page
2. do not redirect when it's not a navigation request (browser loading page) and instead is a AJAX reuest or something similar where the user shouldn't be
   reirected.
3. logout happy path: when `/logout` page is loaded, sign off from cognito, clear cookies, and redirect to a page (which must be white listed from the
   middleware function of course)

Your plan is a list of tests that can be implemented as unit tests of this middleware function. Each will only reqauires calling the function once, and will be
checking the results (the HTTP response) is as expected. Use an explicit title lik "it should redirect to authorisation authority when requesting a home page
without headers". This title should be enough for a human with context to write the test. For agents, adds details of the path, request headers, and the
expected response.

You start with simple and happy path tests, then introduces the non-happy path, error cases, ... But keep the rule: 1 test = 1 call to the function (one input,
and one output to validate). Do not duplicate the tests, the list must be as consise as possible.

Write each of this tests in a MARKDOWN file in `specs/middleware-step{n}.md`. Write each with the relative path of each files that the agents will need to
update and enough context for the agent to start working immediately.

---

Now, plan the work on the middleware function and write the requirement of each of the test in its own file.