VIBE WORKFLOW: WEB FEATURE
=======================================

Step 1 - Requirements collection
---------------------------------------

```
aider --map-tokens 0 --model openrouter/google/gemini-2.5-flash-preview-05-20
#aider --map-tokens 0 --model openrouter/meta-llama/llama-4-maverick
/read-only web/src/core/catalog/language
/ask
```

We're starting a new feature, and you have to produce the Requirement Document that gives strong and detailed vision of what you want to build and how it should
work. The document will be used for planning and must gather the following information:

* What the feature **must do**?
* What are the **user flows** (including what happens before and after each operation) ? How the user interacts with the feature (specific UI component, ...) ?
* How the feature is **integrated to the existing application**, or leveraging **external systems** ?
* What is unknown and would need to be explored first?

The document produced must follow the structure:

1. Feature **Summary** (1-3 sentences)
2. **Ubiquity Language** (only the new terms/entities, skip the section if empty)
3. **Scenarios** (ideally 5, maximum 8): For each, provide a complete, step-by-step user journey from start to finish, showing how the user accomplishes the
   goal using the feature. They should cover permissions or restriction cases, and error handling, data validation.
4. **Technical Context**: bullet points list of what is done by services external of this feature: provided by the platform, by supporting APIs, or by other
   domains or the same application, ... It should also contain an **out of scope** section.
5. List of **Explorations**: questions that couldn't be answered during this chat, and will have an impact on the design of the feature. Do not add question
   related to implementation details (good example: "what are the authorisation capabilities of the Provider X ? Feature will extend it only for the use cases
   not already covered."" ; bad example: "how to use AxiosJS?" because it has not functional impact)

To build this requirement document, ask me one question at a time, focusing on concrete implementation details and user journeys. Each question should build on
my previous answers, and the process should continue until all relevant details are gathered for a detailed and complete vision of the feature to build.

If asked, the requirements should be written in the file: `specs/YYYY-MM_<feature name>.md` in markdown (ask for the current date).

Here's the idea:


Step 2 - Refinement and Story Mapping
---------------------------------------

### Non-interactive

```
aider --model openrouter/meta-llama/llama-4-maverick --map-tokens 0 
/ask
```

Before starting the development of the new feature, you need to break it down into stories following these rules:

* the list of stories must cover the complete requirement: nothing should be left out
* each story must be an interation that bring the project forward, but small enough to be done by an LLM Agent
* each story must represent a vertical slice of the feature: an end-to-end journey, never a technical layer

Iterate as many times as you need until you have a strong and actionable list of stories. You will then write them in a clear and concise way with the following
structure:

