package service

import (
	"golang-api/helper"

	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrTokenMalformed   = jwt.ErrTokenMalformed
	ErrTokenExpired     = jwt.ErrTokenExpired
	ErrTokenNotValidYet = jwt.ErrTokenNotValidYet
)

// JWTService is a contract of what the JWT Service can do
type JWTService interface {
	GenerateToken(userId string) string
	ValidateToken(token string) (*jwt.Token, error)
	MapClaims(token *jwt.Token) jwt.MapClaims
	GetJWTIssuer() string
}

type jwtCustomClaim struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type jwtService struct {
	issuer    string
	secretKey string
	TTL       int64
}

// NewJWTService creates a new instance of JWTService
func NewJWTService() JWTService {
	return &jwtService{
		issuer:    getJWTIssuer(),
		secretKey: getJWTSecretKey(),
		TTL:       getJWTTTL(),
	}
}

func getJWTIssuer() string {
	return os.Getenv(helper.JWT_ISSUER)
}

func getJWTSecretKey() string {
	return os.Getenv("JWT_PRIVATE_KEY")
}

func getJWTTTL() int64 {
	ttl, err := strconv.ParseInt(os.Getenv("TOKEN_TTL"), 10, 0)
	if err != nil {
		panic(err)
	}

	return ttl
}

// GenerateToken generated a JWT token in the userId and additional env details provided
func (s *jwtService) GenerateToken(userId string) string {
	claims := &jwtCustomClaim{
		userId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.TTL) * time.Hour)),
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(s.secretKey))

	if err != nil {
		panic(err)
	}

	return t
}

// ValidateToken validates the provided Token signing method and JWT details.
func (s *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(s.secretKey), nil
	})
}

func (s *jwtService) MapClaims(token *jwt.Token) jwt.MapClaims {
	return token.Claims.(jwt.MapClaims)
}

func (s *jwtService) GetJWTIssuer() string {
	return getJWTIssuer()
}
