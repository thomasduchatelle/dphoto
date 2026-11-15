# 04 — Tag switcher on the Tags page

Status: ready

## Description

Let the owner view the medias of any of their tags on the Tags page, not just `favorite`.

Add a switcher that lists the owner's tags grouped by namespace (with `favorite` surfaced distinctly, from the list-owner-tags endpoint) and, on selection, shows the grid of all medias carrying that tag. The selected tag is reflected in the URL so a tag view is shareable/bookmarkable (the normalised tag is URL-safe).

Load the `nextjs` and `ui-components` skills.

## Dependencies

02.3, 05.

## Acceptance criteria

- The Tags page offers a switcher of the owner's tags grouped by namespace.
- Selecting a tag shows all its medias across albums; the tag appears in the URL.
- Landing on the page still defaults to `favorite`.
- Covered by Storybook/visual and unit tests.

## Out of scope

- Tag counts.
- Filtering by several tags at once; crossing with location.
- Editing tags from this page.
