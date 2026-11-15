# 03 — Tagging mode for namespaced tags

Status: ready

## Description

Implement the general tagging flow for namespaced tags on the album page, using the interaction chosen in the spike (recorded in `specs/2026-11-media-tagging/design.md`, "Chosen interaction"). The description below assumes the working-assumption "paint" model; adjust to the chosen interaction.

The owner enters tagging mode (e.g. via `t` or a dedicated button) and picks a tag with a type-ahead picker that autocompletes existing tags grouped by namespace, offers recently-used tags as quick-picks, and lets them create a new `namespace:value` tag on the fly. Once a tag is active, a header banner shows it and clicking a media toggles that tag on it. `ESC` or a button exits the mode.

Recently-used tags are client-side only (last 5–10, reset on refresh). Favourite remains its own star gesture and is unaffected.

Wire to the real add-tag / remove-tag endpoints and the autocomplete source for existing tags.

Load the `nextjs` and `ui-components` skills.

## Dependencies

01, 02.1, 02.2.

## Acceptance criteria

- Entering/exiting tagging mode works via the chosen trigger and `ESC`/button; the active tag is clearly shown.
- The picker autocompletes existing tags (grouped by namespace) and creates a new normalised `namespace:value` tag on the fly.
- Recently-used tags are offered as quick-picks (client-side, reset on refresh).
- While in mode, clicking a media toggles the active tag; state reflects server truth with optimistic update and revert on failure.
- Covered by Storybook/visual and unit tests.

## Out of scope

- The favourite/star gesture (already shipped).
- The Tags page and its switcher.
- Server-side multi-select.
