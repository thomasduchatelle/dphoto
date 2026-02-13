# Story 1.4 - UI Components Review TODO

## 1. Instructions

### Context

We're implementing `specs/stories/1-4-ui-components-layout.md`. Read it.

All 10 components from Story 1.4 have been created with their Storybook visual tests. The implementation followed the initial story requirements, but we need to
validate them against the **final design direction** specified in:

- Lead dev guidance: you will work iteratively, following the instruction of the user.
- `specs/designs/ux-design-direction-final.html` (PRIMARY reference - edge-to-edge photos, dark blue gradient, text overlays)

### Goals

1. **Keep the good ideas** from the initial implementation
2. **Apply the final design** direction consistently across all components (colors, typography, spacing, layout patterns)
3. **Validate visual quality** through Storybook stories
4. **Document learnings** as we review each component to apply consistently across remaining components

### Workflow for each task

For each task, you must follow this steps:

1. pick the next component (or group of component on the list)
2. Take the initiative ! Read and applies the learnings from the previous components !
3. Then present the component we're working on, what you already updated, and the properties structure required to use it.
4. Wait - apply - wait - apply: the lead dev will ask you to make some modifications, apply them.
    * you are the UX/UI expert, don't simply listen: make propositions, think about alternatives, implement different approach so they can be compared side by
      side !
5. Once the options refined and chosen, cleanup the code to remove anything not required.
6. Propose a list of learning to apply to the next components.
    * they are not always "learnings" to remember, only suggest what is agnostic of this specific component, and is actionable (no general conceptual idea),
7. Update this document: the learning, and check the task list.

---

## 2. Tasks

### Components

* [x] **AppLayout**: `components/layout/AppLayout/`
  _Main application layout wrapper with fixed header and content area. Verify responsive padding, semantic HTML structure, and integration with AppHeader._
  Includes:
    * **AppHeader**: `components/layout/AppHeader/` (Not testing independently)
      _Application header with logo, user profile, and timeline (medias page only, desktop only). The bar is transparent when on the top, blurred when the page
      content pass under, and increase when the mouse goes over it. On Mobile, it is scrolling up with the rest of the page, a small icon to scroll back up
      shows, and on medias page a link back to all albums and "next" are showing._
    * **UserAvatar**: `components/shared/UserAvatar/` (Not testing independently)
      _User avatar with picture or initials fallback in three sizes. Check initials logic, border styling, image loading with Next.js Image component._
    * **Timeline**
      _Lines with dots, each dot represent an album, the label is only the Month+Year. Hovering on the dot show the album card._

* [ ] **AlbumGrid**: `app/(authenticated)/_components/AlbumGrid/`
  _Responsive grid layout for album cards. Verify column configuration (xs:1, sm:2, md:3, lg:4), 32px gap, max-width centering, and semantic HTML._
  Includes:
    * **AlbumCard**: `app/(authenticated)/_components/AlbumCard/`
      _Album display card with name, date range, media count, and density indicator. **CRITICAL:** Validate against final design (edge-to-edge photos, text
      overlay, hover effects). Check typography (22px serif, 13px monospace), density color-coding, owner/sharing status display._
    * **SharedByIndicator**: `components/shared/SharedByIndicator/` (Not testing independently)
      _Displays group of user avatars for sharing status. Verify AvatarGroup overlap, "+N" overflow handling, and tooltip functionality._

* [ ] **EmptyState**: `components/shared/EmptyState/`
  _Should be renamed to be specific to albums: NoAlbum. Also extract the NoMedia. As a note for later, we could add other album cards (next/previous)_

* [ ] **ErrorDisplay**: `components/shared/ErrorDisplay/`
  _Error message display with technical details and recovery actions. Validate collapsible details, ARIA attributes, and action button styling._

