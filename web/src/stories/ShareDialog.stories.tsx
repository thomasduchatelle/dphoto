import React from "react";
import { ComponentMeta, ComponentStory } from "@storybook/react";
import ShareDialog from "../pages/authenticated/albums/ShareDialog";
import { SharingType } from "../core/catalog";

// More on default export: https://storybook.js.org/docs/react/writing-stories/introduction#default-export
export default {
  title: "Albums/ShareDialog",
  component: ShareDialog,
} as ComponentMeta<typeof ShareDialog>;

// More on component templates: https://storybook.js.org/docs/react/writing-stories/introduction#using-args
const Template: ComponentStory<typeof ShareDialog> = (args) => (
  <ShareDialog {...args} />
);

export const Empty = Template.bind({});
Empty.args = {
  open: true,
  sharedWith: [],
};
Empty.parameters = {
  delay: 300,
};

export const WithShares = Template.bind({});
WithShares.args = {
  open: true,
  sharedWith: [
    {
      user: { name: "Morgan Stark", email: "morgan@stark.com" },
      role: SharingType.contributor,
    },
    {
      user: {
        name: "Natasha",
        email: "blackwidow@avenger.com",
        picture: "black-widow-profile.jpg",
      },
      role: SharingType.visitor,
    },
  ],
};
WithShares.parameters = {
  delay: 300,
};

export const WithGenericError = Template.bind({});
WithGenericError.args = {
  open: true,
  error: {
    type: "general",
    message: "A generic issue occurred while doing something",
  },
  sharedWith: [
    {
      user: {
        name: "Natasha",
        email: "blackwidow@avenger.com",
        picture: "black-widow-profile.jpg",
      },
      role: SharingType.visitor,
    },
  ],
};
WithGenericError.parameters = {
  delay: 300,
};

export const WithSaveError = Template.bind({});
WithSaveError.args = {
  open: true,
  error: {
    type: "adding",
    message: "A specific issue occurred while setting the role",
  },
  sharedWith: [
    {
      user: {
        name: "Natasha",
        email: "blackwidow@avenger.com",
        picture: "black-widow-profile.jpg",
      },
      role: SharingType.visitor,
    },
  ],
};
WithSaveError.parameters = {
  delay: 300,
};
