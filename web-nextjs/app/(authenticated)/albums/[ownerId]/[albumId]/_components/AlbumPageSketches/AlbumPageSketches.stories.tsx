/**
 * SKETCHES — Album Page Layout Concepts
 *
 * Three independent layout ideas for the full album page.
 * These are intentional rough sketches: all data is hardcoded, components are
 * self-contained inline, and no project abstractions are reused.
 * The goal is to compare page-level compositions before committing to one.
 *
 * Layout A — Cinematic Header
 *   Full-bleed hero strip built from thumbnail slices, album title + dates centred
 *   on top of it, action buttons (Share / Edit) floating top-right of the hero.
 *   Day-grouped media grid below. Sibling albums in a right sidebar on lg screens.
 *
 * Layout B — Magazine Strip
 *   Compact sticky sub-header (name + dates + actions) below the app bar.
 *   Horizontal scrollable row of sibling album cards above the grid.
 *   Mobile: swappable prev/next album bar pinned at the bottom of the viewport.
 *
 * Layout C — Explorer Rail
 *   Persistent left rail (lg+) lists all sibling albums chronologically.
 *   Wide main column shows the day-grouped grid; album title lives in a slim band
 *   between the app bar and the grid.
 *   Mobile / tablet: collapsed rail replaced by a horizontal album strip at top.
 */

import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {
    AppBar,
    Avatar,
    Box,
    Button,
    Chip,
    Divider,
    Grid,
    IconButton,
    Stack,
    Toolbar,
    Tooltip,
    Typography,
} from '@mui/material';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import ShareIcon from '@mui/icons-material/Share';
import EditIcon from '@mui/icons-material/Edit';
import PlayCircleOutlineIcon from '@mui/icons-material/PlayCircleOutline';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import {AppBackground} from '@/components/AppLayout/AppBackground';

// ---------------------------------------------------------------------------
// Shared sample data
// ---------------------------------------------------------------------------

const ALBUM = {
    name: 'Japan — Cherry Blossom',
    start: new Date(2025, 2, 28),
    end: new Date(2025, 3, 12),
    totalCount: 137,
};

const SIBLINGS = [
    {name: 'Iceland Winter', start: new Date(2025, 0, 3), end: new Date(2025, 0, 15), count: 84},
    {name: 'Paris Spring', start: new Date(2025, 2, 14), end: new Date(2025, 2, 21), count: 62},
    {name: 'Morocco Desert', start: new Date(2025, 4, 1), end: new Date(2025, 4, 10), count: 108},
    {name: 'Norwegian Fjords', start: new Date(2025, 5, 18), end: new Date(2025, 6, 2), count: 95},
    {name: 'Kyoto Autumn', start: new Date(2025, 9, 5), end: new Date(2025, 9, 18), count: 121},
];

const MEDIA_DAYS = [
    {
        day: new Date(2025, 2, 28),
        medias: [
            {id: 'm1', isVideo: false},
            {id: 'm2', isVideo: false},
            {id: 'm3', isVideo: true},
            {id: 'm4', isVideo: false},
            {id: 'm5', isVideo: false},
        ],
    },
    {
        day: new Date(2025, 2, 29),
        medias: [
            {id: 'm6', isVideo: false},
            {id: 'm7', isVideo: false},
            {id: 'm8', isVideo: false},
            {id: 'm9', isVideo: true},
            {id: 'm10', isVideo: false},
            {id: 'm11', isVideo: false},
            {id: 'm12', isVideo: false},
        ],
    },
    {
        day: new Date(2025, 3, 1),
        medias: [
            {id: 'm13', isVideo: false},
            {id: 'm14', isVideo: false},
            {id: 'm15', isVideo: false},
        ],
    },
];

// ---------------------------------------------------------------------------
// Tiny helpers
// ---------------------------------------------------------------------------

const fmt = (d: Date) =>
    d.toLocaleDateString('en-US', {month: 'short', day: 'numeric', year: 'numeric'});

const fmtDay = (d: Date) =>
    d.toLocaleDateString('en-US', {weekday: 'long', month: 'long', day: 'numeric'});

