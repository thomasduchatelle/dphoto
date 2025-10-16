import React from 'react';
import {action, Story} from '@ladle/react';
import {BrowserRouter} from 'react-router-dom';
import AlbumsListActions from "../components/albums/AlbumsListActions/AlbumListActions";

export default {
    title: 'Albums / AlbumsListActions',
};

type Props = React.ComponentProps<typeof AlbumsListActions>;

const AlbumsListActionsWrapper: Story<Partial<Props>> = (args) => (
    <BrowserRouter><AlbumsListActions
        {...args as Props}
        onAlbumFiltered={action('onAlbumFiltered')}/></BrowserRouter>
);

export const NoOptions = (args: Props) => <AlbumsListActionsWrapper {...args} />
NoOptions.args = {
    albumFilter: {
        criterion: {
            owners: []
        },
        avatars: [],
        name: "All Albums",
    },
    displayedAlbumIdIsOwned: true,
    hasAlbumsToDelete: true,
    canCreateAlbums: true,
};

export const AllOptionsOpen = (args: Props) => <AlbumsListActionsWrapper {...args} />
AllOptionsOpen.args = {
    albumFilter: {
        criterion: {
            owners: ['black-widow', 'hulk']
        },
        avatars: ['black-widow-profile.jpg', 'hulk-profile.webp'],
        name: "Avenger Family",
    },
    albumFilterOptions: [
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
    displayedAlbumIdIsOwned: true,
    hasAlbumsToDelete: true,
    canCreateAlbums: true,
};

export const NoOwnAlbums = (args: Props) => <AlbumsListActionsWrapper {...args} />
NoOwnAlbums.args = {
    albumFilter: {
        criterion: {
            owners: ['black-widow', 'hulk']
        },
        avatars: ['black-widow-profile.jpg', 'hulk-profile.webp'],
        name: "Avenger Family",
    },
    albumFilterOptions: [
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
    displayedAlbumIdIsOwned: true,
    hasAlbumsToDelete: true,
    canCreateAlbums: true,
};

export const OnlyOwnAlbums = (args: Props) => <AlbumsListActionsWrapper {...args} />
OnlyOwnAlbums.args = {
    albumFilter: {
        criterion: {
            owners: []
        },
        avatars: ['tonystark-profile.jpg'],
        name: "All Albums",
    },
    albumFilterOptions: [
        {
            criterion: {
                owners: []
            },
            avatars: ['tonystark-profile.jpg'],
            name: "All Albums",
        },
    ],
    displayedAlbumIdIsOwned: true,
    hasAlbumsToDelete: true,
    canCreateAlbums: true,
};

export const EditDatesButtonDisabled = (args: Props) => <AlbumsListActionsWrapper {...args} />
EditDatesButtonDisabled.args = {
    albumFilter: {
        criterion: {
            owners: ['black-widow', 'hulk']
        },
        avatars: ['black-widow-profile.jpg', 'hulk-profile.webp'],
        name: "Avenger Family",
    },
    albumFilterOptions: [
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
    hasAlbumsToDelete: true,
    canCreateAlbums: true,
    displayedAlbumIdIsOwned: false,
};

export const DeleteButtonDisabled = (args: Props) => <AlbumsListActionsWrapper {...args} />
DeleteButtonDisabled.args = {
    albumFilter: {
        criterion: {
            owners: ['black-widow', 'hulk']
        },
        avatars: ['black-widow-profile.jpg', 'hulk-profile.webp'],
        name: "Avenger Family",
    },
    albumFilterOptions: [
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
    displayedAlbumIdIsOwned: true,
    hasAlbumsToDelete: false,
    canCreateAlbums: true,
};

export const CreateButtonDisabled = (args: Props) => <AlbumsListActionsWrapper {...args} />
CreateButtonDisabled.args = {
    albumFilter: {
        criterion: {
            owners: ['black-widow', 'hulk']
        },
        avatars: ['black-widow-profile.jpg', 'hulk-profile.webp'],
        name: "Avenger Family",
    },
    albumFilterOptions: [
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
    displayedAlbumIdIsOwned: true,
    hasAlbumsToDelete: true,
    canCreateAlbums: false,
};
