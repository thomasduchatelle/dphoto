// vitest-dom adds custom matchers for asserting on DOM nodes.
// allows you to do things like:
// expect(element).toHaveTextContent(/react/i)
// learn more: https://github.com/testing-library/jest-dom
import '@testing-library/jest-dom/vitest';
import { vi } from 'vitest';

// Create a shared router state that persists across hook calls
const createRouterState = () => {
    const state = {
        path: '/',
        query: '',
    };
    
    const listeners = new Set<() => void>();
    
    const notify = () => {
        listeners.forEach(listener => listener());
    };
    
    return {
        getPath: () => state.path,
        getQuery: () => state.query,
        push: (path: string) => {
            const [pathname, search] = path.split('?');
            state.path = pathname;
            state.query = search || '';
            if (typeof window !== 'undefined' && window.history) {
                try {
                    window.history.pushState({}, '', path);
                } catch (e) {
                    // Ignore errors in test environment
                }
            }
            notify();
        },
        replace: (path: string) => {
            const [pathname, search] = path.split('?');
            state.path = pathname;
            state.query = search || '';
            if (typeof window !== 'undefined' && window.history) {
                try {
                    window.history.replaceState({}, '', path);
                } catch (e) {
                    // Ignore errors in test environment
                }
            }
            notify();
        },
        subscribe: (listener: () => void) => {
            listeners.add(listener);
            return () => listeners.delete(listener);
        },
        reset: () => {
            state.path = '/';
            state.query = '';
            if (typeof window !== 'undefined' && window.history) {
                try {
                    window.history.replaceState({}, '', '/');
                } catch (e) {
                    // Ignore errors in test environment
                }
            }
        },
    };
};

const routerState = createRouterState();

// Reset router state before each test
beforeEach(() => {
    routerState.reset();
});

// Mock waku router for tests
vi.mock('waku', async () => {
    const actual = await vi.importActual('waku');
    const { useState, useEffect } = await import('react');
    
    return {
        ...actual,
        useRouter: () => {
            const [, forceUpdate] = useState({});
            
            useEffect(() => {
                return routerState.subscribe(() => {
                    forceUpdate({});
                });
            }, []);
            
            return {
                path: routerState.getPath(),
                query: routerState.getQuery(),
                push: routerState.push,
                replace: routerState.replace,
            };
        },
    };
});