package main

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-jwt/jwt/v5"
)

var (
	cognitoRegion      string
	cognitoUserPoolId  string
	cognitoIssuer      string
	jwksCache          *JWKSCache
	jwksCacheMutex     sync.RWMutex
	jwksCacheExpiresAt time.Time
)

type JWKSCache struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type CustomClaims struct {
	Groups        []string `json:"cognito:groups"`
	Username      string   `json:"cognito:username"`
	Email         string   `json:"email"`
	EmailVerified bool     `json:"email_verified"`
	jwt.RegisteredClaims
}

func init() {
	cognitoRegion = getEnv("COGNITO_REGION", "us-east-1")
	cognitoUserPoolId = getEnv("COGNITO_USER_POOL_ID", "")
	if cognitoUserPoolId == "" {
		panic("COGNITO_USER_POOL_ID environment variable is required")
	}
	cognitoIssuer = fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", cognitoRegion, cognitoUserPoolId)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func Handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
	// Extract token from Authorization header or cookie
	token := extractToken(request)
	if token == "" {
		return denyAllResources("No token found"), nil
	}

	// Validate and parse the token
	claims, err := validateToken(token)
	if err != nil {
		fmt.Printf("Token validation failed: %v\n", err)
		return denyAllResources(fmt.Sprintf("Invalid token: %v", err)), nil
	}

	// Determine required group based on path
	requiredGroup := determineRequiredGroup(request.RouteKey)

	// Check if user has required permissions
	if !hasRequiredPermission(claims.Groups, requiredGroup) {
		fmt.Printf("Access denied: user groups %v, required %s\n", claims.Groups, requiredGroup)
		return denyAllResources("Insufficient permissions"), nil
	}

	// Allow the request and pass user context
	return events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: true,
		Context: map[string]interface{}{
			"userId":   claims.Username,
			"email":    claims.Email,
			"groups":   strings.Join(claims.Groups, ","),
		},
	}, nil
}

func extractToken(request events.APIGatewayV2HTTPRequest) string {
	// Try Authorization header first
	if authHeader := request.Headers["authorization"]; authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
	}

	// Fallback to cookie
	if cookies := request.Cookies; len(cookies) > 0 {
		for _, cookie := range cookies {
			if strings.HasPrefix(cookie, "dphoto-access-token=") {
				return strings.TrimPrefix(cookie, "dphoto-access-token=")
			}
		}
	}

	return ""
}

func validateToken(tokenString string) (*CustomClaims, error) {
	// Parse without verification first to get the kid
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get the kid from token header
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid not found in token header")
		}

		// Get the public key
		publicKey, err := getPublicKey(kid)
		if err != nil {
			return nil, err
		}

		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid claims type")
	}

	// Verify issuer
	if claims.Issuer != cognitoIssuer {
		return nil, fmt.Errorf("invalid issuer: %s", claims.Issuer)
	}

	// Verify token is not expired
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token is expired")
	}

	return claims, nil
}

func getPublicKey(kid string) (*rsa.PublicKey, error) {
	// Get JWKS
	jwks, err := getJWKS()
	if err != nil {
		return nil, err
	}

	// Find the key with matching kid
	for _, key := range jwks.Keys {
		if key.Kid == kid {
			return convertJWKToPublicKey(key)
		}
	}

	return nil, fmt.Errorf("key with kid %s not found", kid)
}

func getJWKS() (*JWKSCache, error) {
	jwksCacheMutex.RLock()
	if jwksCache != nil && time.Now().Before(jwksCacheExpiresAt) {
		defer jwksCacheMutex.RUnlock()
		return jwksCache, nil
	}
	jwksCacheMutex.RUnlock()

	// Fetch JWKS
	jwksCacheMutex.Lock()
	defer jwksCacheMutex.Unlock()

	// Check again in case another goroutine updated it
	if jwksCache != nil && time.Now().Before(jwksCacheExpiresAt) {
		return jwksCache, nil
	}

	jwksURL := fmt.Sprintf("%s/.well-known/jwks.json", cognitoIssuer)
	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch JWKS: status %d", resp.StatusCode)
	}

	var newCache JWKSCache
	if err := json.NewDecoder(resp.Body).Decode(&newCache); err != nil {
		return nil, fmt.Errorf("failed to decode JWKS: %w", err)
	}

	jwksCache = &newCache
	jwksCacheExpiresAt = time.Now().Add(1 * time.Hour)

	return jwksCache, nil
}

func convertJWKToPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	// Decode the modulus
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode modulus: %w", err)
	}

	// Decode the exponent
	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode exponent: %w", err)
	}

	// Convert exponent bytes to int
	var eInt int
	for _, b := range eBytes {
		eInt = eInt<<8 + int(b)
	}

	// Create RSA public key
	publicKey := &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: eInt,
	}

	return publicKey, nil
}

func determineRequiredGroup(routeKey string) string {
	// Parse route key format: "METHOD /path"
	parts := strings.SplitN(routeKey, " ", 2)
	if len(parts) != 2 {
		return "admins" // Default to most restrictive
	}
	path := parts[1]

	// Admin-only endpoints
	if strings.HasPrefix(path, "/api/users") ||
		strings.HasPrefix(path, "/api/admin") {
		return "admins"
	}

	// Owner endpoints (album management, media upload, etc.)
	if strings.HasPrefix(path, "/api/albums") && !strings.Contains(path, "/shared") ||
		strings.HasPrefix(path, "/api/medias") ||
		strings.HasPrefix(path, "/api/upload") {
		return "owners"
	}

	// Visitor endpoints (viewing shared albums)
	if strings.HasPrefix(path, "/api/albums") && strings.Contains(path, "/shared") ||
		strings.HasPrefix(path, "/api/timeline") {
		return "visitors"
	}

	// Default to owners for unknown paths
	return "owners"
}

func hasRequiredPermission(userGroups []string, requiredGroup string) bool {
	// Hierarchical permissions: admins > owners > visitors
	groupHierarchy := map[string]int{
		"admins":   3,
		"owners":   2,
		"visitors": 1,
	}

	requiredLevel, ok := groupHierarchy[requiredGroup]
	if !ok {
		requiredLevel = 3 // Default to highest level if unknown
	}

	for _, group := range userGroups {
		if level, ok := groupHierarchy[group]; ok && level >= requiredLevel {
			return true
		}
	}

	return false
}

func denyAllResources(reason string) events.APIGatewayV2CustomAuthorizerSimpleResponse {
	return events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: false,
		Context: map[string]interface{}{
			"reason": reason,
		},
	}
}

func main() {
	lambda.Start(Handler)
}
