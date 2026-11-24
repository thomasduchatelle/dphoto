// PROTOTYPE — throwaway. Spike for issue 01 (album-page tagging interaction).
// Story group: "Prototypes / Tagging". Pure UI, components embedded, minimal state.
// Goal: settle the FINAL LOOK of the album page in "tagging mode".
//
// Requirements captured from the spike review:
//  1. No "Enter tagging mode: …" helper text — a tooltip on the Tag button instead.
//  2. The keyboard shortcut is just `t` (plain, no Ctrl/Cmd — would clash with the browser).
//  3. Type-along uses MUI Autocomplete: ENTER selects the tag as typed; TAB / ARROW
//     selects an existing match (native Autocomplete behaviour, nothing custom).
//  4. Clicking an image toggles the active tag (adds if absent, REMOVES if present).
//  5. A tag must contain a ':' (validation) — except the reserved `favorite`.
//  6. `favorite` is always first in any tag list, shown with a star icon.
//  7. On mobile, tags are NOT shown on images until the per-image TAG icon is tapped.
//  8. Tapping that per-image TAG icon shows: selected tags | — | recommended | — | "+ More…"
//     ("+ More…" opens the type-along modal where a new tag can be created).
//
// Banner (desktop only) is the "paint" affordance; the per-image TAG icon is the
// mobile affordance.

import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {
    Autocomplete,
    Box,
    Button,
    Chip,
    Dialog,
    DialogContent,
    DialogTitle,
    Divider,
    IconButton,
    Paper,
    Popover,
    Stack,
    TextField,
    Tooltip,
    Typography,
} from '@mui/material';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import LocalOfferIcon from '@mui/icons-material/LocalOffer';
import StarIcon from '@mui/icons-material/Star';
import CloseIcon from '@mui/icons-material/Close';
import {useEffect, useMemo, useState} from 'react';
import {AppBackground} from '@/components/AppLayout/AppBackground';

// --- fake data -------------------------------------------------------------

const THUMBS = [
    'clair-obscur-1', 'clair-obscur-2', 'clair-obscur-3', 'clair-obscur-4', 'clair-obscur-5',
    'clair-obscur-6', 'clair-obscur-7', 'clair-obscur-8', 'astro-bot-01', 'astro-bot-02',
    'astro-bot-03', 'death-stranding-1-01', 'death-stranding-1-02', 'death-stranding-1-03',
    'the-witcher-3-01', 'the-witcher-3-02', 'the-witcher-3-03', 'death-stranding-2-01',
];

const KNOWN_TAGS = ['favorite', 'people:luke', 'people:leia', 'place:beach', 'place:home', 'activity:ski', 'activity:hike'];
const PREVIOUS_TAGS = ['people:luke', 'place:beach', 'activity:ski'];
const FAVORITE = 'favorite';

interface FakeMedia {
    id: string;
    url: string;
    tags: string[];
}

function makeMedias(): FakeMedia[] {
    const seeded: Record<number, string[]> = {
        0: ['favorite', 'people:luke'],
        2: ['place:beach'],
        5: ['favorite'],
        8: ['people:leia', 'activity:ski'],
        12: ['place:home'],
    };
    return THUMBS.map((t, i) => ({id: t, url: `/thumbnails/${t}.jpg`, tags: seeded[i] ?? []}));
}

// --- tag helpers -----------------------------------------------------------

function normalise(raw: string): string {
    return raw.trim().toLowerCase().replace(/\s+/g, '');
}

// Req 5: a tag must contain a ':' — except the reserved `favorite`.
function tagError(tag: string): string | null {
    const t = normalise(tag);
    if (!t) return 'Enter a tag';
    if (t === FAVORITE) return null;
    if (!t.includes(':')) return 'A tag must be in the form namespace:value (e.g. people:luke)';
    return null;
}

// Req 6: favorite always first, then the rest in given order (de-duplicated).
function favoriteFirst(tags: string[]): string[] {
    const rest = tags.filter((t) => t !== FAVORITE);
    const out = tags.includes(FAVORITE) ? [FAVORITE, ...rest] : rest;
    return Array.from(new Set(out));
}

// --- a chip that renders `favorite` with a star --------------------------

interface TagChipProps {
    tag: string;
    active?: boolean;
    onClick?: (e: React.MouseEvent) => void;
    onRemove?: () => void;
    size?: 'tiny' | 'small';
}