// Deterministic pastel hue from a string so every thumbnail has a unique colour
const hue = (id: string) => (id.split('').reduce((a, c) => a + c.charCodeAt(0), 0) * 47) % 360;

// A fake thumbnail square
function Thumb({id, isVideo, size = 120}: {id: string; isVideo: boolean; size?: number}) {
    return (
        <Box
            sx={{
                width: size,
                height: size,
                bgcolor: `hsl(${hue(id)}, 40%, 22%)`,
                borderRadius: '2px',
                position: 'relative',
                flexShrink: 0,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
            }}
        >
            {isVideo && (
                <PlayCircleOutlineIcon
                    sx={{
                        color: 'rgba(255,255,255,0.85)',
                        fontSize: 32,
                        position: 'absolute',
                    }}
                />
            )}
        </Box>
    );
}

// A fake album mini-card with 4 thumb squares side by side
function AlbumMiniCard({name, start, end, count, selected = false}: {
    name: string;
    start: Date;
    end: Date;
    count: number;
    selected?: boolean;
}) {
    return (
        <Box
            sx={{
                border: selected ? '2px solid #4a9ece' : '1px solid rgba(255,255,255,0.1)',
                borderRadius: 1,
                overflow: 'hidden',
                cursor: 'pointer',
                bgcolor: selected ? 'rgba(74,158,206,0.08)' : 'rgba(255,255,255,0.03)',
                '&:hover': {bgcolor: 'rgba(74,158,206,0.12)'},
                transition: 'all 0.2s',
                minWidth: 140,
            }}
        >
            {/* 4 thumbnail strip */}
            <Box sx={{display: 'flex', gap: '1px', height: 56, overflow: 'hidden'}}>
                {[name + '0', name + '1', name + '2', name + '3'].map((k) => (
                    <Box
                        key={k}
                        sx={{flex: 1, bgcolor: `hsl(${hue(k)}, 35%, 25%)`}}
                    />
                ))}
            </Box>
            <Box sx={{p: 0.75}}>
                <Typography variant="caption" sx={{fontWeight: 600, display: 'block', lineHeight: 1.2}}>
                    {name}
                </Typography>
                <Typography variant="caption" sx={{color: 'rgba(255,255,255,0.45)', fontSize: '0.65rem'}}>
                    {fmt(start)} · {count} photos
                </Typography>
            </Box>
        </Box>
    );
}

// The day-grouped media grid (shared by all layouts; column count driven by prop)
function MediaGrid({columns = 4}: {columns?: number}) {
    return (
        <Box>
            {MEDIA_DAYS.map(({day, medias}) => (
                <Box key={day.toISOString()} sx={{mb: 4}}>
                    <Typography
                        variant="h2"
                        sx={{mb: 1.5, color: 'rgba(255,255,255,0.6)', fontSize: '0.8rem', letterSpacing: '0.12em', textTransform: 'uppercase'}}
                    >
                        {fmtDay(day)}
                    </Typography>
                    <Box
                        sx={{
                            display: 'grid',
                            gridTemplateColumns: `repeat(${columns}, 1fr)`,
                            gap: '2px',
                        }}
                    >
                        {medias.map((m) => (
                            <Box
                                key={m.id}
                                sx={{
                                    aspectRatio: '1',
                                    bgcolor: `hsl(${hue(m.id)}, 40%, 22%)`,
                                    borderRadius: '1px',
                                    position: 'relative',
                                    display: 'flex',
                                    alignItems: 'center',
                                    justifyContent: 'center',
                                    cursor: 'pointer',
                                    '&:hover': {opacity: 0.85},
                                }}
                            >
                                {m.isVideo && (
                                    <PlayCircleOutlineIcon sx={{color: 'rgba(255,255,255,0.85)', fontSize: 36}}/>
                                )}
                            </Box>
                        ))}
                    </Box>
                </Box>
            ))}
        </Box>
    );
}

