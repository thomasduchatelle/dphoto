import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {fn} from 'storybook/test';
import {AlbumCard} from './index';
import {Box} from '@mui/material';
import {Album, AlbumId} from '@/domains/catalog/language/catalog-state';
import {AppBackground} from "../../../../components/AppLayout/AppBackground";

const createAlbumId = (owner: string, folderName: string): AlbumId => ({owner, folderName});

const meta = {
    title: 'Catalog/AlbumCard',
    component: AlbumCard,
    parameters: {
        layout: 'fullscreen',
    },
    decorators: [
        (Story) => (
            <AppBackground>
                <Box sx={{
                    maxWidth: "500px", p: {md: 3, xs: 1}
                }}>
                    <Story/>
                </Box>
            </AppBackground>
        ),
    ],
    args: {
        onShare: fn(),
    },
} satisfies Meta<typeof AlbumCard>;

export default meta;
type Story = StoryObj<typeof meta>;

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
    args: {
        album: clairObscurAlbum,
    },
};

export const Shared: Story = {
    args: {
        album: {
            ...clairObscurAlbum,
            sharedWith: [{user: {name: 'Hulk', email: 'hulk@avenger.com', picture: '/static/hulk-profile.webp'}}],
        }
    }
};

export const ColdTemperature: Story = {
    args: {
        album: {
            ...clairObscurAlbum,
            temperature: 78.0,
            relativeTemperature: 0.1,
        }
    }
};
export const MidTemperature: Story = {
    args: {
        album: {
            ...clairObscurAlbum,
            temperature: 78.0,
            relativeTemperature: 0.5,
        }
    }
};

export const WithoutThumbnail: Story = {
    args: {
        album: {
            ...clairObscurAlbum,
            thumbnails: [],
        }
    }
}

export const WithOneThumbnail: Story = {
    args: {
        album: {
            ...clairObscurAlbum,
            thumbnails: ['/thumbnails/clair-obscur-7.jpg'],
        }
    }
}

export const WithTwoThumbnails: Story = {
    args: {
        album: {
            ...clairObscurAlbum,
            thumbnails: [
                '/thumbnails/clair-obscur-7.jpg',
                '/thumbnails/clair-obscur-8.jpg',
            ],
        }
    }
}

export const WithThreeThumbnails: Story = {
    args: {
        album: {
            ...clairObscurAlbum,
            thumbnails: [
                '/thumbnails/clair-obscur-7.jpg',
                '/thumbnails/clair-obscur-8.jpg',
                '/thumbnails/clair-obscur-6.jpg',
            ],
        }
    }
}

export const WithErroredThumbnails: Story = {
    args: {
        album: {
            ...clairObscurAlbum,
            thumbnails: [
                '/thumbnails/cannot-be-found.jpg',
                '/thumbnails/death-stranding-2-01.jpg',
                '/thumbnails/astro-bot-01.jpg',
                '/thumbnails/cannot-be-found.jpg',
            ],
        }
    }
}

export const SharedByOne: Story = {
    args: {
        album: {
            ...clairObscurAlbum,
            ownedBy: {
                name: 'Natasha',
                users: [
                    {name: 'Natasha', email: 'blackwidow@avenger.com', picture: '/static/black-widow-profile.jpg'},
                ],
            },
        },
    },
}

export const LongName: Story = {
    args: {
        album: {
            ...clairObscurAlbum,
            name: 'Voyage au bout du monde et retour en passant par les étoiles',
        },
    },
};

export const SharedBySeveral: Story = {
    args: {
        album: {
            ...clairObscurAlbum,
            ownedBy: {
                name: 'Avengers',
                users: [
                    {name: 'Hulk', email: 'hulk@avenger.com', picture: '/static/hulk-profile.webp'},
                    {name: 'Natasha', email: 'blackwidow@avenger.com', picture: '/static/black-widow-profile.jpg'},
                    {name: 'Tony Stark', email: 'ironman@avenger.com', picture: '/static/tonystark-profile.jpg'},
                ],
            },
        },
    },
}