* **title** (as header): use the pattern "As a [user], I want [feature]" (example: "As a user, I can open the delete dialog where I can see a list of deletable
  albums to choose from").
* **Acceptance Criteria**: be thorough, give examples, do not leave any behaviour free for interpretation (all commutations must be explained), write from user
  point of view (not technical), and **stay concise**. Use the BDD-style:
  ```
  GIVEN <description of the initial state>
  WHEN <name of the action dispatched and description of its payload>
  THEN <description of what will return the selector>
  ```

  Place them in an indented code block. Don't overcharge each statement: prefer multiple and simple "given ... when ... then". Add examples where relevant.

* **Out of scope**: reread the title and acceptance criteria and list here what an LLM would be tempted to do but is not in scope of this story (for example "
  validation of the fields is done by the underlying API)

Write each story in its own file named `specs/<requirement filename without extension>_story_<story number in tweo digits>_<few words description>.md`. Write
all the stories before asking for anything. Respect the instructions to write the file as you would do to make code change.

The requirement to refine is written in the file:

### Interactive

```
aider --model openrouter/meta-llama/llama-4-maverick --map-tokens 0 
/ask
```

Before starting the development of the new feature, you need to break it down into stories following these rules:

* the list of stories must cover the complete requirement: nothing should be left out
* each story must be an interation that bring the project forward, but small enough to be done by an LLM Agent
* each story must represent a vertical slice of the feature: an end-to-end journey, never a technical layer

Iterate as many times as you need until you have a strong and actionable enumerated list of stories. Ask your peering partner to review it. Give him the full
list of title at once. Listen and adapt the list based on its comments, feel free to ask clarification questions.

You will then write each story in a clear and concise way with the following structure:

* **title** (as header): use the pattern "As a [user], I want [feature]" (example: "As a user, I can open the delete dialog where I can see a list of deletable
  albums to choose from").
* **Acceptance Criteria**: be thorough, give examples, do not leave any behaviour free for interpretation (all commutations must be explained), write from user
  point of view (not technical), and **stay concise**. Use the BDD-style:
  ```
  GIVEN <description of the initial state>
  WHEN <name of the action dispatched and description of its payload>
  THEN <description of what will return the selector>
  ```

  Place them in an indented code block. Don't overcharge each statement: prefer multiple and simple "given ... when ... then". Add example where relevant.

* **Out of scope**: reread the title and acceptance criteria and list here what an LLM would be tempted to do but is not in scope of this story (for example "
  validation of the fields is done by the underlying API)

Present them one by one to ask for feedback after each one. Upon acceptance, write the story in the file named
`specs/<requirement filename without extension>_story_<story number in tweo digits>_<few words description>.md`, and present the next story.

**Respect the instructions to update files as you would do to make code change!**

The requirement to refine is written in the file:

### Review

```
aider
add specs/
/ask
```

You are the lead engineer of a software team. You are reviewing the **epic** and the **stories** that have been created. You need to make sure that the stories
are clear, that they cover the entirety of the epic, and that they is no details missing.

Give a short list (5 items max) of the most impactful improvements. Look after specifically the **Acceptance Criteria** (no duplicate, examples when
useful, ...).

The Epic is in the file:

The rest are the stories.

Step 3 - Feature Coding
---------------------------------------

### Interactive

```
aider
/add web/src/core/catalog/language 
    web/src/core/catalog/actions.ts 
    web/src/core/catalog/index.ts 
    web/src/core/catalog/thunks.ts  
    web/src/pages/authenticated/albums/CatalogViewerPage.tsx 
    web/src/pages/authenticated/albums/CatalogViewerRoot.tsx 
    web/src/pages/authenticated/albums/AlbumsListActions 
    web/src/pages/authenticated/albums/DeleteAlbumDialog 
    web/src/pages/authenticated/albums/MediasPage/index.tsx
/ask        
```

You are a strong developer prioritizing simple and well tested code. You **strictly follow the coding principles** defined in `docs/principles_web.md`.

You are **presenting to me the guiding principles of your design to implement the story**. You have access of the complete epic to give you context, but you
will **focus ONLY on what is required for the story**.

Using coding principles from the document, define the components involved to deliver the story. Give a thorough description of each of them, using code blocks
for code snippet for clarity:

* name
* if it's "new", or "updated", or "reused" (and not updated)
* if the component is a state or a domain model: give a code snippet of the properties to add or change
* if the component is an event or an action: give the schema of its payload
* if the component is a function: give its signature
* if the component is a UI component: explain what data is rendered, what data can be input, and how the user can interact with it

Your story you are designing is:

It's part of the epic:

### Interactive (v2)

```
aider
/add web/src/core/catalog/language 
    web/src/core/catalog/tests/test-helper-state.ts
    web/src/core/catalog/actions.ts 
    web/src/core/catalog/index.ts 
    web/src/core/catalog/thunks.ts  
    web/src/pages/authenticated/albums/CatalogViewerPage.tsx 
    web/src/pages/authenticated/albums/CatalogViewerRoot.tsx 
    web/src/pages/authenticated/albums/AlbumsListActions 
    web/src/pages/authenticated/albums/DeleteAlbumDialog 
    web/src/pages/authenticated/albums/MediasPage/index.tsx
/ask        
```

You are a strong developer prioritizing simple and well tested code. You **strictly follow the coding principles** defined in `docs/principles_web.md`.

You are **presenting to me the guiding principles of your design to implement the story**. You have access of the complete epic to give you context, but you
will **focus ONLY on what is required for the story**.

Using coding principles from the document, define the components involved to deliver the story. Give a thorough description of each of them, using code blocks
for code snippet for clarity:

1. **UI Component**: what data is required ? what thunks are required ?
   For each component (updated or created), write the interface representing the properties the component(s) need to receive.

2. **Thunks**: what data and ports will be required for the thunk(s) to perform their business logic
   For each thunk, gives the function signature, and the port interface.

3. **Actions**: what is the payload of each actions ?
   For each action, gives the interface representing it.

4. **State**: what changes are required in the state ? Provide a code snippet of the changes

5. **Selectors**: does it need new selector(s), or updating the existing ones ? What are the properties selected ? How are they captured ?
   Write the interface of the selection, add a comment from each comment to describe from what State properties it will come

Your story you are designing is:

It's part of the epic:

#### Then - implementation

```
/code Implement the story. **Strictly follow the coding principles from `docs/principles_web.md`**. Implement the tests to make sure the Acceptance Criteria are
covered. Do not add comments to your code. Once complete, leave a comment to the reviewer of places he needs to bring a special attention.
```

#### Finally - Code review

```
/model bb8
/reset
!git diff HEAD^
/ask You are the senior developer in charge of reviewing the code written by your peer. You use the `docs/principles_web.md`, and your personal knowledge, as 
    references of what a good code looks like. You promote clean code, well tested with a suite robust to refactoring, secure and performant. 
    
    Present your comments like in a Merge Request, with the file name and the code snippet.
```

---

NON VALIDATED DRAFT
=======================================

### Non interactive

```
aider --model openrouter/anthropic/claude-sonnet-4 --map-tokens 0 
/read-only docs/principles_web.md
# ... all the files required to change
```

You are a strong developer prioritizing simple and well tested code. You **strictly follow the coding principles** defined in `docs/principles_web.md`. And you
are now implementing a new story part of a larger epic. You will use the epic to contextualise your changes but will focus to only deliver what is requested in
the story.

**Process**

1. **Design Phase** - using the coding principles from the document:
    1. **data flow**: define the components (by name and type) involved to deliver the story
    2. **re-usability**: identify the _existing components_ that can be leveraged (function, events, ...). Do not be too eager: prioritise
       single-responsibility principle over trying to avoid code duplication.
    3. **technical design**: specify for each component a thorough description:
        * name
        * if it's "new", or "updated", or "reused" (and not updated)
        * if the component is a state or a domain model: give a code snippet of the properties to add or change
        * if the component is an event or an action: give the schema of its payload
        * if the component is a function: give its signature
        * if the component is a UI component: explain what data is rendered, what data can be input, and how the user can interact with it

2. **Implementation** - focus on the implementation of one component at a time. **Do not forget its tests.**
    1. Start with writing the tests following the BDD requirements
    2. Then write an implementation that pass the test
    3. Finally, move on to the next component

Do not wait for confirmation at any stage: write the complete implementation of the story immediately. Do not ask to run the tests.

Your story to implement is:

It's part of the epic:

Step X - Design and Planning
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

    * _Introduction_: "you are implementing ..." ; be specific of the type (or layer) and name of the component(s) the agent have to implement ; make explicit
      that the deliveries must include the tests validating the requirements.
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

Step 4 - Implementation
---------------------------------------

```
aider --model openrouter/anthropic/claude-sonnet-4
/read-only docs/principles_web.md
```

> paste the prompt from the file.