// Shared app header bar stub (not the real component — avoids auth dependency)
function FakeAppBar() {
    return (
        <AppBar
            position="fixed"
            elevation={0}
            sx={{
                bgcolor: 'rgba(0,25,41,0.8)',
                backdropFilter: 'blur(10px)',
                borderBottom: '1px solid rgba(74,158,206,0.2)',
            }}
        >
            <Toolbar sx={{height: {xs: 56, sm: 64}}}>
                <Box sx={{display: 'flex', alignItems: 'center', flexGrow: 1, gap: 1}}>
                    <Box
                        sx={{
                            height: 28,
                            width: 100,
                            bgcolor: 'rgba(74,158,206,0.25)',
                            borderRadius: 0.5,
                        }}
                    />
                </Box>
                <Avatar
                    src="/tonystark-profile.jpg"
                    sx={{width: 32, height: 32}}
                />
            </Toolbar>
        </AppBar>
    );
}

// ---------------------------------------------------------------------------
// Meta — we export a single meta and many named stories
// ---------------------------------------------------------------------------

const meta = {
    title: 'Catalog/AlbumPage — Layout Sketches',
    component: Box,
    parameters: {layout: 'fullscreen'},
    decorators: [
        (Story: () => React.ReactNode) => (
            <AppBackground>
                <Story/>
            </AppBackground>
        ),
    ],
} satisfies Meta<typeof Box>;

export default meta;
type Story = StoryObj<typeof meta>;

// ===========================================================================
// LAYOUT A — Cinematic Header
//
// A full-bleed hero whose background is a mosaic of blurred thumbnails.
// The album title + date range are centred over the hero.
// Action buttons float in the top-right corner of the hero.
// Sibling albums appear in a right-hand sidebar on large screens.
// ===========================================================================

export const LayoutA_CinematicHeader: Story = {
    name: 'A — Cinematic Header',
    render: () => (
        <Box sx={{minHeight: '100vh'}}>
            <FakeAppBar/>

            {/* Hero */}
            <Box
                sx={{
                    mt: {xs: '56px', sm: '64px'},
                    position: 'relative',
                    height: {xs: 220, sm: 280, md: 340},
                    overflow: 'hidden',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                }}
            >
                {/* Blurred thumbnail mosaic background */}
                <Box
                    sx={{
                        position: 'absolute',
                        inset: 0,
                        display: 'grid',
                        gridTemplateColumns: 'repeat(8, 1fr)',
                        gap: 0,
                        filter: 'blur(12px)',
                        transform: 'scale(1.1)',
                    }}
                >
                    {Array.from({length: 16}).map((_, i) => (
                        <Box
                            key={i}
                            sx={{bgcolor: `hsl(${(i * 37 + 20) % 360}, 35%, 25%)`, aspectRatio: '1'}}
                        />
                    ))}
                </Box>
                {/* Dark overlay */}
                <Box
                    sx={{
                        position: 'absolute',
                        inset: 0,
                        background: 'linear-gradient(to bottom, rgba(0,25,41,0.5) 0%, rgba(0,25,41,0.75) 100%)',
                    }}
                />

                {/* Back button */}
                <IconButton
                    sx={{
                        position: 'absolute',
                        top: 16,
                        left: 16,
                        color: 'rgba(255,255,255,0.8)',
                        bgcolor: 'rgba(0,0,0,0.3)',
                        '&:hover': {bgcolor: 'rgba(0,0,0,0.5)'},
                    }}
                >
                    <ArrowBackIcon/>
                </IconButton>

                {/* Action buttons */}
                <Stack
                    direction="row"
                    gap={1}
                    sx={{position: 'absolute', top: 16, right: 16}}
                >
                    <Button
                        variant="outlined"
                        startIcon={<ShareIcon/>}
                        size="small"
                        sx={{
                            borderColor: 'rgba(255,255,255,0.4)',
                            color: 'white',
                            backdropFilter: 'blur(4px)',
                            bgcolor: 'rgba(0,0,0,0.25)',
                        }}
                    >
                        Share
                    </Button>
                    <Button
                        variant="outlined"
                        startIcon={<EditIcon/>}
                        size="small"
                        sx={{
                            borderColor: 'rgba(255,255,255,0.4)',
                            color: 'white',
                            backdropFilter: 'blur(4px)',
                            bgcolor: 'rgba(0,0,0,0.25)',
                        }}
                    >
                        Edit
                    </Button>
                </Stack>

                {/* Title + date centred */}
                <Box sx={{position: 'relative', textAlign: 'center', px: 4}}>
                    <Typography
                        sx={{
                            fontSize: {xs: '1.6rem', sm: '2.2rem', md: '2.8rem'},
                            fontWeight: 300,
                            color: '#ffffff',
                            letterSpacing: '-0.01em',
                            textShadow: '0 2px 20px rgba(0,0,0,0.6)',
                        }}
                    >
                        {ALBUM.name}
                    </Typography>
                    <Typography
                        sx={{
                            mt: 1,
                            color: 'rgba(255,255,255,0.7)',
                            fontSize: '0.9rem',
                            letterSpacing: '0.1em',
                            textTransform: 'uppercase',
                        }}
                    >
                        {fmt(ALBUM.start)} – {fmt(ALBUM.end)} · {ALBUM.totalCount} photos
                    </Typography>
                </Box>
            </Box>

            {/* Content area: grid + optional sidebar */}
            <Box
                sx={{
                    display: 'flex',
                    gap: 3,
                    px: {xs: 1, sm: 2, md: 4},
                    py: 3,
                    maxWidth: 1400,
                    mx: 'auto',
                }}
            >
                {/* Media grid */}
                <Box sx={{flex: 1, minWidth: 0}}>
                    <MediaGrid columns={4}/>
                </Box>

                {/* Right sidebar — sibling albums (lg+) */}
                <Box
                    sx={{
                        display: {xs: 'none', lg: 'flex'},
                        flexDirection: 'column',
                        gap: 1.5,
                        width: 200,
                        flexShrink: 0,
                        pt: 0.5,
                    }}
                >
                    <Typography
                        sx={{
                            fontSize: '0.7rem',
                            letterSpacing: '0.12em',
                            textTransform: 'uppercase',
                            color: 'rgba(255,255,255,0.4)',
                            mb: 0.5,
                        }}
                    >
                        Other albums
                    </Typography>
                    {SIBLINGS.map((s) => (
                        <AlbumMiniCard key={s.name} {...s}/>
                    ))}
                </Box>
            </Box>
        </Box>
    ),
};

