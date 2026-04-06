'use client';

import {Box, Button} from '@mui/material';
import {PageMessage} from '@/components/PageMessage';
import {CollapsibleDetails} from './CollapsibleDetails';
import ErrorOutlineIcon from '@mui/icons-material/ErrorOutline';
import Link from '@/components/Link';

export interface ErrorMessageProps {
    error: Error & { digest?: string };
    title?: string;
}

export const ErrorMessage = ({error, title = 'Something Went Wrong'}: ErrorMessageProps) => {
    return (
        <PageMessage
            variant="error"
            icon={<ErrorOutlineIcon/>}
            title={title}
            message={error.message || 'An unexpected error occurred'}
        >
            {error.stack && (
                <Box sx={{width: '100%'}}>
                    <CollapsibleDetails details={error.stack}/>
                </Box>
            )}
            <Button component={Link} href="/" variant="text" prefetch={false}>
                Home
            </Button>
            <Button variant="text" onClick={() => window.location.reload()}>
                Refresh
            </Button>
        </PageMessage>
    );
};
