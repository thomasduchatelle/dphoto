import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {fn} from 'storybook/test';
import {AlbumFilterControl} from './index';
import {Box} from '@mui/material';
import {AppBackground} from '@/components/AppLayout/AppBackground';
import {AlbumFilterEntry} from '@/domains/catalog/language/catalog-state';

const SAMPLE_OPTIONS: AlbumFilterEntry[] = [
    {
        name: 'All albums',
        criterion: {owners: []},
        avatars: ['/static/black-widow-profile.jpg', '/static/hulk-profile.webp'],
    },
    {
        name: 'My albums',
        criterion: {selfOwned: true, owners: []},
        avatars: ['/static/black-widow-profile.jpg'],
    },
    {
        name: 'Tony Stark',
        criterion: {owners: ['tony']},
        avatars: ['/static/tonystark-profile.jpg'],
    },
];

const MANY_OPTIONS: AlbumFilterEntry[] = [
    ...SAMPLE_OPTIONS,
    {
        name: 'Bruce Banner',
        criterion: {owners: ['bruce']},
        avatars: ['/static/hulk-profile.webp'],
    },
    {
        name: 'Steve Rogers',
        criterion: {owners: ['steve']},
        avatars: [],
    },
    {
        name: 'Thor Odinson',
        criterion: {owners: ['thor']},
        avatars: [],
    },
];

const meta = {
    title: 'Components/AlbumFilterControl',
    component: AlbumFilterControl,
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
        filterOptions: SAMPLE_OPTIONS,
        activeFilter: SAMPLE_OPTIONS[0],
        onFilterChange: fn(),
    },
} satisfies Meta<typeof AlbumFilterControl>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const MyAlbumsSelected: Story = {
    args: {
        activeFilter: SAMPLE_OPTIONS[1],
    },
};

export const SpecificOwnerSelected: Story = {
    args: {
        activeFilter: SAMPLE_OPTIONS[2],
    },
};

export const ManyOptions: Story = {
    args: {
        filterOptions: MANY_OPTIONS,
        activeFilter: MANY_OPTIONS[0],
    },
};

export const SingleOption: Story = {
    args: {
        filterOptions: [SAMPLE_OPTIONS[0]],
        activeFilter: SAMPLE_OPTIONS[0],
    },
};
