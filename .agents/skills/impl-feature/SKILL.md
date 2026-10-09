---
name: impl-feature
description: 'Deliver a whole feature (a spec with stories/issues) by conducting the end-to-end cycle: dispatching coding subagents, resolving conflicts, dispatching reviewer agents, and owning the review-and-rebase loop across all PRs until merged. Use when the user asks to implement a feature made of several stories/issues that can be worked in parallel, or explicitly asks to "implement the feature X", "resume the implementation of the feature", or "spawn subagents for each subtask".'
---

# Delivering a feature via parallel subagents

You are the **conductor** and drive the implementation of the feature - and all its issues - end to end.

Your role is:

* **git** - set up the initial branch, merge the approved PRs, solve conflicts.
* **orchestration** - dispatch the coding agents, and the review agents. Report to the user what's ready to be reviewed.
* **coding** - be pragmatic, some coding task are done more efficiently directly or interactively with the user.
* **wrap-up** - prepare to merge to `main`: PR and review agent.

## Coding orchestration

Read the feature: `specs/<slug>/spec.md` and every `specs/<slug>/issues/NN-*.md`. Extract the dependency graph.

**Important! Prioritise the parallelisation of the coding tasks: one subagent = one story = one PR.**

Coding yourself or using subagents?

* **the first story** that all others depend on: code it yourself, interactively with the user.
* **a batch of stories** that can be implemented in parallel: dispatch coding agents.
* **fixes or review feedback to address**: dispatch coding agent(s).
* **explicitly asked to _bring to PR_** locally and address the review: delete the worktree, check out the PR branch, and work interactively with the user to resolve the review (propose a solution before implementing).
* **explicitly asked to make a change**: apply the requested changes on the current branch; do not commit without an explicit request — this is likely outside the automated workflow and the user will say when to resume it.

When every PR in the batch is open (or re-pushed), stop. Post a table to the user, ordered by priority (number of stories blocked):

| # | Story | PR   | Pending stories                      |
|---|-------|------|--------------------------------------|
| … | …     | link | <list of stories blocked by this PR> |

Flag anything unusual: force-pushes you authorised, deleted tests, unresolved review threads. Wait for the user's next instruction — fixes to make, or the go-ahead to merge and launch the next batch.

Once approved, merge the story into the feature branch. Resolve conflicts if necessary. Do not wait for specific instructions — immediately start any stories that became unblocked.

Keep the issue tracker up to date: mark merged stories as `done`. If the feature is complete, archive it before merging.


## Git branching strategy

Features are developed and delivered through 3 branches:

1. **story branch** - one per story/issue, it is where the coding-test-review loop occurs. A PR is created and must be approved by the user before being merged to the _feature branch_.
2. **feature branch** - one per feature, collection of one or more story/issue. A PR is created upon user request and must be reviewed by the review subagent before being presented to the user for final review.
3. **main** - this is the stable branch, continuous deployment is enabled.

Once merged, the **feature branch** must be recreated or updated for the next batch of stories.

Apply this instruction yourself and pass it to all coding subagents:

````
**Push command — use this exact form, nothing else:**

```
git push origin HEAD:refs/heads/<branch-name>
```

Do not use `git push -u`, do not use bare `git push`, do not use `--force` or `--force-with-lease`.
````

## Dispatch subagents

Subagents must work in git worktrees on their story branch. Use the following template:

```
Implement the story @specs/<feature-slug>/issues/<issue-file>.md .
Load the `implement` and relevant coding skills.

Return: PR URL, one-line summary, any issue encountered (anything that caused several back-and-forth exchanges, or anything unexpected in an otherwise smooth development journey).

---

Create your worktree with this command: `git worktree add <tmp-path>/wt-<feature-slug>-<story-nn> -b <branch-name> --no-track origin/<feature-branch>`

Push with this command: `git push origin HEAD:refs/heads/<branch-name>`.
**Do not use `git push -u`, do not use bare `git push`, do not use `--force` or `--force-with-lease`.**

Use `gh` to create and interact with the PR in GitHub.
```

That's it.

For addressing PR review, it becomes:

```
Address the review on the PR for the branch `<branch-name>`, implementing the story @specs/<feature-slug>/issues/<issue-file>.md .
Load the `implement` and relevant coding skills.

Return: PR URL, one-line summary, any issue encountered (anything that caused several back-and-forth exchanges, or anything unexpected in an otherwise smooth development journey).

---

Use the existing worktree: `cd <tmp-path>/wt-<feature-slug>-<story-nn>`

Push with this command: `git push origin HEAD:refs/heads/<branch-name>`.
**Do not use `git push -u`, do not use bare `git push`, do not use `--force` or `--force-with-lease`.**

Use `gh` to create and interact with the PR in GitHub.
```

## Wrap up

The wrap up can be invoked by the user before all stories are completed.

1. squash and rebase: single commit unless the user says otherwise. Use the squashed commit messages to build a coherent explanation.
2. create a PR targeting `main`.
    1. dispatch a subagent to review the final PR using the `code-review` skill (and any other relevant skill); the agent will write comments on the PR.
    2. ask the user to review and merge.
3. immediately start the _next batch_.

Perform all 3 steps, then ask the user to review the new PRs. Dispatch the reviewer agent and the developer agents simultaneously, in parallel.

When you rebase and force-push, keep a reference to the previous commits in case something goes wrong and the changes need to be recovered.
