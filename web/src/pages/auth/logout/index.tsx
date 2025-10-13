'use client';

import {useEffect} from 'react';

export default function AuthLogout() {
    useEffect(() => {
        // Clear any client-side state
        localStorage.clear();
        sessionStorage.clear();

        // Redirect to home after logout
        setTimeout(() => {
            window.location.href = '/';
        }, 1000);
    }, []);

    return (
        <div style={{textAlign: 'center', marginTop: '100px'}}>
            <h2>Logging out...</h2>
            <p>You have been logged out successfully.</p>
        </div>
    );
}