// ===========================================================================
// LAYOUT B — Magazine Strip
//
// Compact sticky sub-header (album name + dates + actions) after the app bar.
// A horizontal scrollable strip of sibling albums sits between the sub-header
// and the media grid on all screen sizes.
// On mobile a prev/next bar is pinned at the bottom of the viewport.
// ===========================================================================

export const LayoutB_MagazineStrip: Story = {
    name: 'B — Magazine Strip',
    render: () => (
        <Box sx={{minHeight: '100vh', pb: {xs: 10, sm: 0}}}>
            <FakeAppBar/>

            {/* Sticky sub-header */}
            <Box
                sx={{
                    mt: {xs: '56px', sm: '64px'},
                    position: 'sticky',
                    top: {xs: 56, sm: 64},
                    zIndex: 100,
                    bgcolor: 'rgba(10,21,32,0.92)',
                    backdropFilter: 'blur(8px)',
                    borderBottom: '1px solid rgba(74,158,206,0.15)',
                    px: {xs: 2, sm: 3, md: 5},
                    py: 1.25,
                    display: 'flex',
                    alignItems: 'center',
                    gap: 2,
                    flexWrap: 'wrap',
                }}
            >
                <IconButton size="small" sx={{color: 'rgba(255,255,255,0.6)'}}>
                    <ArrowBackIcon fontSize="small"/>
                </IconButton>
                <Box sx={{flex: 1, minWidth: 0}}>
                    <Typography
                        sx={{fontWeight: 500, fontSize: '1.05rem', lineHeight: 1.2, color: '#fff'}}
                        noWrap
                    >
                        {ALBUM.name}
                    </Typography>
                    <Typography sx={{color: 'rgba(255,255,255,0.45)', fontSize: '0.75rem'}}>
                        {fmt(ALBUM.start)} – {fmt(ALBUM.end)} · {ALBUM.totalCount} photos
                    </Typography>
                </Box>
                <Stack direction="row" gap={1}>
                    <Tooltip title="Share">
                        <IconButton size="small" sx={{color: '#4a9ece'}}>
                            <ShareIcon fontSize="small"/>
                        </IconButton>
                    </Tooltip>
                    <Tooltip title="Edit">
                        <IconButton size="small" sx={{color: 'rgba(255,255,255,0.6)'}}>
                            <EditIcon fontSize="small"/>
                        </IconButton>
                    </Tooltip>
                </Stack>
            </Box>

            {/* Sibling albums horizontal strip (all screen sizes) */}
            <Box
                sx={{
                    px: {xs: 1, sm: 2, md: 5},
                    py: 2,
                    borderBottom: '1px solid rgba(255,255,255,0.06)',
                }}
            >
                <Typography
                    sx={{
                        fontSize: '0.7rem',
                        letterSpacing: '0.12em',
                        textTransform: 'uppercase',
                        color: 'rgba(255,255,255,0.35)',
                        mb: 1,
                    }}
                >
                    Other albums
                </Typography>
                <Box
                    sx={{
                        display: 'flex',
                        gap: 1.5,
                        overflowX: 'auto',
                        pb: 0.5,
                        '&::-webkit-scrollbar': {height: 3},
                        '&::-webkit-scrollbar-thumb': {bgcolor: 'rgba(74,158,206,0.3)', borderRadius: 2},
                    }}
                >
                    {SIBLINGS.map((s) => (
                        <AlbumMiniCard key={s.name} {...s}/>
                    ))}
                </Box>
            </Box>

            {/* Media grid */}
            <Box sx={{px: {xs: 1, sm: 2, md: 5}, py: 3}}>
                <MediaGrid columns={4}/>
            </Box>

            {/* Mobile prev/next bar pinned at bottom */}
            <Box
                sx={{
                    display: {xs: 'flex', sm: 'none'},
                    position: 'fixed',
                    bottom: 0,
                    left: 0,
                    right: 0,
                    bgcolor: 'rgba(10,21,32,0.96)',
                    backdropFilter: 'blur(8px)',
                    borderTop: '1px solid rgba(74,158,206,0.2)',
                    px: 2,
                    py: 1.5,
                    gap: 1.5,
                    alignItems: 'center',
                    zIndex: 200,
                }}
            >
                <Button
                    startIcon={<ChevronLeftIcon/>}
                    variant="text"
                    size="small"
                    sx={{color: 'rgba(255,255,255,0.7)', flex: 1, justifyContent: 'flex-start'}}
                >
                    Iceland Winter
                </Button>
                <Divider orientation="vertical" flexItem sx={{bgcolor: 'rgba(255,255,255,0.15)'}}/>
                <Button
                    endIcon={<ChevronRightIcon/>}
                    variant="text"
                    size="small"
                    sx={{color: 'rgba(255,255,255,0.7)', flex: 1, justifyContent: 'flex-end'}}
                >
                    Morocco Desert
                </Button>
            </Box>
        </Box>
    ),
};

