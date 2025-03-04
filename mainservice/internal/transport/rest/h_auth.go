package rest

import (
	"github.com/quietguido/mapnu/mainservice/internal/services/oauth"
	"log"
	"net/http"
	"strings"
)

type OAuthHandler struct {
	service oauth.OAuthService
}

func NewOAuthHandler(service oauth.OAuthService) *OAuthHandler {
	return &OAuthHandler{
		service: service,
	}
}

func (h *OAuthHandler) HandleTokenExchange(w http.ResponseWriter, r *http.Request) {
	log.Println("HandleTokenExchange called")

	var req struct {
		IDToken string `json:"id_token"`
	}

	if err := JsonBodyDecoding(r, &req); err != nil {
		log.Println("Failed to decode JSON", err)
		RespondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	claims, err := h.service.VerifyIDToken(r.Context(), req.IDToken)
	if err != nil {
		log.Println("ID token verification failed", err)
		RespondWithError(w, http.StatusUnauthorized, "Invalid ID token")
		return
	}

	jwtToken, err := h.service.GenerateJWT(claims.Email, claims.GivenName, claims.FamilyName)
	if err != nil {
		log.Println("Failed to generate JWT", err)
		RespondWithError(w, http.StatusInternalServerError, "Error generating JWT")
		return
	}

	RespondWithJson(w, http.StatusOK, map[string]any{"jwt_token": jwtToken})
}

func (h *OAuthHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	log.Println("GetUserProfile called")

	// Step 1: Extract JWT from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Println("Authorization header missing")
		RespondWithError(w, http.StatusUnauthorized, "Missing authorization token")
		return
	}

	// Bearer token format -> "Bearer <token>"
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		log.Println("Invalid Authorization header format")
		RespondWithError(w, http.StatusUnauthorized, "Invalid token format")
		return
	}

	jwtToken := tokenParts[1]

	// Step 2: Verify JWT token
	claims, err := h.service.VerifyJWT(r.Context(), jwtToken)
	if err != nil {
		log.Println("JWT verification failed", err)
		RespondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	// Step 3: Return user profile details
	RespondWithJson(w, http.StatusOK, map[string]any{
		"given_name":  claims.GivenName,
		"family_name": claims.FamilyName,
	})
}
