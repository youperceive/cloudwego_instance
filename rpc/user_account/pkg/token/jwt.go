package token

import (
	"errors"
	"os"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/golang-jwt/jwt/v4"
)

var (
	secretKey      = getSecretKey()
	expireDuration = 2 * time.Hour
)

func getSecretKey() string {
	key := os.Getenv("JWT_SECRETKEY")
	if key == "" {
		klog.Warn("token/jwt", "not find env param $JWT_SECRETKEY, has been replaced by 'temprory key'")
		key = "temprory key"
	}
	if len(key) < 32 {
		klog.Warn("token/jwt", "length of secret_key is less than 32, suggested more than 32")
	}
	return key
}

type CustomClaims struct {
	UserID   int64 `json:"user_id"`
	UserType int8  `json:"user_type"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int64, userType int8) (string, error) {
	claims := CustomClaims{
		UserID:   userID,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "user_account_service",
			Subject:   "login_token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", errors.New("generate jwt token failed: " + err.Error())
	}
	return signedToken, nil
}

func VerifyToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (any, error) {
			alg, ok := token.Header["alg"].(string)
			if !ok {
				return nil, errors.New("invalid alg type in jwt header")
			}
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method: " + alg)
			}
			return []byte(secretKey), nil
		},
	)
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			switch {
			case ve.Errors&jwt.ValidationErrorExpired != 0:
				return nil, errors.New("jwt token expired")
			case ve.Errors&jwt.ValidationErrorSignatureInvalid != 0:
				return nil, errors.New("jwt signature invalid")
			default:
				return nil, errors.New("jwt validation failed: " + ve.Error())
			}
		}
		return nil, errors.New("verify jwt token failed: " + err.Error())
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		if claims.Issuer != "user_account_service" {
			return nil, errors.New("invalid jwt issuer: not from user_account_service")
		}
		return claims, nil
	}
	return nil, errors.New("invalid jwt token")
}
