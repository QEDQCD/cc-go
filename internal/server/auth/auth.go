package auth

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
	"time"

	"github.com/QEDQCD/cc-go/internal/config"
	"github.com/gin-gonic/gin"
)

const (
	SessionCookieName = "cc_go_session"
	SessionDuration   = 7 * 24 * time.Hour
)

type Manager struct {
	cfg      *config.Config
	sessions map[string]time.Time
	mu       sync.RWMutex
}

func NewManager(cfg *config.Config) *Manager {
	m := &Manager{
		cfg:      cfg,
		sessions: make(map[string]time.Time),
	}
	go m.cleanupLoop()
	return m
}

func (m *Manager) RegisterRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	auth.POST("/login", m.handleLogin)
	auth.POST("/logout", m.handleLogout)
	auth.GET("/me", m.handleMe)
}

func (m *Manager) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(SessionCookieName)
		if err != nil || !m.validate(token) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未登录或会话已过期"})
			return
		}
		c.Next()
	}
}

func (m *Manager) handleLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
		return
	}
	if req.Username != m.cfg.Auth.Username || req.Password != m.cfg.Auth.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	token, err := generateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建会话失败"})
		return
	}

	m.mu.Lock()
	m.sessions[token] = time.Now().Add(SessionDuration)
	m.mu.Unlock()

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(SessionCookieName, token, int(SessionDuration.Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"username": req.Username,
	})
}

func (m *Manager) handleLogout(c *gin.Context) {
	if token, err := c.Cookie(SessionCookieName); err == nil {
		m.revoke(token)
	}
	c.SetCookie(SessionCookieName, "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (m *Manager) handleMe(c *gin.Context) {
	token, err := c.Cookie(SessionCookieName)
	if err != nil || !m.validate(token) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": m.cfg.Auth.Username})
}

func (m *Manager) validate(token string) bool {
	if token == "" {
		return false
	}
	m.mu.RLock()
	expires, ok := m.sessions[token]
	m.mu.RUnlock()
	if !ok || time.Now().After(expires) {
		if ok {
			m.revoke(token)
		}
		return false
	}
	return true
}

func (m *Manager) revoke(token string) {
	m.mu.Lock()
	delete(m.sessions, token)
	m.mu.Unlock()
}

func (m *Manager) cleanupLoop() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		m.mu.Lock()
		for token, expires := range m.sessions {
			if now.After(expires) {
				delete(m.sessions, token)
			}
		}
		m.mu.Unlock()
	}
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
