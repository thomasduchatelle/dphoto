// PROTOTYPE — throwaway. Spike for issue 01 (album-page tagging interaction).
// Story group: "Prototypes / Tagging". Pure UI, component embedded, minimal state.
// Goal: settle the FINAL LOOK of a taggable media thumbnail.
//
// A thumbnail is a clickable (link) image, as in the album media grid. It carries
// two overlay icons:
//   - a STAR that toggles `favorite` on click;
//   - a BOOKMARK that, on hover (or click on mobile), reveals a small strip of
//     quick tags to toggle (multi-select; clicking does not dismiss it).

import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {Box, Chip, IconButton, Paper} from '@mui/material';
import StarIcon from '@mui/icons-material/Star';
import StarBorderIcon from '@mui/icons-material/StarBorder';
import BookmarkIcon from '@mui/icons-material/BookmarkBorder';
import {useState} from 'react';
import {AppBackground} from '@/components/AppLayout/AppBackground';

const QUICK_TAGS = ['people:luke', 'people:leia', 'place:beach', 'activity:ski'];

interface TaggingThumbnailProps {
    thumbnailUrl: string;
    initialTags?: string[];
}

function TaggingThumbnail({thumbnailUrl, initialTags = []}: TaggingThumbnailProps) {
    const [tags, setTags] = useState<string[]>(initialTags);
    const [pickerOpen, setPickerOpen] = useState(false);

    const toggle = (tag: string) =>
        setTags((prev) => (prev.includes(tag) ? prev.filter((t) => t !== tag) : [...prev, tag]));

    const isFavorite = tags.includes('favorite');

    return (
        <Box
            sx={{
                aspectRatio: '1',
                position: 'relative',
                overflow: 'hidden',
                bgcolor: 'rgba(255,255,255,0.05)',
                cursor: 'pointer',
                '& .tb-icon': {opacity: {xs: 1, lg: 0}, transition: 'opacity .12s'},
                '&:hover .tb-icon': {opacity: 1},
            }}
        >
            <Box component="img" src={thumbnailUrl} alt="" sx={{width: '100%', height: '100%', objectFit: 'cover', display: 'block'}}/>

            {/* Star — always visible when favourited, otherwise on hover. */}
            <IconButton
                size="small"
                className={isFavorite ? undefined : 'tb-icon'}
                onClick={(e) => {
                    e.stopPropagation();
                    toggle('favorite');
                }}
                sx={{position: 'absolute', top: 2, left: 2, color: isFavorite ? '#ffd54a' : 'common.white', textShadow: '0 1px 2px rgba(0,0,0,.6)'}}
            >
                {isFavorite ? <StarIcon fontSize="small"/> : <StarBorderIcon fontSize="small"/>}
            </IconButton>

            {/* Bookmark — reveals the quick-tag picker. */}
            <IconButton
                size="small"
                className={tags.some((t) => t !== 'favorite') ? undefined : 'tb-icon'}
                onClick={(e) => {
                    e.stopPropagation();
                    setPickerOpen((v) => !v);
                }}
                sx={{position: 'absolute', top: 2, right: 2, color: pickerOpen ? 'primary.light' : 'common.white'}}
            >
                <BookmarkIcon fontSize="small"/>
            </IconButton>

            {/* Quick-tag picker strip, revealed from the bookmark. */}
            {pickerOpen && (
                <Paper
                    onClick={(e) => e.stopPropagation()}
                    sx={{position: 'absolute', top: 34, right: 2, p: 0.75, display: 'flex', flexDirection: 'column', gap: 0.5, bgcolor: 'rgba(20,20,20,0.95)', maxWidth: '85%'}}
                >
                    {QUICK_TAGS.map((t) => {
                        const active = tags.includes(t);
                        return (
                            <Chip
                                key={t}
                                label={t}
                                size="small"
                                onClick={() => toggle(t)}
                                sx={{height: 22, fontSize: '0.65rem', justifyContent: 'flex-start', color: active ? 'primary.contrastText' : 'common.white', bgcolor: active ? 'primary.main' : 'rgba(255,255,255,0.18)'}}
                            />
                        );
                    })}
                </Paper>
            )}

            {/* Current tags, always visible at the bottom. */}
            {tags.filter((t) => t !== 'favorite').length > 0 && (
                <Box sx={{position: 'absolute', left: 0, right: 0, bottom: 0, p: 0.5, display: 'flex', flexWrap: 'wrap', gap: 0.25, bgcolor: 'rgba(0,0,0,0.5)'}}>
                    {tags.filter((t) => t !== 'favorite').map((t) => (
                        <Chip key={t} label={t} size="small" sx={{height: 18, fontSize: '0.6rem', color: 'common.white', bgcolor: 'rgba(255,255,255,0.22)'}}/>
                    ))}
                </Box>
            )}
        </Box>
    );
}

const meta = {
    title: 'Prototypes/Tagging/TaggingButton',
    component: TaggingThumbnail,
    parameters: {layout: 'centered'},
    decorators: [
        (Story) => (
            <AppBackground>
                <Box sx={{width: 260}}>
                    <Story/>
                </Box>
            </AppBackground>
        ),
    ],
    args: {
        thumbnailUrl: '/thumbnails/clair-obscur-1.jpg',
    },
} satisfies Meta<typeof TaggingThumbnail>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Favourited: Story = {
    args: {thumbnailUrl: '/thumbnails/clair-obscur-2.jpg', initialTags: ['favorite']},
};

export const WithTags: Story = {
    args: {thumbnailUrl: '/thumbnails/clair-obscur-3.jpg', initialTags: ['favorite', 'people:luke', 'place:beach']},
};
