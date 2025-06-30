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