---
name: ui-components
description: UI components and visual testing guide for DPhoto's `web-nextjs/` project. Covers Storybook testing strategy, component props principles, composition patterns, and design system implementation. Required skill when creating or updating UI components in `web-nextjs/components/` and `web-nextjs/app/`.
---

# UI Components and Visual Testing Guide

## Component Props Principles (CRITICAL - Read Before Creating Components)

### 1. Use Domain Types (Preferred)

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

### 2. Edge Case: Specific Props Only

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

### 3. Native Types (Only for Simple Components)

Use native types (string, number, boolean) ONLY if:

- Very few props (≤ 3)
- No domain representation exists
- Component is truly generic/reusable

### 4. Red Flag

If you need to pass many properties (>3) and no domain type exists, **this is a red flag**:

- Consider if the component abstraction is correct
- Consider creating a domain type
- Discuss with team before proceeding

---

## Component Structure Principles

### 1. Prevent Props Cascading

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

### 2. Prefer Composition to Configuration

- Components should accept children and compose
- Avoid passing props through multiple layers
- Each component receives only what it needs directly

---

## Directory Structure

**Flat component structure**: All components in `components/` directory without intermediate folders

- `components/AppLayout/` - Main layout component with stories
- `components/AppLayout/AppHeader.tsx` - Header subcomponent (cannot be used outside AppLayout), no stories
- `components/UserAvatar/` - Avatar subcomponent, no stories

---

## Typography Approach

**MUI theme variants**: Customise typography in `web-nextjs/components/theme/theme.ts` instead of custom Typography component. Typography already defined:

- h1: main title of the page
- body1: can be used as description of the page, under the h1
- h2: sub separator within the page content

Usage example:

```tsx
<Typography variant="h1">The title</Typography>
```

---

## Testing Content

Use images and avatars that are from `test/wiremock/__files/api/static/tonystark-profile.jpg`. Copy them if necessary. `web-nextjs/public/tonystark-profile.jpg` is already available.

---

## Writing Storybook Stories - Implementation Guide (CRITICAL - Read Before Creating Stories)

**DO NOT START OR STOP STORYBOOK.** You might run `npm run build` to verify there are no compile errors.

### 1. Single Story File Per Component Hierarchy

**Create stories ONLY for the main/root component** that includes all its subcomponents. Never create separate story files for subcomponents.

**Example**: `AlbumCard.stories.tsx` demonstrates both `AlbumCard` and all its internal subcomponents (thumbnails, sharing indicators, etc.) together. No separate stories for those subcomponents.

**Key rules**:
- **One story file per root component**: Create stories only for root/parent components that demonstrate the assembly of subcomponents
- **No duplication**: Each component should be demonstrated ONCE. Subcomponents are tested through their parent component stories
- **Demonstrate all states**: Show all meaningful states (menu open, button disabled, loading, error, etc.)
- **Mobile variants acceptable**: If a component renders very differently on mobile, create a mobile-specific story for that state

### 2. Story Naming Convention

- **Default story**: Must be named `Default` (not `Desktop`, not `Primary`)
- **State variants**: Name by state, e.g., `MenuOpen`, `Loading`, `Error`, `Disabled`

Example:

```tsx
export const Default: Story = {};  // Basic case with meta.args

export const Shared: Story = {     // Shows sharing feature
    args: { album: { ...clairObscurAlbum, sharedWith: [...] } }
};

export const WithoutThumbnail: Story = {  // Shows edge case
    args: { album: { ...clairObscurAlbum, thumbnails: [] } }
};
```

### 3. Viewport Configuration

Stories must NOT define a viewport, even when specific behavior for mobile is required. The visual tests will automatically test with different viewports.

### 4. Imports - CRITICAL

**CORRECT**:

```tsx
import {fn} from 'storybook/test';
```

**WRONG** (do not use):

```tsx
import {fn} from '@storybook/test'; // ❌ Extra @
```

### 5. Use `meta.decorators` for Component Wrapping

Wrap components in `meta.decorators` to provide necessary context, styling, or layout for realistic demonstrations:

```tsx
const meta = {
    title: 'Components/AlbumCard',
    component: AlbumCard,
    parameters: {
        layout: 'fullscreen',
    },
    decorators: [
        (Story) => (
            <AppBackground>
                <Box sx={{ maxWidth: "500px", p: { md: 3, xs: 1 } }}>
                    <Story/>
                </Box>
            </AppBackground>
        ),
    ],
    // ...
} satisfies Meta<typeof AlbumCard>;
```

**Purpose**: Decorators ensure components render with appropriate backgrounds, padding, containers, or theme providers without polluting the component itself.

### 6. Use `meta.args` for Sensible Defaults

Define default values in `meta.args` that apply to ALL stories unless overridden. These should be realistic, minimal defaults:

```tsx
const meta = {
    // ...
    args: {
        albums: sampleAlbums,  // Realistic sample data
        onShare: fn(),          // Spy function for callbacks
    },
} satisfies Meta<typeof AlbumGrid>;
```

**Each story should override ONLY the specific args relevant to that story's demonstration**:

```tsx
export const Shared: Story = {
    args: {
        album: {
            ...clairObscurAlbum,
            sharedWith: [{user: {name: 'Hulk', email: 'hulk@avenger.com', picture: '/static/hulk-profile.webp'}}],
        }
    }
};
```

### 7. Use `fn()` to Spy on Actions

Use `fn()` from `storybook/test` for callback props that have no effect on the current component's rendering:

```tsx
import {fn} from 'storybook/test';

const meta = {
    args: {
        onShare: fn(),   // Spy on the callback
        onLogout: fn(),  // Track interactions without implementation
    },
} satisfies Meta<typeof Component>;
```

**Purpose**: `fn()` allows Storybook to track and display when callbacks are invoked in the Actions panel, useful for demonstrating interactivity without affecting visual state.

### 8. Use Fake Implementations for Controlled Properties

When a component has **controlled properties** (like `isOpen`/`onClose` pairs), create fake implementations using `useState` within the story decorator:

```tsx
export const MenuOpen: Story = {
    decorators: [
        (Story) => {
            const [isOpen, setIsOpen] = useState(true);
            return <Story args={{ isOpen, onClose: () => setIsOpen(false) }} />;
        },
    ],
};
```

**Use case**: Components that require state management (modals, dialogs, dropdowns) where the open/close state needs to be interactive in Storybook.

### 9. Complete Story File Template

```tsx
import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {fn} from 'storybook/test';
import {MyComponent} from './index';
import {Box} from '@mui/material';
import {AppBackground} from '@/components/AppLayout/AppBackground';

const meta = {
    title: 'Components/MyComponent',
    component: MyComponent,
    parameters: {
        layout: 'fullscreen',
    },
    decorators: [
        (Story) => (
            <AppBackground>
                <Box sx={{p: {xs: 2, md: 6}}}>
                    <Story/>
                </Box>
            </AppBackground>
        ),
    ],
    args: {
        // Sensible defaults for all stories
        data: sampleData,
        onAction: fn(),
    },
} satisfies Meta<typeof MyComponent>;

export default meta;
type Story = StoryObj<typeof meta>;

// Default story inherits meta.args
export const Default: Story = {};

// Override specific args to demonstrate variants
export const VariantA: Story = {
    args: {
        data: { ...sampleData, specificProperty: 'value' }
    }
};
```
