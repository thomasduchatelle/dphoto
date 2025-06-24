import React from 'react';
import {ComponentMeta, ComponentStory} from '@storybook/react';
import {userEvent, within} from "@storybook/testing-library";
import AlbumsListActions from "../pages/authenticated/albums/AlbumsListActions/AlbumListActions";
import {BrowserRouter} from 'react-router-dom';

// More on default export: https://storybook.js.org/docs/react/writing-stories/introduction#default-export
export default {
    title: 'Albums/AlbumsListActions',
    component: AlbumsListActions,
    argTypes: {
        onAlbumFiltered: {action: 'onAlbumFiltered'},
    },
} as ComponentMeta<typeof AlbumsListActions>;

// More on component templates: https://storybook.js.org/docs/react/writing-stories/introduction#using-args
const Template: ComponentStory<typeof AlbumsListActions> = (args) => (<BrowserRouter><AlbumsListActions {...args}/></BrowserRouter>);

export const NoOptions = Template.bind({});
NoOptions.args = {
    selected: {
        criterion: {
            owners: []
        },
        avatars: [],
        name: "All Albums",
    },
};

export const AllOptionsOpen = Template.bind({});
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
AllOptionsOpen.parameters = {
    delay: 300,
};
AllOptionsOpen.play = async ({canvasElement}) => {
    const canvas = within(canvasElement);

    await userEvent.click(canvas.getAllByRole("button")[0]);
};

export const NoOwnAlbums = Template.bind({});
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

export const OnlyOwnAlbums = Template.bind({});
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

export const NoOwnAlbumsOpen = Template.bind({});
NoOwnAlbumsOpen.args = NoOwnAlbums.args;
NoOwnAlbumsOpen.parameters = {
    delay: 300,
};
NoOwnAlbumsOpen.play = async ({canvasElement}) => {
    const canvas = within(canvasElement);

    await userEvent.click(canvas.getAllByRole("button")[0]);
};

export const EditDatesButtonDisabled = Template.bind({});
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
