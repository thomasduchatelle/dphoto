'use client';

import {ErrorMessage} from '@/components/ErrorMessage';

export default function AuthenticatedError({error}: { error: Error & { digest?: string } }) {
    return <ErrorMessage error={error} title="Failed to load the albums"/>;
}
