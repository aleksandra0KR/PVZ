package auth

import (
	"final/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	os.Setenv("SECRET_KEY", "test_secret")
	os.Setenv("JWT_EXP", "1h")

	err := InitAuthFromConfig()
	if err != nil {
		panic(err)
	}

	code := m.Run()
	os.Exit(code)
}
func TestInitAuthFromConfig(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		os.Setenv("SECRET_KEY", "test_secret")
		os.Setenv("JWT_EXP", "1h")

		err := InitAuthFromConfig()
		assert.NoError(t, err)
		assert.Equal(t, "test_secret", string(jwtSecret))
		assert.Equal(t, time.Hour, tokenTTL)
	})

	t.Run("missing secret", func(t *testing.T) {
		os.Unsetenv("SECRET_KEY")

		err := InitAuthFromConfig()
		assert.ErrorIs(t, err, domain.ErrMissingSecret)
	})

	t.Run("invalid expiration", func(t *testing.T) {
		os.Setenv("SECRET_KEY", "test_secret")
		os.Setenv("JWT_EXP", "invalid")

		err := InitAuthFromConfig()
		assert.ErrorIs(t, err, domain.ErrInitAuthConfig)
	})
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("valid token", func(t *testing.T) {
		router := gin.Default()
		router.Use(Middleware())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		token, err := GenerateJWT("user1", "employee")
		assert.NoError(t, err)

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("missing token", func(t *testing.T) {
		router := gin.Default()
		router.Use(Middleware())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid token", func(t *testing.T) {
		router := gin.Default()
		router.Use(Middleware())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid_token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestCheckRole(t *testing.T) {
	t.Run("allowed role", func(t *testing.T) {
		role := "employee"
		assert.True(t, CheckRole(&role))
	})

	t.Run("not allowed role", func(t *testing.T) {
		role := "admin"
		assert.False(t, CheckRole(&role))
	})

	t.Run("nil role", func(t *testing.T) {
		assert.False(t, CheckRole(nil))
	})
}

func TestRequireRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("allowed role", func(t *testing.T) {
		router := gin.Default()
		router.Use(Middleware())
		router.Use(RequireRole("employee"))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		token, err := GenerateJWT("user1", "employee")
		assert.NoError(t, err, "Error generating JWT")
		assert.NotEmpty(t, token, "Generated JWT should not be empty")

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Expected status OK")
	})

	t.Run("not allowed role", func(t *testing.T) {
		router := gin.Default()
		router.Use(RequireRole("employee"))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		token, _ := GenerateJWT("user1", "admin")
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestGenerateJWT(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		os.Setenv("SECRET_KEY", "test_secret")
		os.Setenv("JWT_EXP", "1h")
		InitAuthFromConfig()

		token, err := GenerateJWT("user1", "employee")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		assert.NoError(t, err)
		assert.True(t, parsedToken.Valid)

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.Equal(t, "user1", claims["user_id"])
		assert.Equal(t, "employee", claims["role"])
	})
}
