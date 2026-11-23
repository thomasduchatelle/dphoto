// PROTOTYPE — throwaway. Spike for issue 01 (album-page tagging interaction).
// Story group: "Prototypes / Tagging". Pure UI, components embedded, minimal state.
// Goal: settle the FINAL LOOK of the album page in "tagging mode".
//
// Flow:
//   - Enter tagging mode via the "Tag" button in the album banner, the "play"-style
//     quick link, or Ctrl+T / Cmd+T.
//   - The tagging banner opens with focus on a type-ahead, pre-populated with a few
//     existing / previously-used tags.
//   - Select a tag by pressing Enter in the type-ahead, or by clicking an existing
//     tag chip. The active tag is simply the highlighted (primary) one.
//   - While in mode, every image shows all its tags. Clicking an image adds the
//     active tag. Hovering an existing tag shows a cross; clicking it removes the tag.

import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {Box, Button, Chip, IconButton, InputBase, Paper, Stack, Typography} from '@mui/material';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import LocalOfferIcon from '@mui/icons-material/LocalOffer';
import CloseIcon from '@mui/icons-material/Close';
import {useEffect, useMemo, useRef, useState} from 'react';
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

// --- tagging banner --------------------------------------------------------

interface TaggingBannerProps {
    activeTag: string | null;
    onSelect: (tag: string) => void;
    onExit: () => void;
}

function TaggingBanner({activeTag, onSelect, onExit}: TaggingBannerProps) {
    const [query, setQuery] = useState('');
    const inputRef = useRef<HTMLInputElement>(null);

    useEffect(() => {
        inputRef.current?.focus();
    }, []);

    const suggestions = useMemo(() => {
        const q = query.trim().toLowerCase().replace(/\s+/g, '');
        const pool = query ? KNOWN_TAGS.filter((t) => t.includes(q)) : [...PREVIOUS_TAGS, ...KNOWN_TAGS.filter((t) => !PREVIOUS_TAGS.includes(t))];
        return Array.from(new Set(pool)).slice(0, 7);
    }, [query]);

    const commit = (tag: string) => {
        const t = tag.trim().toLowerCase().replace(/\s+/g, '');
        if (t) onSelect(t);
        setQuery('');
    };

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
                    {activeTag ? 'Click photos to tag them with' : 'Pick a tag to apply'}
                </Typography>

                <InputBase
                    inputRef={inputRef}
                    value={query}
                    placeholder="type a tag — e.g. people:luke"
                    onChange={(e) => setQuery(e.target.value)}
                    onKeyDown={(e) => {
                        if (e.key === 'Enter') commit(query || (suggestions[0] ?? ''));
                        if (e.key === 'Escape') onExit();
                    }}
                    sx={{px: 1.25, py: 0.5, bgcolor: 'rgba(255,255,255,0.08)', borderRadius: 1, color: 'common.white', fontSize: '0.9rem', minWidth: 220}}
                />

                <Button size="small" variant="text" color="secondary" onClick={onExit} sx={{ml: 'auto'}}>
                    Exit (Esc)
                </Button>
            </Box>

            <Box sx={{display: 'flex', flexWrap: 'wrap', gap: 0.5, mt: 1}}>
                {suggestions.map((t) => (
                    <Chip
                        key={t}
                        label={t}
                        size="small"
                        onClick={() => onSelect(t)}
                        color={t === activeTag ? 'primary' : 'default'}
                        variant={t === activeTag ? 'filled' : 'outlined'}
                        sx={{color: t === activeTag ? undefined : 'common.white'}}
                    />
                ))}
            </Box>
        </Paper>
    );
}

// --- media cell ------------------------------------------------------------

interface TaggingCellProps {
    media: FakeMedia;
    taggingMode: boolean;
    activeTag: string | null;
    onAddActive: () => void;
    onRemoveTag: (tag: string) => void;
}

function TaggingCell({media, taggingMode, activeTag, onAddActive, onRemoveTag}: TaggingCellProps) {
    const hasActive = activeTag ? media.tags.includes(activeTag) : false;

    return (
        <Box
            onClick={taggingMode && activeTag ? onAddActive : undefined}
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

            {media.tags.length > 0 && (
                <Box sx={{position: 'absolute', left: 0, right: 0, bottom: 0, p: 0.5, display: 'flex', flexWrap: 'wrap', gap: 0.25, bgcolor: 'rgba(0,0,0,0.5)'}}>
                    {media.tags.map((t) => (
                        <RemovableTag
                            key={t}
                            label={t}
                            active={t === activeTag}
                            removable={taggingMode}
                            onRemove={() => onRemoveTag(t)}
                        />
                    ))}
                </Box>
            )}
        </Box>
    );
}

function RemovableTag({label, active, removable, onRemove}: {label: string; active: boolean; removable: boolean; onRemove: () => void}) {
    const [hover, setHover] = useState(false);
    return (
        <Chip
            label={label}
            size="small"
            onMouseEnter={() => setHover(true)}
            onMouseLeave={() => setHover(false)}
            onDelete={removable && hover ? onRemove : undefined}
            deleteIcon={<CloseIcon/>}
            onClick={removable ? (e) => {e.stopPropagation(); onRemove();} : undefined}
            sx={{
                height: 20,
                fontSize: '0.6rem',
                color: active ? 'primary.contrastText' : 'common.white',
                bgcolor: active ? 'primary.main' : 'rgba(255,255,255,0.22)',
                '& .MuiChip-deleteIcon': {fontSize: '0.8rem', color: 'inherit'},
            }}
        />
    );
}

// --- page ------------------------------------------------------------------

function AlbumPageTaggingMode() {
    const [medias, setMedias] = useState<FakeMedia[]>(() => makeMedias());
    const [taggingMode, setTaggingMode] = useState(false);
    const [activeTag, setActiveTag] = useState<string | null>(null);

    useEffect(() => {
        const onKey = (e: KeyboardEvent) => {
            if ((e.key === 't' || e.key === 'T') && (e.metaKey || e.ctrlKey)) {
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
    }, []);

    const addActiveTo = (id: string) => {
        if (!activeTag) return;
        setMedias((prev) => prev.map((m) => (m.id === id && !m.tags.includes(activeTag) ? {...m, tags: [...m.tags, activeTag]} : m)));
    };
    const removeTagFrom = (id: string, tag: string) =>
        setMedias((prev) => prev.map((m) => (m.id === id ? {...m, tags: m.tags.filter((t) => t !== tag)} : m)));

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
                    <Button
                        variant={taggingMode ? 'contained' : 'outlined'}
                        color={taggingMode ? 'primary' : 'secondary'}
                        startIcon={<LocalOfferIcon/>}
                        onClick={() => setTaggingMode((v) => !v)}
                    >
                        Tag
                    </Button>
                </Stack>
            </Box>

            {taggingMode && (
                <TaggingBanner
                    activeTag={activeTag}
                    onSelect={setActiveTag}
                    onExit={() => {
                        setTaggingMode(false);
                        setActiveTag(null);
                    }}
                />
            )}

            {!taggingMode && (
                <Typography variant="body2" sx={{opacity: 0.5, mb: 1}}>
                    Enter tagging mode: “Tag” button, the Play quick-link, or <kbd>Ctrl</kbd>/<kbd>⌘</kbd>+<kbd>T</kbd>.
                </Typography>
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
                            onAddActive={() => addActiveTo(m.id)}
                            onRemoveTag={(tag) => removeTagFrom(m.id, tag)}
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
