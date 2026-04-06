'use client';

import {ErrorMessage} from '@/components/ErrorMessage';

export default function Error({error}: { error: Error & { digest?: string } }) {
    return <ErrorMessage error={error}/>;
}
