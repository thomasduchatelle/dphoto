import { headers } from 'next/headers';
import { BackendSession } from './constants';

/**
 * Server-side utility to extract BackendSession from request headers.
 * This should only be used in Server Components or Server Actions.
 */
export async function getBackendSession(): Promise<BackendSession | null> {
    const headersList = await headers();
    const backendSessionHeader = headersList.get('x-backend-session');
    
    if (!backendSessionHeader) {
        return null;
    }
    
    try {
        return JSON.parse(backendSessionHeader) as BackendSession;
    } catch (error) {
        console.error('Failed to parse backend session:', error);
        return null;
    }
}
