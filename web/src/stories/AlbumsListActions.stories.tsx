import React from 'react';
import {action, Story} from '@ladle/react';
import {BrowserRouter} from 'react-router-dom';
import AlbumsListActions from "../pages/albums/AlbumsListActions/AlbumListActions";

export default {
    title: 'Albums / AlbumsListActions',
};

type Props = React.ComponentProps<typeof AlbumsListActions>;

const AlbumsListActionsWrapper: Story<Partial<Props>> = (args) => (
    <BrowserRouter><AlbumsListActions {...args as Props} onAlbumFiltered={action('onAlbumFiltered')}/></BrowserRouter>
);

export const NoOptions = (args: Props) => <AlbumsListActionsWrapper {...args} />
NoOptions.args = {
    selected: {
        criterion: {
            owners: []
        },
        avatars: [],
        name: "All Albums",
    },
};

export const AllOptionsOpen = (args: Props) => <AlbumsListActionsWrapper {...args} />
AllOptionsOpen.args = {
    selected: {
        criterion: {
            owners: ['black-widow', 'hulk']
        },
        avatars: ['black-widow-profile.jpg', 'hulk-profile.webp'],
        name: "Avenger Family",
    },
    options: [
        {
            criterion: {
                selfOwned: true,
                owners: [],
            },
            avatars: ['tonystark-profile.jpg'],
            name: "My Albums",
        },
        {
            criterion: {
                owners: []
            },
            avatars: ['black-widow-profile.jpg', 'hulk-profile.webp', 'tonystark-profile.jpg', '4.jpg', '5.jpg'],
            name: "All Albums",
        },
        {
            criterion: {
                owners: ['black-widow', 'hulk']
            },
            avatars: ['black-widow-profile.jpg', 'hulk-profile.webp'],
            name: "Avengers Family",
        },
        {
            criterion: {
                owners: ['black-widow']
            },
            avatars: ['black-widow-profile.jpg'],
            name: "Black Widow",
        },
        {
            criterion: {
                owners: ['hulk']
            },
            avatars: ['hulk-profile.webp'],
            name: "Hulk",
        }
    ],
};

export const NoOwnAlbums = (args: Props) => <AlbumsListActionsWrapper {...args} />
NoOwnAlbums.args = {
    selected: {
        criterion: {
            owners: ['black-widow', 'hulk']
        },
        avatars: ['black-widow-profile.jpg', 'hulk-profile.webp'],
        name: "Avenger Family",
    },
    options: [
        {
            criterion: {
                owners: []
            },
            avatars: ['black-widow-profile.jpg', 'hulk-profile.webp', 'tonystark-profile.jpg', '4.jpg', '5.jpg'],
            name: "All Albums",
        },
        {
            criterion: {
                owners: ['black-widow', 'hulk']
            },
            avatars: ['black-widow-profile.jpg', 'hulk-profile.webp'],
            name: "Avengers Family",
        },
        {
            criterion: {
                owners: ['black-widow']
            },
            avatars: ['black-widow-profile.jpg'],
            name: "Black Widow",
        },
        {
            criterion: {
                owners: ['hulk']
            },
            avatars: ['hulk-profile.webp'],
            name: "Hulk",
        }
    ],
};

export const OnlyOwnAlbums = (args: Props) => <AlbumsListActionsWrapper {...args} />
OnlyOwnAlbums.args = {
    selected: {
        criterion: {
            owners: []
        },
        avatars: ['tonystark-profile.jpg'],
        name: "All Albums",
    },
    options: [
        {
            criterion: {
                owners: []
            },
            avatars: ['tonystark-profile.jpg'],
            name: "All Albums",
        },
    ],
};

export const EditDatesButtonDisabled = (args: Props) => <AlbumsListActionsWrapper {...args} />
EditDatesButtonDisabled.args = {
    selected: {
        criterion: {
            owners: ['black-widow', 'hulk']
        },
        avatars: ['black-widow-profile.jpg', 'hulk-profile.webp'],
        name: "Avenger Family",
    },
    options: [
        {
            criterion: {
                owners: []
            },
            avatars: ['black-widow-profile.jpg', 'hulk-profile.webp', 'tonystark-profile.jpg', '4.jpg', '5.jpg'],
            name: "All Albums",
        },
        {
            criterion: {
                owners: ['black-widow', 'hulk']
            },
            avatars: ['black-widow-profile.jpg', 'hulk-profile.webp'],
            name: "Avengers Family",
        },
    ],
    displayedAlbumIdIsOwned: false,
};
