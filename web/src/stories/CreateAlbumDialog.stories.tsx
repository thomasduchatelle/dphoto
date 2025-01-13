import React from 'react';
import {ComponentMeta, ComponentStory} from '@storybook/react';
import CreateAlbumDialog from "../pages/authenticated/albums/CreateAlbumDialog";

// More on default export: https://storybook.js.org/docs/react/writing-stories/introduction#default-export
export default {
    title: 'Albums/CreateAlbumDialog',
    component: CreateAlbumDialog,
} as ComponentMeta<typeof CreateAlbumDialog>;

// More on component templates: https://storybook.js.org/docs/react/writing-stories/introduction#using-args
const Template: ComponentStory<typeof CreateAlbumDialog> = (args) => (<CreateAlbumDialog {...args}/>
);

export const Empty = Template.bind({});
Empty.args = {
    open: true,
};
Empty.parameters = {
    delay: 300,
};

