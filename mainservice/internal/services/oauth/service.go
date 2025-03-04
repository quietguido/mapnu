package oauth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type OAuthService interface {
	VerifyIDToken(ctx context.Context, idToken string) (*Claims, error)
	GenerateJWT(email, givenName, familyName string) (string, error)
	VerifyJWT(ctx context.Context, tokenString string) (*Claims, error)
}

type Service struct {
	lg          *zap.Logger
	jwtSecret   []byte
	publicKeys  map[string]*rsa.PublicKey
	keyCacheTTL time.Time
	mu          sync.Mutex
	clientID    string
}

func NewOAuthService(lg *zap.Logger) OAuthService {
	jwtSecret, exists := os.LookupEnv("JWT_SECRET")
	if !exists {
		lg.Fatal("JWT_SECRET is missing")
	}

	clientID, exists := os.LookupEnv("GOOGLE_CLIENT_ID")
	if !exists {
		lg.Fatal("GOOGLE_CLIENT_ID is missing")
	}

	return &Service{
		lg:         lg,
		jwtSecret:  []byte(jwtSecret),
		publicKeys: make(map[string]*rsa.PublicKey),
		clientID:   clientID,
	}
}

func (s *Service) VerifyIDToken(ctx context.Context, idToken string) (*Claims, error) {
	s.lg.Info("Verifying ID Token")

	// Step 1: Decode JWT Header to get `kid`
	parsedToken, _, err := new(jwt.Parser).ParseUnverified(idToken, &Claims{})
	if err != nil {
		s.lg.Error("Failed to parse token header", zap.Error(err))
		return nil, errors.New("invalid ID token")
	}

	kid, ok := parsedToken.Header["kid"].(string)
	if !ok {
		s.lg.Error("Invalid token header: missing kid")
		return nil, errors.New("invalid token header: missing kid")
	}

	// Step 2: Fetch and cache Google's public keys
	publicKey, err := s.getGooglePublicKey(kid)
	if err != nil {
		s.lg.Error("Failed to get Google public key", zap.Error(err))
		return nil, err
	}

	// Step 3: Verify JWT
	token, err := jwt.ParseWithClaims(idToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		s.lg.Error("Failed to verify ID token", zap.Error(err))
		return nil, errors.New("failed to verify ID token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		s.lg.Error("Invalid token claims")
		return nil, errors.New("invalid token claims")
	}

	// Step 4: Validate issuer
	if claims.Issuer != "https://accounts.google.com" && claims.Issuer != "accounts.google.com" {
		s.lg.Error("Invalid token issuer", zap.String("issuer", claims.Issuer))
		return nil, errors.New("invalid token issuer")
	}

	// Step 5: Validate audience
	if !contains(claims.Audience, s.clientID) {
		s.lg.Error("Invalid token audience", zap.Strings("audience", claims.Audience))
		return nil, errors.New("invalid token audience")
	}

	// Step 6: Validate expiration
	if claims.ExpiresAt.Time.Before(time.Now()) {
		s.lg.Warn("JWT has expired", zap.Time("exp", claims.ExpiresAt.Time))
		return nil, errors.New("token has expired")
	}

	s.lg.Info("ID Token successfully verified", zap.String("email", claims.Email))
	return claims, nil
}

func contains(audience jwt.ClaimStrings, clientID string) bool {
	for _, aud := range audience {
		if aud == clientID {
			return true
		}
	}
	return false
}

func (s *Service) getGooglePublicKey(kid string) (*rsa.PublicKey, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Refresh keys if cache expired
	if time.Now().After(s.keyCacheTTL) {
		s.lg.Info("Refreshing Google public keys")
		keys, err := fetchGooglePublicKeys()
		if err != nil {
			return nil, err
		}
		s.publicKeys = keys
		s.keyCacheTTL = time.Now().Add(1 * time.Hour)
	}

	// Retrieve key
	publicKey, exists := s.publicKeys[kid]
	if !exists {
		return nil, errors.New("public key not found for given kid")
	}
	return publicKey, nil
}

func fetchGooglePublicKeys() (map[string]*rsa.PublicKey, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v3/certs")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var keys struct {
		Keys []struct {
			Kid string `json:"kid"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&keys); err != nil {
		return nil, err
	}

	publicKeys := make(map[string]*rsa.PublicKey)
	for _, key := range keys.Keys {
		rsaKey, err := parseRSAPublicKey(key.N, key.E)
		if err == nil {
			publicKeys[key.Kid] = rsaKey
		}
	}
	return publicKeys, nil
}

func parseRSAPublicKey(n, e string) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(n)
	if err != nil {
		return nil, err
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(e)
	if err != nil {
		return nil, err
	}

	eInt := 0
	for _, b := range eBytes {
		eInt = eInt<<8 + int(b)
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: eInt,
	}, nil
}

func (s *Service) VerifyJWT(ctx context.Context, tokenString string) (*Claims, error) {
	s.lg.Info("Verifying JWT token")

	// Parse and verify the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			s.lg.Error("Unexpected signing method", zap.String("alg", token.Method.Alg()))
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		s.lg.Error("JWT parsing error", zap.Error(err))
		return nil, errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		s.lg.Error("Invalid JWT claims")
		return nil, errors.New("invalid token claims")
	}

	// Check expiration
	if claims.ExpiresAt.Time.Before(time.Now()) {
		s.lg.Warn("JWT has expired", zap.Time("exp", claims.ExpiresAt.Time))
		return nil, errors.New("token has expired")
	}

	s.lg.Info("JWT successfully verified", zap.String("email", claims.Email))
	return claims, nil
}

func (s *Service) GenerateJWT(email, givenName, familyName string) (string, error) {
	expirationTime := time.Now().Add(20 * time.Minute) // Set expiration to 20 minutes

	claims := &Claims{
		Email:      email,
		GivenName:  givenName,
		FamilyName: familyName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

type Claims struct {
	Email      string `json:"email"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	jwt.RegisteredClaims
}