// ===========================================================================
// LAYOUT C — Explorer Rail
//
// A persistent left rail on lg+ screens lists all sibling albums vertically.
// The current album is highlighted in the rail.
// The main column has a slim album-info band (name, dates, actions) then the grid.
// On sm/md a horizontal strip replaces the rail at the top of the page.
// ===========================================================================

export const LayoutC_ExplorerRail: Story = {
    name: 'C — Explorer Rail',
    render: () => (
        <Box sx={{minHeight: '100vh'}}>
            <FakeAppBar/>

            <Box
                sx={{
                    mt: {xs: '56px', sm: '64px'},
                    display: 'flex',
                    minHeight: 'calc(100vh - 64px)',
                }}
            >
                {/* Left rail — lg+ */}
                <Box
                    sx={{
                        display: {xs: 'none', lg: 'flex'},
                        flexDirection: 'column',
                        width: 220,
                        flexShrink: 0,
                        borderRight: '1px solid rgba(255,255,255,0.07)',
                        bgcolor: 'rgba(0,10,20,0.5)',
                        overflowY: 'auto',
                        position: 'sticky',
                        top: 64,
                        maxHeight: 'calc(100vh - 64px)',
                        p: 2,
                        gap: 1,
                    }}
                >
                    <Typography
                        sx={{
                            fontSize: '0.65rem',
                            letterSpacing: '0.14em',
                            textTransform: 'uppercase',
                            color: 'rgba(255,255,255,0.35)',
                            mb: 0.5,
                        }}
                    >
                        Albums
                    </Typography>
                    {/* Current album highlighted */}
                    <AlbumMiniCard {...ALBUM} count={ALBUM.totalCount} start={ALBUM.start} end={ALBUM.end} selected/>
                    {SIBLINGS.map((s) => (
                        <AlbumMiniCard key={s.name} {...s}/>
                    ))}
                </Box>

                {/* Main column */}
                <Box sx={{flex: 1, minWidth: 0, display: 'flex', flexDirection: 'column'}}>

                    {/* Horizontal album strip — xs/sm/md only */}
                    <Box
                        sx={{
                            display: {xs: 'block', lg: 'none'},
                            px: {xs: 1, sm: 2, md: 4},
                            py: 1.5,
                            borderBottom: '1px solid rgba(255,255,255,0.06)',
                        }}
                    >
                        <Box
                            sx={{
                                display: 'flex',
                                gap: 1.5,
                                overflowX: 'auto',
                                '&::-webkit-scrollbar': {height: 3},
                                '&::-webkit-scrollbar-thumb': {bgcolor: 'rgba(74,158,206,0.3)', borderRadius: 2},
                            }}
                        >
                            <AlbumMiniCard {...ALBUM} count={ALBUM.totalCount} start={ALBUM.start} end={ALBUM.end} selected/>
                            {SIBLINGS.map((s) => (
                                <AlbumMiniCard key={s.name} {...s}/>
                            ))}
                        </Box>
                    </Box>

                    {/* Slim album-info band */}
                    <Box
                        sx={{
                            px: {xs: 2, sm: 3, md: 5},
                            py: {xs: 2, sm: 3},
                            display: 'flex',
                            alignItems: {xs: 'flex-start', sm: 'center'},
                            flexDirection: {xs: 'column', sm: 'row'},
                            gap: {xs: 1.5, sm: 2},
                            borderBottom: '1px solid rgba(255,255,255,0.07)',
                        }}
                    >
                        <IconButton size="small" sx={{color: 'rgba(255,255,255,0.5)', ml: {xs: -0.5, sm: 0}}}>
                            <ArrowBackIcon fontSize="small"/>
                        </IconButton>
                        <Box sx={{flex: 1}}>
                            <Typography variant="h1" sx={{mb: 0.25}}>
                                {ALBUM.name}
                            </Typography>
                            <Typography variant="body1" sx={{fontSize: '0.85rem'}}>
                                {fmt(ALBUM.start)} – {fmt(ALBUM.end)} · {ALBUM.totalCount} photos
                            </Typography>
                        </Box>
                        <Stack direction="row" gap={1} flexShrink={0}>
                            <Button variant="outlined" startIcon={<ShareIcon/>} size="small">
                                Share
                            </Button>
                            <Button variant="outlined" startIcon={<EditIcon/>} size="small">
                                Edit
                            </Button>
                        </Stack>
                    </Box>

                    {/* Media grid */}
                    <Box sx={{px: {xs: 1, sm: 2, md: 5}, py: 3, flex: 1}}>
                        <MediaGrid columns={4}/>
                    </Box>
                </Box>
            </Box>
        </Box>
    ),
};

