package main

import (
	"database/sql"
	"github.com/ayo-ajayi/rest/controllers"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	var db *sql.DB
	controllers.DBinit()
	defer db.Close()
	router := gin.Default()
	router.POST("/choice", controllers.PostChoice)
	router.GET("/choice", controllers.GetChoice)
	router.GET("/choice/:id", controllers.CheckID, controllers.GetChoiceByID)
	router.PUT("/choice/:id", controllers.CheckID, controllers.UpdateChoice)
	router.DELETE("/choice/:id", controllers.CheckID, controllers.DeleteChoice)
	log.Fatal(router.Run(":808"))
}
