import { describe, it, expect } from 'vitest';
import { Cookies } from './Cookies';
import { renderToStaticMarkup } from 'react-dom/server';

describe('Cookies', () => {
  it('should render Set-Cookie meta tags for each cookie', () => {
    const cookies = {
      access_token: 'foobar',
      refresh_token: 'baz',
    };
    
    const html = renderToStaticMarkup(<Cookies cookies={cookies} />);
    
    expect(html).toContain('<meta http-equiv="Set-Cookie"');
    expect(html).toContain('access_token=foobar; Path=/');
    expect(html).toContain('refresh_token=baz; Path=/');
  });

  it('should set correct cookie content for access_token', () => {
    const cookies = {
      access_token: 'foobar',
    };
    
    const html = renderToStaticMarkup(<Cookies cookies={cookies} />);
    
    expect(html).toContain('<meta http-equiv="Set-Cookie" content="access_token=foobar; Path=/"/>');
  });

  it('should set correct cookie content for multiple cookies', () => {
    const cookies = {
      access_token: 'foobar',
      refresh_token: 'baz',
    };
    
    const html = renderToStaticMarkup(<Cookies cookies={cookies} />);
    
    expect(html).toContain('access_token=foobar; Path=/');
    expect(html).toContain('refresh_token=baz; Path=/');
  });

  it('should handle empty cookies object', () => {
    const html = renderToStaticMarkup(<Cookies cookies={{}} />);
    
    expect(html).toBe('');
  });

  it('should handle cookies with special characters', () => {
    const cookies = {
      session: 'abc123-xyz789',
    };
    
    const html = renderToStaticMarkup(<Cookies cookies={cookies} />);
    
    expect(html).toContain('<meta http-equiv="Set-Cookie" content="session=abc123-xyz789; Path=/"/>');
  });
});
