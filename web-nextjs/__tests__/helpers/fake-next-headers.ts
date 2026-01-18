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
 * IMPORTANT: Due to vitest's mock hoisting requirements, you CANNOT import and use this
 * helper's vi.mock() directly. Instead, each test file must define its own vi.mock() inline.
 * 
 * This file serves as a TEMPLATE and REFERENCE for how to structure your test mocks.
 * 
 * For test files that need BOTH headers and cookies mocking (like proxy.test.ts):
 * Copy the pattern from this file with module-level variables and vi.mock() call.
 * 
 * For test files that only need cookies mocking (like access-token-service.test.ts):
 * Use a simplified version without the headers mock and testRequest variable.
 * 
 * Usage pattern (copy this structure to your test file):
 * ```typescript
 * // At the top of your test file, before any imports of code under test:
 * import {NextRequest} from 'next/server';
 * 
 * let testRequest: NextRequest | undefined;
 * let mockCookies: Map<string, string>;
 * 
 * vi.mock('next/headers', () => {
 *     return {
 *         cookies: vi.fn(() => Promise.resolve({
 *             get: vi.fn((key: string) => {
 *                 const value = mockCookies?.get(key) || testRequest?.cookies.get(key)?.value;
 *                 return value ? {value} : undefined;
 *             }),
 *             set: vi.fn((key: string, value: string, options?: any) => {
 *                 if (mockCookies) {
 *                     mockCookies.set(key, value);
 *                 }
 *             }),
 *         })),
 *         headers: vi.fn(() => Promise.resolve({
 *             get: vi.fn((key: string) => {
 *                 if (key === 'host' && testRequest) {
 *                     return new URL(testRequest.url).host;
 *                 }
 *                 return testRequest?.headers.get(key) || null;
 *             }),
 *         })),
 *     };
 * });
 * 
 * // In your test setup:
 * describe("...", () => {
 *   beforeEach(() => {
 *     testRequest = undefined;
 *     mockCookies = new Map();
 *     vi.clearAllMocks();
 *   });
 * 
 *   it("...", () => {
 *     const req = new NextRequest(...);
 *     testRequest = req;
 *     mockCookies.set("cookie-name", "cookie-value");
 * 
 *     // ... test code ...
 * 
 *     expect(mockCookies.get("cookie-name-2")).toBe("expected-value");
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
