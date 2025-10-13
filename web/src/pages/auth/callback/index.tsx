'use client';

import {useEffect, useState} from 'react';

export default function AuthCallback() {
    const [status, setStatus] = useState<'processing' | 'success' | 'error'>('processing');
    const [error, setError] = useState<string>('');

    useEffect(() => {
        const urlParams = new URLSearchParams(window.location.search);
        const code = urlParams.get('code');
        const state = urlParams.get('state');
        const error = urlParams.get('error');

        if (error) {
            setStatus('error');
            setError(error);
            return;
        }

        if (!code || !state) {
            setStatus('error');
            setError('Missing authorization code or state');
            return;
        }

        // The actual token exchange should happen server-side
        // For now, just redirect to home
        setStatus('success');
        setTimeout(() => {
            window.location.href = '/';
        }, 1000);
    }, []);

    if (status === 'processing') {
        return (
            <div style={{textAlign: 'center', marginTop: '100px'}}>
                <h2>Completing authentication...</h2>
                <p>Please wait while we complete your login.</p>
            </div>
        );
    }

    if (status === 'error') {
        return (
            <div style={{textAlign: 'center', marginTop: '100px'}}>
                <h2>Authentication Error</h2>
                <p>{error || 'An error occurred during authentication.'}</p>
                <button onClick={() => window.location.href = '/'}>
                    Return to Home
                </button>
            </div>
        );
    }

    return (
        <div style={{textAlign: 'center', marginTop: '100px'}}>
            <h2>Authentication Successful</h2>
            <p>Redirecting...</p>
        </div>
    );
}
