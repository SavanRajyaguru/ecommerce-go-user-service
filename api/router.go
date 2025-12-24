package api

import (
	"github.com/SavanRajyaguru/ecommerce-go-user-service/api/middleware"
	"github.com/SavanRajyaguru/ecommerce-go-user-service/api/v1/user"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	v1 := r.Group("/v1")
	{
		user.RegisterRoutes(v1)
	}

	return r
}
