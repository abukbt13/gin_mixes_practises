package router

import (
	"github.com/gin-gonic/gin"
	"practise/controllers"
	"practise/controllers/picture"
	"practise/middleware"
)

func MapUrl(r *gin.Engine) {
	r.POST("/user/register", controllers.CreateUser)
	r.POST("/user/login", controllers.LoginUser)
	r.DELETE("/user/:id", controllers.DeleteUser)

	authorized := r.Group("/user").Use(middleware.AuthMiddleware())
	authorized.GET("/", controllers.ProtectedEndpoint)
	authorized.POST("/picture", picture.SavePictures) // Pass function reference here
}