* [ ] **PageLoadingIndicator**: `components/shared/PageLoadingIndicator/`
  _Discrete full-page loading with thin LinearProgress bar at top. Verify 3px height, brand blue color (#185986), and optional message display._

* [ ] **NavigationLoadingIndicator**: `components/shared/NavigationLoadingIndicator/`
  _Small top loading indicator when NextJS links are clicked. Should use https://github.com/TheSGJ/nextjs-toploader ._

### Integration

* [ ] **Error boundaries**: `app/error.tsx` and `app/(authenticated)/error.tsx`
  _Verify ErrorDisplay component integration with proper error prop passing and onRetry handlers._

* [ ] **Not-found page**: `app/not-found.tsx`
  _Validate EmptyState component integration with appropriate icon, message, and "Go Home" action._

* [ ] **Authenticated layout**: `app/(authenticated)/layout.tsx`
  _Confirm AppLayout wrapper integration with user data passing._

---

## 3. Learnings and Rules

**DO NOT START OR STOP STORYBOOK.**

You might run `npm run build` to verify there are no compile error.

### Story Principles (CRITICAL - Read Before Creating Stories)

#### 1. Component Demonstration

- **One story file per root component**: Create stories only for root/parent components that demonstrate the assembly of subcomponents
- **No duplication**: Each component should be demonstrated ONCE. Subcomponents are tested through their parent component stories
- **Demonstrate all states**: Show all meaningful states (menu open, button disabled, loading, error, etc.)
- **Mobile variants acceptable**: If a component renders very differently on mobile, create a mobile-specific story for that state

#### 2. Story Naming Convention

- **Default story**: Must be named `Default` (not `Desktop`, not `Primary`)
- **Mobile story**: Must be named `DefaultMobile` (if mobile rendering differs significantly)
- **State variants**: Name by state, e.g., `MenuOpen`, `Loading`, `Error`, `Disabled`

#### 3. Viewport Configuration

- **Default story**: DO NOT set viewport config - uses Storybook's default viewport. On default export (meta), use the global:
  ```tsx
  globals: { viewport: {} }
  ```
- **Mobile stories**: Use `globals.viewport` only for mobile variants:
  ```tsx
  globals: { viewport: { value: 'mobile2', isRotated: false } }
  ```

#### 4. Imports - CRITICAL

**CORRECT**:

```tsx
import {fn} from 'storybook/test';
```

**WRONG** (do not use):

```tsx
import {fn} from '@storybook/test'; // ❌ Extra @
```

### Component Props Principles (CRITICAL - Read Before Creating Components)

#### 1. Use Domain Types (Preferred)

Components should use types from the domain layer:

- `@/libs/security/session-service`: `AuthenticatedUser`, `AuthenticatedSession`, etc.
- `@/domains/catalog/language`: `Album`, `AlbumId`, `Media`, `UserDetails`, etc.

**Good example**:

```tsx
import {AuthenticatedUser} from '@/libs/security/session-service';

export interface AppHeaderProps {
    user: AuthenticatedUser;
    onLogout?: () => void;
}
```

**Bad example** (inline types):

```tsx
export interface AppHeaderProps {
    user: {
        name: string;
        email: string;
        picture?: string;
    };
    onLogout?: () => void;
}
```

#### 2. Edge Case: Specific Props Only

If a component explicitly shows only 1-2 properties from a domain object, extract them as separate props:

**Good example** (UserAvatar only needs name and picture for initials):

```tsx
export interface UserAvatarProps {
    name: string;
    picture?: string;
    size?: 'small' | 'medium' | 'large';
}
```

**Counter-example**: `AlbumCard` intention is to represent a whole album, and must take the `Album` model even if only few properties are used.

```tsx
import {Album} from '@/domains/catalog/language';

export interface AlbumCardProps {
    album: Album; // ✅ Takes whole Album - component decides what to display
}
```

#### 3. Native Types (Only for Simple Components)

Use native types (string, number, boolean) ONLY if:

- Very few props (≤ 3)
- No domain representation exists
- Component is truly generic/reusable

#### 4. Red Flag

If you need to pass many properties (>3) and no domain type exists, **this is a red flag**:

- Consider if the component abstraction is correct
- Consider creating a domain type
- Discuss with team before proceeding

### Component Structure Principles

#### 1. Prevent Props Cascading

**Bad** (props drilling through layers):

```tsx
<AppLayout user={user} onLogout={onLogout}>
    <Content/>
</AppLayout>
```

**Good** (composable views):

```tsx
<AppLayout>
    <AppHeader user={user} onLogout={onLogout}/>
    <Content/>
</AppLayout>
```

#### 2. Prefer Composition Over Configuration

- Components should accept children and compose
- Avoid passing props through multiple layers
- Each component receives only what it needs directly

### Directory Structure:

**Flat component structure**: All components in `components/` directory without intermediate folders

- `components/AppLayout/` - Main layout component with stories
- `components/AppLayout/AppHeader.tsx` - Header subcomponent (cannot be used without AppLayout), no stories
- `components/UserAvatar/` - Avatar subcomponent, no stories

### Typography Approach:

**MUI theme variants**: Customise typography in `web-nextjs/components/theme/theme.ts` instead of custom Typography component. Typography already defined:

- h1: main title of the page
- body1: can be used as description of the page, under the h1
- h2: sub separator within the page content

Usage example:

```tsx
<Typography variant="h1">The title</Typography>
```

### Testing Content:

Use images and avatar that are from `test/wiremock/__files/api/static/tonystark-profile.jpg`. Copy them if necessary. `web-nextjs/public/tonystark-profile.jpg`
is already available.


---

**Last Updated:** 2026-02-10 (Revision 2 - Compliance fixes in progress)  
**Status:** AppLayout group being updated for compliance  
**Next Action:** After fixes verified, move to AlbumGrid component group
