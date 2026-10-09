# 15 — Re-randomise covers UI on the album page

Status: wontdo
Phase: 3
Layer: web — `web-nextjs`
Depends on: 07, 14

## Description

Add a control on the album page (the grid of all pictures of a single album) that lets the owner
re-randomise the album's covers. On success, the album card's covers reflect the new set without a full
album-list reload.

See `../spec.md` → User journeys → Re-randomise covers.

## Acceptance criteria

- A visible action on the album page (grid-of-all-pictures view). Follow the `ui-components` and
  `nextjs` skills for placement, styling, and the thunk/action/selector pattern.
- Visible to the owner only — hidden or disabled for visitors. The permission context is already
  available in the app shell (follow the same pattern used by other owner-only controls, e.g.
  rename/share).
- On click, calls `POST /api/v1/owners/{owner}/albums/{folderName}/covers/refresh`.
- On success:
  - Updates the client-side state so the album's `covers` field reflects the response body.
  - Shows a brief confirmation (snackbar / toast — follow existing patterns).
- On error, shows an error snackbar; the covers are left unchanged.
- Storybook stories cover:
  - Owner view with the action visible and enabled.
  - Visitor view with the action absent (or disabled).
  - Loading state while the request is in flight.
  - Error state.

## Out of scope

- Cherry-pick / unpick controls — Phase 4+.
- Any change to the album card rendering (covered by issue 07).

## References

- `../spec.md` (Phase 3; User journeys → Re-randomise covers).
- `../design.md` (Phase 3 REST contract).
- Load skills: `nextjs`, `ui-components`.
