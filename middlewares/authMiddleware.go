package middlewares

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var privateKey = []byte(os.Getenv("JWT_PRIVATE_KEY"))

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

var blacklist = struct {
	sync.RWMutex
	tokens map[string]struct{}
}{tokens: make(map[string]struct{})}

func tokenIsBlacklisted(token string) bool {
	blacklist.RLock()
	defer blacklist.RUnlock()
	_, exists := blacklist.tokens[token]
	return exists
}

func GenerateToken(username string, role string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour) //Reset to 30mins
	claims := &Claims{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "Unikorn",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func GenerateRefreshToken(username string, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "Unikorn",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func validateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return privateKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Check if the token has been revoked
		isRevoked := tokenIsBlacklisted(tokenString)
		if isRevoked {
			return nil, errors.New("token is revoked")
		}
		return claims, nil
	}
	return nil, err
}

func InvalidateToken(tokenString string) error {
	blacklist.Lock()
	defer blacklist.Unlock()
	blacklist.tokens[tokenString] = struct{}{}
	return nil
}

func GetAuthUserID(c *gin.Context) string {
	claims, _ := c.Get("claims")
	return claims.(*Claims).Username
}

func GetAuthUserRole(c *gin.Context) string {
	claims, _ := c.Get("claims")
	return claims.(*Claims).Role
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if len(tokenString) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is missing"})
			return
		}
		tokenString = strings.Split(tokenString, "Bearer ")[1]
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is missing"})
			return
		}
		claims, err := validateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}
