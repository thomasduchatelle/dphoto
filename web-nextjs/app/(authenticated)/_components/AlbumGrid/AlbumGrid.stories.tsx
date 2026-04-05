import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {fn} from 'storybook/test';
import {AlbumGrid} from './index';
import {Box} from '@mui/material';
import {AppBackground} from '@/components/AppLayout/AppBackground';
import {Album, AlbumId} from '@/domains/catalog/language/catalog-state';

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
        sharedWith: [
            {user: {name: 'Tony Stark', email: 'ironman@avenger.com', picture: '/static/tonystark-profile.jpg'}},
        ],
        thumbnails: [
            '/thumbnails/astro-bot-01.jpg',
            '/thumbnails/astro-bot-02.jpg',
            '/thumbnails/astro-bot-03.jpg',
        ],
    },
    {
        albumId: createAlbumId('cdprojekt', 'the-witcher-3'),
        name: 'The Witcher 3: Wild Hunt — Blood and Wine',
        start: new Date('2024-05-30'),
        end: new Date('2024-06-14'),
        totalCount: 189,
        temperature: 9.4,
        relativeTemperature: 0.55,
        sharedWith: [],
        thumbnails: [
            '/thumbnails/the-witcher-3-01.jpg',
            '/thumbnails/the-witcher-3-02.jpg',
            '/thumbnails/the-witcher-3-03.jpg',
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
                {name: 'Tony Stark', email: 'ironman@avenger.com', picture: '/static/tonystark-profile.jpg'},
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
    {
        albumId: createAlbumId('kojima', 'death-stranding-2'),
        name: 'Death Stranding 2: On The Beach — A very long album title that should be truncated',
        start: new Date('2025-06-05'),
        end: new Date('2025-06-30'),
        totalCount: 78,
        temperature: 4.1,
        relativeTemperature: 0.22,
        sharedWith: [
            {user: {name: 'Tony Stark', email: 'ironman@avenger.com', picture: '/static/tonystark-profile.jpg'}},
        ],
        thumbnails: [
            '/thumbnails/death-stranding-2-01.jpg',
        ],
    },
    {
        albumId: createAlbumId('sandfall', 'clair-obscur-dlc'),
        name: 'Clair Obscur DLC',
        start: new Date('2025-09-10'),
        end: new Date('2025-09-15'),
        totalCount: 11,
        temperature: 1.2,
        relativeTemperature: 0.06,
        sharedWith: [],
        thumbnails: [],
    },
];

const meta = {
    title: 'Catalog/AlbumGrid',
    component: AlbumGrid,
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
        onShare: fn(),
        onCreateAlbum: fn(),
    },
} satisfies Meta<typeof AlbumGrid>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Empty: Story = {
    args: {
        albums: [],
    },
};
