VIBE WORKFLOW: WEB FEATURE
=======================================

Step 1 - Requirements collection
---------------------------------------

```
aider --model openrouter/anthropic/claude-sonnet-4 --map-tokens 0 
/read-only web/src/core/catalog/language
/ask
```

We're starting a new feature, and you have to produce the Requirement Document. The document will be used for planning and must gather the following
information:

* What the feature **must do**?
* What are the **specific UI components, and user flows** needed?
* What are the interactions with external systems (user, APIs, ...), including what happens before and after each operation?
* What is unknown and would need to be explored first?

The document produced must follow the structure:

1. Feature **Summary** (1-3 sentences)
2. **Ubiquity Language** (only the new terms/entities, skip the section if empty)
3. **Scenarios** (ideally 5, maximum 8): For each, provide a complete, step-by-step user journey from start to finish, showing how the user accomplishes the
   goal using the
   feature. They should cover permissions or restriction cases, and error handling.
4. **Technical Context**: bullet points of the things provided by the platform, supporting APIs, other domains, ...
5. List of **Explorations**: formulated as questions in bullet points

To build this requirement document, ask me one question at a time, focusing on concrete implementation details and user journeys. Each question should build on
my previous answers, and the process should continue until all relevant details are gathered for a complete, actionable requirements document.

Here's the idea:

Step 2 - Planning
---------------------------------------

```
aider --model openrouter/anthropic/claude-sonnet-4 --map-tokens 0 
/read-only web/src/core/catalog/language
/read-only docs/principles_web.md
/read-only docs/feature_edit_album_claude_sonnet_4.md
/ask
```

We've just finished the requirement gathering phase and we are now moving to the planning phase. The objective is to write a detailed and iterative list of
prompts that are actionable and testable individually by an LLM.

Start by breaking down the data flow and requirements, using the Principle Handbook, into the different layers described in the handbook. Make sure each layer
is **never going beyond its responsibilities**, and **never leak its responsibilities** to other layers. **The principles are absolute and must be respected.**

Once you have the flow broken down by layer, for each layer, break it down further into a list of tasks to develop. Each task must be small enough to fit into
the smallest LLM context, but it must progress the project forward by being a **unit of work implementable and testable autonomously**. Order the tasks to start
with structural ones on which the others will be built from.

Then, write a prompt for each of the task. One prompt per task that:

1. describes in details the requirements expected from the task using BDD style `GIVEN ... WHEN ... THEN ..`.
2. describes the names of the functions, interfaces, classes that will be implemented and exported as part of the task, and the filenames where they should be
   implemented.
3. gives the references and description of the previous tasks **relevant** to implement this tasks. Only the references that are expected to be used.
4. quote the principles from the handbook that must be followed to implement the task. Only the part relevant for the layer of the task.
5. insist on the TDD approach: no code should be written if no test make it necessary.

The principle handbook is:
The requirement document is: