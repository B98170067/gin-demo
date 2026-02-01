package main

// import (
// 	"github.com/gin-gonic/gin"
// )

// func main() {
// 	r := gin.Default()

// 	r.GET("/ping", func(c *gin.Context) {
// 		c.JSON(200, gin.H{
// 			"message": "pong",
// 		})
// 	})

// 	r.Run() // :8080
// }

import (
	config "gin-demo/internal/configs"
	router "gin-demo/internal/routes"
)

func main() {
	db := config.InitDB()
	r := router.SetupRouter(db)
	r.Run(":8080")
}
