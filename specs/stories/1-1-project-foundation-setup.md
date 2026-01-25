# Story 1.1: Project Foundation Setup

Status: ready-for-dev

## Story

As a developer,
I want to set up Material UI with the dark theme and remove Tailwind CSS,
So that the project has a consistent design system foundation.

## Acceptance Criteria

**Given** the project currently uses Tailwind CSS
**When** I set up the Material UI foundation
**Then** Tailwind CSS dependencies are completely removed from package.json and configuration files
**And** Material UI (@mui/material ^6.x, @mui/icons-material ^6.x) is installed
**And** Emotion dependencies (@emotion/react ^11.x, @emotion/styled ^11.x) are installed
**And** MUI theme is configured in `components/theme/theme.ts` with:

- Dark mode as default (`mode: 'dark'`)
- Brand color #185986 as primary color
- Background #121212 and surface #1e1e1e
- Text colors (primary: #ffffff, secondary: rgba(255,255,255,0.7))

**And** MUI breakpoint system is configured: xs (<600px), sm (600px), md (960px), lg (1280px)
**And** ThemeProvider wraps the application in root layout
**And** No Tailwind classes remain in any component files
**And** The application builds successfully with `npm run build`
**And** Unit tests pass with `npm run test`

## Tasks / Subtasks

- [ ] Remove Tailwind CSS dependencies (AC: All)
    - [ ] Uninstall `tailwindcss` and `@tailwindcss/postcss` from package.json
    - [ ] Remove tailwind.config.js if exists
    - [ ] Remove Tailwind imports from global CSS files
    - [ ] Remove any Tailwind utility classes from existing components

- [ ] Install Material UI dependencies (AC: All)
    - [ ] Install @mui/material ^6.x
    - [ ] Install @mui/icons-material ^6.x
    - [ ] Install @emotion/react ^11.x
    - [ ] Install @emotion/styled ^11.x

- [ ] Create MUI theme configuration (AC: All)
    - [ ] Create `components/theme/theme.ts` file
    - [ ] Configure dark mode as default
    - [ ] Set brand blue (#185986) as primary color
    - [ ] Configure background colors (#121212 base, #1e1e1e surface)
    - [ ] Configure text colors (white primary, rgba white 0.7 secondary)
    - [ ] Configure breakpoint system (xs, sm, md, lg)

- [ ] Integrate ThemeProvider in root layout (AC: All)
    - [ ] Wrap app content with MUI ThemeProvider in `app/layout.tsx`
    - [ ] Import and apply the theme configuration
    - [ ] Ensure ThemeProvider is client-side compatible

- [ ] Verify build and tests (AC: All)
    - [ ] Run `npm run build` and confirm successful build
    - [ ] Run `npm run test` and confirm all tests pass
    - [ ] Verify no Tailwind references remain in codebase

## Dev Notes

### Architecture Context

This story establishes the foundational design system for the entire web-nextjs application. The architecture decision (Architecture.md, Decision #1) explicitly
requires:

- **Complete Tailwind removal** - no mixing of CSS systems
- **Dark theme as default** - photos as primary visual focus
- **Brand color integration** - #185986 throughout for identity
- **MUI breakpoint system** - standardized responsive approach

### State Management Migration Note

This story does NOT migrate the state management from `web/src/core/catalog/`. That is Story 1.2. This story focuses purely on the visual design system
foundation.

### Technology Stack

**Current NextJS Setup:**

- Next.js 16.1.1 with App Router
- React 19.2.3
- TypeScript 5.x with strict mode
- Currently has Tailwind CSS (to be removed)

**Target Design System:**

- Material UI 6.x (latest major version)
- Emotion for CSS-in-JS
- Dark theme with brand customization

### File Locations

**Create:**

- `web-nextjs/components/theme/theme.ts` - MUI theme configuration

**Modify:**

- `web-nextjs/package.json` - dependency changes
- `web-nextjs/app/layout.tsx` - ThemeProvider integration
- Any existing components with Tailwind classes (if present)

**Remove:**

- `tailwind.config.js` (if exists)
- Tailwind imports from CSS files
- Tailwind classes from components

### Theme Configuration Example

From Architecture.md, the theme should follow this structure:

```typescript
palette: {
    mode: 'dark',
        primary
:
    {
        main: '#185986', // Brand blue
    }
,
    background: {
    default:
        '#121212',
            paper
    :
        '#1e1e1e',
    }
,
    text: {
        primary: '#ffffff',
            secondary
    :
        'rgba(255, 255, 255, 0.7)',
    }
,
}
breakpoints: {
    values: {
        xs: 0,
            sm
    :
        600,
            md
    :
        960,
            lg
    :
        1280,
            xl
    :
        1920,
    }
,
}
```

### Testing Requirements

From nextjs.instructions.md:

- Run `npm run test` for unit tests (~5s expected)
- Run `npm run build` to verify production build
- All tests must pass before completing story

No visual tests (Ladle) are required for this story as it's purely infrastructure setup without new UI components.

### Coding Standards Compliance

**From nextjs.instructions.md:**

1. **No comments in code** - communicate via chat/story notes only
2. **Explicit types** - no `any` types
3. **File structure** - components in own folders with index.tsx export
4. **Testing strategy** - unit tests with vitest, visual tests with Ladle (not applicable for this story)

### Browser Support

From Architecture.md:

- Latest 2 versions of Chrome, Firefox, Safari, Edge
- Modern ES2020+ and CSS features allowed
- No IE11 or legacy browser support needed

### Performance Considerations

- MUI tree-shaking will be configured in next.config.ts (Story 1.2)
- Bundle size monitoring - MUI should be loaded efficiently
- This foundation enables future performance optimizations

### Related Stories

- **Story 1.2 (Next):** State Management Migration - will build on this theme foundation
- **Story 1.3:** Basic Album List Display - will use MUI components with this theme
- **Story 1.4:** Album Card Enhancements - will leverage brand color in components

### References

- **Architecture.md, Decision #1**: Material UI Integration (complete removal of Tailwind, dark theme, brand color)
- **UX Design Specification**: Color System section defines dark theme palette
- **PRD**: Technical Foundation section specifies Material UI requirement
- **nextjs.instructions.md**: File structure, testing strategy, coding standards

### Known Constraints

- **No backend changes** - this is frontend-only work
- **No breaking existing auth** - authentication flow already works, must preserve it
- **NextJS App Router** - theme must work with server/client component split

### Success Validation

**Before marking complete:**

1. `npm install` runs successfully with new dependencies
2. `npm run build` completes without errors
3. `npm run test` shows all tests passing
4. No Tailwind references in codebase (search for "tailwind", "tw-", "@tailwindcss")
5. ThemeProvider wraps app in root layout
6. Theme file exists with correct configuration
7. Application starts with `npm run dev` and displays dark theme

### Anti-Patterns to Avoid

- ❌ Mixing Tailwind and MUI (remove all Tailwind completely)
- ❌ Using generic MUI defaults (must customize with brand color)
- ❌ Inline styles (use MUI `sx` prop pattern)
- ❌ Creating custom breakpoint systems (use MUI's standard breakpoints)

## Dev Agent Record

### Agent Model Used

<!-- Dev agent will fill this in -->

### Debug Log References

<!-- Dev agent will fill in paths to relevant logs -->

### Completion Notes List

<!-- Dev agent will document what was implemented and any decisions made -->

### File List

<!-- Dev agent will list all files created or modified -->
