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

