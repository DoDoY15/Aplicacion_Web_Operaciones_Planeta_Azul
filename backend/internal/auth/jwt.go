package auth

import (
	"crypto/rsa"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/planeta-azul/backend/internal/models"
)

type Service struct {
	privateKey    *rsa.PrivateKey
	publicKey     *rsa.PublicKey
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type AccessClaims struct {
	jwt.RegisteredClaims
	UserID string          `json:"user_id"`
	Email  string          `json:"email"`
	Role   models.UserRole `json:"role"`
	AreaID string          `json:"area_id"`
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
	Type   string `json:"type"`
}

func NewService(privateKeyPath, publicKeyPath string, accessExpiry, refreshExpiry time.Duration) (*Service, error) {
	privBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privBytes)
	if err != nil {
		return nil, err
	}

	pubBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		return nil, err
	}

	return &Service{
		privateKey:    privateKey,
		publicKey:     publicKey,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}, nil
}

// NewServiceFromKeys creates service from in-memory keys (useful for testing)
func NewServiceFromKeys(private *rsa.PrivateKey, public *rsa.PublicKey, accessExpiry, refreshExpiry time.Duration) *Service {
	return &Service{
		privateKey:    private,
		publicKey:     public,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

func (s *Service) GenerateTokenPair(user *models.User) (*TokenPair, error) {
	areaID := ""
	if user.AreaID != nil {
		areaID = user.AreaID.String()
	}

	// Access token
	accessClaims := AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessExpiry)),
		},
		UserID: user.ID.String(),
		Email:  user.Email,
		Role:   user.Role,
		AreaID: areaID,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)
	accessSigned, err := accessToken.SignedString(s.privateKey)
	if err != nil {
		return nil, err
	}

	// Refresh token
	refreshClaims := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshExpiry)),
		},
		UserID: user.ID.String(),
		Type:   "refresh",
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)
	refreshSigned, err := refreshToken.SignedString(s.privateKey)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessSigned,
		RefreshToken: refreshSigned,
	}, nil
}

func (s *Service) ValidateAccessToken(tokenStr string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AccessClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func (s *Service) ValidateRefreshToken(tokenStr string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &RefreshClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid || claims.Type != "refresh" {
		return nil, errors.New("invalid refresh token")
	}
	return claims, nil
}

// ParseUserIDFromRefresh extracts UUID from a validated refresh token
func (s *Service) ParseUserIDFromRefresh(tokenStr string) (uuid.UUID, error) {
	claims, err := s.ValidateRefreshToken(tokenStr)
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(claims.UserID)
}
