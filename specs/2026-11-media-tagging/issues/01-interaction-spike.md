# 01 — Spike: album-page tagging interaction

Status: ready

## Description

Throwaway prototype to decide **how an owner applies tags to medias from the album page**, where speed is the whole point. Explore several radically different interactions and present them for the owner to pick one. Bring options the spec did not consider — do not just build the "paint mode" assumption from the design.

Load the `prototype` skill. This is exploratory: no production code is merged. The deliverable is a decision.

Interactions to contrast should span distinct mental models, e.g.: pick-a-tag-then-click-to-paint; multi-select-then-apply; inline chip editor on a media; a command-palette / quick-switcher; hover quick-tags from recently used; drag medias onto a tag. Include at least one the design did not mention.

Prototype against the existing album media grid (`web-nextjs/app/(authenticated)/albums/[ownerId]/[albumId]/_components/AlbumMediaGrid/`, `MediaThumbnail.tsx`). Favourite (heart) is a separate always-on gesture and is not part of this exploration.

## Dependencies

None.

## Acceptance criteria

- At least three materially different interactions are demonstrated well enough to compare feel and speed.
- One of them is new relative to the spec/design.
- A live, clickable comparison is presented to the user; the user selects one.
- The chosen interaction (and the key reasons) is recorded in `specs/2026-11-media-tagging/design.md` under "Chosen interaction".

## Out of scope

- Any backend, API, or persistence work.
- Merging production code.
- The favourite/heart gesture.
- Namespace taxonomy or tag normalisation rules (already decided in the design).
