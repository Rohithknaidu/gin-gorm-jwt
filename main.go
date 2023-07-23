package main

import (
	"golang/jwt/controllers"
	"golang/jwt/initializers"
	"golang/jwt/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectTodb()
	initializers.SyncData()
}

func main() {
	r := gin.Default()
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate",middleware.RequiredAuth,controllers.Validate)
	r.Run()
}
