import type {Meta, StoryObj} from '@storybook/nextjs-vite';
import {fn} from 'storybook/test';
import {AlbumCard} from './index';
import {Box} from '@mui/material';
import {Album, AlbumCover, AlbumId} from '@/domains/catalog/language/catalog-state';
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

const cover = (mediaId: string, contentPath: string): AlbumCover => ({
    mediaId,
    filename: contentPath.split(/(\\|\/)/g).pop(),
    origin: 'RANDOM',
    contentPath: contentPath,
});

const clairObscurAlbum: Album = {
    albumId: createAlbumId('sandfall', 'clair-obscur'),
    name: 'Clair Obscur',
    start: new Date('2025-04-24'),
    end: new Date('2025-06-01'),
    totalCount: 47,
    temperature: 6.7,
    relativeTemperature: 1,
    sharedWith: [],
    covers: [
        cover('m1', '/thumbnails/clair-obscur-1.jpg'),
        cover('m2', '/thumbnails/clair-obscur-2.jpg'),
        cover('m3', '/thumbnails/clair-obscur-3.jpg'),
        cover('m4', '/thumbnails/clair-obscur-4.jpg'),
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

export const WithoutCovers: Story = {
    args: {
        album: {
            ...clairObscurAlbum,
            covers: [],
        }
    }
}

export const WithOneCover: Story = {
    args: {
        album: {
            ...clairObscurAlbum,
            covers: [cover('m7', '/thumbnails/clair-obscur-7.jpg')],
        }
    }
}

export const WithTwoCovers: Story = {
    args: {
        album: {
            ...clairObscurAlbum,
            covers: [
                cover('m7', '/thumbnails/clair-obscur-7.jpg'),
                cover('m8', '/thumbnails/clair-obscur-8.jpg'),
            ],
        }
    }
}

export const WithThreeCovers: Story = {
    args: {
        album: {
            ...clairObscurAlbum,
            covers: [
                cover('m7', '/thumbnails/clair-obscur-7.jpg'),
                cover('m8', '/thumbnails/clair-obscur-8.jpg'),
                cover('m6', '/thumbnails/clair-obscur-6.jpg'),
            ],
        }
    }
}

export const WithErroredCovers: Story = {
    args: {
        album: {
            ...clairObscurAlbum,
            covers: [
                cover('missing1', '/thumbnails/cannot-be-found.jpg'),
                cover('m-ds2', '/thumbnails/death-stranding-2-01.jpg'),
                cover('m-astro', '/thumbnails/astro-bot-01.jpg'),
                cover('missing2', '/thumbnails/cannot-be-found.jpg'),
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
