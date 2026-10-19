---
name: ui-components
description: UI components and visual testing guide for DPhoto's `web-nextjs/` project. Covers Storybook testing strategy, component props principles, composition patterns, and design system implementation. Required skill when creating or updating UI components in `web-nextjs/components/` and `web-nextjs/app/`.
---

**DO NOT START OR STOP STORYBOOK.** You might run `npm run build` to verify there are no compile errors.
**AND DO NOT RUN THE VISUAL TESTS LOCALLY**, they are only relevant on CI.

# UI Components Architecture

## Component structure

```
web-nextjs/
├── app/
│   └── albums/
│       └── pages.tsx                             # 
│       └── _components/                          # components used by a SINGLE page
│           └── CreateAlbumModal/                 # one folder per component, flat (no nesting)
│               ├── index.tsx                     # the component, always named `index.tsx`
│               ├── AlbumTitle.tsx                # subcomponent extracted for readability, same folder
│               └── CreateAlbumModal.stories.tsx  # 
└── components/                                   # components shared across 2+ pages
    └── AppLayout/                                # same flat structure applies
        └── index.tsx
```

## Design principles

**Components are features, not widgets.** Name them after what the user is doing (`CreateAlbumModal`), not the UI primitive they wrap (`FormDialog`). A feature component owns its full tree — modal, form, buttons — composed directly inside it.

**Props use domain types.** Accept `Album`, `Media`, `AuthenticatedUser` directly. Exception: small leaf components rendering <3 fields of one concept take flat primitives (`UserAvatar({name, picture, size})`).

**Shallow tree.** No middle components that only wrap 2-3 children. No prop cascading. Each component is as deep as possible with only these few exceptions:

* server-side component, fetching data directly, should be separated from their `*Content` component rendering it.
* state management is separated from the view: when a feature component uses non-trivial hooks (`useReducer`, `useThunks`, router, session), split it into two components in the same file: the stateful one keeps the original name and owns hooks/selectors; the pure `XxxView` sibling takes flat props and is the one tested in Storybook.
* too large components: when the component reach a certain size (~100 lines), consider extracting subcomponents (`AlbumTitle`, `StartOrEndDateInput`, ...). The big picture and layout remains completely in the components, only its details are abstracted into meaningfully named subcomponents.

Example of View splitting:

```tsx
export function MyComponent({initialState}: MyComponentProps) {
    const [state, dispatch] = useReducer(reducer, initialState);
    useSessionAwareThunks(thunks, {dispatch}, state);
    return <MyComponentView {...mySelector(state)}/>;
}

export type MyComponentViewProps = MySelection;

export function MyComponentView({field1, field2, onAction}: MyComponentViewProps) {
    // pure JSX
}
```

## Theme

Customise typography in `web-nextjs/components/theme/theme.ts` rather than wrapping `Typography`. Variants already defined: `h1` (page title), `body1` (page description), `h2` (in-page separator).

## Storybook stories

### When to write stories

- **One story file per root component.** Subcomponents (siblings in the same folder) are covered by the root's stories — never give them their own.
- **One story per visual state.** Name the neutral/happy-path one `Default`. Others describe the state: `MenuOpen`, `Loading`, `Error`, `Disabled`.
- **Never set a viewport.** Visual tests run all viewports automatically.
- **When state is split from view, stories target the `View` export.** Evaluate selectors in the stories file to turn state fixtures into View props.

### File shape

```tsx
import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {fn} from 'storybook/test';                        // NOT '@storybook/test'
import {AlbumCard} from './index';
import {AppBackground} from '@/components/AppLayout/AppBackground';

const clairObscurAlbum: Album = { /* one fixture, reused across all stories */ };

const meta = {
    title: 'catalog/AlbumCard',                            // {domain}/{ComponentName}
    component: AlbumCard,
    parameters: {layout: 'fullscreen'},
    decorators: [(Story) => <AppBackground><Story/></AppBackground>],
    args: {                                                // happy-path defaults, used by Default
        album: clairObscurAlbum,
        onSelect: fn(),                                    // spy on every callback
    },
} satisfies Meta<typeof AlbumCard>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};                          // inherits meta.args

export const NoThumbnails: Story = {                       // override ONLY what the story demonstrates
    args: {album: {...clairObscurAlbum, thumbnails: []}},
};
```

### Rules

- **Never use `render`.** Compose via `args` and `decorators` only.
- **Import `fn` from `storybook/test`**, not `@storybook/test`.
- **One fixture, override per story.** Define the dataset once at the top of the file; each story overrides only the fields it demonstrates.
- **Wrap in `AppBackground`** via a decorator, unless the component replaces its own background (modals, full-screen overlays).
- **Controlled props** (`isOpen`/`onClose`, etc.) get a `useState` decorator:

  ```tsx
  export const MenuOpen: Story = {
      decorators: [(Story) => {
          const [isOpen, setIsOpen] = useState(true);
          return <Story args={{isOpen, onClose: () => setIsOpen(false)}}/>;
      }],
  };
  ```

### Test assets

Use `web-nextjs/public/tonystark-profile.jpg` for avatars/images. Copy from `test/wiremock/__files/api/static/` if more are needed.
