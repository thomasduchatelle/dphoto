---
name: code-review
description: "Review the changes since a fixed point (commit, branch, tag, or merge-base): can the code be merged ?"
---

**You must answer a simple - but complex - question: can this code be merged ?**

The criticality of the comments are classified as follows:

1. **risk of data loss** - running this code will cause irrecoverable loss of data.
2. **risk of regression** - existing features and behaviour will fail using this new code. Look after what's not obvious and covered with test, includes misconfiguration, backward compatibility (or impossibility to rollback), storage drift, ...
3. **risk of intrusion** - sensitive information becomes more accessible. Look for any change in the authorisation layers, infrastructure setup, ... that would be reduced or missing and would let someone unauthorised to access data.
4. **risk of poor maintainability** - architecture and standard has been violated, testing strategy hasn't been respected, and that will affect future development. Look for responsibility leakage, poor boundaries, dead code or unnecessary code that could cause confusion... You're not nitpicking here, it's confirmed serious offences.
5. **risk of bugs** - newly delivered features will not work as expected. Look for any smell in the code that has been produced, any accidental behaviour. Explain why the test missed it if possible.

After that line, no point of raising a comment. Ignore every "poor coding standard" nitpicking, strict guidance or skilled not followed, recommended but not necessary changes.

You're the last line of defence. Before you, the code has been pair-programed and reviewed by a lead. After it's straight to production. You must make it count, not making noise.

---

You are to review all the changes for the branch you're on. Load the commits, inspect the code, read the requirements in the issue-tracker to understand the context.

If you've been asked to review a PR, post a single comment on the PR with your recommendation and your findings. The recommendations are:

* `MERGE` - no issue has been found, the PR is safe to merge.
* `RECONSIDER` - while there is no confirmed critical issue (data loss, regression, security), one point needs to be investigated to confirm it's a false-positive.
* `HOLD` - a critical issue has been found: it is not safe to merge.
