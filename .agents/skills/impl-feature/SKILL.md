---
name: impl-feature
description: Deliver a whole feature (a spec with stories/issues) by dispatching one subagent per independent story, one PR per story, and owning the review-and-rebase loop across all PRs until merged. Use when the user asks to implement a feature made of several stories/issues that can be worked in parallel, or explicitly asks to "run the stories" or "spawn subagents for each subtask".
---

# Delivering a feature via parallel subagents

You are the **conductor**. Subagents write the code, one per story; you own dispatch, the review-loop, cross-PR feedback fan-out, and merge conflicts.

## The pipeline

1. Read the feature: `specs/<slug>/spec.md`, `specs/<slug>/stories.md`, every `specs/<slug>/issues/NN-*.md`. Extract the dependency graph.
2. The current git branch is the **base branch**. Every subagent PR targets it. Create it and push it if you're still on `main`.
3. Group stories into **batches**: everything in a batch has all its blockers already merged, and runs in parallel.
4. Dispatch one subagent per story ([recipe below](#dispatching-a-subagent)).
5. When the batch's PRs are open, list them to the user ([handoff](#batch-handoff)) and wait.
6. On user feedback, propagate it ([propagating feedback](#propagating-feedback)).
7. When PR(s) have been approved (or at least one of them), [wrap-up the (potentially partial) batch](#batch-wrap-up)
8. If stories remain: loop back to (3) to find the next batch.
9. If all stories are completed: clean-up GIT and local drive of the branches and worktree that have been merged.

## Dispatching a subagent

Each subagent handles one story end-to-end: worktree → implement → test → commit → push → PR. Fill in the template below (`<...>` placeholders) and pass it as the subagent's prompt. Give it only its own issue file — never the whole spec's list of issues.

### Prompt template

> You are implementing one story of the `<feature-slug>` feature for the DPhoto repo.
>
> ## Read
>
> - `specs/<feature-slug>/issues/<issue-file>.md` — your task.
> - Load the `implement` skill.
> - Load the skill(s) relevant to the paths this story touches (e.g. `go` for `pkg/`, `nextjs` for `web-nextjs/`, ...).
>
> ## Setup
>
> Create a fresh worktree and a fresh branch off `origin/<base-branch>` — `--no-track` is mandatory (see the note below):
>
> ```
> cd <base-worktree-path>
> git fetch origin
> git worktree add <tmp-path>/wt-<feature-slug>-<story-nn> \
>     -b <branch-name> --no-track origin/<base-branch>
> ```
>
> All subsequent commands run inside that worktree.
>
> ## Task
>
> Implement the issue. Follow the design principles from the skills you loaded — no shortcuts on test structure (fakes, and testing strategy) and architecture.
>
> ## Verification
>
> Before committing, from the worktree root:
>
> ```
> <exact test command, e.g. go test ./pkg/<pkg>/...>
> ```
>
> Must be green. Also verify: `<any story-specific check, e.g. "no test file imports internal/mocks">`.
>
> ## Commit and push
>
> One commit, message per the `implement` skill.
>
> **Push command — use this exact form, nothing else:**
>
> ```
> git push origin HEAD:refs/heads/<branch-name>
> ```
>
> Do not use `git push -u`, do not use bare `git push`, do not use `--force` or `--force-with-lease` (force-push is not allowed on a first attempt — the conductor will tell you if it becomes needed).
>
> ## PR
>
> Create the PR with `gh`, base `<base-branch>`, head `<branch-name>`.
> Title with the name of the story and the headline of the change.
> Body must help to reviewer to understand the change: explain how it works (data flow, main interfaces), the intention of some changes (renamed class, ...), and why would the reviewer trust the test (if existing test have been changed: why no regression will be caused anyway, if new tests, how they covers the acceptance criteria). Also reference `specs/<feature-slug>/issues/<issue-file>.md` in the body.
>
> ## Final output
>
> Return: PR URL, one-line summary, whether tests passed, any issue encountered (something that caused several back and forth must be mentioned, or something that wouldn't be expected in a smooth development journey).

## Propagating feedback

Reviewer comments on one PR often reveal a pattern that hides in the sibling PRs. Steps:

1. Classify the nature of the request: local changes (skills, stories, ...) or standard review on the PRs. It can be both.
2. **Local changes**:
    1. update the local branch (new commit) with the requested changes
    2. then dispatch subagents to update all the PRs (same process as standard reviews) ; this is the only case where a rebase and forced push is allowed.
3. **Standard reviews** - user's comments might be directly in the chat, or in the PRs ; and they can affect several PRs (a pattern on several PRs) even if they have been written on the PR.
    1. find what are the PRs are affected.
    2. dispatch a subagent for each affected PR, with the requested change in details. Subagents must use GitHub comment to interact with the user:
        1. close the comments if accepted and fixed as part of the update.
        2. pushback and propose a plan if the request will make the code worst for a reason the user might have missed.
4. Summarise the change for each PR, especially the pushbacks.

### Prompt template — standard feedback

Default template for a reviewer comment or user-requested tweak on a single PR. The subagent adds one follow-up commit; no rebase, no force-push.

> You are addressing review feedback on PR #`<pr-number>` (story `<issue-file>`) of the DPhoto repo.
>
> ## Read
>
> - `specs/<feature-slug>/issues/<issue-file>.md` — the story.
> - Load the `implement` skill.
> - Load the skill(s) relevant to the paths this story touches (e.g. `go` for `pkg/`, `nextjs` for `web-nextjs/`, ...).
>
> ## Setup
>
> Reuse the existing worktree for this story: `<tmp-path>/wt-<feature-slug>-<story-nn>`. It is already on branch `<branch-name>` and tracks `origin/<branch-name>`. All commands run inside that worktree.
>
> Before starting, sync the branch with the remote in case another loop pushed to it: `git pull --ff-only`.
>
> If the worktree path is missing or the branch state looks off (detached HEAD, unexpected local commits, etc.), stop and tell the conductor — do not try to reconstruct it yourself.
>
> ## Feedback to address
>
> `<verbatim reviewer request, or a distilled version. Include the specific files/lines flagged. For each thread you want resolved, include its numeric comment ID and thread ID (PRRT_...).>`
>
> ## Task
>
> Apply the changes. Follow the design principles from the skills you loaded — no shortcuts on test structure (fakes, and testing strategy) and architecture.
>
> For each reviewer thread listed above:
>
> - If you can apply the change as requested, do so, then reply on the thread confirming the fix and **resolve the thread**.
> - If the change would make the code worse (violates a documented principle, requires much more complex or deeper changes, breaks another test, introduces a regression) or you believe the reviewer missed something, **push back**: reply on the thread with a short explanation and a counter-proposal. **Do not resolve the thread** — leave it for the conductor to handle. Do not silently ignore the request.
>
> ## Verification
>
> Before committing, from the worktree root:
>
> ```
> <exact test command, e.g. go test ./pkg/<pkg>/...>
> ```
>
> Must be green. Also verify: `<any story-specific check>`.
>
> ## Commit and push
>
> Add one follow-up commit describing the fix, message per the `implement` skill. Push:
>
> ```
> git push origin HEAD:refs/heads/<branch-name>
> ```
>
> Do not use `git push -u`, do not use bare `git push`.
>
> ## Final output
>
> Return: PR URL, one-line summary of what changed, whether tests passed, which threads you resolved, which threads you pushed back on (and why), any issue encountered.

### Prompt template — rebase-and-force-push

Use this variant **only** when the conductor needs the PR rebased on a moved base branch (e.g. after a policy-driven skill update on the base). Same as the standard template above, with these two sections replaced:

> ## Setup
>
> Reuse the existing worktree for this story: `<tmp-path>/wt-<feature-slug>-<story-nn>`. It is already on branch `<branch-name>` and tracks `origin/<branch-name>`. All commands run inside that worktree.
>
> Before starting, sync with the remote and rebase on the updated base:
>
> ```
> git pull --ff-only
> git fetch origin
> git rebase origin/<base-branch>
> ```
>
> Resolve any conflicts. If a conflict is not obvious, stop and tell the conductor. If the worktree path is missing or the branch state looks off (detached HEAD, unexpected local commits, etc.), stop and tell the conductor — do not try to reconstruct it yourself.
>
> ## Commit and push
>
> `<one of, pick before dispatching:>`
> - Add one follow-up commit describing the fix, message per the `implement` skill.
> - Amend the existing commit (if the conductor asked to keep a single commit).
>
> Push with force-with-lease (rebase makes force required):
>
> ```
> git push origin HEAD:refs/heads/<branch-name> --force-with-lease
> ```
>
> Do not use `git push -u`, do not use bare `git push`, do not use `--force` (without `-with-lease`).


## Batch wrap-up

Once you got the approval from the user for at least one PR:

1. merge approved stories in the current branch ; **do not merge them to `main`, verify the PR before merging it**.
2. update the issue-tracker if necessary, stories completed should be `done` ; if the feature is completed, archive it ; if a feature is not refined (next stories are missing) inform the user.
3. squash and rebase: single commit unless told otherwise by the user. Use the messages of the squash commits to create a consistent explanation.
4. create a PR toward the `main` branch.
   1. dispatch a subagent that will review the final PR using the `code-review` skill and any other relevant skill, and publish its comments on the PR itself.
   2. request the user to review and merge.
5. and, in parallel while the reviewer is working:
    1. start the next batch ! Dispatch subagents for the next stories.

When you use rebase and push force, keep a reference of the previous commits in case something goes wrong, and we need to recover the changes.

This step is slow, do not wait for user approval at each step: the subagents (review and next batch) must run in parallel, and the only decision point of the user will be on the PR directly.


## Batch handoff

When every PR in the batch is open (or re-pushed), stop. Post a table to the user:

| # | Story | PR   |
|---|-------|------|
| … | …     | link |

Flag anything unusual: force-pushes you authorised, deleted tests, review threads left unresolved. Wait for the user's next instruction — fixes to make, or the go-ahead to merge and launch the next batch.

