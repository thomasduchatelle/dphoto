---
name: code-review
description: "Review the code changes on a branch: can the code be merged? It defines the priorities, and the code review standards, to use on this repository. You cannot code and review ; if you coded, dispatch a subagent to do the review."
---

# Code review

**You must answer a simple - but complex - question: can this code be merged ?**

You are to review all the changes for the branch you're on. Load the commits, inspect the code, read the requirements in the issue-tracker to understand the context, load the relevant skills.

Do not run the tests: they are executed independently on Github Action and will block the merge on their own right. Your role is to review statically the code changes.

## Comment only on what matter

The criticality of the comments are classified as follows:

1. **risk of data loss** - running this code will cause irrecoverable loss of data.
2. **risk of regression** - existing features and behaviour will fail using this new code. Look after what's not obvious and covered with test, includes misconfiguration, backward compatibility (or impossibility to rollback), storage drift, ...
3. **risk of intrusion** - sensitive information becomes more accessible. Look for any change in the authorisation layers, infrastructure setup, ... that would be reduced or missing and would let someone unauthorised to access data.
4. **risk of poor maintainability** - architecture and standard has been violated, testing strategy hasn't been respected, and that will affect future development. Look for responsibility leakage, poor boundaries, dead code or unnecessary code that could cause confusion... You're not nitpicking here, it's confirmed serious offences.
5. **risk of bugs** - newly delivered features will not work as expected. Look for any smell in the code that has been produced, any accidental behaviour. Explain why the test missed it if possible.

If asked to review a certain aspect (ex: feature parity), add it to the above list.

**You shall only comment on what would prevent the branch to be merged: no nitpicking, no praising, no explanation of what you've done or what the author did.** You're the last line of defence. Before you, the code has been pair-programed and reviewed by a lead. After it's straight to production. You must make it count, not making noise.

## Assess the mergeability

After a thorough review of the changes, give a clear recommendation:

* `HOLD` - a critical issue has been found: it is not safe to merge.
* `RECONSIDER` - while there is no confirmed critical issue (data loss, regression, security), at least one point needs to be investigated to fix or to mark as false-positive.
* `MERGE` - no issue has been found, the PR is safe to merge.

This recommendation is always about the full PR, never only the latest changes or a reduced scope.

## Where to comment

Use the most appropriate way to report your comments:

* if running as a subagent or reviewing on a PR, **then comment on the PR**.
* if reviewing a local branch and interacting with a user, comment in the chat.

## Follow-up reviews

If reviewing a PR that you have already reviewed (comments from the same account), only review what has changed and:

* raise new comment only if they were introduced by the recent changes.
* resolve or close previous comments that are not relevant any more.
* if your comment was answered but not fixed, review the answer. Then resolve the comment if you can confirm the reason the code wasn't changed, or push-back further if you consider the author is still missing something.

Finally, conclude with a new mergeability recommendation.

## Dispatch a code-reviewer subagent

You cannot code and review ; if you coded, dispatch a subagent to do the review.

When dispatching, keep the prompt minimum: "review my changes on the PR xxxx using code-reviewer skill. Post on the PR.". Don't reduce its scope, don't add custom instructions not explicitly given by the user for the review, and don't explain your changes.
