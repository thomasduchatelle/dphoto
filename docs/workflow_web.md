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

Step 2 - Story Mapping
---------------------------------------

TBD

Delivery: list of stories "As a user, I can open the edit date dialog with for a given AlbumId and I can see on the dialog its name and its dates"

In scope and Out of Scope must be detailed - LLM has tendency to go too far.

Step 3 - Design and Planning
---------------------------------------

```
aider --model openrouter/anthropic/claude-sonnet-4 --map-tokens 0 
/read-only web/src/core/catalog/language
           docs/principles_web.md
           docs/feature_edit_album_claude_sonnet_4.md
/ask

tail -n +853 .aider.chat.history.md > docs/feature_edit_album_plan_0.1-claude.md
```

To develop a User Story, we're going to write a detailed and iterative list of prompts that are actionable and testable individually by an LLM.

**Process**:

1. **Design Phase** - use the Design Pattern concepts from the principle handbook to:
    * describe the data flow of the story: what each concept will require from the underlying layers
    * give the details of each concept required:
        * interfaces: name, with their schema and properties to add/modify
        * functions: name and signature
        * events / actions: name and payload schema
        * UI component: name and purpose

2. **Collaboration** - Present the design and ask for feedback before proceeding ; we will iterate on the design
   * exhaustive list of each component to create and its layer
   * if the component is a state or a domain model: give a code snippet of the changes
   * if the component is an event or an action: give the schema of its payload
   * if the component is a function: give its signature
   * if the component is a UI component: explain what it will contain

3. **Task Breakdown** - Create independently implementable tasks where each task describes **application behavior**, not developer tasks:
    * GOOD: "GIVEN dialog is closed WHEN I dispatch editDatesDialogOpened with AlbumId THEN selectEditDatesDialog returns open dialog with album name"
    * BAD: "GIVEN system needs dialog support WHEN defining state model THEN selector should return properties"

   Each task must be:
    * A unit of work (1 action OR 1 thunk OR 1 component) ; selectors and state change are part of the action task that requires it.
    * Independently testable
    * Described in BDD format focusing on runtime behavior and writen in a code block. Example for an action:
      ```
      GIVEN <description of the initial state>
      WHEN <name of the action dispatched and description of its payload>
      THEN <description of what will return the selector>
      ```

4. **Collaboration** - present the tasks and ask for feedback

5. **Prompt Structure**:

   * _Introduction_: "you are implementing ..." ; be specific of the type (or layer) and name of the component(s) the agent have to implement
   * _Requirements_: the BDD-style requirements defined and reviewed on the previous step
   * _Implementation Details_:
      * in what folder the new components must be created (feature related), insist on the naming convention to be respected
      * add the TDD principle: "Implement the tests first, then implement the code **the simplest and most readable possible**: no behaviour should be
        implemented if it is not required by
        one test"
      * list the general files that must be edited, and for each what's expected ("general" because not specific for the feature: global state, actions/thunks
        register, ...)
      * any recommendation raised during the design
   * _Interface Specification_: data structure and signature that have been decided during design
   * _References_: gives the references and description of the previous tasks **relevant** to implement this tasks. Only the references that are expected to be
     used.

The principle handbook is:

The requirement document is:

The story to work on is:

#### Summary

```
/code The break down is good. Write the prompts in `docs/prompts_edit_dates/task_<number>_<two words summary>.md` in markdown, one prompt per file.
```

#### Review (optional)

```
/reset
/add `docs/prompts_edit_dates`
/ask
```

Review the documents in `docs/prompts_edit_dates`. They are prompts to implement a feature and will be consumed by an LLM agent.

Find and list any inconsistency, and anything that could be misleading for the agent.

Then propose a solution for each.