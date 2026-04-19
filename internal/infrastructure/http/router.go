package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sorrawichyooboon/go-audit-partition-purger/internal/infrastructure/http/handler"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/sorrawichyooboon/go-audit-partition-purger/docs"
)

func SetupRouter(auditHandler *handler.AuditHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		v1.POST("/audit-logs", auditHandler.TrackAuditLog)
	}

	return r
}
