package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/olahol/melody"
)

// private dan public tokennya adalah utf-8 yang di encode ke base64 saat akan membuat atau memvalidasi token maka tokennya di kembalikan ke utf-8 untuk masuk di func jwt.ParseRSAPrivateKeyFromPEM
func CreateToken(ttl time.Duration, content any, privateKey string) (string, error) {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", fmt.Errorf("could not decode key: %w", err)
	}

	// parse the key from base64 back to utf8
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	now := time.Now().UTC()

	claims := jwt.MapClaims{
		"sub": content,
		"exp": now.Add(ttl).Unix(),
		"iat": now.Unix(),
		"nbf": now.Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	// fmt.Println(claims)
	return token, nil
}

func ValidateToken(token string, publicKey string) (map[string]any, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		return nil, fmt.Errorf("validate: parse key: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		// write any validation u want
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		// return the key to parse function
		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	// if cannot convert to jwt mapclaims or the parsed token is invalid
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}
	// fmt.Println(claims["sub"])
	return claims["sub"].(map[string]any), nil
}

func GetToken(ctx *gin.Context, cookieName string, headerName string) (token string) {
	tokenCookie, err := ctx.Cookie(cookieName)
	tokenHeader := ctx.Request.Header.Get(headerName)

	tokenHeaderFieldsSplits := strings.Fields(tokenHeader)

	if len(tokenHeaderFieldsSplits) != 0 && tokenHeaderFieldsSplits[0] == "Bearer" {
		token = tokenHeaderFieldsSplits[1]
	} else if err == nil {
		// if no header token check if there's cookie token
		// if no error cookie token then
		token = tokenCookie
	} else {
		token = ""
	}
	return
}

func GetTokenWS(s *melody.Session, cookieName string, headerName string) (token string) {
	tokenCookie, err := s.Request.Cookie(cookieName)
	tokenHeader := s.Request.Header.Get(headerName)

	tokenHeaderFieldsSplits := strings.Fields(tokenHeader)

	if len(tokenHeaderFieldsSplits) != 0 && tokenHeaderFieldsSplits[0] == "Bearer" {
		token = tokenHeaderFieldsSplits[1]
	} else if err == nil {
		// if no header token check if there's cookie token
		// if no error cookie token then
		token = tokenCookie.Value
	} else {
		token = ""
	}
	return
}

func GetCookie(ctx *gin.Context, cookieName string) (token string, err error) {
	tokenCookie, err := ctx.Cookie(cookieName)

	if err != nil {
		return "", err
	}

	return tokenCookie, nil
}

// generateVerificationToken generates a unique verification token.
func GenerateVerificationToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(token), nil
}

// generateRandomCode generates a random 6-digit code.
func GenerateRandomCode() string {
	return uniuri.NewLen(6)
}

func GenerateRandomCodeLength(length int) string {
	return uniuri.NewLen(length)
}

func GenerateRandomNumCode() (string, error) {
	// Generate a random big integer in the range [0, 999999]
	randomBigInt, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}

	// Format the random number as a 6-digit string with leading zeros
	code := fmt.Sprintf("%06s", randomBigInt.String())

	return code, nil
}

func GenerateRandomNumCodeLength(l int64) (string, error) {
	base := big.NewInt(10)
	exponent := big.NewInt(l)

	result := new(big.Int).Exp(base, exponent, nil)
	randomBigInt, err := rand.Int(rand.Reader, result)
	if err != nil {
		return "", err
	}

	// Format the random number as a 6-digit string with leading zeros
	code := fmt.Sprintf("%06s", randomBigInt.String())

	return code, nil
}
