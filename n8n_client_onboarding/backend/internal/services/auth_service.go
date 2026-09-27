package services

import (
	"backend/internal/config"
	"backend/internal/models"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthService struct {
	config config.Config

	oauthConfig *oauth2.Config

	httpClient *http.Client
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type AuthClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`

	jwt.RegisteredClaims
}

const oauthStateCookieName = "google_oauth_state"

func NewAuthService(cfg config.Config) *AuthService {

	return &AuthService{
		config: cfg,
		oauthConfig: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleRedirectURL,
			Endpoint:     google.Endpoint,
			Scopes:       []string{"openid", "profile", "email"},
		},
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (s *AuthService) GenerateStateToken() (string, error) {
	buffer := make([]byte, 32)

	if _, err := rand.Read(buffer); err != nil {
		return "", fmt.Errorf("generate state token failed")
	}

	return base64.RawURLEncoding.EncodeToString(buffer), nil
}

func (s *AuthService) SetOauthStateCookie(c fiber.Ctx, value string) {

	c.Cookie(&fiber.Cookie{
		Name:     oauthStateCookieName,
		Value:    value,
		Path:     "/",
		HTTPOnly: true,
		Secure:   s.config.CookieSecure,
		SameSite: s.config.CookieSameSite,
		Domain:   s.config.CookieDomain,
		MaxAge:   20 * 60,
	})
}

func (s *AuthService) BuildGoogleAuthUrl(state string) string {

	// AuthCodeURL creates google oauth consent url
	return s.oauthConfig.AuthCodeURL(state)
}

func (s *AuthService) ReadOauthStateCookie(c fiber.Ctx) string {
	return c.Cookies(oauthStateCookieName, "")
}

func (s *AuthService) SetAuthCookie(c fiber.Ctx, token string) {

	maxAge := s.config.JWTExpiresInHours * 60 * 60

	c.Cookie(&fiber.Cookie{
		Name:     s.config.AuthCookieName,
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   s.config.CookieSecure,
		SameSite: s.config.CookieSameSite,
		Domain:   s.config.CookieDomain,
		MaxAge:   maxAge,
	})
}

func (s *AuthService) ClearOauthStateCookie(c fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     oauthStateCookieName,
		Value:    "",
		Path:     "/",
		HTTPOnly: true,
		Secure:   s.config.CookieSecure,
		SameSite: s.config.CookieSameSite,
		Domain:   s.config.CookieDomain,
		// Expire the cookie in the past for better browser compatibility
		Expires: time.Unix(0, 0),
	})
}

func (s *AuthService) ExchangeCodeForGoogleUser(ctx context.Context, code string) (*GoogleUserInfo, error) {

	token, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchange auth code failed")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("creating user info failed")
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching user info failed")
	}

	// prevents resource leaks
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Google userinfo returned failed")
	}

	var userInfo GoogleUserInfo

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("decode user info failed")
	}

	if strings.TrimSpace(userInfo.Email) == "" {
		return nil, fmt.Errorf("google userinfo didnt return any email")
	}

	return &userInfo, nil
}

func (s *AuthService) SignJWT(user *models.User) (string, error) {
	expiresAt := time.Now().Add(time.Duration(s.config.JWTExpiresInHours) * time.Hour)

	claims := AuthClaims{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("Sign jwt error")
	}

	return signed, nil
}

func (s *AuthService) ParseJWT(tokenString string) (*AuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(s.config.JWTSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse jwt failed: %w", err)
	}

	// convert generic claims into our custom authclaims type
	claims, ok := token.Claims.(*AuthClaims)

	if !ok || !token.Valid {
		return nil, fmt.Errorf("Invalid jwt token")
	}

	return claims, nil

}

func (s *AuthService) ClearAuthCookie(c fiber.Ctx) {

	c.Cookie(&fiber.Cookie{
		Name:     s.config.AuthCookieName,
		Value:    "",
		Path:     "/",
		HTTPOnly: true,
		Secure:   s.config.CookieSecure,
		SameSite: s.config.CookieSameSite,
		Domain:   s.config.CookieDomain,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
}
