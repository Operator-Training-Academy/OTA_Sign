package app

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

const keycloakStateCookieName = "otasign_oidc_state"

type KeycloakClient struct {
	config   oauth2.Config
	verifier *oidc.IDTokenVerifier
	cfg      Config
}

type oidcState struct {
	State        string `json:"state"`
	Nonce        string `json:"nonce"`
	CodeVerifier string `json:"code_verifier"`
	ExpiresAt    int64  `json:"expires_at"`
}

func NewKeycloakClient(ctx context.Context, cfg Config) (*KeycloakClient, error) {
	provider, err := oidc.NewProvider(ctx, cfg.KeycloakIssuerURL)
	if err != nil {
		return nil, err
	}

	return &KeycloakClient{
		config: oauth2.Config{
			ClientID:     cfg.KeycloakClientID,
			ClientSecret: cfg.KeycloakClientSecret,
			RedirectURL:  cfg.KeycloakRedirectURL,
			Endpoint:     provider.Endpoint(),
			Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		},
		verifier: provider.Verifier(&oidc.Config{ClientID: cfg.KeycloakClientID}),
		cfg:      cfg,
	}, nil
}

func (c *KeycloakClient) authorizationURL(w http.ResponseWriter) (string, error) {
	state := oidcState{
		State:        randomID(),
		Nonce:        randomID(),
		CodeVerifier: randomID(),
		ExpiresAt:    time.Now().Add(10 * time.Minute).Unix(),
	}
	encoded, err := c.signState(state)
	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     keycloakStateCookieName,
		Value:    encoded,
		Path:     "/auth/keycloak",
		HttpOnly: true,
		Secure:   c.cfg.SessionCookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((10 * time.Minute).Seconds()),
	})

	return c.config.AuthCodeURL(state.State,
		oidc.Nonce(state.Nonce),
		oauth2.S256ChallengeOption(state.CodeVerifier),
	), nil
}

func (c *KeycloakClient) exchange(r *http.Request) (User, error) {
	state, err := c.readState(r)
	if err != nil {
		return User{}, err
	}
	if subtleEqual(state.State, r.URL.Query().Get("state")) == false {
		return User{}, errors.New("OIDC state mismatch")
	}
	if code := strings.TrimSpace(r.URL.Query().Get("code")); code == "" {
		return User{}, errors.New("missing authorization code")
	} else {
		token, exchangeErr := c.config.Exchange(r.Context(), code, oauth2.VerifierOption(state.CodeVerifier))
		if exchangeErr != nil {
			return User{}, exchangeErr
		}
		rawIDToken, ok := token.Extra("id_token").(string)
		if !ok || rawIDToken == "" {
			return User{}, errors.New("Keycloak response did not contain an ID token")
		}
		idToken, verifyErr := c.verifier.Verify(r.Context(), rawIDToken)
		if verifyErr != nil {
			return User{}, verifyErr
		}

		var claims map[string]any
		if claimsErr := idToken.Claims(&claims); claimsErr != nil {
			return User{}, claimsErr
		}
		if claimString(claims, "nonce") != state.Nonce {
			return User{}, errors.New("OIDC nonce mismatch")
		}

		return userFromKeycloakClaims(c.cfg, idToken.Issuer, idToken.Subject, claims)
	}
}

func (c *KeycloakClient) clearState(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     keycloakStateCookieName,
		Value:    "",
		Path:     "/auth/keycloak",
		HttpOnly: true,
		Secure:   c.cfg.SessionCookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

func (c *KeycloakClient) signState(state oidcState) (string, error) {
	payload, err := json.Marshal(state)
	if err != nil {
		return "", err
	}
	encoded := base64.RawURLEncoding.EncodeToString(payload)
	mac := hmac.New(sha256.New, []byte(c.cfg.KeycloakClientSecret))
	_, _ = mac.Write([]byte(encoded))
	return encoded + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), nil
}

func (c *KeycloakClient) readState(r *http.Request) (oidcState, error) {
	cookie, err := r.Cookie(keycloakStateCookieName)
	if err != nil {
		return oidcState{}, errors.New("missing OIDC state cookie")
	}
	encoded, signature, ok := strings.Cut(cookie.Value, ".")
	if !ok || encoded == "" || signature == "" {
		return oidcState{}, errors.New("invalid OIDC state cookie")
	}
	mac := hmac.New(sha256.New, []byte(c.cfg.KeycloakClientSecret))
	_, _ = mac.Write([]byte(encoded))
	expected := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(signature)) {
		return oidcState{}, errors.New("invalid OIDC state signature")
	}
	payload, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return oidcState{}, errors.New("invalid OIDC state payload")
	}
	var state oidcState
	if err := json.Unmarshal(payload, &state); err != nil || state.ExpiresAt < time.Now().Unix() {
		return oidcState{}, errors.New("expired OIDC state")
	}
	return state, nil
}

func userFromKeycloakClaims(cfg Config, issuer, subject string, claims map[string]any) (User, error) {
	dodID := claimString(claims, cfg.KeycloakDoDIDClaim)
	uic := claimString(claims, cfg.KeycloakUICClaim)
	email := claimString(claims, "email")
	fullName := claimString(claims, "name")
	firstName := claimString(claims, "given_name")
	lastName := claimString(claims, "family_name")
	if fullName == "" {
		fullName = strings.TrimSpace(firstName + " " + lastName)
	}
	if issuer == "" || subject == "" || dodID == "" || uic == "" || fullName == "" || email == "" {
		return User{}, errors.New("Keycloak token is missing a required OTA Sign identity claim")
	}

	return User{
		ID:              "keycloak-" + subject,
		MoodleUserID:    "keycloak-" + subject,
		KeycloakIssuer:  issuer,
		KeycloakSubject: subject,
		FullName:        fullName,
		FirstName:       firstName,
		LastName:        lastName,
		Email:           email,
		ArmyEmail:       validArmyEmail(claimString(claims, cfg.KeycloakArmyEmailClaim)),
		DoDID:           dodID,
		Rank:            claimString(claims, cfg.KeycloakRankClaim),
		UIC:             uic,
	}, nil
}

func claimString(claims map[string]any, name string) string {
	value, ok := claims[name]
	if !ok {
		return ""
	}
	stringValue, ok := value.(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(stringValue)
}

func subtleEqual(left, right string) bool {
	return hmac.Equal([]byte(left), []byte(right))
}
