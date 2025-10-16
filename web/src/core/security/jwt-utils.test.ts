import {decodeJWTPayload, isOwnerFromJWT} from './jwt-utils';

describe('jwt-utils', () => {
    describe('decodeJWTPayload', () => {
        it('should decode a valid JWT payload', () => {
            // This is a valid JWT with payload: {"sub":"test@example.com","Scopes":"owner:testuser","iss":"dphoto","aud":["dphoto"],"exp":1234567890,"iat":1234567890,"jti":"test-id"}
            const token = 'eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0QGV4YW1wbGUuY29tIiwiU2NvcGVzIjoib3duZXI6dGVzdHVzZXIiLCJpc3MiOiJkcGhvdG8iLCJhdWQiOlsiZHBob3RvIl0sImV4cCI6MTIzNDU2Nzg5MCwiaWF0IjoxMjM0NTY3ODkwLCJqdGkiOiJ0ZXN0LWlkIn0.signature';
            
            const payload = decodeJWTPayload(token);
            
            expect(payload).not.toBeNull();
            expect(payload?.sub).toBe('test@example.com');
            expect(payload?.Scopes).toBe('owner:testuser');
        });

        it('should return null for invalid JWT', () => {
            const payload = decodeJWTPayload('invalid-token');
            
            expect(payload).toBeNull();
        });

        it('should return null for JWT with wrong number of parts', () => {
            const payload = decodeJWTPayload('header.payload');
            
            expect(payload).toBeNull();
        });
    });

    describe('isOwnerFromJWT', () => {
        it('should return true for owner scope', () => {
            // JWT with payload: {"sub":"owner@example.com","Scopes":"owner:testowner","iss":"dphoto","aud":["dphoto"]}
            const token = 'eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJvd25lckBleGFtcGxlLmNvbSIsIlNjb3BlcyI6Im93bmVyOnRlc3Rvd25lciIsImlzcyI6ImRwaG90byIsImF1ZCI6WyJkcGhvdG8iXX0.signature';
            
            expect(isOwnerFromJWT(token)).toBe(true);
        });

        it('should return true for multiple scopes including owner', () => {
            // JWT with payload: {"sub":"owner@example.com","Scopes":"api:admin owner:testowner","iss":"dphoto","aud":["dphoto"]}
            const token = 'eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJvd25lckBleGFtcGxlLmNvbSIsIlNjb3BlcyI6ImFwaTphZG1pbiBvd25lcjp0ZXN0b3duZXIiLCJpc3MiOiJkcGhvdG8iLCJhdWQiOlsiZHBob3RvIl19.signature';
            
            expect(isOwnerFromJWT(token)).toBe(true);
        });

        it('should return false for visitor scope', () => {
            // JWT with payload: {"sub":"visitor@example.com","Scopes":"visitor","iss":"dphoto","aud":["dphoto"]}
            const token = 'eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ2aXNpdG9yQGV4YW1wbGUuY29tIiwiU2NvcGVzIjoidmlzaXRvciIsImlzcyI6ImRwaG90byIsImF1ZCI6WyJkcGhvdG8iXX0.signature';
            
            expect(isOwnerFromJWT(token)).toBe(false);
        });

        it('should return false for invalid token', () => {
            expect(isOwnerFromJWT('invalid-token')).toBe(false);
        });

        it('should return false for token with no scopes', () => {
            // JWT with payload: {"sub":"test@example.com","iss":"dphoto","aud":["dphoto"]}
            const token = 'eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0QGV4YW1wbGUuY29tIiwiaXNzIjoiZHBob3RvIiwiYXVkIjpbImRwaG90byJdfQ.signature';
            
            expect(isOwnerFromJWT(token)).toBe(false);
        });
    });
});
