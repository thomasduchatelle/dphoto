import {NextRequest} from 'next/server';
import {vi} from 'vitest';

let testRequest: NextRequest | undefined;
let mockCookies: Map<string, string>;

// This mock must be defined at module level for vitest to properly intercept the imports
vi.mock('next/headers', () => {
    return {
        cookies: vi.fn(() => Promise.resolve({
            get: vi.fn((key: string) => {
                const value = mockCookies?.get(key) || testRequest?.cookies.get(key)?.value;
                return value ? {value} : undefined;
            }),
            set: vi.fn((key: string, value: string, options?: any) => {
                if (mockCookies) {
                    mockCookies.set(key, value);
                }
            }),
        })),
        headers: vi.fn(() => Promise.resolve({
            get: vi.fn((key: string) => {
                if (key === 'host' && testRequest) {
                    return new URL(testRequest.url).host;
                }
                return testRequest?.headers.get(key) || null;
            }),
        })),
    };
});

export interface FakeNextHeaders {
    /**
     * Set the current request context for the mock
     */
    withRequest(request: NextRequest): void;

    /**
     * Set a cookie value in the mock store
     */
    setCookie(key: string, value: string): void;

    /**
     * Get a cookie value that was set by the code under test
     */
    getSetCookie(key: string): string | undefined;

    /**
     * Reset the mock state
     */
    reset(): void;
}

/**
 * Creates a reusable state manager for Next.js headers() and cookies() mocks.
 * 
 * IMPORTANT: This must be called AFTER the imports so the vi.mock() at module level takes effect.
 * 
 * Usage:
 * ```typescript
 * import {fakeNextHeaders} from '@/__tests__/helpers/fake-next-headers';
 * 
 * const fakeHeaders = fakeNextHeaders();
 * 
 * describe("...", () => {
 *   afterEach(() => fakeHeaders.reset());
 * 
 *   it("...", () => {
 *     const req = new NextRequest(...);
 *     fakeHeaders.withRequest(req);
 *     fakeHeaders.setCookie("cookie-name", "cookie-value");
 * 
 *     // ... test code ...
 * 
 *     expect(fakeHeaders.getSetCookie("cookie-name-2")).toBe("expected-value");
 *   });
 * });
 * ```
 */
export function fakeNextHeaders(): FakeNextHeaders {
    // Initialize on first call
    if (!mockCookies) {
        mockCookies = new Map();
    }

    return {
        withRequest(request: NextRequest): void {
            testRequest = request;
        },

        setCookie(key: string, value: string): void {
            mockCookies.set(key, value);
        },

        getSetCookie(key: string): string | undefined {
            return mockCookies.get(key);
        },

        reset(): void {
            testRequest = undefined;
            mockCookies.clear();
            vi.clearAllMocks();
        },
    };
}