function TagChip({tag, active, onClick, onRemove, size = 'small'}: TagChipProps) {
    const isFav = tag === FAVORITE;
    return (
        <Chip
            icon={isFav ? <StarIcon sx={{fontSize: size === 'tiny' ? '0.8rem' : '1rem', color: 'inherit !important'}}/> : undefined}
            label={tag}
            size="small"
            onClick={onClick}
            onDelete={onRemove}
            deleteIcon={<CloseIcon/>}
            sx={{
                height: size === 'tiny' ? 20 : 26,
                fontSize: size === 'tiny' ? '0.6rem' : '0.75rem',
                cursor: onClick ? 'pointer' : 'default',
                color: active ? 'primary.contrastText' : 'common.white',
                bgcolor: active ? 'primary.main' : 'rgba(255,255,255,0.22)',
                '& .MuiChip-deleteIcon': {fontSize: size === 'tiny' ? '0.8rem' : '1rem', color: 'inherit'},
            }}
        />
    );
}

// --- type-along modal (Req 3 + 8) -----------------------------------------

interface TagPickerModalProps {
    open: boolean;
    onClose: () => void;
    onPick: (tag: string) => void;
}

function TagPickerModal({open, onClose, onPick}: TagPickerModalProps) {
    const [value, setValue] = useState('');
    const error = value ? tagError(value) : null;

    const submit = (raw: string) => {
        const t = normalise(raw);
        if (!t || tagError(t)) return;
        onPick(t);
        setValue('');
        onClose();
    };

    return (
        <Dialog open={open} onClose={onClose} fullWidth maxWidth="xs">
            <DialogTitle>Pick or create a tag</DialogTitle>
            <DialogContent sx={{pt: 1}}>
                <Autocomplete
                    freeSolo
                    autoHighlight
                    openOnFocus
                    options={favoriteFirst(KNOWN_TAGS)}
                    inputValue={value}
                    onInputChange={(_, v) => setValue(v)}
                    // ENTER-as-typed vs TAB/ARROW-selects-existing is native Autocomplete.
                    onChange={(_, v) => typeof v === 'string' && submit(v)}
                    renderInput={(params) => (
                        <TextField
                            {...params}
                            autoFocus
                            placeholder="people:luke"
                            error={Boolean(error)}
                            helperText={error ?? 'Enter to use as typed · Tab/↓ to pick an existing one'}
                            onKeyDown={(e) => {
                                if (e.key === 'Enter' && value && !tagError(value)) submit(value);
                            }}
                        />
                    )}
                />
            </DialogContent>
        </Dialog>
    );
}

// --- desktop banner (Req 1 tooltip owner, paint affordance) ---------------

interface TaggingBannerProps {
    activeTag: string | null;
    onSelect: (tag: string) => void;
    onExit: () => void;
}

function TaggingBanner({activeTag, onSelect, onExit}: TaggingBannerProps) {
    const recommended = useMemo(() => favoriteFirst([...PREVIOUS_TAGS, ...KNOWN_TAGS]).slice(0, 7), []);
    const [modalOpen, setModalOpen] = useState(false);

    return (
        <Paper
            sx={{
                p: 1.5,
                mb: 2,
                bgcolor: 'rgba(25,25,30,0.96)',
                border: '1px solid',
                borderColor: 'primary.main',
                position: 'sticky',
                top: 8,
                zIndex: 20,
            }}
        >
            <Box sx={{display: 'flex', alignItems: 'center', gap: 1, flexWrap: 'wrap'}}>
                <LocalOfferIcon fontSize="small" sx={{color: 'primary.light'}}/>
                <Typography variant="body2" sx={{opacity: 0.8, mr: 0.5}}>
                    {activeTag ? 'Click photos to toggle' : 'Pick a tag to apply'}
                </Typography>

                {activeTag && <TagChip tag={activeTag} active/>}

                <Box sx={{flex: 1}}/>

                <Button size="small" variant="text" color="secondary" onClick={onExit}>
                    Exit (Esc)
                </Button>
            </Box>

            {/* selected active | — | recommended | — | + More… */}
            <Box sx={{display: 'flex', flexWrap: 'wrap', alignItems: 'center', gap: 0.5, mt: 1}}>
                {recommended.map((t) => (
                    <TagChip key={t} tag={t} active={t === activeTag} onClick={() => onSelect(t)}/>
                ))}
                <Divider orientation="vertical" flexItem sx={{mx: 0.5, borderColor: 'rgba(255,255,255,0.2)'}}/>
                <Chip
                    label="+ More…"
                    size="small"
                    onClick={() => setModalOpen(true)}
                    sx={{height: 26, fontSize: '0.75rem', color: 'common.white', bgcolor: 'transparent', border: '1px dashed rgba(255,255,255,0.4)'}}
                />
            </Box>

            <TagPickerModal open={modalOpen} onClose={() => setModalOpen(false)} onPick={onSelect}/>
        </Paper>
    );
}

