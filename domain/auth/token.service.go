package auth

import (
	"errors"
	"fmt"
	"grf/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// TokenService encapsula a lógica de JWT
type TokenService struct {
	Config *config.Config
	DB     *gorm.DB
}

// CustomClaims define as reivindicações do JWT
type CustomClaims struct {
	UserID uint64 `json:"user_id"`
	Type   string `json:"type"` // "access" ou "refresh"
	jwt.RegisteredClaims
}

// NewTokenService cria um novo serviço de token
func NewTokenService(db *gorm.DB, config *config.Config) *TokenService {
	return &TokenService{Config: config, DB: db}
}

// GenerateTokenPair cria um par de tokens de acesso e refresh
func (s *TokenService) GenerateTokenPair(user *User) (accessToken string, refreshToken string, err error) {
	// 1. Criar Access Token
	accessExp := time.Now().Add(time.Minute * time.Duration(s.Config.JWTExpiresInMinutes))
	accessClaims := CustomClaims{
		UserID: user.ID,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).
		SignedString([]byte(s.Config.JWTSecret))
	if err != nil {
		return "", "", err
	}

	// 2. Criar Refresh Token
	refreshExp := time.Now().Add(time.Hour * 24 * time.Duration(s.Config.JWTRefreshExpiresInDays))
	refreshClaims := CustomClaims{
		UserID: user.ID,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).
		SignedString([]byte(s.Config.JWTSecret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ValidateToken analisa um token e retorna o usuário (se válido e ativo)
func (s *TokenService) ValidateToken(tokenString string, expectedType string) (*User, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return []byte(s.Config.JWTSecret), nil
	})

	if err != nil {
		return nil, errors.New("token inválido ou expirado")
	}

	if !token.Valid || claims.UserID == 0 {
		return nil, errors.New("token inválido")
	}

	if claims.Type != expectedType {
		return nil, fmt.Errorf("tipo de token inválido: esperado '%s', obteve '%s'", expectedType, claims.Type)
	}

	// Token é válido, buscar usuário no DB
	var user User
	if err := s.DB.First(&user, claims.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuário do token não encontrado")
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("usuário inativo")
	}

	return &user, nil
}
