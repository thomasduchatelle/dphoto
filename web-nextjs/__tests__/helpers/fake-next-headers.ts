import {NextRequest} from 'next/server';
import {vi} from 'vitest';

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
    getSetCookie(key: string): { value: string, options?: any } | undefined;

    /**
     * Reset the mock state
     */
    reset(): void;

    /**
     * Get the mock implementation of next/headers
     */
    mock(): HeadersAndCookies;
}

export interface HeadersAndCookies {
    cookies: () => Promise<{
        get: (key: string) => { value: string } | undefined;
        set: (key: string, value: string, options?: any) => void;
        getAll: () => Array<{ name: string, value: string, options?: any }>;
    }>;
    headers: () => Promise<{
        get: (key: string) => string | null;
    }>;
}

export interface SetCookieValue {
    value: string
    options?: any
}

class FakeHeader {

    constructor(
        private testRequest: NextRequest | undefined = undefined,
        private readonly requestCookies: Map<string, string> = new Map(),
        private readonly setCookies: Map<string, SetCookieValue> = new Map(),
    ) {
    }

    public withRequest(request: NextRequest): void {
        this.testRequest = request;
    }

    public setCookie(key: string, value: string): void {
        this.requestCookies.set(key, value);
    }

    public getSetCookie(key: string): SetCookieValue | undefined {
        return this.setCookies.get(key);
    }

    public reset(): void {
        this.testRequest = undefined;
        this.requestCookies.clear();
        this.setCookies.clear();
        vi.clearAllMocks();
    }

    public mock(): HeadersAndCookies {
        return {
            cookies: () => Promise.resolve({
                get: (key: string): { value: string } | undefined => {
                    const value = this.requestCookies?.get(key) || this.testRequest?.cookies.get(key)?.value;
                    return value ? {value} : undefined;
                },
                set: (key: string, value: string, options?: any): void => {
                    if (this.setCookies) {
                        this.setCookies.set(key, {value, options});
                    }
                },
                getAll: (): Array<{ name: string, value: string, options?: any }> => {
                    const result: Array<{ name: string, value: string, options?: any }> = [];
                    for (const [name, {value, options}] of this.setCookies) {
                        result.push({name, value, options});
                    }
                    return result;
                },
            }),
            headers: () => Promise.resolve({
                get: (key: string): string | null => {
                    if (key === 'host' && this.testRequest) {
                        return new URL(this.testRequest.url).host;
                    }
                    return this.testRequest?.headers.get(key) || null;
                }
            }),
        }
    }
}

/**
 * Usage:
 *
 * const fakeHeaders = fakeNextHeaders();
 *
 * vi.mock('next/headers', () => {
 *     return {
 *         cookies: vi.fn(() => fakeHeaders.mock().cookies()),
 *         headers: vi.fn(() => fakeHeaders.mock().headers()),
 *     };
 * });
 *
 * // In test:
 * fakeHeaders.withRequest(yourTestRequest);
 * fakeHeaders.setCookie('key', 'value');
 *
 * // After test:
 * fakeHeaders.reset();
 */
export function fakeNextHeaders(): FakeNextHeaders {
    return new FakeHeader()
}
