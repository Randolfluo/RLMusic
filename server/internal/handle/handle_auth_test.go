package handle_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	g "server/internal/global"
	"server/internal/handle"
	"server/internal/middleware"
	"server/internal/model"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// SetupRouter 初始化路由和数据库
func SetupRouter() (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Setup Memory DB
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&model.User{})

	// Inject DB
	r.Use(middleware.WithGormDB(db))

	// Mock Config for JWT
	g.Conf = &g.Config{}
	g.Conf.JWT.Secret = "test_secret"
	g.Conf.JWT.Issuer = "test_issuer"
	g.Conf.JWT.Expire = 24

	// Register Routes
	userAuthAPI := &handle.UserAuth{}
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", userAuthAPI.Register)
		auth.POST("/login", userAuthAPI.Login)
	}

	// Register Auth Routes
	protected := r.Group("/api")
	protected.Use(func(c *gin.Context) {
		// Mock Auth Middleware for testing
		db := c.MustGet(g.CtxDB).(*gorm.DB)
		user, _ := model.GetUserAuthInfoByName(db, "newuser")
		if user != nil {
			c.Set(g.CtxUserAuth, user)
		}
		c.Next()
	})
	protected.DELETE("/user", userAuthAPI.DeleteUser)

	return r, db
}

func TestRegisterAndLogin(t *testing.T) {
	router, _ := SetupRouter()

	// 1. Test Register
	t.Run("Register Success", func(t *testing.T) {
		reqBody := handle.RegisterReq{
			Username: "newuser",
			Password: "password123",
			Email:    "newuser@example.com",
		}
		jsonValue, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response handle.Response[any]
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, 1000, response.Code) // Success Code
	})

	// 2. Test Login
	t.Run("Login Success", func(t *testing.T) {
		reqBody := handle.LoginReq{
			Username: "newuser",
			Password: "password123",
		}
		jsonValue, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response handle.Response[handle.LoginVO]
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, 1000, response.Code)
		assert.NotEmpty(t, response.Data.Token)
		assert.Equal(t, "newuser", response.Data.Username)
	})

	// 3. Test Login Failed
	t.Run("Login Wrong Password", func(t *testing.T) {
		reqBody := handle.LoginReq{
			Username: "newuser",
			Password: "wrongpassword",
		}
		jsonValue, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var response handle.Response[any]
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.NotEqual(t, 0, response.Code)
	})

	// 4. Test Delete User
	t.Run("Delete User Success", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/user", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response handle.Response[any]
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, 1000, response.Code)
	})
}
