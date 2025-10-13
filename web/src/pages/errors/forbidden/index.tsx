'use client';

export default function Forbidden() {
    return (
        <div style={{textAlign: 'center', marginTop: '100px', padding: '20px'}}>
            <h1>Access Denied</h1>
            <p>You do not have permission to access this resource.</p>
            <button 
                onClick={() => window.location.href = '/'}
                style={{
                    marginTop: '20px',
                    padding: '10px 20px',
                    fontSize: '16px',
                    cursor: 'pointer'
                }}
            >
                Return to Home
            </button>
        </div>
    );
}
