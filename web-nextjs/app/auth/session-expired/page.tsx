import {Button} from '@mui/material';
import LockClockOutlinedIcon from '@mui/icons-material/LockClockOutlined';
import {PageMessage} from '@/components/PageMessage';
import Link from '@/components/Link';

export default function SessionExpiredPage() {
    return (
        <PageMessage
            variant="info"
            icon={<LockClockOutlinedIcon/>}
            title="Session Expired"
            message="For your security, your session has ended. Sign in again to continue."
        >
            <Button
                component={Link}
                href="/"
                prefetch={false}
                variant="contained"
            >
                Sign In Again
            </Button>
        </PageMessage>
    );
}
