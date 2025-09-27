VIBE WORKFLOW: MIGRATION
=======================================

Step 1: Evaluation & Scoping
---------------------------------------

We're starting a new technological migration, and you have to produce the Technical Planning Document that details the plan, as a series of tasks, that will be
followed and implemented by agents to completely migrate from the legacy technology.

You need to ensure that the project is never far (in time and effort) to be shippable to production by leveraging side-by-side deployments, feature flags, and
pre-migration refactors whenever possible.

Start by listing the changes that are required to perform the migration:

* use your knowledge of the legacy and the new technology: how they differ, and how migrations are recommended
* inspect the code, and ask targeted questions to identify precisely, for this project, what needs to be done

Then, write 3 list of tasks:

* **forward compatible refactoring**: tasks that can be done before the new technology is introduced to make the code forward compatible. Be creative:
  remove/replace non-compatible libraries, extract business logic from code too coupled with legacy technology, ...
* **swap technology and stabilise**: tasks to introduce the new technology and get a project shippable to production as quickly as possible even if it's
  degraded: side deployment not feature complete yet, missing testing capabilities, ...
* **completion and clean up**: tasks to get the new migrated codebase at the same level the previous one was (features and testing capabilities), and remove all
  trace of the previous technology

You should prefer using your knowledge and inspecting the code but you might ask me some questions if you need details, one at a time so you can build on my
previous answer to ask the next question.

---

Now, we're migrating from Create-React-App to Waku framework. 

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