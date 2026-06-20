package api

import (
	"github.com/QEDQCD/cc-go/internal/bridge"
	"github.com/QEDQCD/cc-go/internal/config"
	"github.com/QEDQCD/cc-go/internal/server/auth"
	"github.com/QEDQCD/cc-go/internal/server/ws"
	"github.com/QEDQCD/cc-go/internal/store"
	"github.com/QEDQCD/cc-go/internal/wechat"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, cfg *config.Config, st *store.Store, br *bridge.Bridge, wc *wechat.Client, hub *ws.Hub, authMgr *auth.Manager) {
	api := r.Group("/api/v1")
	authMgr.RegisterRoutes(api)

	protected := api.Group("")
	protected.Use(authMgr.Middleware())
	registerWechatRoutes(protected, cfg, wc, br)
	registerWechatBotRoutes(protected, wc, br)
	registerClaudeRoutes(protected, st, br)
	registerSessionRoutes(protected, st, br)
	registerPermissionRoutes(protected, br)
	registerPushRoutes(protected, cfg)
	registerSettingsRoutes(protected, cfg)
	r.GET("/ws/events", authMgr.Middleware(), hub.HandleWS)
}