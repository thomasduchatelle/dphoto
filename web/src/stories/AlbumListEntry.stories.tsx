import React from "react";
import { ComponentMeta, ComponentStory } from "@storybook/react";
import { List } from "@mui/material";
import { AlbumListEntry } from "../pages/authenticated/albums/AlbumsList/AlbumListEntry";
import { StoriesContext } from "./StoriesContext";
import { SharingType } from "../core/catalog";

// More on default export: https://storybook.js.org/docs/react/writing-stories/introduction#default-export
export default {
  title: "Albums/AlbumListEntry",
  component: AlbumListEntry,
} as ComponentMeta<typeof AlbumListEntry>;

// More on component templates: https://storybook.js.org/docs/react/writing-stories/introduction#using-args
const Template: ComponentStory<typeof AlbumListEntry> = (args) => (
  <StoriesContext maxWidth={450}>
    <List>
      <AlbumListEntry {...args} />
    </List>
  </StoriesContext>
);

export const Default = Template.bind({});
Default.args = {
  album: {
    albumId: { owner: "tony@stark.com", folderName: "2010_Avenger" },
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

export const Selected = Template.bind({});
Selected.args = {
  album: {
    albumId: { owner: "tony@stark.com", folderName: "2010_Avenger" },
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

export const SharedBySomeoneElse = Template.bind({});
SharedBySomeoneElse.args = {
  album: {
    albumId: { owner: "tony@stark.com", folderName: "2010_Avenger" },
    name: "Avenger 2010",
    start: new Date(2023, 3, 22, 8, 41, 0),
    end: new Date(2023, 4, 23, 8, 41, 0),
    totalCount: 214,
    temperature: 25,
    relativeTemperature: 0.6,
    ownedBy: {
      name: "Stark friends",
      users: [
        {
          name: "Black Widow",
          email: "blckwidow@avenger.com",
          picture: "black-widow-profile.jpg",
        },
        {
          name: "Hulk",
          email: "hulk@avenger.com",
          picture: "hulk-profile.webp",
        },
      ],
    },
    sharedWith: [],
  },
  selected: false,
};

export const SharedToOthers = Template.bind({});
SharedToOthers.args = {
  album: {
    albumId: { owner: "tony@stark.com", folderName: "2010_Avenger" },
    name: "Avenger 2010",
    start: new Date(2023, 3, 22, 8, 41, 0),
    end: new Date(2023, 4, 23, 8, 41, 0),
    totalCount: 214,
    temperature: 25,
    relativeTemperature: 0.6,
    sharedWith: [
      {
        user: { name: "Pepper Stark", email: "pepper@stark.com" },
        role: SharingType.visitor,
      },
    ],
  },
  selected: false,
};
