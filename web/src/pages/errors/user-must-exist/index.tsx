'use client';

export default function UserMustExist() {
    return (
        <div style={{textAlign: 'center', marginTop: '100px', padding: '20px'}}>
            <h1>Access Required</h1>
            <p>Access must be granted by an administrator before you can use this application.</p>
            <p>Please contact your administrator to request access.</p>
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
