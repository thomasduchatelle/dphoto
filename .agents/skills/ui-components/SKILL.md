---
name: ui-components
description: UI components and visual testing guide for DPhoto's `web-nextjs/` project. Covers Storybook testing strategy, component props principles, composition patterns, and design system implementation. Required skill when creating or updating UI components in `web-nextjs/components/` and `web-nextjs/app/`.
---

**DO NOT START OR STOP STORYBOOK.** You might run `npm run build` to verify there are no compile errors. And do not run the visual tests locally, they are only
relevant on CI.

# UI Components Architecture

**Before implementing, you need to define the UI components architecture to respect the following rules.** Your architecture needs to be composed of:

* **components created / updated / deleted / renamed**: list the UI components that will be modified as part of you tasks.
* **signature**: describe the signature of each component (input properties).
* **stories**: for each component, list the stories names that you will write or update, with a brief description of how the component should look like.
  Describe in plain English, no code.
* **integration / dependency**: brief explanation on where the component will be used, and what other component it will depend on (if any). Explain in plain
  English, no code.

If asked to plan your work, you shall present your architecture to the user.

## Directory Structure

Components can be found in two places:

* `web-nextjs/app/{path to the page}/_components/{ComponentName}/index.tsx`: components only used by `web-nextjs/app/{path to the page}/page.tsx` (example:
  `web-nextjs/app/(authenticated)/_components`).
* `web-nextjs/components/{ComponentName}/index.tsx`: components shared and used in different places (layout, ...).

Both `app/.../_components` and `web-nextjs/components/` adopt a **flat structure**: all the components are directly children of them.

A component is always a folder with an index file:

* `{ComponentName}/`
    * `index.tsx`: required, the main component code.
    * `{ComponentName}.stories.tsx`: required with few exceptions, the stories validating the component is rendering appropriately.
    * `{SubComponent}.tsx`: optional, a component extracted from the main `index.tsx` for readability purpose.

## Component signatures

You must read the model of the domain you're working on before defining the signature of the components. They can be found, for example:

* **auth**: `import {AuthenticatedUser, AuthenticatedSession} from '@/libs/security/session-service`'.
* **catalog**: in `web-nextjs/domains/catalog/language`, example: `import {Album, AlbumId, Media, UserDetails) from '@/domains/catalog/language'`.

### Use the domain types (Preferred)

Components should use the domain types. Good example:

```tsx
import {AuthenticatedUser} from '@/libs/security/session-service';

export interface AppHeaderProps {
    user: AuthenticatedUser;
    onLogout?: () => void;
}
```

Bad example (inline types):

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

### Exception: simple and specialised rendering component

If a component is specialised to render a single property, or a small set of properties (less than 3), from a domain object, then use flat properties with
native types:

```tsx
export interface UserAvatarProps { // the name is explicit, only the avatar will be shown, with a fallback to the name. No need to pass the complete user.
    name: string;
    picture?: string;
    size?: 'small' | 'medium' | 'large';
}

export interface CatalogNameProps {
    name: string;
    temperature?: number;
}
```

## Prioritise simple and single-responsibility components

This is a **good example** to reproduce:

```tsx
export function HomeContent() {
    return <>
        <PageHeader actions={(
            <>
                <CreateAlbumButton onCreateAlbumButton={handleCreateButton} disabled={!canCreateAlbum}/>
                <DeleteAlbumButton onCreateAlbumButton={handleCreateButton} disabled={!canDeleteCurrentAlbum}/>
            </>
        )}>
            <AlbumFilter filter={currentFilter} onFilterChange={handleFilterChange}/>
        </PageHeader>
        <AlbumGrid albums={albums}/>
        <PageFooter>
            <BackToTopButton/>
        </PageFooter>
    </>
}  
```

It is a good because it has a **shallow dependency tree of components** which is characterised by:

* No middle-components that only aggregate few components together.
* No properties cascading between the components.
* Only few properties by components.

### A BAD example TO AVOID

This is a bad example, to demonstrate the contrast, with the following dependency structure:

* `AlbumsFilterAndGrid` - the component only calls two other components but needs to pass MANY properties down.
    * `AlbumsHeader` - same, useless
        * `AlbumButtons` - shall not be introduced unless we need the EXACT same list of buttons in a different page
            * `CreateAlbumButton`
            * `DeleteAlbumButton`
        * `AlbumFilter`
    * `AlbumGrid`

## Reusable components

Most components that are in DPhoto are specific for a business use case and are only used in a SINGLE place. There are few exceptions, especially to encapsulate
the layout and themes.

* Avoid conditional logic, and prefer composition (with `children` property), on these components.
* Auto-style the children using CSS rules.
* Control the overall layout by placing and boxing the different elements, including the children.

Good example:

