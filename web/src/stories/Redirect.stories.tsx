import React from 'react';
import { Story } from '@ladle/react';
import { Redirect } from '../components/Redirect';

export default {
  title: 'Components/Redirect',
};

type RedirectProps = React.ComponentProps<typeof Redirect>;

export const ToExternalURL: Story<RedirectProps> = (args) => <Redirect {...args} />;
ToExternalURL.args = {
  url: 'https://example.com',
};

export const ToRelativeURL: Story<RedirectProps> = (args) => <Redirect {...args} />;
ToRelativeURL.args = {
  url: '/albums',
};

export const ToURLWithQueryParams: Story<RedirectProps> = (args) => <Redirect {...args} />;
ToURLWithQueryParams.args = {
  url: '/albums?filter=recent&sort=date',
};

export const ToAlbumPage: Story<RedirectProps> = (args) => <Redirect {...args} />;
ToAlbumPage.args = {
  url: '/albums/owner1/album1',
};
