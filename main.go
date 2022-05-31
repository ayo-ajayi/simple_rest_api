package main

import (
	"database/sql"
	"github.com/ayo-ajayi/rest/controllers"
	"github.com/ayo-ajayi/rest/route"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	var db *sql.DB
	controllers.DBinit()
	defer db.Close()
	router := gin.Default()
	route.GroupChoice(router.Group("/choice"))
	route.GroupChoiceByID(router.Group("/choice"))
	log.Fatal(router.Run(":808"))
}
