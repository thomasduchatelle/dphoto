# Media Tagging

## Intent

Let an owner attach **tags** to their medias so they can find and act on medias across albums — starting with marking favourites, growing into a full tag-based browsing experience, and ultimately downloading the originals of every media carrying a given tag.

Tagging must be **extremely quick** from the album page, where the owner already spends time looking at the media grid.

## Vocabulary

- **Tag**: a normalised, lowercase, URL-safe label attached to a media (e.g. `people:luke`, `place:school`, `activity:ski`, `favorite`). Owner-scoped: tags belong to the owner, alongside the media.
- **Namespace**: the part of a tag before the first colon (`people`, `place`, `activity`). Every tag has a namespace except the reserved `favorite`. Namespaces are free-form — the owner invents them on the fly.
- **Value**: the part after the first colon (`luke`, `school`, `ski`).
- **Favorite**: a reserved tag with no namespace. Treated as a normal tag by the backend, but surfaced in the UI as a dedicated star gesture.
- **Tagging mode**: a UI state where a single tag is "active" and clicking a media toggles that tag on it ("paint" model).

## Who can do what

- Only the **owner** can add or remove tags on their medias (security boundary).
- Visitor visibility of tags on shared albums is **out of scope for this feature** — acceptable if visitors see no tags for now.

## User journeys

### 1. Mark a favourite (fastest path)
On the album page, each media shows a star on hover. One click toggles `favorite` on that media.

### 2. Tag medias in "tagging mode"
On the album page, the owner presses `t` (or clicks a dedicated button) and picks a tag via a type-ahead picker that autocompletes existing tags (grouped by namespace), offers recently-used tags as quick-picks, and lets them create a new `namespace:value` tag on the fly. Once a tag is picked, the page enters **tagging mode**: a header banner shows the active tag, and clicking any media toggles that tag on it. `ESC` or a button exits the mode.

Recently-used tags are a client-side convenience (last 5–10, reset on page refresh).

### 3. Browse photos of a tag
A top-level header navigation (`Albums | Tags`, with `Map | Settings` as future siblings) gives access to a **Tags page**. It opens on `favorite` by default and offers a switcher to pick any other of the owner's tags (grouped by namespace). It displays the grid of all medias carrying the selected tag, across all albums.

### 4. Download originals of a tag
From the CLI: `dphoto download --tag <tag> --to <dir>` downloads the originals (from S3) of every media carrying the tag, across all albums. This is the end goal the data model must support; it is specified here only to anchor the design, and will be delivered in a later feature.

## Scope of this feature

Delivered in this order, driven by the interaction model being the risky, subjective part:

1. **Interaction spike** — explore radically different tagging interactions and pick one (decision only, no production code).
2. **Favourites, end-to-end** — tag a media and see it. One story, built in parallel once the REST contract is pinned in the design: backend (domain + persistence + GSI + endpoints), the favourite star on the album grid, and the Tags page + header navigation showing favourites.
3. **Tagging mode** — the chosen interaction wired in for general namespaced tags.
4. **Tag switcher** on the Tags page.
5. **List the owner's distinct tags** — the enumeration the switcher relies on.
6. **CLI download by tag** — the north-star goal.

## Out of scope

- Tag counts on the Tags page.
- Visitor visibility of tags on shared albums (owner-only in this feature).
- Filtering by several tags at once, and crossing tags with location — a separate future feature.
