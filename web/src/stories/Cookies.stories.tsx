import React from 'react';
import { Story } from '@ladle/react';
import { Cookies } from '../components/Cookies';

export default {
  title: 'Components/Cookies',
};

type CookiesProps = React.ComponentProps<typeof Cookies>;

export const SingleCookie: Story<CookiesProps> = (args) => (
  <html>
    <head>
      <Cookies {...args} />
    </head>
    <body>
      <p>View the HTML source to see the Set-Cookie meta tags</p>
    </body>
  </html>
);
SingleCookie.args = {
  cookies: {
    session_id: 'abc123xyz',
  },
};

export const MultipleCookies: Story<CookiesProps> = (args) => (
  <html>
    <head>
      <Cookies {...args} />
    </head>
    <body>
      <p>View the HTML source to see the Set-Cookie meta tags</p>
    </body>
  </html>
);
MultipleCookies.args = {
  cookies: {
    access_token: 'foobar',
    refresh_token: 'baz',
  },
};

export const AuthenticationTokens: Story<CookiesProps> = (args) => (
  <html>
    <head>
      <Cookies {...args} />
    </head>
    <body>
      <p>View the HTML source to see the Set-Cookie meta tags</p>
    </body>
  </html>
);
AuthenticationTokens.args = {
  cookies: {
    access_token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9',
    refresh_token: 'rt_abc123xyz789',
    user_id: 'user-12345',
  },
};

export const NoCookies: Story<CookiesProps> = (args) => (
  <html>
    <head>
      <Cookies {...args} />
    </head>
    <body>
      <p>No cookies to set</p>
    </body>
  </html>
);
NoCookies.args = {
  cookies: {},
};
