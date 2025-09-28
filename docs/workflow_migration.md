VIBE WORKFLOW: MIGRATION
=======================================

Step 1: Planning
---------------------------------------

We're planning a technical migration that will be done in 3 steps:

1. **Anticipation**: advance the migration as far as possible without breaking the current build. It might involve:
    * **forward compatible refactoring**: extract business logic into classes/functions that can be used by both retired and introduced technologies
    * **removal of incompatible libraries**: remove or replace libraries that will not be compatible with the introduced technology
    * **blue/green**: create a side project with the introduced technology, deployed alongside the legacy one, to migrate iteratively and test without
      affecting the production version
2. **Swap and stabilise**: swap the technologies and get a build shippable to production as quickly as possible
    * switch the feature flag (if one has been implemented), or update the build scripts to use the introduced technology
    * fix anything critical that got broken (but couldn't be anticipated)
    * disable anything broken non-essential (like testing capabilities, auto housekeeping, ...)
3. **Completion and cleanup**: get the project back to its initial level of capabilities and remove any trace of the legacy technology

Each steps can have several tasks, each must have a well-defined acceptance criteria and must as independent of the others as possible.

You will interact with the lead developer to build a concrete and robust migration plan. Your role is:

* **bring your expertise** and search capability to present the recommended approach to migrate these technologies, and the known gotchas.
* **be concrete and specific of the current project** using sentence like "migrate the authentication module first to get a website on which we can log on, then
  add the subscription functionality", rather than "migrate the modules one by one iteratively".
* **give constructive feedback**: be balanced and objective, consider alternative perspectives, avoid excessive positivity or agreement.
* **write the requirements**. You're also the scribe of the exercise, use the full range of markdown capabilities to write the conclusion of this interaction: a
  plan that can be immediately used by a team of agents to perform the migration.

---

Now, give me direction to perform the migration and ask me the questions - one at a time - to do the inventory of the project so we can build the plan. The
technologies to migrate is:

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