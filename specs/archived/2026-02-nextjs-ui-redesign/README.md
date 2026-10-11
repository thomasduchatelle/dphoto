# 2026-02 NextJS UI Redesign — archived

Archived on 2026-10-11.

This feature kicked off the `web-nextjs` rebuild. Issues 01–08 shipped the foundation (MUI, state migration,
album list page, album filtering, album page with day-grouped media grid). Mid-way through, the scope was
tightened: the remaining ambition (random photos, timeline navigation, visual date-picker, polish) was dropped
in favour of **reaching feature parity with `web/` and cutting over**.

The continuation lives in **`specs/2026-10-nextjs-cutover/`**: it carries the viewer stories (issues 09 and 10
here, preserved verbatim as 02 and 03 there) and adds the missing modals (create / share / rename / dates /
delete) plus the actual cutover (delete `web/`, drop the `/nextjs` basepath).

Deferred items (random photos, timeline navigation, visual date-picker, scroll-restore on viewer close,
album preview on hover) are tracked in `TODO.md`.

## Issue status at archive time

- `01` – `08`: **done**
- `09` – `10`: **wontdo** — moved verbatim to `2026-10-nextjs-cutover`
- `11`: **wontdo** — dropped; deferred to `TODO.md`
