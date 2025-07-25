VIBE WORKFLOW: MIGRATION
=======================================

Step 1: Evaluation & Scoping
---------------------------------------

We're starting a new technological migration, and you have to produce the Requirement Document that gives strong and detailed vision of the target and how it
should be delivered. **We prioritise iterative migration where small chunks can be deployed and immediately bring value**, instead of waiting the project to be
complete and deployed as a major release.

Your role is to use your knowledge, experience, and problem-solving skill to understand what I want to get to and help me to structure a complex migration into
milestones. Each milestone will be mergeable into the main code base to bring a value (even minor). Think out of the box to solve the problem: pre-requisites
that can be done before introducing the new technology, side-by-side deployments, feature flags, ... everything should be considered.

You will then produce a document to report our discussions and summarise our plan, which will then be used for planning. We should have answered the following
questions:

* What the technology(ies) or practice(s) being retired, and what's the one getting introduced ?
* How will the application look like once migrated ?
* What are the migration milestones to continuously merge the chunks to migrate iteratively the application ?
* What are the most complex part of this migration ? How can they be mitigated ?

The document produced must follow the structure:

1. Migration **Summary** (1-3 sentences)
2. **Technology landscape**: technologies retired, upgraded, introduced ; and the impact on other components
3. **Milestones**: for each, write:
    * a **one-sentence outcome** summary using clear, factual, and concise declarative statements (examples: "Terraform 0.x is not used
      anywhere on the application", "A developer can run locally the whole stack")
    * a **list of the main tasks** to achieve this outcome, not the one that are obvious from the outcome summary, only the one that brings better understanding
      of the milestone.
    * any notes or consideration we had during our discussion

Ask me one question at a time. Each question should build on my previous answers, and the process should continue until all relevant details are gathered for a
detailed and complete vision of the migration to execute.

Here's the migration idea:

---

Write the document in the file: `specs/YYYY-MM_<feature name>.md`. Write it completely with everything we spoke about. We will review it and edited it after.

Step 2: Restructure existing code
---------------------------------------

### Pinning IDs

```
aider
```

The CDK project do not respect the principles from `docs/principles_cdk.md` and we're starting a migration of the exiting code to re-align the code.

The first step is to pin the logical IDs of the resource managed by the `InfrastructureStack`. Write a test asserting them, if you don't know the IDs, leave a
placeholder I'll complete them.

### Design

```
/ask
```

You're in charge of updating the CDK project structure to **strictly follow the principles `docs/principles_cdk.md`**. You are to present the high level target
design
that will respect the principles. Do not describe _how_ to migrate, only _where_ you are aiming: file structure with short description of what's expected in
each file.

The main domains are:

* **archive**: storage and access to raw medias and cached transformations
* **catalog**: management and listing of albums and medias
* **access**: access control to the service, user management, authentication, CLI access

### Execution

```
/code
```

Execute the plan we discussed. Do not leave anything behind. Do not edit the files to be deleted or moved: create the new file and provide a script to delete
old files. Only make the strictly necessary changes on tests to follow the new structure.

Step 3: Implementation
---------------------------------------

### Non-interactive - CDK

```
/read-only docs/principles_cdk.md
```

You're a senior engineer with strong knowledge of AWS and CDK who **strictly follow the principles in `docs/principles_cdk.md`**.

You are charged of implementing the "Milestone 2: CDK Infrastructure Parity" from `specs/2025-06_CDK-migration.md`.

* create or update any file required to **fully complete this task**
  Provide a shell script to move or delete files
* cover your changes by writing tests **strictly following the testing strategy**
* **implement only what is explicitly defined in the milestone scope** - do not add features, properties, or changes beyond the specified requirements

Once complete, leave a comment to the reviewer of places he needs to bring a special attention, and a suggestion of next steps.