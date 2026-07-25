import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {fn} from 'storybook/test';
import {HomePageContent} from './index';
import {Box} from '@mui/material';
import {AppBackground} from '@/components/AppLayout/AppBackground';
import {Album, AlbumFilterEntry, AlbumId} from '@/domains/catalog/language/catalog-state';

const createAlbumId = (owner: string, folderName: string): AlbumId => ({owner, folderName});

const sampleAlbums: Album[] = [
    {
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
    },
    {
        albumId: createAlbumId('sony', 'astro-bot'),
        name: 'Astro Bot',
        start: new Date('2024-09-06'),
        end: new Date('2024-09-12'),
        totalCount: 23,
        temperature: 12.1,
        relativeTemperature: 0.72,
        sharedWith: [],
        thumbnails: [
            '/thumbnails/astro-bot-01.jpg',
            '/thumbnails/astro-bot-02.jpg',
            '/thumbnails/astro-bot-03.jpg',
        ],
    },
    {
        albumId: createAlbumId('kojima', 'death-stranding-1'),
        name: 'Death Stranding',
        start: new Date('2023-11-08'),
        end: new Date('2023-11-22'),
        totalCount: 312,
        temperature: 17.5,
        relativeTemperature: 1.0,
        ownedBy: {
            name: 'Kojima',
            users: [
                {name: 'Hideo Kojima', email: 'hideo@kojima.com', picture: '/static/tonystark-profile.jpg'},
            ],
        },
        sharedWith: [],
        thumbnails: [
            '/thumbnails/death-stranding-1-01.jpg',
            '/thumbnails/death-stranding-1-02.jpg',
            '/thumbnails/death-stranding-1-03.jpg',
            '/thumbnails/death-stranding-1-04.jpg',
        ],
    },
];

const sampleFilterOptions: AlbumFilterEntry[] = [
    {
        name: 'All albums',
        criterion: {owners: []},
        avatars: [],
    },
    {
        name: 'My albums',
        criterion: {selfOwned: true, owners: []},
        avatars: [],
    },
    {
        name: 'Kojima',
        criterion: {owners: ['kojima']},
        avatars: ['/static/tonystark-profile.jpg'],
    },
];

const meta = {
    title: 'Pages/HomePageContent',
    component: HomePageContent,
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
        albums: sampleAlbums,
        filterOptions: sampleFilterOptions,
        activeFilter: sampleFilterOptions[0],
        onFilterChange: fn(),
    },
} satisfies Meta<typeof HomePageContent>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const WithMyAlbumsSelected: Story = {
    args: {
        activeFilter: sampleFilterOptions[1],
        albums: sampleAlbums.filter(a => a.ownedBy === undefined),
    },
};

export const WithOwnerSelected: Story = {
    args: {
        activeFilter: sampleFilterOptions[2],
        albums: sampleAlbums.filter(a => a.albumId.owner === 'kojima'),
    },
};

export const EmptyFilterResult: Story = {
    args: {
        activeFilter: sampleFilterOptions[1],
        albums: [],
    },
};

export const EmptyCatalog: Story = {
    args: {
        albums: [],
        filterOptions: [{name: 'All albums', criterion: {owners: []}, avatars: []}],
        activeFilter: {name: 'All albums', criterion: {owners: []}, avatars: []},
    },
};

export const SingleFilterOption: Story = {
    args: {
        filterOptions: [sampleFilterOptions[0]],
        activeFilter: sampleFilterOptions[0],
    },
};

export const WithError: Story = {
    args: {
        error: new Error('Failed to load albums from the server'),
    },
};