// --- media cell ------------------------------------------------------------

interface TaggingCellProps {
    media: FakeMedia;
    taggingMode: boolean;
    activeTag: string | null;
    onToggleActive: () => void;
    onToggleTag: (tag: string) => void;
}

function TaggingCell({media, taggingMode, activeTag, onToggleActive, onToggleTag}: TaggingCellProps) {
    const hasActive = activeTag ? media.tags.includes(activeTag) : false;
    const [sheetAnchor, setSheetAnchor] = useState<HTMLElement | null>(null);
    const [modalOpen, setModalOpen] = useState(false);

    const recommended = useMemo(() => favoriteFirst([...PREVIOUS_TAGS, ...KNOWN_TAGS]).slice(0, 6), []);
    const selected = favoriteFirst(media.tags);

    return (
        <Box
            onClick={taggingMode && activeTag ? onToggleActive : undefined}
            sx={{
                aspectRatio: '1',
                position: 'relative',
                overflow: 'hidden',
                bgcolor: 'rgba(255,255,255,0.05)',
                cursor: taggingMode && activeTag ? 'pointer' : 'default',
                outline: hasActive ? '3px solid' : 'none',
                outlineColor: 'primary.main',
                outlineOffset: '-3px',
            }}
        >
            <Box component="img" src={media.url} alt="" sx={{width: '100%', height: '100%', objectFit: 'cover', display: 'block'}}/>

            {/* Per-image TAG icon — the mobile affordance (Req 7/8). Shown on xs only. */}
            {taggingMode && (
                <IconButton
                    size="small"
                    onClick={(e) => {
                        e.stopPropagation();
                        setSheetAnchor(e.currentTarget);
                    }}
                    sx={{position: 'absolute', top: 2, right: 2, display: {xs: 'flex', md: 'none'}, color: 'common.white', bgcolor: 'rgba(0,0,0,0.5)'}}
                >
                    <LocalOfferIcon sx={{fontSize: '1rem'}}/>
                </IconButton>
            )}

            {/* Tags on the image. Req 7: on xs, hidden until the TAG icon is tapped. */}
            {media.tags.length > 0 && (
                <Box sx={{position: 'absolute', left: 0, right: 0, bottom: 0, p: 0.5, display: {xs: 'none', md: 'flex'}, flexWrap: 'wrap', gap: 0.25, bgcolor: 'rgba(0,0,0,0.5)'}}>
                    {selected.map((t) => (
                        <TagChip
                            key={t}
                            tag={t}
                            size="tiny"
                            active={t === activeTag}
                            onRemove={taggingMode ? () => onToggleTag(t) : undefined}
                        />
                    ))}
                </Box>
            )}

            {/* Mobile per-image popover: selected | — | recommended | — | + More… */}
            <Popover
                open={Boolean(sheetAnchor)}
                anchorEl={sheetAnchor}
                onClose={() => setSheetAnchor(null)}
                anchorOrigin={{vertical: 'bottom', horizontal: 'center'}}
                transformOrigin={{vertical: 'top', horizontal: 'center'}}
                onClick={(e) => e.stopPropagation()}
            >
                <Box sx={{p: 1.5, width: 260, bgcolor: 'background.paper'}}>
                    <Typography variant="caption" sx={{opacity: 0.6}}>On this photo</Typography>
                    <Box sx={{display: 'flex', flexWrap: 'wrap', gap: 0.5, mt: 0.5, mb: 1}}>
                        {selected.length ? (
                            selected.map((t) => <TagChip key={t} tag={t} size="tiny" active onRemove={() => onToggleTag(t)}/>)
                        ) : (
                            <Typography variant="caption" sx={{opacity: 0.5}}>None yet</Typography>
                        )}
                    </Box>
                    <Divider sx={{my: 1}}/>
                    <Typography variant="caption" sx={{opacity: 0.6}}>Recommended</Typography>
                    <Box sx={{display: 'flex', flexWrap: 'wrap', gap: 0.5, mt: 0.5, mb: 1}}>
                        {recommended.map((t) => (
                            <TagChip key={t} tag={t} size="tiny" active={media.tags.includes(t)} onClick={() => onToggleTag(t)}/>
                        ))}
                    </Box>
                    <Divider sx={{my: 1}}/>
                    <Chip
                        label="+ More…"
                        size="small"
                        onClick={() => setModalOpen(true)}
                        sx={{height: 26, fontSize: '0.75rem', color: 'text.primary', border: '1px dashed', borderColor: 'divider'}}
                    />
                </Box>
            </Popover>

            <TagPickerModal open={modalOpen} onClose={() => setModalOpen(false)} onPick={(t) => onToggleTag(t)}/>
        </Box>
    );
}

