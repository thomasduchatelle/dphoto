import {beforeEach, describe, expect, it} from 'vitest';
import {act, renderHook} from '@testing-library/react';
import {useClientRouter} from './ClientRouter';
import {ReactNode} from 'react';
import {Router} from "waku/router/client";

// Test wrapper component
const TestWrapper = ({children}: { children: ReactNode }) => (
    <Router>{children}</Router>
);

describe('useClientRouter', () => {
    beforeEach(() => {
        // Reset window location before each test
        window.history.pushState({}, '', '/');
    });

    it('should return current path', () => {
        const {result} = renderHook(() => useClientRouter(), {wrapper: TestWrapper});
        expect(result.current.path).toBe('/');
    });

    it('should navigate to new path without reload', () => {
        const {result} = renderHook(() => useClientRouter(), {wrapper: TestWrapper});

        act(() => {
            result.current.navigate('/albums');
        });

        expect(window.location.pathname).toBe('/albums');
        expect(result.current.path).toBe('/albums');
    });

    it('should replace current path without reload', () => {
        const {result} = renderHook(() => useClientRouter(), {wrapper: TestWrapper});

        act(() => {
            result.current.navigate('/albums');
        });

        act(() => {
            result.current.replace('/albums/owner/album');
        });

        expect(window.location.pathname).toBe('/albums/owner/album');
        expect(result.current.path).toBe('/albums/owner/album');
    });

    it('should parse album params from path', () => {
        const {result} = renderHook(() => useClientRouter(), {wrapper: TestWrapper});

        act(() => {
            result.current.navigate('/albums/owner1/album1');
        });

        expect(result.current.params).toEqual({
            owner: 'owner1',
            album: 'album1',
        });
    });

    it('should parse media params from path', () => {
        const {result} = renderHook(() => useClientRouter(), {wrapper: TestWrapper});

        act(() => {
            result.current.navigate('/albums/owner1/album1/encoded123/photo.jpg');
        });

        expect(result.current.params).toEqual({
            owner: 'owner1',
            album: 'album1',
            encodedId: 'encoded123',
            filename: 'photo.jpg',
        });
    });

    it('should parse query parameters', () => {
        const {result} = renderHook(() => useClientRouter(), {wrapper: TestWrapper});

        act(() => {
            result.current.navigate('/albums?filter=recent');
        });

        expect(result.current.query.get('filter')).toBe('recent');
    });

    it('should update when query params change', () => {
        const {result, rerender} = renderHook(() => useClientRouter(), {wrapper: TestWrapper});

        act(() => {
            result.current.navigate('/albums?filter=old');
        });

        expect(result.current.query.get('filter')).toBe('old');

        act(() => {
            result.current.navigate('/albums?filter=new');
        });

        rerender();
        expect(result.current.query.get('filter')).toBe('new');
    });
});
