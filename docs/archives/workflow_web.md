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

### Non-interactive

```
aider --model six
/add web/src/core/catalog/language 
    web/src/core/catalog/tests/test-helper-state.ts
    web/src/core/catalog/actions.ts 
    web/src/core/catalog/index.ts 
    web/src/core/catalog/thunks.ts
    
/add <existing UI components, actions or thunk to update or use, ...>
```

You are a strong developer prioritizing simple and well tested code. You **strictly follow the coding principles** defined in `docs/principles_web.md`.

You are implementing the user story linked below. You have access of the complete epic to give you context, but you will **focus ONLY on what is required for
this story**.

Follow the important design principles below:

1. Do not update the existing tests: create new ones.
2. Do no add comments in your code.
3. Only add on the UI components the properties they **require** for the **current story**
    * Data must come from a selector
4. Keep the payload of action and thunk **minimum**:
    * _actions_ payload is only what's new on the state, or the Identifiers required to update it
    * _thunks_ arguments is only what it is required to perform the business logic (REST requests, making decisions, ...) ; only extract from the state what's
      needed ; only take as parameters what cannot be derived from the state
5. Write tests that are both readable and robust against refactoring
    * use the tests to make sure **the story's Acceptance Criteria are covered**
    * **always test together** actions and selectors: state is considered as private (robust: changing it won't affect the tests). The result of a selection is
      tested as a whole to prevent unexpected regressions.
    * **use pre-defined constants** of the states and the selections in a known situation: when adding new properties, only these constants are updated so it
      doesn't affect the tests (robust)

Once your listing complete, provide a list of next step considerations (only if you have any):

* ask questions that would help to refine the behaviour and improve user experience, security, or error handling
* suggest simplifications or improvements that would involve files out of the current context (and couldn't be done otherwise)
* suggest behaviour tests or acceptance tests to be added if you identify a case that would require them, justify them thoroughly: what sequence they will test and why they are required.

The story you are implementing is:

It's part of the epic:


### Interactive

```

aider
/add web/src/core/catalog/language
web/src/core/catalog/tests/test-helper-state.ts
web/src/core/catalog/actions.ts
web/src/core/catalog/index.ts
web/src/core/catalog/thunks.ts
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

Present your design, ask your question(s) **only if you have any**, and request the file(s) you will need to update to be added on the chat (**only the
files requiring to be edited**).

Your story you are designing is:

It's part of the epic:

#### Then - implementation

```

/code Implement the story. **Strictly follow the coding principles from `docs/principles_web.md`**. Implement the tests to make sure the Acceptance Criteria are
covered. Once complete, leave a comment to the reviewer of places he needs to bring a special attention.

Do not update the existing tests: add new ones. Do not add comments to your code.

```


---

NON VALIDATED DRAFT
=======================================

#### Finally - Code review

```

/model bb8
/reset
!git diff HEAD^
/ask You are the senior developer in charge of reviewing the code written by your peer. You use the `docs/principles_web.md`, and your personal knowledge, as
references of what a good code looks like. You promote clean code, well tested with a suite robust to refactoring, secure and performant.

    Present your comments like in a Merge Request, with the file name and the code snippet.

```