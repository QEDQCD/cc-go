package api

import (
	"github.com/gin-gonic/gin"
	"github.com/QEDQCD/cc-go/internal/bridge"
	"github.com/QEDQCD/cc-go/internal/config"
	"github.com/QEDQCD/cc-go/internal/server/ws"
	"github.com/QEDQCD/cc-go/internal/store"
	"github.com/QEDQCD/cc-go/internal/wechat"
)

func RegisterRoutes(r *gin.Engine, cfg *config.Config, st *store.Store, br *bridge.Bridge, wc *wechat.Client, hub *ws.Hub) {
	api := r.Group("/api/v1")
	registerWechatRoutes(api, cfg, wc, br)
	registerWechatBotRoutes(api, wc, br)
	registerClaudeRoutes(api, st, br)
	registerSessionRoutes(api, st, br)
	registerPermissionRoutes(api, br)
	registerPushRoutes(api, cfg)
	registerSettingsRoutes(api, cfg)
	r.GET("/ws/events", hub.HandleWS)
}