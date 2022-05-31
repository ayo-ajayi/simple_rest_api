package route

import (
	"github.com/ayo-ajayi/rest/controllers"
	"github.com/gin-gonic/gin"
)

var GroupChoice = func(r *gin.RouterGroup) {
	p := r.Group("/")
	p.POST("/", controllers.PostChoice)
	p.GET("/", controllers.GetChoice)
}
var GroupChoiceByID = func(r *gin.RouterGroup) {
	p := r.Group("/:id")
	p.Use(controllers.CheckID)
	p.GET("/", controllers.GetChoiceByID)
	p.PUT("/", controllers.UpdateChoice)
	p.DELETE("/", controllers.DeleteChoice)

}
