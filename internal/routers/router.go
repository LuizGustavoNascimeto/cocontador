package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func WelcomeMessage(c *gin.Context) {
	var initialData = []int{0, 1, 2, 3, 4, 5}
	c.JSON(http.StatusOK, initialData)
}

func SetupRouter() *gin.Engine {

	router := gin.Default()

	router.GET("/", WelcomeMessage)
	// router.GET("/", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, initialData)
	// })

	// router.POST("/createStudent", student.CreateStudent(initialData))

	return router
}