// ===========================================================================
// LAYOUT D — Immersive Full-Width
//
// No traditional header for the album — the app bar stays but the album name
// overlays the very first row of the grid as a large translucent watermark.
// A floating action pill (Share + Edit) lives in the bottom-right corner.
// Prev/next album arrows peek in from the left and right edges of the screen.
// Sibling albums appear in a bottom drawer that expands on hover/tap.
// ===========================================================================

export const LayoutD_ImmersiveFullWidth: Story = {
    name: 'D — Immersive Full-Width',
    render: () => (
        <Box sx={{minHeight: '100vh', position: 'relative'}}>
            <FakeAppBar/>

            <Box sx={{mt: {xs: '56px', sm: '64px'}, px: {xs: 0.5, sm: 1, md: 2}, pt: 2, position: 'relative'}}>

                {/* Album name watermark overlay — sits above the grid */}
                <Box
                    sx={{
                        position: 'relative',
                        mb: 2,
                        display: 'flex',
                        alignItems: 'flex-end',
                        gap: 2,
                        px: {xs: 1, md: 2},
                        pb: 2,
                        borderBottom: '1px solid rgba(255,255,255,0.06)',
                    }}
                >
                    <IconButton sx={{color: 'rgba(255,255,255,0.5)'}}>
                        <ArrowBackIcon/>
                    </IconButton>
                    <Box sx={{flex: 1}}>
                        <Typography
                            sx={{
                                fontSize: {xs: '1.8rem', sm: '2.5rem'},
                                fontWeight: 200,
                                color: 'rgba(255,255,255,0.9)',
                                lineHeight: 1,
                                letterSpacing: '-0.02em',
                            }}
                        >
                            {ALBUM.name}
                        </Typography>
                        <Stack direction="row" gap={1} mt={0.75} flexWrap="wrap">
                            <Chip
                                label={`${fmt(ALBUM.start)} – ${fmt(ALBUM.end)}`}
                                size="small"
                                sx={{
                                    bgcolor: 'rgba(74,158,206,0.15)',
                                    color: '#4a9ece',
                                    fontSize: '0.7rem',
                                    height: 22,
                                }}
                            />
                            <Chip
                                label={`${ALBUM.totalCount} photos`}
                                size="small"
                                sx={{
                                    bgcolor: 'rgba(255,255,255,0.07)',
                                    color: 'rgba(255,255,255,0.55)',
                                    fontSize: '0.7rem',
                                    height: 22,
                                }}
                            />
                        </Stack>
                    </Box>
                </Box>

                {/* Gutterless grid — edge to edge */}
                {MEDIA_DAYS.map(({day, medias}) => (
                    <Box key={day.toISOString()} sx={{mb: 3}}>
                        <Typography
                            sx={{
                                px: {xs: 1, md: 2},
                                mb: 0.75,
                                fontSize: '0.7rem',
                                letterSpacing: '0.14em',
                                textTransform: 'uppercase',
                                color: 'rgba(255,255,255,0.4)',
                            }}
                        >
                            {fmtDay(day)}
                        </Typography>
                        <Box
                            sx={{
                                display: 'grid',
                                gridTemplateColumns: 'repeat(5, 1fr)',
                                gap: '2px',
                            }}
                        >
                            {medias.map((m) => (
                                <Box
                                    key={m.id}
                                    sx={{
                                        aspectRatio: '1',
                                        bgcolor: `hsl(${hue(m.id)}, 40%, 20%)`,
                                        position: 'relative',
                                        display: 'flex',
                                        alignItems: 'center',
                                        justifyContent: 'center',
                                        cursor: 'pointer',
                                        '&:hover': {opacity: 0.82},
                                        transition: 'opacity 0.15s',
                                    }}
                                >
                                    {m.isVideo && (
                                        <PlayCircleOutlineIcon sx={{color: 'rgba(255,255,255,0.85)', fontSize: 40}}/>
                                    )}
                                </Box>
                            ))}
                        </Box>
                    </Box>
                ))}
            </Box>

            {/* Prev/next album peeking arrows */}
            <Box
                sx={{
                    display: {xs: 'none', md: 'flex'},
                    position: 'fixed',
                    left: 0,
                    top: '50%',
                    transform: 'translateY(-50%)',
                    flexDirection: 'column',
                    alignItems: 'flex-start',
                    zIndex: 300,
                }}
            >
                <Box
                    sx={{
                        bgcolor: 'rgba(10,21,32,0.85)',
                        backdropFilter: 'blur(6px)',
                        border: '1px solid rgba(255,255,255,0.1)',
                        borderLeft: 'none',
                        borderRadius: '0 8px 8px 0',
                        p: 1,
                        cursor: 'pointer',
                        '&:hover': {bgcolor: 'rgba(74,158,206,0.2)'},
                    }}
                >
                    <ChevronLeftIcon sx={{color: 'rgba(255,255,255,0.7)'}}/>
                    <Typography sx={{fontSize: '0.6rem', color: 'rgba(255,255,255,0.5)', px: 0.25}}>
                        Iceland
                    </Typography>
                </Box>
            </Box>
            <Box
                sx={{
                    display: {xs: 'none', md: 'flex'},
                    position: 'fixed',
                    right: 0,
                    top: '50%',
                    transform: 'translateY(-50%)',
                    flexDirection: 'column',
                    alignItems: 'flex-end',
                    zIndex: 300,
                }}
            >
                <Box
                    sx={{
                        bgcolor: 'rgba(10,21,32,0.85)',
                        backdropFilter: 'blur(6px)',
                        border: '1px solid rgba(255,255,255,0.1)',
                        borderRight: 'none',
                        borderRadius: '8px 0 0 8px',
                        p: 1,
                        cursor: 'pointer',
                        '&:hover': {bgcolor: 'rgba(74,158,206,0.2)'},
                    }}
                >
                    <ChevronRightIcon sx={{color: 'rgba(255,255,255,0.7)'}}/>
                    <Typography sx={{fontSize: '0.6rem', color: 'rgba(255,255,255,0.5)', px: 0.25}}>
                        Morocco
                    </Typography>
                </Box>
            </Box>

            {/* Floating action pill — bottom right */}
            <Box
                sx={{
                    position: 'fixed',
                    bottom: {xs: 24, sm: 32},
                    right: {xs: 16, sm: 24},
                    display: 'flex',
                    gap: 1,
                    zIndex: 300,
                }}
            >
                <Button
                    variant="contained"
                    startIcon={<ShareIcon/>}
                    size="small"
                    sx={{
                        bgcolor: 'rgba(24,89,134,0.9)',
                        backdropFilter: 'blur(6px)',
                        '&:hover': {bgcolor: '#185986'},
                        boxShadow: '0 4px 24px rgba(0,0,0,0.5)',
                    }}
                >
                    Share
                </Button>
                <Button
                    variant="outlined"
                    startIcon={<EditIcon/>}
                    size="small"
                    sx={{
                        borderColor: 'rgba(255,255,255,0.3)',
                        color: 'white',
                        backdropFilter: 'blur(6px)',
                        bgcolor: 'rgba(0,0,0,0.3)',
                        '&:hover': {bgcolor: 'rgba(74,158,206,0.15)', borderColor: '#4a9ece'},
                        boxShadow: '0 4px 24px rgba(0,0,0,0.4)',
                    }}
                >
                    Edit
                </Button>
            </Box>

            {/* Collapsed sibling drawer — sm and below */}
            <Box
                sx={{
                    display: {xs: 'flex', md: 'none'},
                    position: 'fixed',
                    bottom: 0,
                    left: 0,
                    right: 0,
                    bgcolor: 'rgba(5,15,25,0.96)',
                    backdropFilter: 'blur(10px)',
                    borderTop: '1px solid rgba(74,158,206,0.2)',
                    px: 1.5,
                    py: 1.5,
                    gap: 1.5,
                    overflowX: 'auto',
                    zIndex: 200,
                    '&::-webkit-scrollbar': {display: 'none'},
                }}
            >
                {[...SIBLINGS.slice(0, 2), ...SIBLINGS.slice(0, 2)].map((s, i) => (
                    <AlbumMiniCard key={s.name + i} {...s}/>
                ))}
            </Box>
        </Box>
    ),
};
