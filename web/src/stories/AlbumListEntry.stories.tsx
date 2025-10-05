import React from 'react';
import {Story} from '@ladle/react';
import {List} from "@mui/material";
import {AlbumListEntry} from "../components/albums/AlbumsList/AlbumListEntry";
import {StoriesContext} from "./StoriesContext";

export default {
    title: 'Albums / AlbumListEntry',
};

type Props = React.ComponentProps<typeof AlbumListEntry>;

const AlbumListEntryWrapper: Story<Partial<Props>> = (args) => (
    <StoriesContext maxWidth={450}><List><AlbumListEntry {...args as Props} /></List></StoriesContext>
);

export const Default = (args: Props) => <AlbumListEntryWrapper {...args} />
Default.args = {
    album: {
        albumId: {owner: "tony@stark.com", folderName: "2010_Avenger"},
        name: "Avenger 2010",
        start: new Date(2023, 3, 22, 8, 41, 0),
        end: new Date(2023, 4, 23, 8, 41, 0),
        totalCount: 214,
        temperature: 25,
        relativeTemperature: 0.6,
        sharedWith: [],
    },
    selected: false,
};

export const Selected = (args: Props) => <AlbumListEntryWrapper {...args} />
Selected.args = {
    album: {
        albumId: {owner: "tony@stark.com", folderName: "2010_Avenger"},
        name: "Avenger 2010",
        start: new Date(2023, 3, 22, 8, 41, 0),
        end: new Date(2023, 4, 23, 8, 41, 0),
        totalCount: 214,
        temperature: 25,
        relativeTemperature: 0.6,
        sharedWith: [],
    },
    selected: true,
};

export const SharedBySomeoneElse = (args: Props) => <AlbumListEntryWrapper {...args} />
SharedBySomeoneElse.args = {
    album: {
        albumId: {owner: "tony@stark.com", folderName: "2010_Avenger"},
        name: "Avenger 2010",
        start: new Date(2023, 3, 22, 8, 41, 0),
        end: new Date(2023, 4, 23, 8, 41, 0),
        totalCount: 214,
        temperature: 25,
        relativeTemperature: 0.6,
        ownedBy: {
            name: "Stark friends",
            users: [
                {name: "Black Widow", email: "blckwidow@avenger.com", picture: "black-widow-profile.jpg"},
                {name: "Hulk", email: "hulk@avenger.com", picture: "hulk-profile.webp"},
            ]
        },
        sharedWith: [],
    },
    selected: false,
};

export const SharedToOthers = (args: Props) => <AlbumListEntryWrapper {...args} />
SharedToOthers.args = {
    album: {
        albumId: {owner: "tony@stark.com", folderName: "2010_Avenger"},
        name: "Avenger 2010",
        start: new Date(2023, 3, 22, 8, 41, 0),
        end: new Date(2023, 4, 23, 8, 41, 0),
        totalCount: 214,
        temperature: 25,
        relativeTemperature: 0.6,
        sharedWith: [
            {user: {name: "Pepper Stark", email: "pepper@stark.com"}},
        ],
    },
    selected: false,
};
