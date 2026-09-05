I'm hiring a UX designer that will redesign the entire website (currently in `web/`).

You have to write two documents.

The first is `specs/2026-01-ux-functionnal.md`. you need to write in plain English a precise description of what DPhoto is, and each interaction a user can have
with it.
Do not mention _how_ (what page, or what button), only _what_ a user can do (ex: a user can see his albums. he can filter the one he owns. he can share with
someone else. An album has a name, start date, end date, ...). This is the only document the UX designer will have access to in order to make the design of the
whole website. You must be precise and exhaustive, without any implementation details, only functional.

The second is `specs/2026-01-ux-reference.md`. You will write it using TypeScript code snippets to describe the operations, and their signature, which the
frontend
developer will use to write the callbacks of his components. The frontend developer will implement the UI components, focused on the visual and without any
logic and behaviour, to fulfil the UX design. It's critical to give him a precise descrition of what behaviours will be attached to his compoenents: they will
be his compoenents' properties ( data to display and callback function, or thunks).

You will only use the deprecated website code (in `web/`) as source to generate this two documents.

1. first, read the following files (and directory), don't be lazy, read them thoroughly
    1. `.github/instructions/nextjs.instructions.md`: the coding standards apply on the old website and explain the concepts of the thunks. They
       are the key of the user interactions that the UX designer must know, and what the front developer will have to use.
    2. `web/src/core/catalog/language/**`: domain model used by the UI, critical to understand what we're speaking about
    3. `web/src/core/catalog/**`: everything else about the domain, you can identify the callback, the reads, ...
2. then summarise all your findings into the two documents:
    1. `specs/2026-01-ux-functionnal.md`: plain English to describe what a user can do using the UI that the UX designer must design
    2. `specs/2026-01-ux-reference.md`: api reference of what will be made available to the frontend developer to power his UI component, focused on the visual.
