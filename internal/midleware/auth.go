package auth

import (
	"final/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"time"
)

var jwtSecret []byte
var tokenTTL time.Duration

func InitAuthFromConfig() error {
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		return domain.ErrMissingSecret
	}
	jwtSecret = []byte(secret)

	expiration := os.Getenv("JWT_EXP")
	if expiration == "" {
		expiration = "24h"
	}
	ttl, err := time.ParseDuration(expiration)
	if err != nil {
		return domain.ErrInitAuthConfig
	} else if ttl == 0 {
		ttl = 72 * time.Hour
	}
	tokenTTL = ttl

	return nil
}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if !strings.HasPrefix(tokenStr, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, domain.ErrorResponse{Message: domain.ErrMissingToken.Error()})
			return
		}
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, domain.ErrInvalidToken
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, domain.ErrorResponse{Message: domain.ErrInvalidToken.Error()})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, domain.ErrorResponse{Message: domain.ErrInvalidToken.Error()})
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])
		c.Next()
	}
}

var allowedRoles = map[string]struct{}{
	"employee":  {},
	"moderator": {},
}

func CheckRole(role *string) bool {
	if role == nil {
		return false
	}
	_, ok := allowedRoles[*role]
	return ok
}

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		for _, r := range allowedRoles {
			if role == r {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusForbidden, domain.ErrorResponse{Message: domain.ErrInvalidRole.Error()})
		c.Abort()
	}
}

func GenerateJWT(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(tokenTTL).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Error(err)
		return "", domain.ErrGenerateToken
	}
	return tokenString, nil
}
