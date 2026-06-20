package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/QEDQCD/cc-go/internal/config"
	"github.com/gin-gonic/gin"
)

func TestLoginAndSession(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := config.DefaultConfig()
	mgr := NewManager(cfg)
	r := gin.New()
	api := r.Group("/api/v1")
	mgr.RegisterRoutes(api)

	loginBody := `{"username":"admin","password":"admin123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("login status = %d, body = %s", w.Code, w.Body.String())
	}

	cookie := w.Result().Cookies()
	if len(cookie) == 0 || cookie[0].Name != SessionCookieName {
		t.Fatal("expected session cookie")
	}

	meReq := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	meReq.AddCookie(cookie[0])
	meW := httptest.NewRecorder()
	r.ServeHTTP(meW, meReq)
	if meW.Code != http.StatusOK {
		t.Fatalf("me status = %d", meW.Code)
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := config.DefaultConfig()
	mgr := NewManager(cfg)
	r := gin.New()
	api := r.Group("/api/v1")
	mgr.RegisterRoutes(api)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(`{"username":"admin","password":"wrong"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestMiddlewareBlocksUnauthenticated(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := config.DefaultConfig()
	mgr := NewManager(cfg)
	r := gin.New()
	protected := r.Group("/api/v1")
	protected.Use(mgr.Middleware())
	protected.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