```tsx
export const PageMessage = ({variant = 'info', icon, title, message, children}: PageMessageProps) => {
    const colors = variantColors[variant];

    return (
        <Box>
            <Typography
                variant="h4"
                component="h2"
            >
                {title}
            </Typography>
            <Typography
                variant="body1"
            >
                {message}
            </Typography>
            {children && (
                <Box sx={{
                    display: 'flex',
                    gap: 2,
                    flexWrap: 'wrap',
                    justifyContent: 'center',
                    '& .MuiButton-contained': {
                        bgcolor: '#185986',
                        color: '#ffffff',
                        px: 4,
                        py: 1.5,
                        textTransform: 'uppercase',
                        letterSpacing: '0.1em',
                        fontSize: '14px',
                        fontWeight: 400,
                        '&:hover': {
                            bgcolor: '#206ba8',
                            boxShadow: '0 0 24px rgba(24, 89, 134, 0.6)',
                        },
                    },
                    '& .MuiButton-outlined': {
                        borderColor: 'rgba(74, 158, 206, 0.4)',
                        color: 'rgba(255, 255, 255, 0.9)',
                        px: 4,
                        py: 1.5,
                        textTransform: 'uppercase',
                        letterSpacing: '0.1em',
                        fontSize: '14px',
                        fontWeight: 400,
                        '&:hover': {
                            borderColor: '#4a9ece',
                            bgcolor: 'rgba(74, 158, 206, 0.1)',
                        },
                    },
                }}>
                    {children}
                </Box>
            )}
        </Box>
    )
}
```

### Typography components

**MUI theme variants**: Customise typography in `web-nextjs/components/theme/theme.ts` instead of custom Typography component. Typography already defined:

- h1: main title of the page
- body1: can be used as description of the page, under the h1
- h2: sub separator within the page content

Usage example:

```tsx
<Typography variant="h1">The title</Typography>
```

---

## Writing Storybook Stories - CRITICAL to get PR accepted

### Visual Strategy: when to create Stories ?

The strategy is what you need to define and present when your planning a task, and before implementing it.

* **Create stories ONLY for the main/root component** that includes all its subcomponents. Never create separate story files for subcomponents.

  Examples:

    ```
    * AppLayout/
      * index.tsx
      * AppBackground.tsx
      * AppLayout.stories.tsx     -> this is stories to showcase the AppLayout AND AppBackground.tsx. There will NOT be a AppBackground.stories.tsx.
    ```

* Create a story for each state the component can take. The most complete and neutral state is called `Default` (not `Desktop`, not `Primary`).

  Examples: `MenuOpen`, `Loading`, `Error`, `Disabled`.

* Stories must NOT define a viewport, even when specific behaviour for mobile is required. The visual tests will automatically test with different viewports.

### Coding conversion: How to write a story both readable and maintainable ?

Once the strategy for the component has been defined, the implementation must strictly follow these coding conventions. A complete example follow the
templates (see below for descriptions on each section):

```tsx
import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {fn} from 'storybook/test';
import {MyComponent} from './index';
import {AppBackground} from '@/components/AppLayout/AppBackground';

const meta = {
    title: 'Layout/MyComponent',
    component: MyComponent,
    parameters: {
        layout: 'fullscreen',
    },
    decorators: [
        (Story) => (
            <AppBackground>
                <Story/>
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
        data: {...sampleData, specificProperty: 'value'}
    }
};
```

**DO NOT USE `render` for any story or meta! It is strictly PROHIBITED!**

#### Imports

Import `fn` from:

```tsx
import {fn} from 'storybook/test';
```

(not `'@storybook/test'` which contains an extra `@`).

#### Meta

* **title** convention is `{domain}/{ComponentName`, example: `catalog/DeleteAlbumModal`.
* **decorators**:
    * use the `AppBackground` unless the component will be used in a component likely to redefine its background (modals, ...).
    * When a component has **controlled properties** (like `isOpen`/`onClose` pairs), create fake implementations using `useState` within the story decorator:

  ```tsx
  export const MenuOpen: Story = {
      decorators: [
          (Story) => {
              const [isOpen, setIsOpen] = useState(true);
              return <Story args={{isOpen, onClose: () => setIsOpen(false)}}/>;
          },
      ],
  };
  ```

* **args**:
    * **actions / callbacks**: use the Storybook spy `fn()` to capture any interactions.
    * **props / data structure**: fulfil all required argument in their most neutral -- happy path -- possible. They will be used as it by the `Default` story.

#### Default story

The default story is the component with only the required parameters, defined in the `meta.args`.

#### Other stories

A dataset must be created to be used across every story so **only the values demonstrated are overridden for the story**.

Examples:

```tsx
const clairObscurAlbum: Album = {
    albumId: createAlbumId('sandfall', 'clair-obscur'),
    name: 'Clair Obscur',
    start: new Date('2025-04-24'),
    end: new Date('2025-06-01'),
    totalCount: 47,
    temperature: 6.7,
    relativeTemperature: 1,
    sharedWith: [],
    thumbnails: [
        '/thumbnails/clair-obscur-1.jpg',
        '/thumbnails/clair-obscur-2.jpg',
        '/thumbnails/clair-obscur-3.jpg',
        '/thumbnails/clair-obscur-4.jpg',
    ],
};

export const Default: Story = {
    args: {}
};

export const NoThumbnails: Story = {
    args: {
        album: {...clairObscurAlbum, thumbnails: []}
    }
};

export const Cold: Story = {
    args: {
        album: {...clairObscurAlbum, temperature: 1.42}
    }
};
```

### Testing Content

Use images and avatars that are from `test/wiremock/__files/api/static/tonystark-profile.jpg`. Copy them if necessary. `web-nextjs/public/tonystark-profile.jpg`
is already available.
