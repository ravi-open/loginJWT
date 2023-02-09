package main

import (
	"JWTAUTH/controllers"
	"JWTAUTH/initializers"
	"JWTAUTH/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVar()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate",middleware.Requireuth, controllers.Validate)
	r.Run()

}
