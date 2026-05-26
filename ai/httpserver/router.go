package httpserver

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vpa/quanlynhahang-backend/ai/config"
	"github.com/vpa/quanlynhahang-backend/ai/httpserver/handlers"
)

type Services struct {
	FileStore handlers.FileStore
	Vector    handlers.VectorStore
	Gemini    handlers.Gemini
	RAG       handlers.RAG
}

func NewRouter(cfg config.Config, svc *Services) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	corsCfg := cors.Config{
		AllowOrigins:     cfg.CorsOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Retry-After"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			origin = strings.TrimSpace(origin)
			for _, o := range cfg.CorsOrigins {
				if o == origin {
					return true
				}
			}
			return false
		},
	}
	r.Use(cors.New(corsCfg))

	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	api := r.Group("/api")
	{
		chat := handlers.NewChatHandler(svc.FileStore, svc.RAG, svc.Gemini)
		api.POST("/chat", chat.Chat)

		kn := handlers.NewKnowledgeHandler(svc.FileStore, svc.Vector, svc.Gemini)
		api.POST("/knowledge/menu/import", kn.ImportMenu)
		api.POST("/knowledge/menu/items", kn.AddMenuItem)
		api.GET("/knowledge/menu/items", kn.ListMenuItems)
		api.DELETE("/knowledge/menu/items/:item_id", kn.DeleteMenuItem)
		api.PUT("/knowledge/restaurant", kn.UpdateRestaurant)
		api.GET("/knowledge/restaurant", kn.GetRestaurant)
	}

	return r
}
