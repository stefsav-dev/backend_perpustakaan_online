package config

import (
	"backend_perpustakaan_online/models"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	SecretKey      string
	AdminSecretKey string
	UserSecretKey  string
	ExpiresIn      time.Duration
}

var JWT *JWTConfig

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func InitJWT() {
	expiresIn, _ := strconv.Atoi(os.Getenv("JWT_EXPIRES_IN"))
	if expiresIn == 0 {
		expiresIn = 24
	}

	JWT = &JWTConfig{
		SecretKey:      getEnv("JWT_SECRET_KEY", "default-secret-key"),
		AdminSecretKey: getEnv("JWT_SECRET_KEY", "admin-secret-key"),
		UserSecretKey:  getEnv("JWT_SECRET_KEY", "user-secret-key"),
		ExpiresIn:      time.Hour * time.Duration(expiresIn),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (j *JWTConfig) GenerateToken(userID uint, email string, role models.UserRole) (string, error) {
	var secretKey string

	switch role {
	case models.RoleAdmin:
		secretKey = j.AdminSecretKey
	case models.RoleUser:
		secretKey = j.UserSecretKey
	default:
		secretKey = j.SecretKey
	}

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   string(role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func (j *JWTConfig) ValidateToken(tokenString string, role models.UserRole) (*Claims, error) {
	var secretKey string

	switch role {
	case models.RoleAdmin:
		secretKey = j.AdminSecretKey
	case models.RoleUser:
		secretKey = j.UserSecretKey
	default:
		secretKey = j.SecretKey
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

func (j *JWTConfig) ValidateAnyToken(tokenString string) (*Claims, error) {

	if claims, err := j.ValidateToken(tokenString, models.RoleAdmin); err == nil {
		return claims, nil
	}

	if claims, err := j.ValidateToken(tokenString, models.RoleUser); err == nil {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
