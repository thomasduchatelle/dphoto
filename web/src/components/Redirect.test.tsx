import { describe, it, expect } from 'vitest';
import { Redirect } from './Redirect';
import { renderToStaticMarkup } from 'react-dom/server';

describe('Redirect', () => {
  it('should render meta refresh tag with the provided URL', () => {
    const testUrl = 'https://example.com/path';
    const html = renderToStaticMarkup(<Redirect url={testUrl} />);
    
    expect(html).toContain('<meta http-equiv="refresh" content="0;url=https://example.com/path"/>');
  });

  it('should render redirecting message', () => {
    const html = renderToStaticMarkup(<Redirect url="https://example.com" />);
    
    expect(html).toContain('Redirecting...');
  });

  it('should handle URLs with query parameters', () => {
    const testUrl = 'https://example.com/path?param=value&other=test';
    const html = renderToStaticMarkup(<Redirect url={testUrl} />);
    
    expect(html).toContain('content="0;url=https://example.com/path?param=value&amp;other=test"');
  });

  it('should handle relative URLs', () => {
    const testUrl = '/albums/owner/album';
    const html = renderToStaticMarkup(<Redirect url={testUrl} />);
    
    expect(html).toContain('content="0;url=/albums/owner/album"');
  });
});
