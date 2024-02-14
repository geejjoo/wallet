package handler

import (
	"github.com/geejjoo/task/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api/v1")
	{
		lists := api.Group("/wallet")
		{
			lists.POST("/", h.createWallet)
			lists.POST("/:id/send", h.updateWallet)
			lists.GET("/:id/history", h.getHistory)
			lists.GET("/:id", h.getWallet)
		}
	}

	return router
}
