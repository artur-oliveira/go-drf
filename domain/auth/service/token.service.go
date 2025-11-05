package service

import (
	"errors"
	"fmt"
	"grf/core/config"
	"grf/domain/auth/model"
	"grf/domain/auth/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type TokenService struct {
	Config   *config.Config
	UserRepo *repository.UserRepository
}

type CustomClaims struct {
	UserID uint64 `json:"user_id"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

func NewTokenService(db *gorm.DB, config *config.Config) *TokenService {
	return &TokenService{
		Config: config, UserRepo: repository.NewUserRepository(db),
	}
}

func (s *TokenService) GenerateTokenPair(user *model.User) (accessToken string, refreshToken string, err error) {
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

func (s *TokenService) ValidateToken(tokenString string, expectedType string) (*model.User, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signature method: %v", token.Header["alg"])
		}
		return []byte(s.Config.JWTSecret), nil
	})

	if err != nil {
		return nil, errors.New("token invalid or expired")
	}

	if !token.Valid || claims.UserID == 0 {
		return nil, errors.New("invalid token")
	}

	if claims.Type != expectedType {
		return nil, fmt.Errorf("invalid token type: expected '%s', received '%s'", expectedType, claims.Type)
	}

	var user *model.User
	if user, err = s.UserRepo.FindById(claims.UserID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("user is not active")
	}

	return user, nil
}
