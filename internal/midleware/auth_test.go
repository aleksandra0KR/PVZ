package auth

import (
	"final/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	err := os.Setenv("SECRET_KEY", "test_secret")
	if err != nil {
		log.Error(err)
		return
	}
	err = os.Setenv("JWT_EXP", "1h")
	if err != nil {
		log.Error(err)
		return
	}

	err = InitAuthFromConfig()
	if err != nil {
		panic(err)
	}

	code := m.Run()
	os.Exit(code)
}
func TestInitAuthFromConfig(t *testing.T) {
	t.Run("Successful_Auth", func(t *testing.T) {
		err := os.Setenv("SECRET_KEY", "test_secret")
		if err != nil {
			log.Error(err)
			return
		}
		err = os.Setenv("JWT_EXP", "1h")
		if err != nil {
			log.Error(err)
			return
		}

		err = InitAuthFromConfig()
		assert.NoError(t, err)
		assert.Equal(t, "test_secret", string(jwtSecret))
		assert.Equal(t, time.Hour, tokenTTL)
	})

	t.Run("Failure_Missing_Secret", func(t *testing.T) {
		err := os.Unsetenv("SECRET_KEY")
		if err != nil {
			log.Error(err)
			return
		}

		err = InitAuthFromConfig()
		assert.ErrorIs(t, err, domain.ErrMissingSecret)
	})

	t.Run("Failure_Invalid_Expiration", func(t *testing.T) {
		err := os.Setenv("SECRET_KEY", "test_secret")
		if err != nil {
			log.Error(err)
			return
		}
		err = os.Setenv("JWT_EXP", "invalid")
		if err != nil {
			log.Error(err)
			return
		}

		err = InitAuthFromConfig()
		assert.ErrorIs(t, err, domain.ErrInitAuthConfig)
	})
}

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success_Valid_Token", func(t *testing.T) {
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

	t.Run("Failure_Missing_Token", func(t *testing.T) {
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

	t.Run("Failure_Invalid_Token", func(t *testing.T) {
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
	t.Run("Success_Allowed_Role", func(t *testing.T) {
		role := "employee"
		assert.True(t, CheckRole(&role))
	})

	t.Run("Failure_Not_Allowed_Role", func(t *testing.T) {
		role := "admin"
		assert.False(t, CheckRole(&role))
	})

	t.Run("Failure_Nil_Role", func(t *testing.T) {
		assert.False(t, CheckRole(nil))
	})
}

func TestRequireRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success_Allowed_Role", func(t *testing.T) {
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

	t.Run("Failure_Not_Allowed_Role", func(t *testing.T) {
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
	t.Run("Success", func(t *testing.T) {
		err := os.Setenv("SECRET_KEY", "test_secret")
		if err != nil {
			log.Error(err)
			return
		}
		err = os.Setenv("JWT_EXP", "1h")
		if err != nil {
			log.Error(err)
			return
		}
		err = InitAuthFromConfig()
		if err != nil {
			log.Error(err)
			return
		}

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
