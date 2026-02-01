import Link from '@/components/Link';
import {Box, Button, Paper, Typography} from '@mui/material';
import ErrorOutlineIcon from '@mui/icons-material/ErrorOutline';

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
        <Box
            sx={{
                display: 'flex',
                minHeight: '100vh',
                alignItems: 'center',
                justifyContent: 'center',
            }}
        >
            <Paper
                component="main"
                sx={{
                    display: 'flex',
                    width: '100%',
                    maxWidth: '28rem',
                    flexDirection: 'column',
                    alignItems: 'center',
                    gap: 4,
                    p: 4,
                }}
            >
                <Box
                    sx={{
                        display: 'flex',
                        height: '4rem',
                        width: '4rem',
                        alignItems: 'center',
                        justifyContent: 'center',
                        borderRadius: '50%',
                        bgcolor: 'error.dark',
                    }}
                >
                    <ErrorOutlineIcon
                        sx={{
                            fontSize: '2rem',
                            color: 'error.light',
                        }}
                    />
                </Box>

                <Box sx={{textAlign: 'center'}}>
                    <Typography
                        variant="h4"
                        sx={{
                            fontWeight: 600,
                        }}
                    >
                        {errorInfo.title}
                    </Typography>
                    <Typography
                        variant="body1"
                        sx={{
                            mt: 2,
                            color: 'text.secondary',
                        }}
                    >
                        {errorInfo.description}
                    </Typography>
                </Box>

                    <Button
                        href="/"
                        component={Link}
                        prefetch={false}
                        variant="contained"
                        fullWidth
                        sx={{
                            height: '48px',
                            borderRadius: '24px',
                        }}
                    >
                        Try Again
                    </Button>
            </Paper>
        </Box>
    );
}
