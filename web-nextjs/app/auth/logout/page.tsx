import {Button} from '@mui/material';
import CheckCircleOutlineIcon from '@mui/icons-material/CheckCircleOutline';
import {PageMessage} from '@/components/PageMessage';
import Link from '@/components/Link';

export default function LogoutPage() {
    return (
        <PageMessage
            variant="success"
            icon={<CheckCircleOutlineIcon/>}
            title="Successfully Logged Out"
            message="You have been successfully logged out of your account."
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
