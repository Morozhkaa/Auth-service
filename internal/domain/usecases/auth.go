package usecases

import (
	"auth/internal/config"
	"auth/internal/domain/models"
	"auth/internal/ports"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"time"

	jwt4 "github.com/golang-jwt/jwt/v4"
	"github.com/xdg-go/pbkdf2"
)

type Auth struct {
	userStorage ports.UserStorage
}

func New(userStorage ports.UserStorage) (*Auth, error) {
	return &Auth{
		userStorage: userStorage,
	}, nil
}

const accessToken_expiration_time = 1 * time.Minute
const refreshToken_expiration_time = 60 * time.Minute

func (a *Auth) EncodePassword(password string) string {
	cfg, _ := config.GetConfig()
	dk := pbkdf2.Key([]byte(password), []byte(cfg.Salt), 1000, 128, sha1.New)
	return base64.StdEncoding.EncodeToString([]byte(dk))
}

func (a *Auth) generateToken(login string, email string, expiredIn time.Duration) (string, error) {
	cfg, _ := config.GetConfig()
	now := time.Now()
	claims := &jwt4.RegisteredClaims{
		ExpiresAt: jwt4.NewNumericDate(now.Add(expiredIn)),
		Issuer:    login,
		Subject:   email,
	}

	token := jwt4.NewWithClaims(jwt4.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", models.ErrGenerateToken
	}
	return ss, nil
}

func (a *Auth) verifyToken(tokenString string) (string, string, error) {
	cfg, _ := config.GetConfig()
	var claims jwt4.RegisteredClaims

	token, err := jwt4.ParseWithClaims(tokenString, &claims, func(token *jwt4.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})
	if token == nil || !token.Valid {
		return "", "", fmt.Errorf("parse token unexpected error: %w", err)
	}
	return claims.Issuer, claims.Subject, nil
}

// Login checks that the user with the given login and password is registered, and generates JWT access, refresh tokens.
func (a *Auth) Login(ctx context.Context, login string, password string) (string, string, error) {
	user, err := a.userStorage.GetUser(ctx, login)
	if err != nil {
		return "", "", models.ErrNotFound
	}
	if user.Password != a.EncodePassword(password) {
		return "", "", models.ErrForbidden
	}

	access, err := a.generateToken(login, user.Email, accessToken_expiration_time)
	if err != nil {
		return "", "", err
	}
	refresh, err := a.generateToken(login, user.Email, refreshToken_expiration_time)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

// Verify validates tokens and returns user data. If the access token has not expired, returns the user data.
// Otherwise, if the refresh token has not expired, returns the user data and new access and refresh tokens.
func (a *Auth) Verify(ctx context.Context, access string, refresh string) (r models.VerifyResponse, err error) {
	r.Login, r.Email, err = a.verifyToken(access)
	if err == nil {
		r.AccessToken = access
		r.RefreshToken = refresh
		return r, nil
	}

	r.Login, r.Email, err = a.verifyToken(refresh)
	if err == nil {
		r.AccessToken, err = a.generateToken(r.Login, r.Email, accessToken_expiration_time)
		if err != nil {
			return r, err
		}
		r.RefreshToken, err = a.generateToken(r.Login, r.Email, refreshToken_expiration_time)
		if err != nil {
			return r, err
		}
		return r, nil
	}
	return r, models.ErrTokenExpired
}
