import {describe, it, expect, beforeEach, vi} from 'vitest';
import {renderHook, act} from '@testing-library/react';
import {useClientRouter, ClientLink, RouterProvider} from './ClientRouter';
import {render, screen, fireEvent} from '@testing-library/react';
import {ReactNode} from 'react';

// Test wrapper component
const TestWrapper = ({children}: {children: ReactNode}) => (
    <RouterProvider>{children}</RouterProvider>
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

describe('ClientLink', () => {
    beforeEach(() => {
        window.history.pushState({}, '', '/');
    });

    it('should navigate on click without reload', () => {
        const mockNavigate = vi.fn();
        
        render(
            <RouterProvider>
                <ClientLink to="/albums">
                    Go to Albums
                </ClientLink>
            </RouterProvider>
        );

        const link = screen.getByText('Go to Albums');
        fireEvent.click(link);

        // Verify the href is set correctly
        expect(link).toHaveAttribute('href', '/albums');
        
        // After click, the URL should be updated
        expect(window.location.pathname).toBe('/albums');
    });

    it('should call onClick handler if provided', () => {
        const mockOnClick = vi.fn();
        
        render(
            <RouterProvider>
                <ClientLink to="/albums" onClick={mockOnClick}>
                    Go to Albums
                </ClientLink>
            </RouterProvider>
        );

        const link = screen.getByText('Go to Albums');
        fireEvent.click(link);

        expect(mockOnClick).toHaveBeenCalled();
    });

    it('should apply className', () => {
        render(
            <RouterProvider>
                <ClientLink to="/albums" className="test-class">
                    Go to Albums
                </ClientLink>
            </RouterProvider>
        );

        const link = screen.getByText('Go to Albums');
        expect(link).toHaveClass('test-class');
    });
});
