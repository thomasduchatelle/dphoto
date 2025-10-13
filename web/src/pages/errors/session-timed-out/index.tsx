'use client';

export default function SessionTimedOut() {
    return (
        <div style={{textAlign: 'center', marginTop: '100px', padding: '20px'}}>
            <h1>Session Timed Out</h1>
            <p>Your session has timed out. Please log in again.</p>
            <button 
                onClick={() => window.location.href = '/'}
                style={{
                    marginTop: '20px',
                    padding: '10px 20px',
                    fontSize: '16px',
                    cursor: 'pointer'
                }}
            >
                Log In Again
            </button>
        </div>
    );
}
