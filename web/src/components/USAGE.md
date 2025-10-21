# Usage Examples for Redirect and Cookies Components

## Redirect Component

The `Redirect` component uses a meta refresh tag to redirect to a different URL.

### Basic Usage

```tsx
import { Redirect } from './components/Redirect';

// Redirect to an external URL
<Redirect url="https://example.com" />

// Redirect to a relative URL
<Redirect url="/albums" />

// Redirect with query parameters
<Redirect url="/albums?filter=recent" />
```

### Full HTML Example

```tsx
import { Redirect } from './components/Redirect';

export default function RedirectPage() {
  return <Redirect url="/dashboard" />;
}
```

## Cookies Component

The `Cookies` component sets cookie headers using meta tags.

### Basic Usage

```tsx
import { Cookies } from './components/Cookies';

// Set a single cookie
<Cookies cookies={{ session_id: 'abc123' }} />

// Set multiple cookies
<Cookies cookies={{
  access_token: 'foobar',
  refresh_token: 'baz'
}} />
```

### Combined Example (Redirect with Cookies)

This is useful for authentication flows where you need to set cookies and redirect:

```tsx
import { Redirect } from './components/Redirect';
import { Cookies } from './components/Cookies';

export default function AuthCallback() {
  const tokens = {
    access_token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...',
    refresh_token: 'rt_abc123xyz789'
  };

  return (
    <html>
      <head>
        <meta httpEquiv="refresh" content="0;url=/dashboard" />
        <Cookies cookies={tokens} />
      </head>
      <body>
        <p>Authentication successful. Redirecting...</p>
      </body>
    </html>
  );
}
```

Or more concisely using both components:

```tsx
import { Cookies } from './components/Cookies';

export default function AuthCallback() {
  return (
    <html>
      <head>
        <meta httpEquiv="refresh" content="0;url=/dashboard" />
        <Cookies cookies={{
          access_token: 'foobar',
          refresh_token: 'baz'
        }} />
      </head>
      <body>
        <p>Redirecting...</p>
      </body>
    </html>
  );
}
```

## Notes

- The `Redirect` component creates a full HTML document structure
- The `Cookies` component only renders meta tags and should be placed within a `<head>` element
- Cookie paths are automatically set to `/`
- The redirect happens immediately (0 second delay)