// --- page ------------------------------------------------------------------

function AlbumPageTaggingMode() {
    const [medias, setMedias] = useState<FakeMedia[]>(() => makeMedias());
    const [taggingMode, setTaggingMode] = useState(false);
    const [activeTag, setActiveTag] = useState<string | null>(null);

    useEffect(() => {
        const onKey = (e: KeyboardEvent) => {
            const el = document.activeElement;
            const typing = el && (el.tagName === 'INPUT' || el.tagName === 'TEXTAREA');
            // Req 2: plain `t` (no modifier) enters the mode; ignore while typing.
            if (e.key === 't' && !typing && !taggingMode) {
                e.preventDefault();
                setTaggingMode(true);
            }
            if (e.key === 'Escape') {
                setTaggingMode(false);
                setActiveTag(null);
            }
        };
        window.addEventListener('keydown', onKey);
        return () => window.removeEventListener('keydown', onKey);
    }, [taggingMode]);

    // Req 4: toggle — add if absent, remove if present.
    const toggleTagOn = (id: string, tag: string) =>
        setMedias((prev) =>
            prev.map((m) =>
                m.id === id ? {...m, tags: m.tags.includes(tag) ? m.tags.filter((t) => t !== tag) : [...m.tags, tag]} : m,
            ),
        );

    return (
        <Box sx={{p: {xs: 2, md: 4}}}>
            {/* Album banner */}
            <Box sx={{py: 2, display: 'flex', alignItems: 'center', gap: 2}}>
                <Box sx={{flex: 1, minWidth: 0}}>
                    <Typography variant="h1" noWrap>2024-07 — Summer holiday</Typography>
                    <Typography variant="body1" sx={{fontSize: '0.85rem'}}>Sun, 14 July 2024 · {medias.length} photos</Typography>
                </Box>
                <Stack direction="row" gap={1} alignItems="center" sx={{flexShrink: 0}}>
                    <Button variant="outlined" color="secondary" startIcon={<PlayArrowIcon/>} onClick={() => setTaggingMode(true)}>
                        Play
                    </Button>
                    {/* Req 1: tooltip explains the shortcut; no inline helper text. */}
                    <Tooltip title="Tag photos — shortcut: t">
                        <Button
                            variant={taggingMode ? 'contained' : 'outlined'}
                            color={taggingMode ? 'primary' : 'secondary'}
                            startIcon={<LocalOfferIcon/>}
                            onClick={() => setTaggingMode((v) => !v)}
                        >
                            Tag
                        </Button>
                    </Tooltip>
                </Stack>
            </Box>

            {/* Banner is desktop only (md+). */}
            {taggingMode && (
                <Box sx={{display: {xs: 'none', md: 'block'}}}>
                    <TaggingBanner
                        activeTag={activeTag}
                        onSelect={setActiveTag}
                        onExit={() => {
                            setTaggingMode(false);
                            setActiveTag(null);
                        }}
                    />
                </Box>
            )}

            {/* Media grid */}
            <Box sx={{borderTop: '1px solid rgba(255,255,255,0.07)', pt: 2}}>
                <Typography variant="h2" sx={{mb: 1.5}}>Sun, 14 July 2024</Typography>
                <Box
                    sx={{
                        display: 'grid',
                        gridTemplateColumns: {xs: 'repeat(2, 1fr)', sm: 'repeat(3, 1fr)', md: 'repeat(4, 1fr)', lg: 'repeat(5, 1fr)'},
                        gap: '2px',
                    }}
                >
                    {medias.map((m) => (
                        <TaggingCell
                            key={m.id}
                            media={m}
                            taggingMode={taggingMode}
                            activeTag={activeTag}
                            onToggleActive={() => activeTag && toggleTagOn(m.id, activeTag)}
                            onToggleTag={(tag) => toggleTagOn(m.id, tag)}
                        />
                    ))}
                </Box>
            </Box>
        </Box>
    );
}

const meta = {
    title: 'Prototypes/Tagging/AlbumPageTaggingMode',
    component: AlbumPageTaggingMode,
    parameters: {layout: 'fullscreen'},
    decorators: [
        (Story) => (
            <AppBackground>
                <Story/>
            </AppBackground>
        ),
    ],
} satisfies Meta<typeof AlbumPageTaggingMode>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};
