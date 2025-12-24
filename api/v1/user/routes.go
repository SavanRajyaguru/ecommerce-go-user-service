package user

import (
	"github.com/SavanRajyaguru/ecommerce-go-user-service/internal/auth"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup) {
	handler := NewUserHandler()

	userGroup := router.Group("/user")
	{
		userGroup.POST("/register", handler.Register)
		userGroup.POST("/login", handler.Login)

		protected := userGroup.Group("")
		protected.Use(auth.AuthMiddleware())
		{
			protected.GET("/profile", handler.GetProfile)
		}
	}
}
