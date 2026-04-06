'use client';

import {Button} from '@mui/material';
import {SearchOff as SearchOffIcon} from '@mui/icons-material';
import {PageMessage} from '@/components/PageMessage';

export default function NotFound() {
    return (
        <PageMessage
            icon={<SearchOffIcon/>}
            title="Page Not Found"
            message="The page you're looking for doesn't exist."
        >
            <Button
                variant="contained"
                onClick={() => {
                    window.location.href = '/';
                }}
            >
                Go Home
            </Button>
        </PageMessage>
    );
}
