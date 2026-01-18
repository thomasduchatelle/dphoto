import Link from 'next/link';
import {basePath} from '@/libs/requests';

interface ErrorPageProps {
    searchParams: Promise<{
        error?: string;
        error_description?: string;
    }>;
}

interface ErrorInfo {
    title: string;
    description: string;
}

function getErrorInfo(error: string, errorDescription?: string): ErrorInfo {
    const errorMap: Record<string, ErrorInfo> = {
        'invalid_request': {
            title: 'Invalid Request',
            description: errorDescription || 'The authentication request was invalid. Please try again.',
        },
        'unauthorized_client': {
            title: 'Unauthorized Client',
            description: 'The application is not authorized to perform this request.',
        },
        'access_denied': {
            title: 'Access Denied',
            description: 'You have denied access to your account. Authentication cannot proceed without your permission.',
        },
        'unsupported_response_type': {
            title: 'Unsupported Response Type',
            description: 'The requested response type is not supported.',
        },
        'invalid_scope': {
            title: 'Invalid Scope',
            description: 'The requested scope is invalid or not supported.',
        },
        'server_error': {
            title: 'Server Error',
            description: 'An error occurred on the authentication server. Please try again later.',
        },
        'temporarily_unavailable': {
            title: 'Service Unavailable',
            description: 'The authentication service is temporarily unavailable. Please try again later.',
        },
        'state-mismatch': {
            title: 'State Mismatch',
            description: 'The authentication state does not match. This could indicate a security issue or an expired session. Please try logging in again.',
        },
        'missing-authentication-cookies': {
            title: 'Missing Authentication Cookies',
            description: 'Required authentication cookies are missing. This may happen if cookies expired or were deleted. Please try logging in again.',
        },
        'token-exchange-failed': {
            title: 'Token Exchange Failed',
            description: 'Failed to exchange the authorization code for tokens. Please try logging in again.',
        },
    };

    return errorMap[error] || {
        title: 'Authentication Error',
        description: 'An unexpected error occurred during authentication. Please try logging in again.',
    };
}

export default async function ErrorPage({searchParams}: ErrorPageProps) {
    const params = await searchParams;
    const error = params.error || 'unknown';
    const errorDescription = params.error_description;
    const errorInfo = getErrorInfo(error, errorDescription);

    return (
        <div className="flex min-h-screen items-center justify-center bg-zinc-50 font-sans dark:bg-black">
            <main className="flex w-full max-w-md flex-col items-center gap-8 rounded-lg bg-white p-8 shadow-lg dark:bg-zinc-900">
                <div className="flex h-16 w-16 items-center justify-center rounded-full bg-red-100 dark:bg-red-900">
                    <svg
                        className="h-8 w-8 text-red-600 dark:text-red-400"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                        />
                    </svg>
                </div>

                <div className="text-center">
                    <h1 className="text-2xl font-semibold text-zinc-900 dark:text-zinc-50">
                        {errorInfo.title}
                    </h1>
                    <p className="mt-4 text-zinc-600 dark:text-zinc-400">
                        {errorInfo.description}
                    </p>
                </div>

                <Link
                    href={`${basePath}/auth/login`}
                    prefetch={false}
                    className="flex h-12 w-full items-center justify-center rounded-full bg-zinc-900 px-5 text-white transition-colors hover:bg-zinc-700 dark:bg-zinc-100 dark:text-zinc-900 dark:hover:bg-zinc-300"
                >
                    Try Again
                </Link>
            </main>
        </div>
    );
}
